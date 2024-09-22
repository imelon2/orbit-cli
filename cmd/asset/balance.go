package cmd

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	prompt "github.com/imelon2/orbit-cli/prompt"
	"github.com/spf13/cobra"
)

var BalanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		getAccountBalance()
	},
}

func init() {

}

func getAccountBalance() {
	selectedWallet, err := prompt.SelectWallet()
	if err != nil {
		log.Fatal(err)
	}

	provider, err := prompt.SelectProvider()
	if err != nil {
		log.Fatal(err)
	}
	client, err := ethclient.Dial(provider)
	if err != nil {
		log.Fatal(err)
	}

	account := common.HexToAddress(selectedWallet)
	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	if err != nil {
		log.Fatal(err)
	}

	ethValue2 := new(big.Float).SetInt(pendingBalance)
	ethValue2.Quo(ethValue2, big.NewFloat(1e18))

	fmt.Printf("\nBalance: %.18f ETH\n", ethValue2)
}
