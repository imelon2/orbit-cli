/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	arbnetwork "github.com/imelon2/orbit-cli/arbNetwork"
	"github.com/imelon2/orbit-cli/common/logs"
	"github.com/imelon2/orbit-cli/common/utils"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/spf13/cobra"
)

// maxTimeVariationCmd represents the maxTimeVariation command
var MaxTimeVariationCmd = &cobra.Command{
	Use:   "maxTimeVariation",
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

		network, err := arbnetwork.GetNetworkInfo(childClient)
		if err != nil {
			log.Fatal(err)
		}

		sequencerInbox, err := network.NewSequencerInbox(parentClient)
		if err != nil {
			log.Fatal(err)
		}

		Callopts := &bind.CallOpts{
			Pending: false,
			Context: nil,
		}
		delayBlocks, futureBlocks, delaySeconds, futureSeconds, err := sequencerInbox.SequencerInboxCaller.MaxTimeVariation(Callopts)
		if err != nil {
			log.Fatal(err)
		}
		maxTimeVariation := arbnetwork.MaxTimeVariation{
			DelayBlocks:   delayBlocks,
			FutureBlocks:  futureBlocks,
			DelaySeconds:  delaySeconds,
			FutureSeconds: futureSeconds,
		}
		logs.PrintFromatter(utils.ConvertBytesToHex(maxTimeVariation))
	},
}
