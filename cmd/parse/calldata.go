/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imelon2/orbit-cli/common/logs"
	"github.com/imelon2/orbit-cli/common/tx"
	"github.com/imelon2/orbit-cli/common/utils"
	"github.com/imelon2/orbit-cli/parse"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/spf13/cobra"
)

// calldataCmd represents the calldata command
var CalldataCmd = &cobra.Command{
	Use:   "calldata",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		var calldataBytes []byte

		chains, err := prompt.SelectChains(prompt.LAST_PROVIDER_STRING, prompt.LAST_CALLDATA_STRING)
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
			transaction, _, err := txLib.GetTransactionByHash()
			if err != nil {
				log.Fatal(err)
			}
			calldataBytes = transaction.Data()
			// Enter Bytes Calldata
		} else if chains == prompt.LAST_CALLDATA_STRING {
			calldata, err := prompt.EnterBytes()
			if err != nil {
				log.Fatal(err)
			}
			calldata = utils.Unhexlify(calldata)
			calldataBytes, err = hex.DecodeString(calldata)
			if err != nil {
				log.Fatalf("fail string calldata decode to hex: %v", err)
			}
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
			transaction, _, err := txLib.GetTransactionByHash()
			if err != nil {
				log.Fatal(err)
			}
			calldataBytes = transaction.Data()
		}

		parse, err := parse.NewParse()
		if err != nil {
			log.Fatal(err)
		}

		resultJson, err := parse.ParseCalldata(calldataBytes)
		if err != nil {
			log.Fatal(err)
		}
		logs.PrintFromatter(resultJson)
	},
}

func init() {

}
