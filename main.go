/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"NextDomain-Utils/cmd"
	"NextDomain-Utils/utils"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(utils.MyHandler{})
	slog.SetDefault(logger)
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
