/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imelon2/orbit-cli/common/logs"
	"github.com/imelon2/orbit-cli/common/utils"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/solgen/go/precompilesgen"
	"github.com/spf13/cobra"
)

type NetworkAccountLogs struct {
	NetworkOwners     []common.Address
	NetworkFeeAccount common.Address
	InfraFeeAccount   common.Address
	L1RewardRecipient common.Address
}

// NetworkAccountCmd represents the NetworkFeeAccount command
var NetworkAccountCmd = &cobra.Command{
	Use:   "NetworkAccount",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		chains, err := prompt.SelectChains()
		if err != nil {
			log.Fatal(err)
		}

		_, child, err := prompt.SelectCrossChainProviders(chains)
		if err != nil {
			log.Fatal(err)
		}

		childClient, err := ethclient.Dial(child)
		if err != nil {
			log.Fatal(err)
		}

		ArbOwnerPublic, err := precompilesgen.NewArbOwnerPublic(types.ArbOwnerPublicAddress, childClient)
		if err != nil {
			log.Fatal(err)
		}

		ArbGasInfo, err := precompilesgen.NewArbGasInfo(types.ArbGasInfoAddress, childClient)
		if err != nil {
			log.Fatal(err)
		}

		Callopts := &bind.CallOpts{
			Pending: false,
			Context: nil,
		}

		owners, err := ArbOwnerPublic.GetAllChainOwners(Callopts)
		if err != nil {
			log.Fatal(err)
		}

		NetworkFeeAccount, err := ArbOwnerPublic.GetNetworkFeeAccount(Callopts)
		if err != nil {
			log.Fatal(err)
		}

		InfraFeeAccount, err := ArbOwnerPublic.GetInfraFeeAccount(Callopts)
		if err != nil {
			log.Fatal(err)
		}

		L1RewardRecipient, err := ArbGasInfo.GetL1RewardRecipient(Callopts)
		if err != nil {
			log.Fatal(err)
		}

		logs.PrintFromatter(utils.ConvertBytesToHex(NetworkAccountLogs{
			NetworkOwners:     owners,
			NetworkFeeAccount: NetworkFeeAccount,
			InfraFeeAccount:   InfraFeeAccount,
			L1RewardRecipient: L1RewardRecipient,
		}))
	},
}
