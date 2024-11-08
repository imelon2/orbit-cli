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
	"github.com/spf13/cobra"
)

// SetSequencerReportedSubMessageCountCmd represents the SetSequencerReportedSubMessageCount command
var SetSequencerReportedSubMessageCountCmd = &cobra.Command{
	Use:   "SetSequencerReportedSubMessageCount",
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
		sequencerReportedSubMessageCount, err := UpgradeExecutor.Bridge.SequencerReportedSubMessageCount(Callopts)
		if err != nil {
			log.Fatal(err)
		}

		newSequencerReportedSubMessageCount, err := prompt.EnterInt(0, fmt.Sprintf("new SequencerReportedSubMessageCount "+logs.GrayString("(current value: %d)"), sequencerReportedSubMessageCount))
		if err != nil {
			log.Fatal(err)
		}

		response, err := UpgradeExecutor.SetSequencerReportedSubMessageCount(auto, big.NewInt(int64(*newSequencerReportedSubMessageCount)))
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
