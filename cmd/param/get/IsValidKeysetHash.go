/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	arbnetwork "github.com/imelon2/orbit-cli/arbNetwork"
	"github.com/imelon2/orbit-cli/common/logs"
	"github.com/imelon2/orbit-cli/common/utils"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/spf13/cobra"
)

// IsValidKeysetHashCmd represents the IsValidKeysetHash command
var IsValidKeysetHashCmd = &cobra.Command{
	Use:   "IsValidKeysetHash",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

		// INIT PROVIDER
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

		// INIT CONTRACT
		network, err := arbnetwork.GetNetworkInfo(childClient)
		if err != nil {
			log.Fatal(err)
		}

		sequencerInbox, err := network.NewSequencerInbox(parentClient)
		if err != nil {
			log.Fatal(err)
		}

		Callopts := &bind.CallOpts{
			Pending: false,
			Context: nil,
		}

		// GET PARAM
		keyHashStr, err := prompt.EnterBytes32()
		if err != nil {
			log.Fatal(err)
		}

		bytes, err := hex.DecodeString(utils.Unhexlify(keyHashStr))
		if err != nil {
			log.Fatal(err)
		}

		var fixedBytes [32]byte
		copy(fixedBytes[:], bytes)

		isVaild, err := sequencerInbox.IsValidKeysetHash(Callopts, fixedBytes)
		if err != nil {
			log.Fatal(err)
		}
		logs.PrintFromatter(utils.ConvertBytesToHex(map[string]bool{
			"IsVaild KeyHash?": isVaild,
		}))
	},
}
