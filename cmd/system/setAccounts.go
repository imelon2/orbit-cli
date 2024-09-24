/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/solgen/go/precompilesgen"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
)

// setAccountsCmd represents the setAccounts command
var SetAccountsCmd = &cobra.Command{
	Use:   "setAccounts",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		provider, err := prompt.SelectProvider()
		if err != nil {
			log.Fatal(err)
		}
		client := utils.GetClient(provider)

		_, ks, account, err := prompt.SelectWalletForSign()

		if err != nil {
			log.Fatal(err)
		}

		ArbOwner, err := precompilesgen.NewArbOwner(types.ArbOwnerAddress, client)
		if err != nil {
			log.Fatal(err)
		}
		// wallet.SignTx()
		chainID, err := client.NetworkID(context.Background())
		nonce, err := client.PendingNonceAt(context.Background(), account.Address)
		auth, err := bind.NewKeyStoreTransactorWithChainID(ks, account, chainID)
		auth.Nonce = big.NewInt(int64(nonce + 1))

		signedTx, err := ArbOwner.SetInfraFeeAccount(auth, common.HexToAddress("0x10012d9D7365bD937d5c28f786045D7C93EDc7eC"))
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
