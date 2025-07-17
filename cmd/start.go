package cmd

import (
	"gecho/server"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the echo server",
	Long: `The start command launches an echo server that listens for incoming TCP connections 
	and echoes back any received data to the client. This server is particularly useful for 
	testing network connectivity and validating communication paths in Kubernetes clusters.

	The server can be configured with various flags to customize its behavior, including 
	port number, interface binding, and connection handling options.`,
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("interface")
		port, _ := cmd.Flags().GetInt("port")
		server := server.NewEchoServer(host, port)
		err := server.Start()
		if err != nil {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().IntP("port", "p", 8080, "The port number to listen on")
	startCmd.Flags().StringP("interface", "i", "0.0.0.0", "The interface to bind to")
	startCmd.Flags().BoolP("verbose", "v", false, "Enable verbose logging")
}
