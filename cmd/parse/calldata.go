/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/imelon2/orbit-toolkit/prompt"
	"github.com/imelon2/orbit-toolkit/utils"
	"github.com/spf13/cobra"
)

// calldataCmd represents the calldata command
var CalldataCmd = &cobra.Command{
	Use:   "calldata",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

		abiPath := utils.GetAbiDir()
		_abi, _ := os.ReadFile(abiPath)

		// // 파일 내용을 ABI 파싱에 전달
		parsedABI, err := abi.JSON(strings.NewReader(string(_abi)))
		if err != nil {
			log.Fatalf("Failed to get ABI: %v", err)
		}
		provider, err := prompt.SelectProvider()
		if err != nil {
			log.Fatal(err)
		}

		txHash, err := prompt.EnterTransactionHash()
		if err != nil {
			log.Fatal(err)
		}

		client := utils.GetClient(provider)
		tx, _, err := client.TransactionByHash(context.Background(), txHash)

		if err != nil {
			log.Fatal(err)
		}

		data := tx.Data()

		method, err := parsedABI.MethodById(data[:4] /* function selector */)
		if err != nil {
			log.Fatalf("Failed to get method from calldata: %v", err)
		}

		inter, err := method.Inputs.Unpack(data[4:] /* data without selector */)

		if err != nil {
			log.Fatalf("Failed to unpack calldata: %v", err)
		}

		jsonCalldata := make(map[string]interface{})
		for i, data := range inter {
			jsonCalldata[method.Inputs[i].Name] = data
		}

		jsonData, err := json.MarshalIndent(jsonCalldata, "", "  ")
		if err != nil {
			log.Fatalf("Failed to MarshalIndent calldata: %v", err)
		}
		fmt.Printf("Function : %s\n", parsedABI.Methods[method.Name])
		fmt.Println(string(jsonData))
	},
}

func init() {
}
