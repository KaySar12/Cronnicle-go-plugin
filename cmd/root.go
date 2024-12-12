/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "NextDomain-Utils",
	Short: "Cronicle Plugin written in golang that interact with NextDomain",
}

func init() {
}
