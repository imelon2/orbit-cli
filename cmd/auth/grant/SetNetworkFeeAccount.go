/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imelon2/orbit-cli/common/logs"
	"github.com/imelon2/orbit-cli/common/tx"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/solgen/go/precompilesgen"
	"github.com/spf13/cobra"
)

// SetNetworkFeeAccountCmd represents the SetNetworkFeeAccount command
var SetNetworkFeeAccountCmd = &cobra.Command{
	Use:   "SetNetworkFeeAccount",
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

		_, ks, account, err := prompt.SelectWalletForSign()
		if err != nil {
			log.Fatal(err)
		}

		auth, err := tx.GetAuthByKeystore(ks, *account, childClient)
		if err != nil {
			log.Fatal(err)
		}

		ArbOwnerPublic, err := precompilesgen.NewArbOwnerPublic(types.ArbOwnerPublicAddress, childClient)
		if err != nil {
			log.Fatal(err)
		}

		Callopts := &bind.CallOpts{
			Pending: false,
			Context: nil,
		}
		NetworkFeeAccount, err := ArbOwnerPublic.GetNetworkFeeAccount(Callopts)
		if err != nil {
			log.Fatal(err)
		}

		newNetworkFeeAccount, err := prompt.EnterAddress(fmt.Sprintf("enter new Network Fee Account "+logs.GrayString("(current address: %s)"), NetworkFeeAccount.Hex()))
		if err != nil {
			log.Fatal(err)
		}

		ArbOwner, err := precompilesgen.NewArbOwner(types.ArbOwnerAddress, childClient)
		if err != nil {
			log.Fatal(err)
		}

		response, err := ArbOwner.SetNetworkFeeAccount(auth, common.HexToAddress(*newNetworkFeeAccount))
		if err != nil {
			log.Fatal(err)
		}

		receipt, err := bind.WaitMined(context.Background(), childClient, response)
		if err != nil {
			log.Fatal(err)
		}

		logs.PrintReceiptFromatter(receipt)
	},
}
