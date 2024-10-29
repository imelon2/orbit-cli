package config

import (
	"sort"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

type Provider struct {
	L1 *ethclient.Client
	L2 *ethclient.Client
	L3 *ethclient.Client
}

func GetChainsKeys() []string {
	configChains := viper.GetStringMap("providers")
	var chains []string
	for key, _ := range configChains {
		chains = append(chains, key)
	}

	sort.Strings(chains)
	return chains
}

func GetProviders(key string) []string {
	return viper.GetStringSlice("providers." + key)
}
