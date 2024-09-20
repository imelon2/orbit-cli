/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	account "github.com/imelon2/orbit-cli/cmd/account"
	getBalance "github.com/imelon2/orbit-cli/cmd/getBalance"
	parse "github.com/imelon2/orbit-cli/cmd/parse"
	system "github.com/imelon2/orbit-cli/cmd/system"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "orbit",
	Short: "Orbit 테스트를 위한 CLI Tool",
	RunE: func(cmd *cobra.Command, args []string) error {

		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			log.Fatal("bad path")
		}

		root := utils.GetRootDir(filename)

		selected, err := prompt.SelectCommand(root)
		if err != nil {
			log.Fatal("bad SelectCommand")
		}

		nextCmd, _, err := cmd.Root().Find([]string{selected})
		if err != nil {
			return err
		}
		nextCmd.Run(nextCmd, args)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(account.AccountCmd)
	rootCmd.AddCommand(getBalance.GetBalanceCmd)
	rootCmd.AddCommand(parse.ParseCmd)
	rootCmd.AddCommand(system.SystemCmd)

	var cfgFile string
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config.yml)")
	// rootCmd.PersistentFlags().IntP("port", "p", 0, "Port to run the application")

	// flag := "flag"
	// rootCmd.PersistentFlags().StringVarP(&flag, "pfoo", "p", "pvar", "This is Command's Persistent Flag")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// rootCmd.Flags().UintP("name", "n", 0, "Help message for name")

	// Viper 설정 초기화
	_, filename, _, _ := runtime.Caller(0)
	parsent := utils.GetParentRootDir(filename)

	cfgFile = filepath.Join(parsent, "config.yml")
	cobra.OnInitialize(func() {
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else {
			viper.SetConfigName("config")
			viper.SetConfigType("yaml")
			viper.AddConfigPath(".")
		}

		// 설정 파일 읽기
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %s", err)
		}

		// CLI 플래그를 설정에 바인딩
		// if err := viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port")); err != nil {
		// 	log.Fatalf("Error binding flag: %s", err)
		// }
	})
}
