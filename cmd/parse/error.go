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

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	ethlib "github.com/imelon2/orbit-cli/ethLib"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
)

type ErrorDataLog struct {
	System  string         `json:"system,omitempty"`
	Message string         `json:"message,omitempty"` // 빈 값 허용
	Custom  *CustomMessage `json:"custom,omitempty"`  // nil 값 허용
	Hex     string         `json:"hex"`
}

type CustomMessage struct {
	Name   string      `json:"name,omitempty"`
	Params interface{} `json:"params,omitempty"`
}

// Revert reason ID for "Error(string)"
var revertReasonID = "08c379a0"

// errorCmd represents the error command
var ErrorCmd = &cobra.Command{
	Use:   "error",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

		abiPath := utils.GetAbiDir()
		_abi, err := os.ReadFile(abiPath)
		if err != nil {
			log.Fatal(err)
		}

		abiParse, err := ethlib.GetAbi(strings.NewReader(string(_abi)))
		if err != nil {
			log.Fatal(err)
		}

		providerOrError, isProvider, err := prompt.SelectProviderOrBytes()
		if err != nil {
			log.Fatal(err)
		}

		errRuslt := ErrorDataLog{}

		if isProvider {
			client, err := ethclient.Dial(providerOrError)
			if err != nil {
				log.Fatal(err)
			}

			hash, err := prompt.EnterTransactionHash()
			if err != nil {
				log.Fatal(err)
			}

			txHash := ethlib.NewTransactionHash(client, hash)
			tx, receipt, sender, err := txHash.GetTransactionSender()
			if err != nil {
				log.Fatal(err)
			}

			if receipt.Status == 1 {
				fmt.Printf("transaction hash %s is SUCCESS", receipt.TxHash)
				return
			}

			callMsg := ethereum.CallMsg{
				From:     *sender,
				To:       tx.To(),
				Gas:      tx.Gas(),
				GasPrice: tx.GasPrice(),
				Data:     tx.Data(),
				// GasFeeCap: tx.GasFeeCap(),
				// GasTipCap: tx.GasTipCap(),
				Value: tx.Value(),
			}

			_, errData := client.CallContract(context.Background(), callMsg, receipt.BlockNumber)

			rpcErr, ok := errData.(rpc.DataError)
			if ok {
				errRuslt.System = rpcErr.Error()
				errRuslt.Hex, _ = rpcErr.ErrorData().(string)

				if errRuslt.Hex == "" {
					errRuslt.Hex = "0x"
					errRuslt.Message = "NULL"
				} else if errRuslt.Hex[2:10] == revertReasonID {
					dataBytes, _ := hex.DecodeString(errRuslt.Hex[10:])

					stringType, err := abi.NewType("string", "", nil)
					arguments := abi.Arguments{
						{Type: stringType},
					}
					decodeAbiString, err := arguments.Unpack(dataBytes)
					if err != nil {
						log.Fatal(err)
					}

					for _, s := range decodeAbiString {
						_s, _ := s.(string)
						errRuslt.Message = _s
						break
					}

				} else {
					// var revertMessage string
					dataBytes, _ := hex.DecodeString(errRuslt.Hex[2:])
					var sigdata [4]byte
					for i, data := range dataBytes[:4] {
						sigdata[i] = data
					}

					errorAbi, err := abiParse.ErrorByID(sigdata)
					if err != nil {
						log.Fatal(err)
					}

					decodedError, err := errorAbi.Inputs.Unpack(dataBytes[4:])
					if err != nil {
						log.Fatal(err)
					}
					jsonError := make(map[string]interface{})
					for i, data := range decodedError {
						jsonError[errorAbi.Inputs[i].Name] = utils.ConvertBytesToHex(data)
					}

					newCustom := CustomMessage{
						Name:   errorAbi.Sig,
						Params: jsonError,
					}
					errRuslt.Custom = &newCustom
				}
			}

		} else {

			errRuslt.Hex = providerOrError
			if errRuslt.Hex[2:10] == revertReasonID {
				dataBytes, _ := hex.DecodeString(errRuslt.Hex[10:])

				stringType, err := abi.NewType("string", "", nil)
				arguments := abi.Arguments{
					{Type: stringType},
				}
				decodeAbiString, err := arguments.Unpack(dataBytes)
				if err != nil {
					log.Fatal(err)
				}

				for _, s := range decodeAbiString {
					_s, _ := s.(string)
					errRuslt.Message = _s
					break
				}

			} else {
				// var revertMessage string
				dataBytes, _ := hex.DecodeString(errRuslt.Hex[2:])
				var sigdata [4]byte
				for i, data := range dataBytes[:4] {
					sigdata[i] = data
				}
				errorAbi, err := abiParse.ErrorByID(sigdata)
				if err != nil {
					log.Fatal(err)
				}
				decodedError, err := errorAbi.Inputs.Unpack(dataBytes[4:])
				if err != nil {
					log.Fatal(err)
				}
				jsonError := make(map[string]interface{})
				for i, data := range decodedError {
					jsonError[errorAbi.Inputs[i].Name] = utils.ConvertBytesToHex(data)
				}

				newCustom := CustomMessage{
					Name:   errorAbi.Sig,
					Params: jsonError,
				}
				errRuslt.Custom = &newCustom
			}
		}

		utils.PrintPrettyJson(errRuslt)
	},
}

func init() {
}
