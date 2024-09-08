/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Get account list",
	Run: func(cmd *cobra.Command, args []string) {
		wallets := viper.GetStringSlice("wallets")
		for index, wallet := range wallets {
			fmt.Printf("[%d] %s\n", index, wallet)
		}
	},
}

func init() {
	// accountCmd.AddCommand(listCmd)
}
