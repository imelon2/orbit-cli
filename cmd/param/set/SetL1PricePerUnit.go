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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imelon2/orbit-cli/common/logs"
	"github.com/imelon2/orbit-cli/common/tx"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/solgen/go/precompilesgen"
	"github.com/spf13/cobra"
)

// SetL1PricePerUnitCmd represents the SetL1PricePerUnit command
var SetL1PricePerUnitCmd = &cobra.Command{
	Use:   "SetL1PricePerUnit",
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

		ArbGasInfo, err := precompilesgen.NewArbGasInfo(types.ArbGasInfoAddress, childClient)
		if err != nil {
			log.Fatal(err)
		}

		Callopts := &bind.CallOpts{
			Pending: false,
			Context: nil,
		}
		currentPricePerUnit, err := ArbGasInfo.GetL1BaseFeeEstimate(Callopts)
		if err != nil {
			log.Fatal(err)
		}

		pricePerUnit, err := prompt.EnterInt(0, fmt.Sprintf("enter new pricePerUnit "+logs.GrayString("(current: %d)"), currentPricePerUnit))
		if err != nil {
			log.Fatal(err)
		}

		ArbOwner, err := precompilesgen.NewArbOwner(types.ArbOwnerAddress, childClient)
		if err != nil {
			log.Fatal(err)
		}

		response, err := ArbOwner.SetL1PricePerUnit(auth, big.NewInt(int64(*pricePerUnit)))
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
