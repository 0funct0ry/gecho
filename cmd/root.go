package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gecho",
	Short: "An echo server",
	Long:  `An echo server that can be used to test connectivity to a Kubernetes cluster over tcp.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
