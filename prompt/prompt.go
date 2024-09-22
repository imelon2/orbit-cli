package prompt

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/viper"
)

type ProviderUrl struct {
	Name string
	Url  string
}

const (
	LAST_WALLET_STRING   = "< Enter Wallet Address >"
	LAST_PROVIDER_STRING = "< Enter Provider URL >"
	LAST_CALLDATA_STRING = "< Enter Calldata type of bytes >"
)

func SelectWallet() (string, error) {
	path := utils.GetKeystoreDir()
	ks := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)

	var addressList []string
	accounts := ks.Accounts()

	for _, wallet := range accounts {
		addressList = append(addressList, wallet.Address.Hex())
	}
	addressList = append(addressList, LAST_WALLET_STRING)
	var qs = &survey.Select{
		Message: "Select Wallet: ",
		Options: addressList,
	}

	var selected string
	err := survey.AskOne(qs, &selected)
	if err != nil {
		return "", fmt.Errorf("Prompt failed %v\n", err)
	}

	selectedWallet := selected

	if selected == LAST_WALLET_STRING {

		var validationQs = []*survey.Question{
			{
				Name:   "Address",
				Prompt: &survey.Input{Message: "Enter the Address: "},
				Validate: func(val interface{}) error {
					// if the input matches the expectation
					if str := val.(string); !utils.IsAddress(str) {
						return errors.New("Invalid Address")
					}
					// nothing was wrong
					return nil
				},
			},
		}

		var selected string
		err := survey.Ask(validationQs, &selected)
		if err != nil {
			return "", fmt.Errorf("Prompt failed %v\n", err)
		}
		selectedWallet = selected
	}

	return selectedWallet, nil
}

func SelectWalletForSign() (accounts.Wallet, string, error) {
	path := utils.GetKeystoreDir()
	ks := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)

	var addressList []string
	accounts := ks.Accounts()

	for _, wallet := range accounts {
		addressList = append(addressList, wallet.Address.Hex())
	}

	var qs = &survey.Select{
		Message: "Select Wallet: ",
		Options: addressList,
	}

	answerIndex := 0
	err := survey.AskOne(qs, &answerIndex)
	if err != nil {
		return nil, "", fmt.Errorf("Prompt failed %v\n", err)
	}

	pw := ""
	err = ks.Unlock(accounts[answerIndex], pw)

	if err == keystore.ErrDecrypt {
		pw, err = EnterPassword()
		err = ks.Unlock(accounts[answerIndex], pw)
	} else if err != nil {
		return nil, "", fmt.Errorf("Open Wallet failed %v\n", err)
	}

	wallet := ks.Wallets()[answerIndex]

	address := accounts[answerIndex].Address.Hex()

	return wallet, address, nil
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

	var qs = &survey.Select{
		Message: "Select Command: ",
		Options: directories,
	}

	var selected string
	err = survey.AskOne(qs, &selected)

	if err != nil {
		return "", fmt.Errorf("Prompt failed %v\n", err)
	}

	return selected, nil
}

func SelectProvider() (string, error) {
	var selectedChain string
	var selectedProvider string

	var _chains []string
	chains := viper.GetStringMap("providers")
	for key, _ := range chains {
		_chains = append(_chains, key)
	}
	_chains = append(_chains, LAST_PROVIDER_STRING)
	var selectQs = &survey.Select{
		Message: "Select Chain: ",
		Options: _chains,
	}

	err := survey.AskOne(selectQs, &selectedChain)
	if err != nil {
		return "", fmt.Errorf("Prompt failed %v\n", err)
	}

	if selectedChain == LAST_PROVIDER_STRING {
		inputQs := &survey.Input{
			Message: "Enter the Provider URL: ",
		}

		err := survey.AskOne(inputQs, &selectedProvider)
		if err != nil {
			return "", fmt.Errorf("Prompt failed %v\n", err)
		}

	} else {
		_providers := viper.GetStringMapString("providers." + selectedChain)
		providers := make([]ProviderUrl, 0)
		for k, v := range _providers {
			providers = append(providers, ProviderUrl{k, v})
		}
		titles := make([]string, len(providers))
		for i, m := range providers {
			titles[i] = m.Name
		}
		var qs = &survey.Select{
			Message: "Select Provider: ",
			Options: titles,
			Description: func(value string, index int) string {
				return providers[index].Url
			},
		}

		answerIndex := 0
		err := survey.AskOne(qs, &answerIndex)
		if err != nil {
			return "", fmt.Errorf("Prompt failed %v\n", err)
		}
		selectedProvider = providers[answerIndex].Url

		if selectedProvider == "" {
			errM := selectedChain + "-" + providers[answerIndex].Name + " Chain No Provider"
			log.Fatal(errM)
		}
	}

	return selectedProvider, nil
}

func SelectProviderOrCalldata() (string, bool, error) {
	var selectedChain string
	var selectedProvider string

	var _chains []string
	chains := viper.GetStringMap("providers")
	for key, _ := range chains {
		_chains = append(_chains, key)
	}
	_chains = append(_chains, LAST_PROVIDER_STRING)
	_chains = append(_chains, LAST_CALLDATA_STRING)
	var selectQs = &survey.Select{
		Message: "Select Chain: ",
		Options: _chains,
	}

	err := survey.AskOne(selectQs, &selectedChain)
	if err != nil {
		return "", false, fmt.Errorf("Prompt failed %v\n", err)
	}

	if selectedChain == LAST_CALLDATA_STRING {
		inputQs := &survey.Input{
			Message: "Enter the calldata type of bytes: ",
		}
		err := survey.AskOne(inputQs, &selectedProvider)

		if err != nil {
			return "", false, fmt.Errorf("Prompt failed %v\n", err)
		}

		return selectedProvider, false, nil

	} else if selectedChain == LAST_PROVIDER_STRING {
		inputQs := &survey.Input{
			Message: "Enter the Provider URL: ",
		}

		err := survey.AskOne(inputQs, &selectedProvider)
		if err != nil {
			return "", false, fmt.Errorf("Prompt failed %v\n", err)
		}

	} else {
		_providers := viper.GetStringMapString("providers." + selectedChain)
		providers := make([]ProviderUrl, 0)
		for k, v := range _providers {
			providers = append(providers, ProviderUrl{k, v})
		}
		titles := make([]string, len(providers))
		for i, m := range providers {
			titles[i] = m.Name
		}

		var qs = &survey.Select{
			Message: "Select Provider: ",
			Options: titles,
			Description: func(value string, index int) string {
				return providers[index].Url
			},
		}

		answerIndex := 0
		err := survey.AskOne(qs, &answerIndex)
		if err != nil {
			return "", false, fmt.Errorf("Prompt failed %v\n", err)
		}
		selectedProvider = providers[answerIndex].Url

		if selectedProvider == "" {
			errM := selectedChain + "-" + providers[answerIndex].Name + " Chain No Provider"
			log.Fatal(errM)
		}
	}

	return selectedProvider, true, nil
}

func EnterTransactionHash() (common.Hash, error) {

	var validationQs = []*survey.Question{
		{
			Name:   "Hash",
			Prompt: &survey.Input{Message: "Enter the transaction hash: "},
			Validate: func(val interface{}) error {
				// if the input matches the expectation
				if str := val.(string); !utils.IsTransaction(str) {
					return errors.New("Invalid transaction hash")
				}
				// nothing was wrong
				return nil
			},
		},
	}
	var selected string
	err := survey.Ask(validationQs, &selected)
	if err != nil {
		return common.HexToHash(""), fmt.Errorf("EnterTransactionHash Prompt failed %v\n", err)
	}

	return common.HexToHash(selected), nil
}

func EnterTransactionHashOrBytes() (string, error) {
	var validationQs = []*survey.Question{
		{
			Name:   "HashOrBytes",
			Prompt: &survey.Input{Message: "Enter the transaction hash or calldata: "},
			Validate: func(val interface{}) error {
				// if the input matches the expectation
				if str := val.(string); len(str) < 10 {
					return errors.New("Invalid transaction hash")
				}
				// nothing was wrong
				return nil
			},
		},
	}
	var selected string
	err := survey.Ask(validationQs, &selected)

	if err != nil {
		return "", fmt.Errorf("EnterTransactionHash Prompt failed %v\n", err)
	}

	return selected, nil
}

func EnterPrivateKey() (string, error) {
	var validationQs = []*survey.Question{
		{
			Name:   "PrivateKey",
			Prompt: &survey.Input{Message: "Enter the private key: "},
			Validate: func(val interface{}) error {
				// if the input matches the expectation
				if str := val.(string); !utils.IsPrivateKey(str) {
					return errors.New("Invalid private key")
				}
				// nothing was wrong
				return nil
			},
		},
	}
	var privateKey string
	err := survey.Ask(validationQs, &privateKey)

	if err != nil {
		return "", fmt.Errorf("EnterPrivateKey Prompt failed %v\n", err)
	}

	return privateKey, nil
}

func EnterPassword() (string, error) {
	var passwordQs = &survey.Password{Message: "Enter the password [for skip <ENTER>] : "}

	var password string
	err := survey.AskOne(passwordQs, &password, survey.WithValidator(func(val interface{}) error {
		// if the input matches the expectation
		if str := val.(string); utils.IsWithSpace(str) {
			return errors.New("Invalid Password : remove space")
		}
		// nothing was wrong
		return nil
	}))

	if err != nil {
		return "", fmt.Errorf("EnterPassword Prompt failed %v\n", err)
	}

	return password, nil
}
