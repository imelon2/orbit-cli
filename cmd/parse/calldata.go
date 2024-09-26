/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	ethlib "github.com/imelon2/orbit-cli/ethLib"
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
		_abi, err := os.ReadFile(abiPath)
		if err != nil {
			log.Fatal(err)
		}

		abi, err := ethlib.GetAbi(strings.NewReader(string(_abi)))
		if err != nil {
			log.Fatal(err)
		}

		providerOrCalldata, isProvider, err := prompt.SelectProviderOrBytes()
		if err != nil {
			log.Fatal(err)
		}

		var data []byte
		if isProvider {
			client, err := ethclient.Dial(providerOrCalldata)
			if err != nil {
				log.Fatal(err)
			}

			hash, err := prompt.EnterTransactionHash()
			if err != nil {
				log.Fatal(err)
			}

			txHash := ethlib.NewTransactionHash(client, hash)
			tx, _, err := txHash.GetTransactionByHash()
			if err != nil {
				log.Fatal(err)
			}

			data = tx.Data()
		} else {
			providerOrCalldata = providerOrCalldata[2:]
			calldataBytes, err := hex.DecodeString(providerOrCalldata)
			if err != nil {
				log.Fatalf("fail string calldata decode to hex: %v", err)
			}
			data = calldataBytes
		}

		calldata, err := ethlib.NewCalldata(&abi, data).GetMethodById()
		if err != nil {
			log.Fatal(err)
		}

		inter, err := calldata.GetUnpackedHexdata()
		if err != nil {
			log.Fatal(err)
		}

		var resultJson CalldataLog
		resultJson.Function = calldata.Method.RawName
		jsonCalldata := make(map[string]interface{})
		for i, data := range inter {
			jsonCalldata[calldata.Method.Inputs[i].Name] = utils.ConvertBytesToHex(data)
		}

		resultJson.Params = jsonCalldata
		utils.PrintPrettyJson(resultJson)
	},
}

func init() {
}
