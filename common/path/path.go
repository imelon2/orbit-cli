package path

import (
	"path/filepath"
	"runtime"
)

func GetConfigPath() string {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Dir(filename)
	depth1 := filepath.Dir(root)
	depth2 := filepath.Dir(depth1)
	return filepath.Join(depth2, "config.yml")
}

func GetAbiPath() string {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Dir(filename)
	depth1 := filepath.Dir(root)
	depth2 := filepath.Dir(depth1)
	return filepath.Join(depth2, "solgen", "abi", "aggregateAbi.json")
}

func GetContractNetworkDir() string {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Dir(filename)
	depth1 := filepath.Dir(root)
	depth2 := filepath.Dir(depth1)
	networkPath := filepath.Join(depth2, "arbNetwork", "networks")
	return networkPath
}

func GetKeystoreDir(tag string) string {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Dir(filename)
	depth1 := filepath.Dir(root)
	depth2 := filepath.Dir(depth1)
	path := filepath.Join(depth2, "keystore", tag)
	return path
}
