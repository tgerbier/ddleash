package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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
	fmt.Println("list: hello world!")
}
