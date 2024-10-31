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
	"github.com/imelon2/orbit-cli/common/path"
	"github.com/imelon2/orbit-cli/config"
)

const (
	LAST_WALLET_STRING   = "< Enter Wallet Address >"
	LAST_PROVIDER_STRING = "< Enter Provider URL >"
	LAST_CALLDATA_STRING = "< Enter data type of bytes >"
)

var ProvidersList = []string{"ethereum", "arbitrum", "orbit"}

type ProviderUrl struct {
	Name string
	Url  string
}

func SelectNextCmd(rootFileName string) (string, error) {
	root := filepath.Dir(rootFileName)
	files, err := os.ReadDir(root)

	if err != nil {
		return "", fmt.Errorf("failed read dir : %v", err)
	}

	var options []string
	for _, file := range files {
		if file.IsDir() {
			options = append(options, file.Name())
		} else {
			name := file.Name()[:len(file.Name())-3]
			_, fileName := filepath.Split(root)

			if name != fileName && name != "root" {
				options = append(options, name)
			}
		}
	}

	var qs = &survey.Select{
		Message: "Select Command: ",
		Options: options,
	}

	var selected string
	err = survey.AskOne(qs, &selected)

	if err != nil {
		return "", fmt.Errorf("failed select command : %v", err)
	}

	return selected, nil
}

func SelectChains(add ...string) (string, error) {
	chainKeys := config.GetChainsKeys()
	chainKeys = append(chainKeys, add...)

	selectQs := &survey.Select{
		Message: "Select Chains: ",
		Options: chainKeys,
	}

	var selected string
	err := survey.AskOne(selectQs, &selected)
	if err != nil {
		return "", fmt.Errorf("failed select chains: %v", err)
	}

	return selected, nil
}

func SelectProviders(chain string) (string, error) {
	providers := config.GetProviders(chain)
	providerName := make([]ProviderUrl, len(providers))
	for i, v := range providers {
		providerName[i] = ProviderUrl{ProvidersList[i], v}
	}

	titles := make([]string, len(providerName))
	for i, m := range providerName {
		titles[i] = m.Name
	}

	selectQs := &survey.Select{
		Message: "Select provider: ",
		Options: titles,
		Description: func(value string, index int) string {
			return providerName[index].Url
		},
	}

	selected := 0
	err := survey.AskOne(selectQs, &selected)
	if err != nil {
		return "", fmt.Errorf("failed select provider : %v", err)
	}

	return providerName[selected].Url, nil
}

func SelectCrossChainProviders(chain string) (string, string, error) {
	providers := config.GetProviders(chain)
	childProviders := providers[1:]
	providerName := make([]ProviderUrl, len(childProviders))
	for i, v := range childProviders {
		providerName[i] = ProviderUrl{ProvidersList[i+1], v}
	}

	titles := make([]string, len(providerName))
	for i, m := range providerName {
		titles[i] = m.Name
	}

	selectQs := &survey.Select{
		Message: "Select provider: ",
		Options: titles,
		Description: func(value string, index int) string {
			return providerName[index].Url
		},
	}

	selected := 0
	err := survey.AskOne(selectQs, &selected)
	if err != nil {
		return "", "", fmt.Errorf("failed select provider : %v", err)
	}

	return providers[selected], providers[selected+1], nil
}

type KeystoreTag struct {
	Tag      string
	Keystore *keystore.KeyStore
}

func SelectWalletForSign() (*accounts.Wallet, *keystore.KeyStore, *accounts.Account, error) {
	keystorePath := path.GetKeystoreDir("")
	files, err := os.ReadDir(keystorePath)
	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 {
		return nil, nil, nil, fmt.Errorf("No keystore was created.")
	}

	var KeystoreList []KeystoreTag
	for _, file := range files {
		pathTag := path.GetKeystoreDir(file.Name())
		ks := keystore.NewKeyStore(pathTag, keystore.StandardScryptN, keystore.StandardScryptP)
		KeystoreList = append(KeystoreList, KeystoreTag{
			Tag:      file.Name(),
			Keystore: ks,
		})
	}

	var addressList []string
	for _, wallet := range KeystoreList {
		accounts := wallet.Keystore.Accounts()
		addressList = append(addressList, accounts[0].Address.Hex())
	}

	var qs = &survey.Select{
		Message: "Select Wallet: ",
		Options: addressList,
		Description: func(value string, index int) string {
			return KeystoreList[index].Tag
		},
	}

	answerIndex := 0
	err = survey.AskOne(qs, &answerIndex)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("fail select wallet : %v", err)
	}

	var pw string = ""
	var wallet accounts.Wallet

	ks := KeystoreList[answerIndex].Keystore
	account := ks.Accounts()
	err = ks.Unlock(account[answerIndex], pw)

	if err == keystore.ErrDecrypt {
		var validationQs = []*survey.Question{
			{
				Name:   "Password",
				Prompt: &survey.Password{Message: "Enter the password [for skip <ENTER>]: "},
				Validate: func(val interface{}) error {
					err = ks.Unlock(account[answerIndex], val.(string))

					if err == keystore.ErrDecrypt {
						return errors.New("Invaild Password :" + err.Error() + "\n")
					}
					// nothing was wrong
					return nil
				},
			},
		}
		err = survey.Ask(validationQs, &pw)
	} else if err != nil {
		return nil, nil, nil, fmt.Errorf("failed select wallet for sign: %v", err)
	}

	wallet = ks.Wallets()[answerIndex]

	return &wallet, ks, &account[answerIndex], nil
}
