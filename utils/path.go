package utils

import (
	"path/filepath"
	"runtime"
)

func GetRootDir(filename string) string {
	root := filepath.Dir(filename)
	return root
}

func GetParentRootDir(filename string) string {
	root := filepath.Dir(filename)
	parent := filepath.Dir(root)
	return parent
}

func GetAbiDir() string {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Dir(filename)
	parent := filepath.Dir(root)
	abiPath := filepath.Join(parent, "solgen", "abi", "aggregateAbi.json")
	return abiPath
}

func GetKeystoreDir() string {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Dir(filename)
	parent := filepath.Dir(root)
	keystorePath := filepath.Join(parent, "keystore", "accounts")
	return keystorePath
}

func GetContractNetworkDir() string {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Dir(filename)
	parent := filepath.Dir(root)
	networkPath := filepath.Join(parent, "contractgen", "network")
	return networkPath
}
