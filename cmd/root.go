package cmd

import (
	"github.com/comoyi/steam-server-monitor/client"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "steam-server-monitor",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client.Start()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
