package prompt

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
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
