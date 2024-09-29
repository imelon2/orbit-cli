package prompt

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/imelon2/orbit-cli/utils"
)

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

func EnterTransactionHashOrBytes(name string) (string, error) {
	var validationQs = []*survey.Question{
		{
			Name:   "HashOrBytes",
			Prompt: &survey.Input{Message: "Enter the " + name + ": "},
			Validate: func(val interface{}) error {
				// if the input matches the expectation
				if str := val.(string); !utils.IsBytes(str) {
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
