/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethLib "github.com/imelon2/orbit-cli/ethLib"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/solgen/go/precompilesgen"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
)

const (
	setL1PricingRewardRecipient = iota
	setInfraFeeAccount
	setNetworkFeeAccount
)

var setAccountsCommand = []string{"SetL1PricingRewardRecipient", "SetInfraFeeAccount", "SetNetworkFeeAccount"}

// setAccountsCmd represents the setAccounts command
var SetAccountsCmd = &cobra.Command{
	Use:   "setAccounts",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

		var qs = &survey.Select{
			Message: "Select Command: ",
			Options: setAccountsCommand,
		}

		answerIndex := 0
		err := survey.AskOne(qs, &answerIndex)
		if err != nil {
			log.Fatal(err)
		}

		ArbOwnerLibs := NewArbOwnerLibsPrompt()
		client := ArbOwnerLibs.client

		var response *types.Transaction
		switch answerIndex {
		case setL1PricingRewardRecipient:
			response = ArbOwnerLibs.SetL1PricingRewardRecipient()
		case setInfraFeeAccount:
			response = ArbOwnerLibs.SetInfraFeeAccount()
		case setNetworkFeeAccount:
			response = ArbOwnerLibs.SetNetworkFeeAccount()
		}

		fmt.Print("\n\nTransaction Response: \n")
		utils.PrintPrettyJson(response)
		fmt.Print("\n\nWait Mined Transaction ... \n\n")

		receipt, err := bind.WaitMined(context.Background(), client, response)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print("Transaction receipt: \n")
		utils.PrintPrettyJson(receipt)
	},
}

func init() {

}

type ArbOwnerLibs struct {
	client   *ethclient.Client
	auth     *bind.TransactOpts
	ArbOwner *precompilesgen.ArbOwner
}

func NewArbOwnerLibs(client *ethclient.Client) *ArbOwnerLibs {
	new := new(ArbOwnerLibs)
	ArbOwner, err := precompilesgen.NewArbOwner(types.ArbOwnerAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	new.ArbOwner = ArbOwner
	return new
}

func NewArbOwnerLibsPrompt() *ArbOwnerLibs {
	new := new(ArbOwnerLibs)

	client, auth, err := ethLib.GenerateAuth()
	if err != nil {
		log.Fatal(err)
	}

	new.client = client
	new.auth = auth
	ArbOwner, err := precompilesgen.NewArbOwner(types.ArbOwnerAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	new.ArbOwner = ArbOwner

	return new
}

func (lib *ArbOwnerLibs) SetL1PricingRewardRecipient() *types.Transaction {
	newRewarderAccount, err := prompt.EnterAddress("new L1 Rewarder account")
	if err != nil {
		log.Fatal(err)
	}

	response, err := lib.ArbOwner.SetL1PricingRewardRecipient(lib.auth, common.HexToAddress(newRewarderAccount))
	if err != nil {
		log.Fatal(err)
	}

	return response
}

func (lib *ArbOwnerLibs) SetInfraFeeAccount() *types.Transaction {
	newInfraAccount, err := prompt.EnterAddress("new infra account")
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := lib.ArbOwner.SetInfraFeeAccount(lib.auth, common.HexToAddress(newInfraAccount))

	if err != nil {
		log.Fatal(err)
	}

	return signedTx
}

func (lib *ArbOwnerLibs) SetNetworkFeeAccount() *types.Transaction {
	newInfraAccount, err := prompt.EnterAddress("new Network Fee account")
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := lib.ArbOwner.SetNetworkFeeAccount(lib.auth, common.HexToAddress(newInfraAccount))

	if err != nil {
		log.Fatal(err)
	}

	return signedTx
}
