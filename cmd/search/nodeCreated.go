/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	arbnetwork "github.com/imelon2/orbit-cli/arbNetwork"
	"github.com/imelon2/orbit-cli/common/logs"
	"github.com/imelon2/orbit-cli/common/tx"
	"github.com/imelon2/orbit-cli/common/utils"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/spf13/cobra"
)

// NodeCreatedCmd represents the NodeCreated command
var NodeCreatedCmd = &cobra.Command{
	Use:   "NodeCreated",
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
		RollupCore, err := network.NewRollupCore(parentClient)
		if err != nil {
			log.Fatal(err)
		}

		Callopts := &bind.CallOpts{
			Pending: false,
			Context: nil,
		}
		max, err := RollupCore.RollupCoreCaller.LatestNodeCreated(Callopts)
		if err != nil {
			log.Fatal(err)
		}
		count, err := prompt.EnterInt(int(max), "NodeCreated events")
		if err != nil {
			log.Fatal(err)
		}

		eventFunc := func(opt bind.FilterOpts) ([]interface{}, error) {
			iterator, err := RollupCore.RollupCoreFilterer.FilterNodeCreated(&opt, nil, nil, nil)
			if err != nil {
				return nil, fmt.Errorf("fail FilterNodeCreated : %s", err)
			}

			_events := make([]interface{}, 0)
			for iterator.Next() {
				event := iterator.Event
				e := arbnetwork.NodeCreatedEvent{
					NodeNum:        event.NodeNum,
					ParentNodeHash: event.ParentNodeHash,
					NodeHash:       event.NodeHash,
					ExecutionHash:  event.ExecutionHash,
					// Assertion:          event.Assertion,
					AfterInboxBatchAcc: event.AfterInboxBatchAcc,
					WasmModuleRoot:     event.WasmModuleRoot,
					InboxMaxCount:      event.InboxMaxCount,
					TransactionHash:    &event.Raw.TxHash,
				}

				_events = append(_events, e)
			}
			_events = utils.Reverse(_events)
			return _events, nil
		}
		result, err := tx.SearchEvent(*count, parentClient, eventFunc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print("\n\n")
		logs.PrintFromatter(utils.ConvertBytesToHex(result))
	},
}
