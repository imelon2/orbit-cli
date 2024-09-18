/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
)

// eventCmd represents the event command
var EventCmd = &cobra.Command{
	Use:   "event",
	Short: "Parse calldata by transaction hash",
	Run: func(cmd *cobra.Command, args []string) {
		abiPath := utils.GetAbiDir()
		_abi, _ := os.ReadFile(abiPath)
		parsedABI, err := abi.JSON(strings.NewReader(string(_abi)))

		if err != nil {
			log.Fatalf("Failed to get ABI: %v", err)
		}
		provider, err := prompt.SelectProvider()
		if err != nil {
			log.Fatalf("Failed to get provider: %v", err)
		}
		hash, err := prompt.EnterTransactionHash()
		if err != nil {
			log.Fatalf("Failed to get hash: %v", err)
		}

		client := utils.GetClient(provider)
		tx, err := client.TransactionReceipt(context.Background(), txHash)
	},
}

func init() {
}
