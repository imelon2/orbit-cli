/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethLib "github.com/imelon2/orbit-cli/ethLib"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/solgen/go/precompilesgen"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
)

// gasInfoCmd represents the gasInfo command
var GasCmd = &cobra.Command{
	Use:   "gas",
	Short: "Print network gas variable from ArbGasInfo, ArbSys",
	Run: func(cmd *cobra.Command, args []string) {
		provider, err := prompt.SelectProvider()
		if err != nil {
			log.Fatal(err)
		}
		client, err := ethclient.Dial(provider)
		if err != nil {
			log.Fatal(err)
		}

		ArbGasInfo, err := precompilesgen.NewArbGasInfo(types.ArbGasInfoAddress, client)
		ArbSys, err := precompilesgen.NewArbSys(types.ArbSysAddress, client)
		if err != nil {
			log.Fatal(err)
		}

		ArbOSVersion, err := ArbSys.ArbOSVersion(ethLib.Callopts)

		fmt.Printf("\nArbOSVersion                 : %d\n", ArbOSVersion.Int64()-55)
		if ArbOSVersion.Int64()-55 >= 20 {

		}
		L1FeesAvailable, err := ArbGasInfo.GetL1FeesAvailable(ethLib.Callopts)
		L1PricingFundsDueForRewards, err := ArbGasInfo.GetL1PricingFundsDueForRewards(ethLib.Callopts) // 20
		L1PricingSurplus, err := ArbGasInfo.GetL1PricingSurplus(ethLib.Callopts)                       // 20
		LastL1PricingSurplus, err := ArbGasInfo.GetLastL1PricingSurplus(ethLib.Callopts)               // 20

		L1PricingUnitsSinceUpdate, err := ArbGasInfo.GetL1PricingUnitsSinceUpdate(ethLib.Callopts) // 20

		L1BaseFeeEstimate, err := ArbGasInfo.GetL1BaseFeeEstimate(ethLib.Callopts)
		LastL1PricingUpdateTime, err := ArbGasInfo.GetLastL1PricingUpdateTime(ethLib.Callopts) // 20

		MinimumGasPrice, err := ArbGasInfo.GetMinimumGasPrice(ethLib.Callopts)
		AmortizedCostCapBips, err := ArbGasInfo.GetAmortizedCostCapBips(ethLib.Callopts)
		L1PricingEquilibrationUnits, err := ArbGasInfo.GetL1PricingEquilibrationUnits(ethLib.Callopts) // 20
		L1BaseFeeEstimateInertia, err := ArbGasInfo.GetL1BaseFeeEstimateInertia(ethLib.Callopts)
		PerBatchGasCharge, err := ArbGasInfo.GetPerBatchGasCharge(ethLib.Callopts)
		L1RewardRate, err := ArbGasInfo.GetL1RewardRate(ethLib.Callopts)
		L1RewardRecipient, err := ArbGasInfo.GetL1RewardRecipient(ethLib.Callopts)
		CurrentTxL1GasFees, err := ArbGasInfo.GetCurrentTxL1GasFees(ethLib.Callopts)
		GasBacklog, err := ArbGasInfo.GetGasBacklog(ethLib.Callopts)
		GasBacklogTolerance, err := ArbGasInfo.GetGasBacklogTolerance(ethLib.Callopts)
		PricingInertia, err := ArbGasInfo.GetPricingInertia(ethLib.Callopts)

		header, err := client.HeaderByNumber(context.Background(), nil /* Latest */)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print("\n")
		fmt.Print(utils.BoldGreenString("#### L1 L2 Gas Price Info ####\n"))
		fmt.Printf("Current L2 Base Fee (block)  : %d\n", header.BaseFee.Uint64())
		fmt.Printf("L1BaseFeeEstimate            : %d\n", L1BaseFeeEstimate)
		fmt.Printf("LastL1PricingUpdateTime      : %d\n", LastL1PricingUpdateTime)
		fmt.Print("\n\n")
		fmt.Print(utils.BoldGreenString("#### Dynamic Data ####\n"))
		fmt.Printf("L1PricingUnitsSinceUpdate    : %d\n", L1PricingUnitsSinceUpdate)
		fmt.Printf("GasBacklog                   : %d\n", GasBacklog)
		fmt.Printf("CurrentTxL1GasFees           : %d\n", CurrentTxL1GasFees)
		fmt.Print("\n\n")
		fmt.Print(utils.BoldGreenString("#### Vault Info ####\n"))
		fmt.Printf("L1FeesAvailable              : %d\n", L1FeesAvailable)
		fmt.Printf("L1PricingFundsDueForRewards  : %d\n", L1PricingFundsDueForRewards)
		fmt.Printf("L1PricingSurplus(d)          : %d\n", L1PricingSurplus)
		fmt.Printf("LastL1PricingSurplus         : %d\n", LastL1PricingSurplus)
		fmt.Print("\n\n")
		fmt.Print(utils.BoldGreenString("#### L1 Gas Constant ####\n"))
		fmt.Printf("PerBatchGasCharge            : %d\n", PerBatchGasCharge)
		fmt.Printf("L1RewardRate                 : %d\n", L1RewardRate)
		fmt.Printf("L1RewardRecipient            : %s\n", L1RewardRecipient.Hex())
		fmt.Printf("AmortizedCostCapBips         : %d\n", AmortizedCostCapBips)
		fmt.Printf("L1PricingEquilibrationUnits  : %d\n", L1PricingEquilibrationUnits)
		fmt.Printf("L1BaseFeeEstimateInertia     : %d\n", L1BaseFeeEstimateInertia)
		fmt.Print("\n\n")
		fmt.Print(utils.BoldGreenString("#### L2 Gas Constant ####\n"))
		fmt.Printf("MinimumGasPrice              : %d\n", MinimumGasPrice)
		fmt.Printf("GasBacklogTolerance          : %d\n", GasBacklogTolerance)
		fmt.Printf("PricingInertia               : %d\n", PricingInertia)

		fmt.Print("\n\n")
	},
}

func init() {
}
