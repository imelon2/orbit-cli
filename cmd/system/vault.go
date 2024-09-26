/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethLib "github.com/imelon2/orbit-cli/ethLib"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/solgen/go/precompilesgen"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// vaultCmd represents the vault command
var VaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		provider, err := prompt.SelectProvider()
		if err != nil {
			log.Fatal(err)
		}
		client, err := ethclient.Dial(provider)
		if err != nil {
			log.Fatal(err)
		}

		/* -------- ArbOS -------- */
		ArbOwnerPublic, err := precompilesgen.NewArbOwnerPublic(types.ArbOwnerPublicAddress, client)
		if err != nil {
			log.Fatal(err)
		}

		ArbGasInfo, err := precompilesgen.NewArbGasInfo(types.ArbGasInfoAddress, client)
		if err != nil {
			log.Fatal(err)
		}

		/* -------- get account address -------- */
		networkFeeAccount, err := ArbOwnerPublic.GetNetworkFeeAccount(ethLib.Callopts)
		if err != nil {
			log.Fatal(err)
		}

		infraFeeAccount, err := ArbOwnerPublic.GetInfraFeeAccount(ethLib.Callopts)
		if err != nil {
			log.Fatal(err)
		}

		l1RewardRecipient, err := ArbGasInfo.GetL1RewardRecipient(ethLib.Callopts)
		if err != nil {
			log.Fatal(err)
		}

		/* -------- get balance -------- */
		networkFeeBalance, err := client.PendingBalanceAt(context.Background(), networkFeeAccount)
		if err != nil {
			log.Fatal(err)
		}
		networkFee := new(big.Float).SetInt(networkFeeBalance)
		networkFee.Quo(networkFee, big.NewFloat(1e18))

		infraFeeAccountBalance, err := client.PendingBalanceAt(context.Background(), infraFeeAccount)
		if err != nil {
			log.Fatal(err)
		}
		infraFee := new(big.Float).SetInt(infraFeeAccountBalance)
		infraFee.Quo(infraFee, big.NewFloat(1e18))

		l1RewardRecipientBalance, err := client.PendingBalanceAt(context.Background(), l1RewardRecipient)
		if err != nil {
			log.Fatal(err)
		}
		l1Reward := new(big.Float).SetInt(l1RewardRecipientBalance)
		l1Reward.Quo(l1Reward, big.NewFloat(1e18))

		L1PricerFundsPoolBalance, err := client.PendingBalanceAt(context.Background(), ethLib.L1PricerFundsPool)
		if err != nil {
			log.Fatal(err)
		}
		L1PricerFundsPool := new(big.Float).SetInt(L1PricerFundsPoolBalance)
		L1PricerFundsPool.Quo(L1PricerFundsPool, big.NewFloat(1e18))

		data := [][]string{
			{fmt.Sprintf("%.18f ETH", networkFee), fmt.Sprintf("%.18f ETH", infraFee), fmt.Sprintf("%.18f ETH", L1PricerFundsPool), fmt.Sprintf("%.18f ETH", l1Reward)},
			{networkFeeAccount.Hex(), infraFeeAccount.Hex(), ethLib.L1PricerFundsPool.Hex(), l1RewardRecipient.Hex()},
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Network Fee", "Infra Fee", "L1 PricerFunds Pool", "l1 Rewarder"})
		table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlueColor},
			tablewriter.Colors{tablewriter.FgBlueColor, tablewriter.Bold},
			tablewriter.Colors{tablewriter.FgBlueColor, tablewriter.Bold},
			tablewriter.Colors{tablewriter.FgBlueColor, tablewriter.Bold},
		)
		table.SetRowLine(true)
		table.AppendBulk(data) // Add Bulk Data
		table.Render()
	},
}

func init() {
}
