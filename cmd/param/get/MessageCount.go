/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	arbnetwork "github.com/imelon2/orbit-cli/arbNetwork"
	"github.com/imelon2/orbit-cli/common/logs"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/spf13/cobra"
)

type MessageCountLog struct {
	DelayedMessage   DelayedMessageCountLog
	SequencerMessage SequencerMessageCountLog
}

type DelayedMessageCountLog struct {
	TotalDelayedMessages *big.Int
	DelayedMessageCount  *big.Int
}

type SequencerMessageCountLog struct {
	SequencerMessageCount            *big.Int
	SequencerReportedSubMessageCount *big.Int
}

// MessageCountCmd represents the MessageCount command
var MessageCountCmd = &cobra.Command{
	Use:   "MessageCount",
	Short: "A brief description of your command",
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

		sequencerInbox, err := network.NewSequencerInbox(parentClient)
		if err != nil {
			log.Fatal(err)
		}

		bridge, err := network.NewBridge(parentClient)
		if err != nil {
			log.Fatal(err)
		}

		Callopts := &bind.CallOpts{
			Pending: false,
			Context: nil,
		}

		TotalDelayedMessages, err := sequencerInbox.TotalDelayedMessagesRead(Callopts)
		if err != nil {
			log.Fatal(err)
		}
		DelayedMessageCount, err := bridge.DelayedMessageCount(Callopts)
		if err != nil {
			log.Fatal(err)
		}
		SequencerMessageCount, err := bridge.SequencerMessageCount(Callopts)
		if err != nil {
			log.Fatal(err)
		}
		sequencerReportedSubMessageCount, err := bridge.SequencerReportedSubMessageCount(Callopts)
		if err != nil {
			log.Fatal(err)
		}

		logs.PrintFromatter(MessageCountLog{
			DelayedMessageCountLog{
				TotalDelayedMessages: TotalDelayedMessages,
				DelayedMessageCount:  DelayedMessageCount,
			},
			SequencerMessageCountLog{
				SequencerMessageCount:            SequencerMessageCount,
				SequencerReportedSubMessageCount: sequencerReportedSubMessageCount,
			},
		})

	},
}
