/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
)

type CalldataLog struct {
	Function string      `json:"function"`
	Params   interface{} `json:"params"`
}

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
		providerOrCalldata, isProvider, err := prompt.SelectProviderOrBytes()
		if err != nil {
			log.Fatal(err)
		}

		var data []byte
		if isProvider {
			client := utils.GetClient(providerOrCalldata)

			txHash, err := prompt.EnterTransactionHash()
			if err != nil {
				log.Fatalf("Failed to get txHash: %v", err)
			}

			tx, _, err := client.TransactionByHash(context.Background(), txHash)
			if err != nil {
				fmt.Errorf("Failed to get TransactionByHash: %v", err)
				return
			}
			data = tx.Data()
		} else {
			providerOrCalldata = providerOrCalldata[2:]
			// hex 문자열을 []byte로 변환
			calldataBytes, err := hex.DecodeString(providerOrCalldata)
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

		var resultJson CalldataLog
		resultJson.Function = method.RawName
		jsonCalldata := make(map[string]interface{})
		for i, data := range inter {
			jsonCalldata[method.Inputs[i].Name] = utils.ConvertBytesToHex(data)
		}

		resultJson.Params = jsonCalldata

		utils.PrintPrettyJson(resultJson)
	},
}

func init() {
}
