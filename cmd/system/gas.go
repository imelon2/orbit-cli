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
		if err != nil {
			log.Fatalf("fail bind ArbGasInfo: %v", err)
		}
		ArbSys, err := precompilesgen.NewArbSys(types.ArbSysAddress, client)
		if err != nil {
			log.Fatalf("fail bind ArbSys: %v", err)
		}

		ArbOSVersion, err := ArbSys.ArbOSVersion(ethLib.Callopts)
		if err != nil {
			log.Fatalf("fail call ArbOSVersion: %v", err)
		}

		fmt.Printf("\nArbOSVersion                 : %d\n", ArbOSVersion.Int64()-55)
		// if ArbOSVersion.Int64()-55 >= 20 {

		// }
		L1FeesAvailable, _ := ArbGasInfo.GetL1FeesAvailable(ethLib.Callopts)
		// if err != nil {
		// 	log.Fatalf("fail call GetL1FeesAvailable: %v", err)
		// }

		L1PricingFundsDueForRewards, _ := ArbGasInfo.GetL1PricingFundsDueForRewards(ethLib.Callopts) // 20
		// if err != nil {
		// 	log.Fatalf("fail call GetL1PricingFundsDueForRewards: %v", err)
		// }

		L1PricingSurplus, _ := ArbGasInfo.GetL1PricingSurplus(ethLib.Callopts) // 20
		// if err != nil {
		// 	log.Fatalf("fail call GetL1PricingSurplus: %v", err)
		// }

		LastL1PricingSurplus, _ := ArbGasInfo.GetLastL1PricingSurplus(ethLib.Callopts) // 20
		// if err != nil {
		// 	log.Fatalf("fail call GetLastL1PricingSurplus: %v", err)
		// }

		L1PricingUnitsSinceUpdate, _ := ArbGasInfo.GetL1PricingUnitsSinceUpdate(ethLib.Callopts) // 20
		// if err != nil {
		// 	log.Fatalf("fail call GetL1PricingUnitsSinceUpdate: %v", err)
		// }

		L1BaseFeeEstimate, _ := ArbGasInfo.GetL1BaseFeeEstimate(ethLib.Callopts)
		// if err != nil {
		// 	log.Fatalf("fail call GetL1BaseFeeEstimate: %v", err)
		// }

		LastL1PricingUpdateTime, _ := ArbGasInfo.GetLastL1PricingUpdateTime(ethLib.Callopts) // 20
		// if err != nil {
		// 	log.Fatalf("fail call GetLastL1PricingUpdateTime: %v", err)
		// }

		MinimumGasPrice, _ := ArbGasInfo.GetMinimumGasPrice(ethLib.Callopts)
		// if err != nil {
		// 	log.Fatalf("fail call GetMinimumGasPrice: %v", err)
		// }

		AmortizedCostCapBips, _ := ArbGasInfo.GetAmortizedCostCapBips(ethLib.Callopts)
		// if err != nil {
		// 	log.Fatalf("fail call GetAmortizedCostCapBips: %v", err)
		// }

		L1PricingEquilibrationUnits, _ := ArbGasInfo.GetL1PricingEquilibrationUnits(ethLib.Callopts) // 20
		// if err != nil {
		// 	log.Fatalf("fail call GetL1PricingEquilibrationUnits: %v", err)
		// }

		L1BaseFeeEstimateInertia, _ := ArbGasInfo.GetL1BaseFeeEstimateInertia(ethLib.Callopts)
		// if err != nil {
		// 	log.Fatalf("fail call GetL1BaseFeeEstimateInertia: %v", err)
		// }

		PerBatchGasCharge, _ := ArbGasInfo.GetPerBatchGasCharge(ethLib.Callopts)
		// if err != nil {
		// 	log.Fatalf("fail call GetPerBatchGasCharge: %v", err)
		// }

		L1RewardRate, _ := ArbGasInfo.GetL1RewardRate(ethLib.Callopts)
		// if err != nil {
		// 	log.Fatalf("fail call GetL1RewardRate: %v", err)
		// }

		L1RewardRecipient, _ := ArbGasInfo.GetL1RewardRecipient(ethLib.Callopts)
		// if err != nil {
		// 	log.Fatalf("fail call GetL1RewardRecipient: %v", err)
		// }

		CurrentTxL1GasFees, _ := ArbGasInfo.GetCurrentTxL1GasFees(ethLib.Callopts)
		// if err != nil {
		// 	log.Fatalf("fail call GetCurrentTxL1GasFees: %v", err)
		// }

		GasBacklog, _ := ArbGasInfo.GetGasBacklog(ethLib.Callopts)
		// if err != nil {
		// 	log.Fatalf("fail call GetGasBacklog: %v", err)
		// }

		GasBacklogTolerance, _ := ArbGasInfo.GetGasBacklogTolerance(ethLib.Callopts)
		// if err != nil {
		// 	log.Fatalf("fail call GetGasBacklogTolerance: %v", err)
		// }

		PricingInertia, _ := ArbGasInfo.GetPricingInertia(ethLib.Callopts)
		// if err != nil {
		// 	log.Fatalf("fail call GetPricingInertia: %v", err)
		// }

		block, err := client.BlockByNumber(context.Background(), nil /* Latest */)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print("\n")
		fmt.Print(utils.BoldGreenString("#### L1 L2 Gas Price Info ####\n"))
		fmt.Printf("Current L2 Base Fee (block)  : %d\n", block.BaseFee().Uint64())
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
