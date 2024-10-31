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
	"github.com/ethereum/go-ethereum/ethclient"
	arbnetwork "github.com/imelon2/orbit-cli/arbNetwork"
	"github.com/imelon2/orbit-cli/common/logs"
	"github.com/imelon2/orbit-cli/common/tx"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/solgen/go/bridgegen"
	"github.com/spf13/cobra"
)

// SetMaxTimeVariationCmd represents the SetMaxTimeVariation command
var SetMaxTimeVariationCmd = &cobra.Command{
	Use:   "SetMaxTimeVariation",
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
		delayBlocks, futureBlocks, delaySeconds, futureSeconds, err := UpgradeExecutor.SequencerInbox.MaxTimeVariation(Callopts)
		if err != nil {
			log.Fatal(err)
		}

		newDelayBlocks, err := prompt.EnterInt(0, fmt.Sprintf("new delayBlocks"+logs.GrayString("(current value: %d)"), delayBlocks))
		if err != nil {
			log.Fatal(err)
		}
		newFutureBlocks, err := prompt.EnterInt(0, fmt.Sprintf("new futureBlocks"+logs.GrayString("(current value: %d)"), futureBlocks))
		if err != nil {
			log.Fatal(err)
		}
		newDelaySeconds, err := prompt.EnterInt(0, fmt.Sprintf("new delaySeconds"+logs.GrayString("(current value: %d)"), delaySeconds))
		if err != nil {
			log.Fatal(err)
		}
		newFutureSeconds, err := prompt.EnterInt(0, fmt.Sprintf("new futureSeconds"+logs.GrayString("(current value: %d)"), futureSeconds))
		if err != nil {
			log.Fatal(err)
		}

		response, err := UpgradeExecutor.SetMaxTimeVariation(auto, bridgegen.ISequencerInboxMaxTimeVariation{
			DelayBlocks:   big.NewInt(int64(*newDelayBlocks)),
			FutureBlocks:  big.NewInt(int64(*newFutureBlocks)),
			DelaySeconds:  big.NewInt(int64(*newDelaySeconds)),
			FutureSeconds: big.NewInt(int64(*newFutureSeconds)),
		})
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
