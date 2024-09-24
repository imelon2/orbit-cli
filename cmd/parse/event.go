/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hokaccha/go-prettyjson"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
)

type EventLog struct {
	Name      string      `json:"name"`
	Signature string      `json:"signature"`
	Topic     string      `json:"topic"`
	Params    interface{} `json:"params"`
}

// eventCmd represents the event command
var EventCmd = &cobra.Command{
	Use:   "event",
	Short: "Parse calldata by transaction hash",
	Run: func(cmd *cobra.Command, args []string) {

		abiPath := utils.GetAbiDir()
		_abi, _ := os.ReadFile(abiPath)
		parsedABI, err := abi.JSON(strings.NewReader(string(_abi)))

		if err != nil {
			log.Fatalf("Failed to get ABI: %v", err)
		}
		provider, err := prompt.SelectProvider()
		if err != nil {
			log.Fatalf("Failed to get provider: %v", err)
		}
		txHash, err := prompt.EnterTransactionHash()
		if err != nil {
			log.Fatalf("Failed to get txHash: %v", err)
		}

		client, err := ethclient.Dial(provider)
		if err != nil {
			log.Fatal(err)
		}
		tx, err := client.TransactionReceipt(context.Background(), txHash)

		resultJson := make([]EventLog, 0)
		for _, data := range tx.Logs {
			var eventHashTopic = data.Topics[0]
			var eventJson EventLog

			event, _ := parsedABI.EventByID(eventHashTopic)
			data.Topics = data.Topics[1:]

			eventJson.Name = event.RawName
			eventJson.Signature = event.Sig
			eventJson.Topic = eventHashTopic.Hex()

			eventData := make(map[string]interface{})

			lastIndex := 0
			for j, topic := range data.Topics {
				data, err := decodeIndexedInput(event.Inputs[j], topic.Bytes())
				if err != nil {
					log.Fatalf("Failed to unpack calldata: %v", err)
				}
				eventData[event.Inputs[j].Name] = data

				lastIndex++
			}

			inter, err := event.Inputs.Unpack(data.Data)
			if err != nil {
				log.Fatalf("Failed to unpack calldata: %v", err)
			}

			for k, data := range inter {
				eventData[event.Inputs[lastIndex+k].Name] = data
			}
			eventJson.Params = utils.ConvertBytesToHex(eventData)

			resultJson = append(resultJson, eventJson)
		}

		formatter := prettyjson.NewFormatter()
		formatter.Indent = 2
		// formatter.KeyColor.DisableColor()

		coloredJson, err := formatter.Marshal(resultJson)
		if err != nil {
			log.Fatalf("Failed to Marshal decoded event: %v", err)
		}

		fmt.Println(string(coloredJson))
	},
}

func init() {
}

func decodeIndexedInput(input abi.Argument, data []byte) (interface{}, error) {
	switch input.Type.T {
	case abi.AddressTy:
		return common.BytesToAddress(data[12:]), nil // 주소는 마지막 20바이트 사용
	case abi.UintTy, abi.IntTy:
		return new(big.Int).SetBytes(data), nil
	case abi.BoolTy:
		return data[len(data)-1] == 1, nil
	case abi.BytesTy, abi.FixedBytesTy:
		return data, nil
	default:
		return nil, fmt.Errorf("지원되지 않는 인덱스 타입: %s", input.Type.String())
	}
}
