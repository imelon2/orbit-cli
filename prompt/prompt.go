package prompt

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sort"

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

const (
	PROMPT_SELECT_WALLET_ERROR   = "Failed Select Wallet"
	PROMPT_SELECT_PROVIDER_ERROR = "Failed Select Provider"
	PROMPT_SELECT_CHAIN_ERROR    = "Failed Select Chain"

	PROMPT_ENTER_PROVIDER_URL_ERROR     = "Failed Enter Provider"
	PROMPT_ENTER_TRANSACTION_HASH_ERROR = "Failed Enter Transaction Hash"
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
		return "", fmt.Errorf("%v : %v\n", PROMPT_SELECT_WALLET_ERROR, err)
	}

	selectedWallet := selected

	if selected == LAST_WALLET_STRING {

		selected, err := EnterAddress("wanted")
		if err != nil {
			return "", err
		}
		selectedWallet = selected
	}

	return selectedWallet, nil
}

func SelectWalletForSign() (accounts.Wallet, *keystore.KeyStore, accounts.Account, error) {
	path := utils.GetKeystoreDir()
	ks := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)

	var addressList []string
	_accounts := ks.Accounts()

	for _, wallet := range _accounts {
		addressList = append(addressList, wallet.Address.Hex())
	}

	var qs = &survey.Select{
		Message: "Select Wallet: ",
		Options: addressList,
	}

	answerIndex := 0
	err := survey.AskOne(qs, &answerIndex)
	if err != nil {
		return nil, nil, accounts.Account{}, fmt.Errorf("%v : %v\n", PROMPT_SELECT_WALLET_ERROR, err)
	}

	var pw string = ""
	var wallet accounts.Wallet

	err = ks.Unlock(_accounts[answerIndex], pw)
	if err == keystore.ErrDecrypt {
		var validationQs = []*survey.Question{
			{
				Name:   "Password",
				Prompt: &survey.Password{Message: "Enter the password [for skip <ENTER>]: "},
				Validate: func(val interface{}) error {
					err = ks.Unlock(_accounts[answerIndex], val.(string))

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
		return nil, nil, accounts.Account{}, fmt.Errorf("Failed Select Wallet For Sign :  %v\n", err)
	}

	wallet = ks.Wallets()[answerIndex]

	return wallet, ks, _accounts[answerIndex], nil
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
		return "", fmt.Errorf("Failed Select Command : %v\n", err)
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

	sort.Strings(_chains)

	_chains = append(_chains, LAST_PROVIDER_STRING)
	var selectQs = &survey.Select{
		Message: "Select Chain: ",
		Options: _chains,
	}

	err := survey.AskOne(selectQs, &selectedChain)
	if err != nil {
		return "", fmt.Errorf("%v : %v\n", PROMPT_SELECT_PROVIDER_ERROR, err)
	}

	if selectedChain == LAST_PROVIDER_STRING {
		inputQs := &survey.Input{
			Message: "Enter the Provider URL: ",
		}

		err := survey.AskOne(inputQs, &selectedProvider)
		if err != nil {
			return "", fmt.Errorf("%v : %v\n", PROMPT_ENTER_PROVIDER_URL_ERROR, err)
		}

	} else {
		_providers := viper.GetStringMapString("providers." + selectedChain)

		providers := make([]ProviderUrl, 0)
		for k, v := range _providers {
			providers = append(providers, ProviderUrl{k, v})
		}

		sort.Slice(providers, func(i, j int) bool {
			return providers[i].Name < providers[j].Name
		})

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
			return "", fmt.Errorf("%v : %v\n", PROMPT_SELECT_PROVIDER_ERROR, err)
		}
		selectedProvider = providers[answerIndex].Url

		if selectedProvider == "" {
			errM := selectedChain + "-" + providers[answerIndex].Name + " Chain No Provider"
			return "", fmt.Errorf("%v\n", errM)
		}
	}

	return selectedProvider, nil
}

func SelectProviderOrBytes() (string, bool, error) {
	var selectedChain string
	var selectedProvider string

	var _chains []string
	chains := viper.GetStringMap("providers")
	for key, _ := range chains {
		_chains = append(_chains, key)
	}

	sort.Strings(_chains)

	_chains = append(_chains, LAST_PROVIDER_STRING)
	_chains = append(_chains, LAST_CALLDATA_STRING)
	var selectQs = &survey.Select{
		Message: "Select Chain: ",
		Options: _chains,
	}

	err := survey.AskOne(selectQs, &selectedChain)
	if err != nil {
		return "", false, fmt.Errorf("%v : %v\n", PROMPT_SELECT_CHAIN_ERROR, err)
	}

	if selectedChain == LAST_CALLDATA_STRING {
		inputQs := &survey.Input{
			Message: "Enter the calldata type of bytes: ",
		}
		err := survey.AskOne(inputQs, &selectedProvider)

		if err != nil {
			return "", false, fmt.Errorf("Failed Enter Bytes %v\n", err)
		}

		return selectedProvider, false, nil

	} else if selectedChain == LAST_PROVIDER_STRING {
		inputQs := &survey.Input{
			Message: "Enter the Provider URL: ",
		}

		err := survey.AskOne(inputQs, &selectedProvider)
		if err != nil {
			return "", false, fmt.Errorf("%v : %v\n", PROMPT_ENTER_PROVIDER_URL_ERROR, err)
		}

	} else {
		_providers := viper.GetStringMapString("providers." + selectedChain)
		providers := make([]ProviderUrl, 0)
		for k, v := range _providers {
			providers = append(providers, ProviderUrl{k, v})
		}

		sort.Slice(providers, func(i, j int) bool {
			return providers[i].Name < providers[j].Name
		})

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
			return "", false, fmt.Errorf("%v : %v\n", PROMPT_SELECT_PROVIDER_ERROR, err)
		}
		selectedProvider = providers[answerIndex].Url

		if selectedProvider == "" {
			errM := selectedChain + "-" + providers[answerIndex].Name + " Chain No Provider"
			return "", false, fmt.Errorf("%v\n", errM)
		}
	}

	return selectedProvider, true, nil
}

// true = Parent -> Child ||
// false = Child -> Parent
func SelectChainTo() (bool, error) {

	inputQs := &survey.Select{
		Message: "Select forward direction",
		Options: []string{
			"Parent -> Child",
			"Parent <- Child",
		},
	}
	answerIndex := 0
	err := survey.AskOne(inputQs, &answerIndex)

	if err != nil {
		return false, fmt.Errorf("Failed Select Direction About Chain : %v\n", err)
	}

	if answerIndex == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func EnterAddress(name string) (string, error) {
	var validationQs = []*survey.Question{
		{
			Name:   "Address",
			Prompt: &survey.Input{Message: "Enter the " + name + " Address: "},
			Validate: func(val interface{}) error {
				// if the input matches the expectation
				if str := val.(string); !utils.IsAddress(str) {
					return errors.New("Invalid Address")
				}
				return nil
			},
		},
	}

	var selectedAddress string
	err := survey.Ask(validationQs, &selectedAddress)
	if err != nil {
		return "", fmt.Errorf("Failed Enter Address :  %v\n", err)
	}

	return selectedAddress, nil
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
		return common.HexToHash(""), fmt.Errorf("%v : %v\n", PROMPT_ENTER_TRANSACTION_HASH_ERROR, err)
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
		return "", fmt.Errorf("%v : %v\n", PROMPT_ENTER_TRANSACTION_HASH_ERROR, err)
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
		return "", fmt.Errorf("Failed Enter PrivateKey %v\n", err)
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
		return "", fmt.Errorf("Failed Enter Password %v\n", err)
	}

	return password, nil
}

func EnterRecipient() (string, error) {
	var recipientQs = &survey.Input{Message: "Enter the recipient address : "}

	var to string
	err := survey.AskOne(recipientQs, &to, survey.WithValidator(func(val interface{}) error {
		if str := val.(string); !utils.IsAddress(str) {
			return errors.New("Invalid Address")
		}
		return nil
	}))

	if err != nil {
		return "", fmt.Errorf("Failed Enter Recipient : %v\n", err)
	}

	return to, nil
}

func EnterValue(name string) (*big.Int, error) {
	var valueQs = &survey.Input{Message: "Enter the " + name + " Value(float) [Set 0 value <ENTER>] : "}

	var value string
	err := survey.AskOne(valueQs, &value)

	if err != nil {
		return nil, fmt.Errorf("Failed Enter Value : %v\n", err)
	}

	etherInWei := new(big.Float)
	etherInWei.SetString(value)

	wei := utils.FloatToWei(etherInWei)

	return wei, nil
}
