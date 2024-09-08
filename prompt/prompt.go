package prompt

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/imelon2/orbit-toolkit/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

var LAST_WALLET_STRING = "Enter"
var LAST_PROVIDER_STRING = "Enter Provider URL"

func SelectWallet() (string, error) {
	wallets := viper.GetStringSlice("wallets")
	wallets = append(wallets, LAST_WALLET_STRING)
	promptSelect := promptui.Select{
		Label: "Select Wallet",
		Items: wallets,
	}

	_, selected, err := promptSelect.Run()

	if err != nil {
		return "", fmt.Errorf("Prompt failed %v\n", err)
	}

	selectedWallet := selected

	if selected == LAST_WALLET_STRING {
		validate := func(input string) error {
			isAddress := utils.IsAddress(input)
			if !isAddress {
				return errors.New("Invalid Address")
			}
			return nil
		}

		promptPrompt := promptui.Prompt{
			Label:    "Enter the wallet address",
			Validate: validate,
		}

		selected, err := promptPrompt.Run()
		if err != nil {
			return "", fmt.Errorf("Prompt failed %v\n", err)
		}
		selectedWallet = selected
	}

	return selectedWallet, nil
}

func SelectProvider() (string, error) {
	var _chains []string
	var selectedProvider string
	var selectedChain string

	chains := viper.GetStringMap("providers")
	for key, _ := range chains {
		_chains = append(_chains, key)
	}
	_chains = append(_chains, LAST_PROVIDER_STRING)

	promptSelectChain := promptui.Select{
		Label: "Select Chain",
		Items: _chains,
	}

	_, selectedChain, err := promptSelectChain.Run()
	if err != nil {
		return "", fmt.Errorf("Prompt failed %v\n", err)
	}
	if selectedChain == LAST_PROVIDER_STRING {
		validate := func(input string) error {
			return nil
		}

		promptPrompt := promptui.Prompt{
			Label:    "Enter the Provider URL",
			Validate: validate,
		}

		selected, err := promptPrompt.Run()
		if err != nil {
			return "", fmt.Errorf("Prompt failed %v\n", err)
		}
		selectedProvider = selected
	} else {
		providers := viper.GetStringMapString("providers." + selectedChain)
		var _providers []string
		for key := range providers {
			_providers = append(_providers, key)
		}

		promptSelectProvider := promptui.Select{
			Label: "Select Provider",
			Items: _providers,
		}

		_, selectedProviderName, err := promptSelectProvider.Run()
		if err != nil {
			return "", fmt.Errorf("Prompt failed %v\n", err)
		}
		selectedProvider = providers[selectedProviderName]

		if selectedProvider == "" {
			errM := selectedChain + "-" + selectedProviderName + " Chain No Provider"
			log.Fatal(errM)
		}
	}

	return selectedProvider, nil
}

func SelectCommand(dirPath string) (string, error) {

	var directories []string
	// 디렉토리 안의 파일 및 디렉토리 목록 읽기
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	// 각 파일이 디렉토리인지 확인
	for _, file := range files {
		if file.IsDir() {
			directories = append(directories, file.Name())
		} else {
			name := file.Name()[:len(file.Name())-3]
			_, fileName := filepath.Split(dirPath)

			if name != fileName && name != "root" {
				directories = append(directories, name)
			}
		}
	}

	promptSelect := promptui.Select{
		Label: "Select Command",
		Items: directories,
	}

	_, selected, err := promptSelect.Run()
	if err != nil {
		return "", err
	}

	return selected, nil
}

func EnterTransactionHash() (common.Hash, error) {
	validate := func(input string) error {
		isAddress := utils.IsTransaction(input)
		if !isAddress {
			return errors.New("Invalid transaction hash")
		}
		return nil
	}

	promptPrompt := promptui.Prompt{
		Label:    "Enter the transaction hash",
		Validate: validate,
	}

	selected, err := promptPrompt.Run()
	if err != nil {
		return common.HexToHash(""), fmt.Errorf("EnterTransactionHash Prompt failed %v\n", err)
	}

	return common.HexToHash(selected), nil
}
