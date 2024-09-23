/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Get account list",
	Run: func(cmd *cobra.Command, args []string) {
		path := utils.GetKeystoreDir()
		ks := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)

		accounts := ks.Accounts()
		if len(accounts) == 0 {
			fmt.Print("\nNo keystore was created.\n")
			fmt.Print("Execute the " + utils.BoldString("{ orbit-cli account new }") + " to create a keystore.\n")
			return
		}

		for i, wallet := range accounts {
			fmt.Printf("[%d] %s\n", i, wallet.Address.Hex())
		}
	},
}

func init() {
}
