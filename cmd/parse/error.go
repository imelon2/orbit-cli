/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	ethlib "github.com/imelon2/orbit-cli/ethLib"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/spf13/cobra"
)

// errorCmd represents the error command
var ErrorCmd = &cobra.Command{
	Use:   "error",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("error called")

		// abiPath := utils.GetAbiDir()
		// _abi, err := os.ReadFile(abiPath)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// abi, err := ethlib.GetAbi(strings.NewReader(string(_abi)))
		// if err != nil {
		// 	log.Fatal(err)
		// }

		providerOrError, isProvider, err := prompt.SelectProviderOrBytes()
		if err != nil {
			log.Fatal(err)
		}

		var data []byte
		if isProvider {
			client, err := ethclient.Dial(providerOrError)
			if err != nil {
				log.Fatal(err)
			}

			hash, err := prompt.EnterTransactionHash()
			if err != nil {
				log.Fatal(err)
			}

			txHash := ethlib.NewTransactionHash(client, hash)
			tx, receipt, sender, err := txHash.GetTransactionSender()
			if err != nil {
				log.Fatal(err)
			}

			callMsg := ethereum.CallMsg{
				From:      *sender,
				To:        tx.To(),
				Gas:       tx.Gas(),
				GasPrice:  tx.GasPrice(),
				Data:      tx.Data(),
				GasFeeCap: tx.GasFeeCap(),
				GasTipCap: tx.GasTipCap(),
				Value:     tx.Value(),
			}

			bytes, err := client.CallContract(context.Background(), callMsg, receipt.BlockNumber)
			if err != nil {
				log.Fatal(err)
			}

			data = bytes
		} else {
			providerOrError = providerOrError[2:]
			ErrorBytes, err := hex.DecodeString(providerOrError)
			if err != nil {
				log.Fatalf("fail string calldata decode to hex: %v", err)
			}
			data = ErrorBytes
		}

		fmt.Print(data)
	},
}

func init() {
}
