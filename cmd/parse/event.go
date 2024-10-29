/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imelon2/orbit-cli/common/logs"
	"github.com/imelon2/orbit-cli/common/tx"
	"github.com/imelon2/orbit-cli/parse"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/spf13/cobra"
)

// eventCmd represents the event command
var EventCmd = &cobra.Command{
	Use:   "event",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		var txLogs []*types.Log

		chains, err := prompt.SelectChains(prompt.LAST_PROVIDER_STRING)
		if err != nil {
			log.Fatal(err)
		}

		// Enter Provider
		if chains == prompt.LAST_PROVIDER_STRING {
			url, err := prompt.EnterProviderUrl()
			if err != nil {
				log.Fatal(err)
			}

			client, err := ethclient.Dial(url)
			if err != nil {
				log.Fatal(err)
			}

			sHash, err := prompt.EnterTransactionHash()
			if err != nil {
				log.Fatal(err)
			}

			hash := common.HexToHash(sHash)
			txLib := tx.NewTxLib(client, &hash)
			receipt, err := txLib.GetTransactionReceipt()
			if err != nil {
				log.Fatal(err)
			}

			txLogs = receipt.Logs
		} else {
			// Select Chain and Provider
			provider, err := prompt.SelectProviders(chains)
			if err != nil {
				log.Fatal(err)
			}

			client, err := ethclient.Dial(provider)
			if err != nil {
				log.Fatal(err)
			}

			sHash, err := prompt.EnterTransactionHash()
			if err != nil {
				log.Fatal(err)
			}

			hash := common.HexToHash(sHash)
			txLib := tx.NewTxLib(client, &hash)
			receipt, err := txLib.GetTransactionReceipt()
			if err != nil {
				log.Fatal(err)
			}

			txLogs = receipt.Logs
		}

		parse, err := parse.NewParse()
		if err != nil {
			log.Fatal(err)
		}

		resultJson, err := parse.ParseEvent(txLogs)
		if err != nil {
			log.Fatal(err)
		}
		logs.PrintFromatter(resultJson)
	},
}

func init() {
}
