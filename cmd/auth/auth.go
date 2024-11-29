/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"runtime"

	grant "github.com/imelon2/orbit-cli/cmd/auth/grant"
	query "github.com/imelon2/orbit-cli/cmd/auth/query"
	revoke "github.com/imelon2/orbit-cli/cmd/auth/revoke"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			log.Fatal("bad path")
		}

		prompt.RootCmdNavigation(filename, cmd, args)
	},
}

func init() {
	AuthCmd.AddCommand(grant.GrantCmd)
	AuthCmd.AddCommand(query.QueryCmd)
	AuthCmd.AddCommand(revoke.RevokeCmd)
}
