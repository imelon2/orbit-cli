/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/imelon2/orbit-cli/prompt"
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var SendCmd = &cobra.Command{
	Use:   "send",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("send called")

		_, address, err := prompt.SelectWalletForSign()

		if err != nil {
			log.Fatal("bad nextCmd : ", err)
		}

		fmt.Print(address)
	},
}

func init() {
}
