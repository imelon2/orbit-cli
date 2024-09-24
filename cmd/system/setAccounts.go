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
	ethlib "github.com/imelon2/orbit-cli/ethLib"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/solgen/go/precompilesgen"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
)

// setAccountsCmd represents the setAccounts command
var SetAccountsCmd = &cobra.Command{
	Use:   "setAccounts",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		provider, err := prompt.SelectProvider()
		if err != nil {
			log.Fatal(err)
		}
		client, err := ethclient.Dial(provider)
		if err != nil {
			log.Fatal(err)
		}

		_, ks, account, err := prompt.SelectWalletForSign()

		if err != nil {
			log.Fatal(err)
		}

		ArbOwner, err := precompilesgen.NewArbOwner(types.ArbOwnerAddress, client)

		if err != nil {
			log.Fatal(err)
		}

		auth := ethlib.GenerateAuth(client, ks, account)
		signedTx, err := ArbOwner.SetInfraFeeAccount(auth, common.HexToAddress("0xd7464B89f726EcE721B4fcB7a90732387b23E6fc"))

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
