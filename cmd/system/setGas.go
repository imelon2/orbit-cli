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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethLib "github.com/imelon2/orbit-cli/ethLib"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/solgen/go/precompilesgen"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
)

const (
	setL2BaseFee = iota
	setL1PricePerUnit
	setMinimumL2BaseFee
)

var setGasCommand = []string{"SetL2BaseFee", "SetL1PricePerUnit", "SetMinimumL2BaseFee"}

var SetGasCmd = &cobra.Command{
	Use:   "setGas",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		var qs = &survey.Select{
			Message: "Select Command: ",
			Options: setGasCommand,
		}

		answerIndex := 0
		err := survey.AskOne(qs, &answerIndex)
		if err != nil {
			log.Fatal(err)
		}

		var client *ethclient.Client
		var signedTx *types.Transaction

		switch answerIndex {
		case setL2BaseFee:
			client, signedTx = SetL2BaseFee()
		case setL1PricePerUnit:
			client, signedTx = SetL1PricePerUnit()
		case setMinimumL2BaseFee:
			client, signedTx = SetMinimumL2BaseFee()
		}

		fmt.Print("\n\nTransaction Response: \n")
		utils.PrintPrettyJson(signedTx)
		fmt.Print("\n\nWait Mined Transaction ... \n\n")

		receipt, err := bind.WaitMined(context.Background(), client, signedTx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Transaction receipt: \n")
		utils.PrintPrettyJson(receipt)
	},
}

func SetL2BaseFee() (*ethclient.Client, *types.Transaction) {
	newL2BaseFeeWei, err := prompt.EnterValue("new L2 base fee")

	if err != nil {
		log.Fatal(err)
	}

	client, auth, err := ethLib.GenerateAuth()
	if err != nil {
		log.Fatal(err)
	}

	ArbOwner, err := precompilesgen.NewArbOwner(types.ArbOwnerAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := ArbOwner.SetL2BaseFee(auth, newL2BaseFeeWei)

	if err != nil {
		log.Fatal(err)
	}

	return client, signedTx
}

func SetL1PricePerUnit() (*ethclient.Client, *types.Transaction) {
	newL1PricePerUnit, err := prompt.EnterValue("new L1 Price Per Unit")

	if err != nil {
		log.Fatal(err)
	}

	client, auth, err := ethLib.GenerateAuth()
	if err != nil {
		log.Fatal(err)
	}

	ArbOwner, err := precompilesgen.NewArbOwner(types.ArbOwnerAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := ArbOwner.SetL1PricePerUnit(auth, newL1PricePerUnit)

	if err != nil {
		log.Fatal(err)
	}

	return client, signedTx
}

func SetMinimumL2BaseFee() (*ethclient.Client, *types.Transaction) {
	newMinBaseFee, err := prompt.EnterValue("new min L2 base fee")

	if err != nil {
		log.Fatal(err)
	}

	client, auth, err := ethLib.GenerateAuth()
	if err != nil {
		log.Fatal(err)
	}

	ArbOwner, err := precompilesgen.NewArbOwner(types.ArbOwnerAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := ArbOwner.SetMinimumL2BaseFee(auth, newMinBaseFee)

	if err != nil {
		log.Fatal(err)
	}

	return client, signedTx
}
