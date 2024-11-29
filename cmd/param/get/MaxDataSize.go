/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	arbnetwork "github.com/imelon2/orbit-cli/arbNetwork"
	"github.com/imelon2/orbit-cli/common/logs"
	"github.com/imelon2/orbit-cli/common/utils"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/spf13/cobra"
)

// MaxDataSizeCmd represents the MaxDataSize command
var MaxDataSizeCmd = &cobra.Command{
	Use:   "MaxDataSize",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		MaxDataSize, err := sequencerInbox.MaxDataSize(Callopts)
		if err != nil {
			log.Fatal(err)
		}
		logs.PrintFromatter(utils.ConvertBytesToHex(map[string]*big.Int{
			"MaxDataSize": MaxDataSize,
		}))

	},
}
