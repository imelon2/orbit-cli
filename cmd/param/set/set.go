/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"runtime"

	"github.com/imelon2/orbit-cli/prompt"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var SetCmd = &cobra.Command{
	Use:   "set",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			log.Fatal("bad path")
		}

		prompt.RootCmdNavigation(filename, cmd, args)
	},
}

func init() {
	SetCmd.AddCommand(SetMinimumAssertionPeriodCmd)
	SetCmd.AddCommand(SetConfirmPeriodBlocksCmd)
	SetCmd.AddCommand(SetMaxTimeVariationCmd)
	SetCmd.AddCommand(SetSequencerReportedSubMessageCountCmd)
	SetCmd.AddCommand(SetL1PricePerUnitCmd)
}
