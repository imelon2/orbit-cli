/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	arblib "github.com/imelon2/orbit-cli/arbLib"
	"github.com/imelon2/orbit-cli/contractgen"
	ethlib "github.com/imelon2/orbit-cli/ethLib"
	"github.com/spf13/cobra"
)

// bridgeCmd represents the bridge command
var BridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, auth, provider, err := ethlib.GenerateAuth()
		if err != nil {
			log.Fatal(err)
		}

		network, err := contractgen.GetNetworkInfoByParent(client)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Router : %s\n", network.TokenBridge.L1GatewayRouter)
		fmt.Printf("Router : %s\n", provider)
		router, err := arblib.NewRouter(client, network.TokenBridge.L1GatewayRouter)
		if err != nil {
			log.Fatal(err)
		}

		router.Auth = auth

		data, err := hex.DecodeString("37c6145a00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000002247cd5c2910000000000000000000000000000000000000000000000000000000000000020000000000000000000000000382ffce2287252f930e1c8dc9328dac5bf282ba1000000000000000000000000eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee0000000000000000000000000000000000000000000000000001791cda050f67000000000000000000000000000000000000000000000000000000000000019f00000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000124d025dec00000000000000000000000003f4123905aba29e3118a51d25e45acad3e1ffc59000000000000000000000000b8901acb165ed027e32754e0ffe830802919727f000000000000000000000000710bda329b2a6224e4b44833de30f38e7f81d564000000000000000000000000000000000000000000000000000000000000a4b10000000000000000000000000000000000000000000000000091d6284ff4f47500000000000000000000000000000000000000000000000000911b7c973cb0ac0000000000000000000000000000000000000000000000000000886c98b760000000000000000000000000000000000000000000000000000000019261b6371d00000000000000000000000000000000000000000000000000000000000008f10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		if err != nil {
			log.Fatalf("fail string calldata decode to hex: %v", err)
		}

		params := arblib.OutboundTransferPrams{
			Token:       common.Address{},
			To:          common.Address{},
			Amount:      big.NewInt(10),
			MaxGas:      big.NewInt(10),
			GasPriceBid: big.NewInt(10),
			Data:        data,
		}

		tx, err := router.DepositFunc(&params)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(tx)
	},
}

func init() {

}
