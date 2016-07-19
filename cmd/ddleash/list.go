package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/matthauglustaine/ddleash"
)

const (
	window = 3600
)

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List a Datadog object (metric, ...)",
	Long: `List a Datadog object. For now, only the "metric" object can be
listed.`,
	Run: runListCmd,
}

func init() {
	RootCmd.AddCommand(ListCmd)
}

func runListCmd(cmd *cobra.Command, args []string) {
	object := "metrics"
	listArgs := []string{}
	if len(args) > 0 {
		object = args[0]
		listArgs = args[1:]
	}

	listFunc, ok := map[string]func(*ddleash.Client, []string) error{
		"metrics": listMetrics,
	}[object]

	if !ok {
		fmt.Printf("Unknown object to list: %q\n", object)
		os.Exit(-1)
	}

	client := ddleash.New(ddleash.Account{
		Team:     viper.GetString("datadog.team"),
		User:     viper.GetString("datadog.user"),
		Password: viper.GetString("datadog.password"),
	})

	if err := listFunc(client, listArgs); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func listMetrics(client *ddleash.Client, _ []string) error {
	if err := client.Login(); err != nil {
		return err
	}

	metrics, err := client.FetchAllMetricNames(window)
	if err != nil {
		return err
	}

	fmt.Println(strings.Join(metrics, "\n"))
	return nil
}
