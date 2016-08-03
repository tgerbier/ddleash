package main

import (
	"fmt"
	"os"
	"strings"

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
	cobra.OnInitialize(assertConfig)

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
	viper.AddConfigPath(".")        // adding current directory as first search path
	viper.AutomaticEnv()            // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// viper returns an `unsupported config type ""` error
		// if it can't find a file. We just ignore it.
		// https://github.com/spf13/viper/issues/210
		if !strings.HasSuffix(err.Error(), `Type ""`) {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
}

// assertConfig returns an error if the mandatory config values are
// not set.
func assertConfig() {
	for _, key := range []string{
		"datadog.team",
		"datadog.user",
		"datadog.password",
		"dogstatsd.url",
	} {
		if !viper.IsSet(key) {
			fmt.Println(fmt.Sprintf(
				"Invalid configuration: %q must be set.",
				key,
			))
			os.Exit(-1)
		}

	}
}
