/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/imelon2/orbit-cli/common/logs"
	"github.com/imelon2/orbit-cli/common/tx"
	"github.com/imelon2/orbit-cli/parse"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/spf13/cobra"
)

// errorCmd represents the error command
var ErrorCmd = &cobra.Command{
	Use:   "error",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		var errorData *rpc.DataError

		chains, err := prompt.SelectChains(prompt.LAST_PROVIDER_STRING, prompt.LAST_CALLDATA_STRING)
		if err != nil {
			log.Fatal(err)
		}

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
			_, errorRpc, status, err := txLib.GetTransactionReturn()
			if err != nil {
				log.Fatal(err)
			}

			if *status {
				fmt.Printf("transaction hash %s is SUCCESS", hash.Hex())
				return
			}

			errorData = errorRpc
		} else if chains == prompt.LAST_CALLDATA_STRING {
			sError, err := prompt.EnterBytes()
			if err != nil {
				log.Fatal(err)
			}

			parse, err := parse.NewParse()
			if err != nil {
				log.Fatal(err)
			}

			errorJson, err := parse.ParseErrorByBytes(sError)
			if err != nil {
				log.Fatal(err)
			}
			logs.PrintFromatter(errorJson)
			return
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
			_, errorRpc, status, err := txLib.GetTransactionReturn()
			if err != nil {
				log.Fatal(err)
			}

			if *status {
				fmt.Printf("transaction hash %s is SUCCESS", hash.Hex())
				return
			}

			errorData = errorRpc
		}

		parse, err := parse.NewParse()
		if err != nil {
			log.Fatal(err)
		}

		errorJson, err := parse.ParseError(*errorData)
		if err != nil {
			log.Fatal(err)
		}

		logs.PrintFromatter(errorJson)
	},
}

func init() {
}
