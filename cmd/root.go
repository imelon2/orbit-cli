/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"runtime"

	account "github.com/imelon2/orbit-cli/cmd/account"
	auth "github.com/imelon2/orbit-cli/cmd/auth"
	param "github.com/imelon2/orbit-cli/cmd/param"
	parse "github.com/imelon2/orbit-cli/cmd/parse"
	search "github.com/imelon2/orbit-cli/cmd/search"
	"github.com/imelon2/orbit-cli/common/path"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "orbit-cli",
	Short: "A brief description of your application",
	Run: func(cmd *cobra.Command, args []string) {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			log.Fatal("bad path")
		}

		prompt.RootCmdNavigation(filename, cmd, args)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(parse.ParseCmd)
	rootCmd.AddCommand(search.SearchCmd)
	rootCmd.AddCommand(param.ParamCmd)
	rootCmd.AddCommand(account.AccountCmd)
	rootCmd.AddCommand(auth.AuthCmd)

	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	configPath := path.GetConfigPath()
	cobra.OnInitialize(func() {
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %s", err)
		}
	})
}
