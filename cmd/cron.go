package cmd

import "github.com/spf13/cobra"

func init() {
	RootCmd.AddCommand(cronCmd)
}

var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "dns related commands",
}
