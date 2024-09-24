/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
)

// aliasCmd represents the alias command
var AliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Print aliasing address ",
	Run: func(cmd *cobra.Command, args []string) {
		isParentToChild, err := prompt.SelectChainTo()
		if err != nil {
			log.Fatalf("Failed to get SelectChainTo: %v", err)
		}
		address, err := prompt.EnterAddress()
		if err != nil {
			log.Fatalf("Failed to get EnterAddress: %v", err)
		}

		aliasing := utils.Alias(common.HexToAddress(address), isParentToChild)

		fmt.Printf("\n\nAlising Address : %s\n\n", aliasing.Hex())
	},
}

func init() {
}
