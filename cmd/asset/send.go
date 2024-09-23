/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var SendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send ETH from select wallet",
	Run: func(cmd *cobra.Command, args []string) {

		value, err := prompt.EnterValue()
		to, err := prompt.EnterRecipient()
		toAddress := common.HexToAddress(to)

		wallet, account, err := prompt.SelectWalletForSign()

		if err != nil {
			log.Fatal(err)
		}

		provider, err := prompt.SelectProvider()
		if err != nil {
			log.Fatal(err)
		}

		client, err := ethclient.Dial(provider)

		nonce, err := client.PendingNonceAt(context.Background(), account.Address)
		gasLimit := uint64(21000)
		gasPrice, err := client.SuggestGasPrice(context.Background())

		tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil /* calldata */)
		chainID, err := client.NetworkID(context.Background())

		if err != nil {
			log.Fatal(err)
		}

		signedTx, err := wallet.SignTx(account, tx, chainID)
		if err != nil {
			log.Fatal(err)
		}

		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			fmt.Println("SendTransaction")
			log.Fatal(err)
		}

		txResponse, _, err := client.TransactionByHash(context.Background(), signedTx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print("\n\nTransaction Response: \n")
		utils.PrintPrettyJson(txResponse)

		fmt.Print("\n\nWait Mined Transaction ... \n\n")

		receipt, err := bind.WaitMined(context.Background(), client, signedTx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Transaction receipt: ")
		utils.PrintPrettyJson(receipt)
	},
}

func init() {
}
