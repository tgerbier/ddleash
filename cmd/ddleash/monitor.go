package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/matryer/try.v1"

	"github.com/matthauglustaine/ddleash"
)

const (
	defaultNumFetchers = 20
	defaultWindow      = 3600
	maxRetries         = 5
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor a value (metrics-count, ...) ",
	Long: `Monitor an value by sending it to Datadog. For now, only the
"metrics-count" object can be monitored.`,
	Run: runMonitorCmd,
}

func init() {
	RootCmd.AddCommand(monitorCmd)
}

func runMonitorCmd(cmd *cobra.Command, args []string) {
	item := "metrics-count"
	monitorArgs := []string{}
	if len(args) > 0 {
		item = args[0]
		monitorArgs = args[1:]
	}

	monitorFunc, ok := map[string]func(*ddleash.Client, []string) error{
		"metrics-count": monitorMetricsCount,
	}[item]

	if !ok {
		fmt.Printf("Unknown value to monitor: %q\n", item)
		os.Exit(-1)
	}

	client := ddleash.New(ddleash.Account{
		Team:     viper.GetString("datadog.team"),
		User:     viper.GetString("datadog.user"),
		Password: viper.GetString("datadog.password"),
	})

	if err := monitorFunc(client, monitorArgs); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func monitorMetricsCount(client *ddleash.Client, _ []string) error {
	statsdClient, err := statsd.New(viper.GetString("dogstatsd.url"))
	if err != nil {
		return err
	}

	if err := client.Login(); err != nil {
		return err
	}

	errs := make(chan error)
	done := make(chan struct{})
	defer close(done)

	names := genMetricNames(done, errs, client)
	metricHostTags := fetchMetrics(
		done, errs, names, defaultNumFetchers, client,
	)
	out := sumMetricContexts(done, metricHostTags)

	select {
	case err := <-errs:
		return err
	case sum := <-out:
		// Send this sum to Datadog
		err = statsdClient.Gauge(
			"foo.bar.baz",
			float64(sum),
			nil,
			1,
		)
		if err != nil {
			return err
		}
		return nil
	}
}

func genMetricNames(
	done <-chan struct{},
	errs chan<- error,
	client *ddleash.Client,
) <-chan string {
	names := make(chan string)

	go func() {
		defer close(names)

		allNames, err := client.FetchAllMetricNames(defaultWindow)
		if err != nil {
			errs <- err
			return
		}

		for _, name := range allNames {
			select {
			case names <- name:
			case <-done:
				return
			}
		}
	}()

	return names
}

func fetcher(
	done <-chan struct{},
	errs chan<- error,
	names <-chan string,
	metricHostTags chan<- *ddleash.MetricHostsTags,
	client *ddleash.Client,
) {
	for name := range names {
		var metricHostTag *ddleash.MetricHostsTags

		err := try.Do(func(attempt int) (bool, error) {
			var err error
			metricHostTag, err = client.FetchMetricHostsTags(
				name, defaultWindow,
			)
			return attempt < maxRetries, err
		})

		if err != nil {
			errs <- err
			return
		}

		select {
		case metricHostTags <- metricHostTag:
		case <-done:
			return
		}
	}
}

func fetchMetrics(
	done <-chan struct{},
	errs chan<- error,
	names <-chan string,
	numFetchers int,
	client *ddleash.Client,
) <-chan *ddleash.MetricHostsTags {
	var wg sync.WaitGroup

	metricHostTags := make(chan *ddleash.MetricHostsTags)

	wg.Add(numFetchers)
	for i := 0; i < numFetchers; i++ {
		go func() {
			fetcher(done, errs, names, metricHostTags, client)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(metricHostTags)
	}()

	return metricHostTags
}

func sumMetricContexts(
	done <-chan struct{},
	metricHostTags <-chan *ddleash.MetricHostsTags,
) <-chan int {
	var sum int

	out := make(chan int)

	go func() {
		defer close(out)

		for metricHostTag := range metricHostTags {
			sum += metricHostTag.NumContexts
		}

		select {
		case out <- sum:
		case <-done:
			return
		}
	}()

	return out
}
