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

// searchCmd represents the search command
var SearchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
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
	SearchCmd.AddCommand(BatchCmd)
	SearchCmd.AddCommand(NodeConfirmedCmd)
	SearchCmd.AddCommand(NodeCreatedCmd)
}
