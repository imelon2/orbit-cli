/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	arblib "github.com/imelon2/orbit-cli/arbLib"
	"github.com/imelon2/orbit-cli/contractgen"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/spf13/cobra"
)

// barchCmd represents the barch command
var BatchCmd = &cobra.Command{
	Use:   "batch",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		providers, err := prompt.SelectProviders()
		if err != nil {
			log.Fatal(err)
		}

		l2ProviderUrl := providers[1]
		l3ProviderUrl := providers[2]

		parentClient, err := ethclient.Dial(l2ProviderUrl)
		if err != nil {
			log.Fatal(err)
		}
		childClient, err := ethclient.Dial(l3ProviderUrl)
		if err != nil {
			log.Fatal(err)
		}
		network, err := contractgen.GetNetworkInfo(childClient)
		if err != nil {
			log.Fatal(err)
		}

		arb := arblib.NewContractLib(&network, parentClient)
		sequencerInbox, err := arb.NewSequencerInbox()
		if err != nil {
			log.Fatal(err)
		}
		erc20Bridge, err := arb.NewERC20Bridge()
		if err != nil {
			log.Fatal(err)
		}
		count, err := erc20Bridge.SequencerMessageCount()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Max Batch Count %d\n", count)

		events, err := sequencerInbox.GetBatchData(count)
		if err != nil {
			log.Fatal(err)
		}
		// sequencerInbox.GetBatchData(big.NewInt(7))
	},
}

func init() {
}
