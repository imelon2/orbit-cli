/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		provider, err := prompt.SelectProvider()
		if err != nil {
			log.Fatal(err)
		}
		client := utils.GetClient(provider)

		ArbGasInfo, err := precompilesgen.NewArbGasInfo(types.ArbGasInfoAddress, client)
		if err != nil {
			log.Fatal(err)
		}
		L1FeesAvailable, err := ArbGasInfo.GetL1FeesAvailable(opts)
		L1BaseFeeEstimate, err := ArbGasInfo.GetL1BaseFeeEstimate(opts)
		L1PricingUnitsSinceUpdate, err := ArbGasInfo.GetL1PricingUnitsSinceUpdate(opts)
		L1PricingFundsDueForRewards, err := ArbGasInfo.GetL1PricingFundsDueForRewards(opts)
		LastL1PricingUpdateTime, err := ArbGasInfo.GetLastL1PricingUpdateTime(opts)
		LastL1PricingSurplus, err := ArbGasInfo.GetLastL1PricingSurplus(opts)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Print("\n")
		fmt.Printf("L1FeesAvailable             : %d\n", L1FeesAvailable)
		fmt.Printf("LastL1PricingSurplus        : %d\n", LastL1PricingSurplus)
		fmt.Printf("L1BaseFeeEstimate           : %d\n", L1BaseFeeEstimate)
		fmt.Printf("L1PricingUnitsSinceUpdate   : %d\n", L1PricingUnitsSinceUpdate)
		fmt.Print("\n")
		fmt.Printf("L1PricingFundsDueForRewards : %d\n", L1PricingFundsDueForRewards)
		fmt.Printf("LastL1PricingUpdateTime     : %d\n", LastL1PricingUpdateTime)

		return nil
	},
}

func init() {
}
