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

// systemCmd represents the system command
var SystemCmd = &cobra.Command{
	Use:   "system",
	Short: "Manage system network Gas, access account, vault",
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
	SystemCmd.AddCommand(GasCmd)
	SystemCmd.AddCommand(AccountsCmd)
	SystemCmd.AddCommand(SetAccountsCmd)
	SystemCmd.AddCommand(VaultCmd)
	SystemCmd.AddCommand(SetGasCmd)
}
