/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"log"

	arblib "github.com/imelon2/orbit-cli/arbLib"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
)

// retryableCmd represents the retryable command
var RetryableCmd = &cobra.Command{
	Use:   "retryable",
	Short: "Decode Retryable Data form event InboxMessageDelivered",
	Run: func(cmd *cobra.Command, args []string) {
		bytes, err := prompt.EnterTransactionHashOrBytes("Bytes Data")
		if err != nil {
			log.Fatal(err)
		}

		bytes = bytes[2:] // remove 0x
		retryableBytes, err := hex.DecodeString(bytes)
		if err != nil {
			log.Fatalf("failed decode bytes data: %v", err)
		}

		data := retryableBytes

		retry := arblib.ParseRetryableMessage(data)
		utils.PrintPrettyJson(retry)
	},
}

func init() {
}
