/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/imelon2/orbit-cli/common/logs"
	"github.com/imelon2/orbit-cli/common/path"
	"github.com/spf13/cobra"
)

type KeystoreTag struct {
	Tag     string
	Address common.Address
}

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		keystorePath := path.GetKeystoreDir("")

		files, err := os.ReadDir(keystorePath)
		if err != nil {
			log.Fatal(err)
		}

		if len(files) == 0 {
			fmt.Print("\nNo keystore was created.\n")
			fmt.Print("Execute the " + logs.BoldString("{ orbit-cli account new }") + " to create a keystore.\n")
			return
		}

		for i, file := range files {
			pathTag := path.GetKeystoreDir(file.Name())
			ks := keystore.NewKeyStore(pathTag, keystore.StandardScryptN, keystore.StandardScryptP)

			accounts := ks.Accounts()
			fmt.Printf("[%d] %s | %s\n", i, accounts[0].Address.Hex(), file.Name())
		}
	},
}
