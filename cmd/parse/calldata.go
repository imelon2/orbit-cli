/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
)

// calldataCmd represents the calldata command
var CalldataCmd = &cobra.Command{
	Use:   "calldata",
	Short: "Parse calldata by transaction hash or bytes",
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
		txHashOrCalldata, err := prompt.EnterTransactionHashOrBytes()

		client := utils.GetClient(provider)

		var data []byte
		if utils.IsTransaction(txHashOrCalldata) {
			tx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(txHashOrCalldata))
			if err != nil {
				log.Fatal(err)
			}
			data = tx.Data()
		} else {
			txHashOrCalldata = txHashOrCalldata[2:]

			// hex 문자열을 []byte로 변환
			calldataBytes, err := hex.DecodeString(txHashOrCalldata)
			if err != nil {
				log.Fatalf("calldata 변환 에러: %v", err)
			}
			data = calldataBytes
		}

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
			jsonCalldata[method.Inputs[i].Name] = utils.ConvertBytesToHex(data)
		}

		jsonData, err := json.MarshalIndent(jsonCalldata, "", "  ")
		if err != nil {
			log.Fatalf("Failed to MarshalIndent calldata: %v", err)
		}
		fmt.Printf("Function : %s\n", parsedABI.Methods[method.Name])
		fmt.Printf("Calldata Length : %d\n\n", len(data))
		fmt.Fprintln(os.Stdout, string(jsonData))
	},
}

func init() {
}
