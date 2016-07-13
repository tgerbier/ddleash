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
	Short: "List all known metrics",
	Long:  "List all metric names known by Datadog, sorted alphabetically.",
	Run:   runListCmd,
}

func init() {
	RootCmd.AddCommand(ListCmd)
}

func runListCmd(cmd *cobra.Command, args []string) {
	item := "metrics"
	listArgs := []string{}
	if len(args) > 0 {
		item = args[0]
		listArgs = args[1:]
	}

	listFunc, ok := map[string]func(*ddleash.Client, []string) error{
		"metrics": listMetrics,
	}[item]

	if !ok {
		fmt.Printf("Unknown item to list: %q\n", item)
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
