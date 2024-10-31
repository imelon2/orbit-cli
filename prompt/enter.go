package prompt

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/imelon2/orbit-cli/common/utils"
)

func EnterTransactionHash() (string, error) {
	validationQs := []*survey.Question{
		{

			Prompt: &survey.Input{Message: "Enter transaction hash: "},
			Validate: func(val interface{}) error {
				// if the input matches the expectation
				if str := val.(string); !utils.IsTransactionHash(str) {
					return errors.New("invalid hash")
				}
				// nothing was wrong
				return nil
			},
		},
	}

	var selected string
	err := survey.Ask(validationQs, &selected)

	if err != nil {
		return "", fmt.Errorf("failed enter transaction hash: %v", err)
	}

	return selected, nil
}

func EnterProviderUrl() (string, error) {
	inputQs := &survey.Input{
		Message: "Enter the Provider URL: ",
	}

	var selected string
	err := survey.AskOne(inputQs, &selected)
	if err != nil {
		return "", fmt.Errorf("failed enter provider url: %v", err)
	}
	return selected, nil
}

func EnterBytes() (string, error) {
	validationQs := []*survey.Question{
		{

			Prompt: &survey.Input{Message: "Enter bytes data: "},
			Validate: func(val interface{}) error {
				// if the input matches the expectation
				if str := val.(string); !utils.IsBytes(str) {
					return errors.New("invalid bytes")
				}
				// nothing was wrong
				return nil
			},
		},
	}

	var selected string
	err := survey.Ask(validationQs, &selected)

	if err != nil {
		return "", fmt.Errorf("failed enter bytes: %v", err)
	}

	return selected, nil
}

func EnterInt(max int, name string) (*int, error) {
	msg := "Enter " + name + ": "
	if max != 0 {
		msg += fmt.Sprintf("(max : %d)", max)
	}

	validationQs := []*survey.Question{
		{

			Prompt: &survey.Input{Message: msg},
			Validate: func(val interface{}) error {
				inputStr := val.(string)
				inputInt, err := strconv.Atoi(inputStr)
				if err != nil {
					return errors.New("invalid number format")
				}

				// 입력값이 max보다 큰지 확인
				if max != 0 && inputInt > max {
					return errors.New("number too large")
				}
				return nil
			},
		},
	}

	selected := 0
	err := survey.Ask(validationQs, &selected)

	if err != nil {
		return nil, fmt.Errorf("failed enter count: %v", err)
	}

	return &selected, nil
}

func EnterPassword() (string, error) {
	var passwordQs = &survey.Password{Message: "Enter the password [for skip <ENTER>] : "}

	var password string
	err := survey.AskOne(passwordQs, &password, survey.WithValidator(func(val interface{}) error {
		// if the input matches the expectation
		if str := val.(string); utils.IsWithSpace(str) {
			return errors.New("invalid password : remove space")
		}
		// nothing was wrong
		return nil
	}))

	if err != nil {
		return "", fmt.Errorf("failed enter password %v", err)
	}

	return password, nil
}

func EnterString(msg string) (string, error) {
	message := "Enter " + msg + ": "
	inputQs := &survey.Input{
		Message: message,
	}

	var selected string
	err := survey.AskOne(inputQs, &selected)

	if err != nil {
		return "", fmt.Errorf("failed enter %v: %v", msg, err)
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
					return errors.New("invalid private key")
				}
				// nothing was wrong
				return nil
			},
		},
	}
	var privateKey string
	err := survey.Ask(validationQs, &privateKey)

	if err != nil {
		return "", fmt.Errorf("failed enter privateKey %v", err)
	}

	return privateKey, nil
}
