/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"runtime"

	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
)

// parseCmd represents the parse command
var ParseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse Calldata, Event, Retrya",
	Run: func(cmd *cobra.Command, args []string) {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			log.Fatal("bad path")
		}

		root := utils.GetRootDir(filename)
		selected, err := prompt.SelectCommand(root)

		if err != nil {
			log.Fatal(err)
		}

		nextCmd, _, err := cmd.Find([]string{selected})
		if err != nil {
			log.Fatal(err)
		}

		nextCmd.Run(nextCmd, args)
	},
}

func init() {
	ParseCmd.AddCommand(CalldataCmd)
	ParseCmd.AddCommand(EventCmd)
	ParseCmd.AddCommand(EventCmd)
	ParseCmd.AddCommand(AliasCmd)
	ParseCmd.AddCommand(RetryableCmd)
}
