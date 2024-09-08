package utils

import (
	"path/filepath"
	"runtime"
)

func GetRootDir(filename string) string {
	root := filepath.Dir(filename)
	return root
}

func GetAbiDir() string {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Dir(filename)
	parent := filepath.Dir(root)
	abiPath := filepath.Join(parent, "solgen", "abi", "aggregateAbi.json")
	return abiPath
}
