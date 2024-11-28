/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	arbnetwork "github.com/imelon2/orbit-cli/arbNetwork"
	"github.com/imelon2/orbit-cli/common/logs"
	"github.com/imelon2/orbit-cli/common/utils"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/spf13/cobra"
)

// SequencerEventCmd represents the SequencerEvent command
var SequencerEventCmd = &cobra.Command{
	Use:   "SequencerEvent",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		chains, err := prompt.SelectChains()
		if err != nil {
			log.Fatal(err)
		}

		parent, child, err := prompt.SelectCrossChainProviders(chains)
		if err != nil {
			log.Fatal(err)
		}

		parentClient, err := ethclient.Dial(parent)
		if err != nil {
			log.Fatal(err)
		}

		childClient, err := ethclient.Dial(child)
		if err != nil {
			log.Fatal(err)
		}

		network, err := arbnetwork.GetNetworkInfo(childClient)
		if err != nil {
			log.Fatal(err)
		}
		SequencerInbox, err := network.NewSequencerInbox(parentClient)
		if err != nil {
			log.Fatal(err)
		}

		query := ethereum.FilterQuery{
			Addresses: []common.Address{network.EthBridge.SequencerInbox},
		}

		logsChan := make(chan types.Log)
		sub, err := parentClient.SubscribeFilterLogs(context.Background(), query, logsChan)
		if err != nil {
			log.Fatalf("Failed to subscribe to logs: %v", err)
		}

		fmt.Println("Subscribed to Transfer events...")

		for {
			select {
			case err := <-sub.Err():
				sub.Unsubscribe()
				log.Fatalf("Error with subscription: %v", err)
			case vLog := <-logsChan:
				if e, err := SequencerInbox.SequencerInboxFilterer.ParseInboxMessageDelivered(vLog); err == nil {
					fmt.Printf("ParseInboxMessageDelivered: %s\n", e.Raw.TxHash.Hex())
					// logs.PrintFromatter(utils.ConvertBytesToHex(e))
				}
				if e, err := SequencerInbox.SequencerInboxFilterer.ParseInboxMessageDeliveredFromOrigin(vLog); err == nil {
					fmt.Printf("ParseInboxMessageDeliveredFromOrigin: %s\n", e.Raw.TxHash.Hex())
					// logs.PrintFromatter(utils.ConvertBytesToHex(e))
				}
				if e, err := SequencerInbox.SequencerInboxFilterer.ParseOwnerFunctionCalled(vLog); err == nil {
					fmt.Printf("\rParseOwnerFunctionCalled: %s\n", e.Raw.TxHash.Hex())
					// txLib := tx.NewTxLib(parentClient, &e.Raw.TxHash)
					// receipt, _ := txLib.GetTransactionReceipt()

					block, _ := parentClient.BlockByNumber(context.Background(), big.NewInt(int64(e.Raw.BlockNumber)))
					logs.PrintFromatter(utils.ConvertBytesToHex(block.Time()))
					// logs.PrintFromatter(utils.ConvertBytesToHex(e))
				}
				if e, err := SequencerInbox.SequencerInboxFilterer.ParseSequencerBatchDelivered(vLog); err == nil {
					fmt.Printf("ParseSequencerBatchDelivered: %s\n", e.Raw.TxHash.Hex())
					// logs.PrintFromatter(utils.ConvertBytesToHex(e))
				}

				if err != nil {
					log.Fatalf("Undefined: %v", err)
				}

			}
		}
	},
}
