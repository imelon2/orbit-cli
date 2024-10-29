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

		cmdName, err := prompt.SelectNextCmd(filename)
		if err != nil {
			log.Fatal(err)
		}

		nextCmd, _, err := cmd.Find([]string{cmdName})
		if err != nil {
			log.Fatal(err)
		}
		nextCmd.Run(nextCmd, args)
	},
}

func init() {
	GetCmd.AddCommand(MaxTimeVariationCmd)
}
