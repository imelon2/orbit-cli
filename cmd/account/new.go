/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var NewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new account with Priavet Key",
	Run: func(cmd *cobra.Command, args []string) {

		pw, err := prompt.EnterPassword()
		if err != nil {
			log.Fatal(err)
		}

		path := utils.GetKeystoreDir()
		account, err := keystore.StoreKey(path, pw, keystore.StandardScryptN, keystore.StandardScryptP)
		if err == keystore.ErrAccountAlreadyExists {
			fmt.Printf("\n%v\n\n", err)
			fmt.Printf("Public address of the key:   %s\n", account.Address.Hex())
			fmt.Printf("Path of the secret key file: %s\n\n", account.URL.Path)
			return
		} else if err != err {
			log.Fatalf("Failed to create account: %v", err)
		}

		fmt.Printf("\nYour new key was generated\n\n")
		fmt.Printf("Public address of the key:   %s\n", account.Address.Hex())
		fmt.Printf("Path of the secret key file: %s\n\n", account.URL.Path)
		fmt.Printf("- You can share your public address with anyone. Others need it to interact with you.\n")
		fmt.Printf("- You must NEVER share the secret key with anyone! The key controls access to your funds!\n")
		fmt.Printf("- You must BACKUP your key file! Without the key, it's impossible to access account funds!\n")
		fmt.Printf("- You must REMEMBER your password! Without the password, it's impossible to decrypt the key!\n\n")
	},
}

func init() {
}
