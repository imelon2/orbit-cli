/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethLib "github.com/imelon2/orbit-cli/ethLib"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/solgen/go/precompilesgen"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var AccountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Print network account from ArbAggregator",
	Run: func(cmd *cobra.Command, args []string) {
		provider, err := prompt.SelectProvider()
		if err != nil {
			log.Fatal(err)
		}

		client, err := ethclient.Dial(provider)
		if err != nil {
			log.Fatal(err)
		}

		ArbAggregator, err := precompilesgen.NewArbAggregator(types.ArbAggregatorAddress, client)
		if err != nil {
			log.Fatal(err)
		}

		ArbOwnerPublic, err := precompilesgen.NewArbOwnerPublic(types.ArbOwnerPublicAddress, client)
		if err != nil {
			log.Fatal(err)
		}

		ArbGasInfo, err := precompilesgen.NewArbGasInfo(types.ArbGasInfoAddress, client)
		if err != nil {
			log.Fatal(err)
		}

		owners, err := ArbOwnerPublic.GetAllChainOwners(ethLib.Callopts)

		networkFees := make([]common.Address, 0)
		if networkFee, err := ArbOwnerPublic.GetNetworkFeeAccount(ethLib.Callopts); err == nil {
			networkFees = append(networkFees, networkFee)
		}

		infraFeeAccounts := make([]common.Address, 0)
		if infraFeeAccount, err := ArbOwnerPublic.GetInfraFeeAccount(ethLib.Callopts); err == nil {
			infraFeeAccounts = append(infraFeeAccounts, infraFeeAccount)
		}

		l1RewardRecipients := make([]common.Address, 0)
		if l1RewardRecipient, err := ArbGasInfo.GetL1RewardRecipient(ethLib.Callopts); err == nil {
			l1RewardRecipients = append(l1RewardRecipients, l1RewardRecipient)
		}

		batchPosters, err := ArbAggregator.GetBatchPosters(ethLib.Callopts)

		feeCollectors := make([]common.Address, 0)
		for _, poster := range batchPosters {
			feeCollector, err := ArbAggregator.GetFeeCollector(ethLib.Callopts, poster)
			if err != nil {
				log.Fatal(err)
			}

			feeCollectors = append(feeCollectors, feeCollector)
		}

		length := utils.SafeGetLongestArray(owners, networkFees, infraFeeAccounts, l1RewardRecipients, batchPosters, feeCollectors)

		data := make([][]string, 0)
		for i := 0; i < length; i++ {
			data = append(data, []string{
				"",
				utils.SafeGetAddressString(owners, i),
				utils.SafeGetAddressString(networkFees, i),
				utils.SafeGetAddressString(infraFeeAccounts, i),
				utils.SafeGetAddressString(l1RewardRecipients, i),
				utils.SafeGetAddressString(batchPosters, i),
				utils.SafeGetAddressString(feeCollectors, i),
			})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Index", "Owners", "Network Fee", "Infra Fee", "l1 Rewarder", "Poster", "Poster Fee Collector"})
		table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlueColor},
			tablewriter.Colors{tablewriter.FgBlueColor, tablewriter.Bold},
			tablewriter.Colors{tablewriter.FgBlueColor, tablewriter.Bold},
			tablewriter.Colors{tablewriter.FgBlueColor, tablewriter.Bold},
			tablewriter.Colors{tablewriter.FgBlueColor, tablewriter.Bold},
			tablewriter.Colors{tablewriter.FgBlueColor, tablewriter.Bold},
			tablewriter.Colors{tablewriter.FgBlueColor, tablewriter.Bold},
		)
		table.AppendBulk(data) // Add Bulk Data
		table.Render()
	},
}

func init() {
}
