/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
)

// retryableCmd represents the retryable command
var RetryableCmd = &cobra.Command{
	Use:   "retryable",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		providerOrBytes, isProvider, err := prompt.SelectProviderOrBytes()
		if err != nil {
			log.Fatal(err)
		}

		var data []byte
		if isProvider {
			fmt.Print("\n\nTODO\n\n")
			return
			client, err := ethclient.Dial(providerOrBytes)
			if err != nil {
				log.Fatal(err)
			}

			txHash, err := prompt.EnterTransactionHash()
			if err != nil {
				log.Fatalf("Failed to get txHash: %v", err)
			}

			tx, _, err := client.TransactionByHash(context.Background(), txHash)
			if err != nil {
				log.Fatalf("Failed to get TransactionByHash: %v", err)
			}
			data = tx.Data()

		} else {
			providerOrBytes = providerOrBytes[2:] // remove 0x
			calldataBytes, err := hex.DecodeString(providerOrBytes)
			if err != nil {
				log.Fatalf("failed decode bytes data: %v", err)
			}

			data = calldataBytes
		}

		retry := utils.ParseRetryableMessage(data)
		utils.PrintPrettyJson(retry)
	},
}

func init() {
}
