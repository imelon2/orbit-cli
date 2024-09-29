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

		data, err := hex.DecodeString("0000000000000000000000000000000000000000000000000013be304992d88000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000000")
		if err != nil {
			log.Fatalf("fail string calldata decode to hex: %v", err)
		}

		params := arblib.OutboundTransferPrams{
			Token:       common.HexToAddress("0xe2ef69e4af84dbefb0a75f8491f27a52bf047b01"),
			To:          common.HexToAddress("0xea9035a97722c1fde906a17184f558794e4a9141"),
			Amount:      big.NewInt(10000000000000000),
			MaxGas:      big.NewInt(10),
			GasPriceBid: big.NewInt(600000000),
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
