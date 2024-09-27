/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hokaccha/go-prettyjson"
	ethlib "github.com/imelon2/orbit-cli/ethLib"
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

		abi, err := ethlib.GetAbi(strings.NewReader(string(_abi)))
		if err != nil {
			log.Fatal(err)
		}

		provider, err := prompt.SelectProvider()
		if err != nil {
			log.Fatal(err)
		}

		hash, err := prompt.EnterTransactionHash()
		if err != nil {
			log.Fatal(err)
		}

		client, err := ethclient.Dial(provider)
		if err != nil {
			log.Fatal(err)
		}

		txHash := ethlib.NewTransactionHash(client, hash)

		tx, err := txHash.GetTransactionReceipt()
		if err != nil {
			log.Fatal(err)
		}

		resultJson := make([]EventLog, 0)
		for _, eLog := range tx.Logs {
			var eventHashTopic = eLog.Topics[0]
			var eventJson EventLog

			event, _ := abi.EventByID(eventHashTopic)
			eLog.Topics = eLog.Topics[1:]

			eventJson.Name = event.RawName
			eventJson.Signature = event.Sig
			eventJson.Topic = eventHashTopic.Hex()

			logDatas, err := event.Inputs.Unpack(eLog.Data)

			if err != nil {
				log.Fatalf("failed to decode unpack event: %v", err)
			}

			eventData := make(map[string]interface{})
			dataIndex := 0
			topicIndex := 0
			for _, eventInput := range event.Inputs {
				if eventInput.Indexed {
					topic := eLog.Topics[topicIndex]
					data, err := decodeIndexedInput(eventInput, topic.Bytes())
					if err != nil {
						log.Fatalf("Failed to decode unpack topics: %v", err)
					}
					eventData[eventInput.Name] = data
					topicIndex++
				} else {
					eventData[eventInput.Name] = logDatas[dataIndex]
					dataIndex++
				}
			}

			eventJson.Params = utils.ConvertBytesToHex(eventData)
			resultJson = append(resultJson, eventJson)
		}

		formatter := prettyjson.NewFormatter()
		formatter.Indent = 2
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
