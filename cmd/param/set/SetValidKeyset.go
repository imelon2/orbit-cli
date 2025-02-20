/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/hex"
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

// SetValidKeysetCmd represents the SetValidKeyset command
var SetValidKeysetCmd = &cobra.Command{
	Use:   "SetValidKeyset",
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

		// INIT SIGNER
		_, ks, account, err := prompt.SelectWalletForSign()
		if err != nil {
			log.Fatal(err)
		}

		auto, err := tx.GetAuthByKeystore(ks, *account, parentClient)
		if err != nil {
			log.Fatal(err)
		}

		// INIT CONTRACTs
		network, err := arbnetwork.GetNetworkInfo(childClient)
		if err != nil {
			log.Fatal(err)
		}
		UpgradeExecutor, err := network.NewUpgradeExecutor(parentClient)
		if err != nil {
			log.Fatal(err)
		}

		// GET PARAM
		newVaildKeysetStr, err := prompt.EnterBytes()
		if err != nil {
			log.Fatal(err)
		}

		bytes, err := hex.DecodeString(utils.Unhexlify(newVaildKeysetStr))
		if err != nil {
			log.Fatal(err)
		}

		// SEND TX
		response, err := UpgradeExecutor.SetValidKeyset(auto, bytes)
		if err != nil {
			log.Fatal(err)
		}

		receipt, err := bind.WaitMined(context.Background(), parentClient, response)
		if err != nil {
			log.Fatal(err)
		}

		logs.PrintReceiptFromatter(receipt)
	},
}
