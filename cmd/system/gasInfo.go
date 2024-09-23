/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/solgen/go/precompilesgen"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
)

var opts = &bind.CallOpts{
	Pending: false, // 트랜잭션이 확정된 상태를 조회
	Context: nil,   // 컨텍스트가 필요한 경우 (예: 시간 초과)
}

// gasInfoCmd represents the gasInfo command
var GasInfoCmd = &cobra.Command{
	Use:   "gasInfo",
	Short: "Get & Set Network Gas variable from ArbOS",
	Run: func(cmd *cobra.Command, args []string) {
		provider, err := prompt.SelectProvider()
		if err != nil {
			log.Fatal(err)
		}
		client := utils.GetClient(provider)

		ArbGasInfo, err := precompilesgen.NewArbGasInfo(types.ArbGasInfoAddress, client)
		ArbSys, err := precompilesgen.NewArbSys(types.ArbSysAddress, client)
		if err != nil {
			log.Fatal(err)
		}

		ArbOSVersion, err := ArbSys.ArbOSVersion(opts)

		fmt.Printf("\nArbOSVersion                 : %d\n", ArbOSVersion.Int64()-55)
		if ArbOSVersion.Int64()-55 >= 20 {

		}
		L1FeesAvailable, err := ArbGasInfo.GetL1FeesAvailable(opts)
		L1PricingFundsDueForRewards, err := ArbGasInfo.GetL1PricingFundsDueForRewards(opts) // 20
		L1PricingSurplus, err := ArbGasInfo.GetL1PricingSurplus(opts)                       // 20
		LastL1PricingSurplus, err := ArbGasInfo.GetLastL1PricingSurplus(opts)               // 20

		L1PricingUnitsSinceUpdate, err := ArbGasInfo.GetL1PricingUnitsSinceUpdate(opts) // 20

		L1BaseFeeEstimate, err := ArbGasInfo.GetL1BaseFeeEstimate(opts)
		LastL1PricingUpdateTime, err := ArbGasInfo.GetLastL1PricingUpdateTime(opts) // 20

		MinimumGasPrice, err := ArbGasInfo.GetMinimumGasPrice(opts)
		AmortizedCostCapBips, err := ArbGasInfo.GetAmortizedCostCapBips(opts)
		L1PricingEquilibrationUnits, err := ArbGasInfo.GetL1PricingEquilibrationUnits(opts) // 20
		L1BaseFeeEstimateInertia, err := ArbGasInfo.GetL1BaseFeeEstimateInertia(opts)
		PerBatchGasCharge, err := ArbGasInfo.GetPerBatchGasCharge(opts)
		L1RewardRate, err := ArbGasInfo.GetL1RewardRate(opts)
		L1RewardRecipient, err := ArbGasInfo.GetL1RewardRecipient(opts)
		CurrentTxL1GasFees, err := ArbGasInfo.GetCurrentTxL1GasFees(opts)
		GasBacklog, err := ArbGasInfo.GetGasBacklog(opts)
		GasBacklogTolerance, err := ArbGasInfo.GetGasBacklogTolerance(opts)
		PricingInertia, err := ArbGasInfo.GetPricingInertia(opts)

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
