package main

import (
	"log"

	"github.com/DataDog/datadog-go/statsd"

	"github.com/MattHauglustaine/ddleash"
)

const (
	parallelFetches = 20

	dogstatsdUrl = "dogstatsd:8125"
)

func initDDLeash() (*ddleash.DDLeash, error) {
	leash, err := ddleash.New(ddleash.Config{
		Team:     "XXX",
		Username: "YYY",
		Password: "ZZZ",
	})
	if err != nil {
		return nil, err
	}

	if err := leash.Login(); err != nil {
		return nil, err
	}

	return leash, nil
}

func produceMetrics(
	leash *ddleash.DDLeash,
	names chan<- string,
	done chan<- bool,
	errors chan<- error,
) {
	fetchedNames, err := leash.FetchAllMetricNames()
	if err != nil {
		errors <- err
	}

	for _, name := range fetchedNames {
		names <- name
	}
	done <- true
}

func consumeMetrics(
	leash *ddleash.DDLeash,
	names <-chan string,
	metricNumContexts chan<- int,
	errors chan<- error,
) {
	for {
		name := <-names
		hostsTags, err := leash.FetchMetricHostsTags(name)
		if err != nil {
			errors <- err
			break
		}

		metricNumContexts <- hostsTags.NumContexts
	}
}

func computeContextsSum(leash *ddleash.DDLeash) (int, error) {
	sum := 0

	names := make(chan string)
	metricNumContexts := make(chan int)
	done := make(chan bool)
	errors := make(chan error)

	go produceMetrics(leash, names, done, errors)

	for i := 0; i < parallelFetches; i++ {
		go consumeMetrics(leash, names, metricNumContexts, errors)
	}

	for {
		select {
		case numContexts := <-metricNumContexts:
			sum += numContexts
		case err := <-errors:
			return 0, err
		case <-done:
			return sum, nil
		}
	}
}

func main() {
	// Init DDLeash & our DD StatsD connection
	leash, err := initDDLeash()
	if err != nil {
		log.Fatal(err)
	}

	statsdClient, err := statsd.New(dogstatsdUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Compute the number of contexts on Datadog
	contextsSum, err := computeContextsSum(leash)
	if err != nil {
		log.Fatal(err)
	}

	// Send this sum to Datadog
	err = statsdClient.Gauge("foo.bar.baz", float64(contextsSum), nil, 1)
	if err != nil {
		log.Fatal(err)
	}
}
