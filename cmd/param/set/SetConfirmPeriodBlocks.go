/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	arbnetwork "github.com/imelon2/orbit-cli/arbNetwork"
	"github.com/imelon2/orbit-cli/common/logs"
	"github.com/imelon2/orbit-cli/common/tx"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/spf13/cobra"
)

// SetConfirmPeriodBlocksCmd represents the SetConfirmPeriodBlocks command
var SetConfirmPeriodBlocksCmd = &cobra.Command{
	Use:   "SetConfirmPeriodBlocks",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		chains, err := prompt.SelectChains()
		if err != nil {
			log.Fatal(err)
		}

		parent, child, err := prompt.SelectCrossChainProviders(chains)
		if err != nil {
			log.Fatal(err)
		}

		parentClient, err := ethclient.Dial(parent)
		if err != nil {
			log.Fatal(err)
		}

		childClient, err := ethclient.Dial(child)
		if err != nil {
			log.Fatal(err)
		}

		_, ks, account, err := prompt.SelectWalletForSign()
		if err != nil {
			log.Fatal(err)
		}

		auto, err := tx.GetAuthByKeystore(ks, *account, parentClient)
		if err != nil {
			log.Fatal(err)
		}

		network, err := arbnetwork.GetNetworkInfo(childClient)
		if err != nil {
			log.Fatal(err)
		}
		UpgradeExecutor, err := network.NewUpgradeExecutor(parentClient)
		if err != nil {
			log.Fatal(err)
		}

		Callopts := &bind.CallOpts{
			Pending: false,
			Context: nil,
		}
		confirmPeriodBlocks, err := UpgradeExecutor.RollupUserLogic.ConfirmPeriodBlocks(Callopts)
		if err != nil {
			log.Fatal(err)
		}

		newConfirmPeriod, err := prompt.EnterInt(0, fmt.Sprintf("new ConfirmPeriodBlocks(block)"+logs.GrayString("(current value: %d)"), confirmPeriodBlocks))
		if err != nil {
			log.Fatal(err)
		}

		response, err := UpgradeExecutor.SetConfirmPeriodBlocks(auto, uint64(*newConfirmPeriod))
		if err != nil {
			log.Fatal(err)
		}

		receipt, err := bind.WaitMined(context.Background(), parentClient, response)
		if err != nil {
			log.Fatal(err)
		}

		logs.PrintReceiptFromatter(receipt)
	},
}
