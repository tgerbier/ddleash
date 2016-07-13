package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base commad when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ddleash",
	Short: "Manipulate Datadog metrics to keep in on a leash.",
	Long: `ddleash exposes simple commands to show, manipulate and
monitor Datadog metrics.

Complete documentation is available at
https://github.com/MattHauglustaine/ddleash.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(assertDatadogAccountConfig)

	RootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		"config file (default is $HOME/.ddleash.yaml)",
	)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".ddleash") // name of config file (without extension)
	viper.AddConfigPath("$HOME")    // adding home directory as first search path
	viper.AutomaticEnv()            // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// assertDatadogAccountConfig returns an error if no Datadog account
// info (team, user, and password) are set.
func assertDatadogAccountConfig() {
	for _, subKey := range []string{"team", "user", "password"} {
		key := "datadog." + subKey
		if !viper.IsSet(key) {
			fmt.Println(fmt.Sprintf(
				"Invalid configuration: %q must be set.",
				key,
			))
			os.Exit(-1)
		}

	}
}
