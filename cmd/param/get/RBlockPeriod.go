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

// RBlockPeriodCmd represents the RBlockPeriod command
var RBlockPeriodCmd = &cobra.Command{
	Use:   "RBlockPeriod",
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

		RollupCore, err := network.NewRollupCore(parentClient)
		if err != nil {
			log.Fatal(err)
		}

		Callopts := &bind.CallOpts{
			Pending: false,
			Context: nil,
		}
		minimumAssertionPeriod, err := RollupCore.RollupCoreCaller.MinimumAssertionPeriod(Callopts)
		if err != nil {
			log.Fatal(err)
		}
		confirmPeriodBlocks, err := RollupCore.RollupCoreCaller.ConfirmPeriodBlocks(Callopts)
		if err != nil {
			log.Fatal(err)
		}
		maxTimeVariation := arbnetwork.RBlockPeriod{
			MinimumAssertionPeriod: minimumAssertionPeriod,
			ConfirmPeriodBlocks:    big.NewInt(int64(confirmPeriodBlocks)),
		}
		logs.PrintFromatter(utils.ConvertBytesToHex(maxTimeVariation))
	},
}
