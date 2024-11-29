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

// getCmd represents the get command
var GetCmd = &cobra.Command{
	Use:   "get",
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
	GetCmd.AddCommand(MaxTimeVariationCmd)
	GetCmd.AddCommand(RBlockPeriodCmd)
	GetCmd.AddCommand(MessageCountCmd)
	GetCmd.AddCommand(MaxDataSizeCmd)
}
