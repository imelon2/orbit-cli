package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type HardHatArtifact struct {
	Format       string        `json:"_format"`
	ContractName string        `json:"contractName"`
	SourceName   string        `json:"sourceName"`
	Abi          []interface{} `json:"abi"`
	Bytecode     string        `json:"bytecode"`
}

type FoundryBytecode struct {
	Object string `json:"object"`
}

type FoundryArtifact struct {
	Abi      []interface{}   `json:"abi"`
	Bytecode FoundryBytecode `json:"bytecode"`
}

type moduleInfo struct {
	contractNames []string
	abis          []string
	bytecodes     []string
}

func (m *moduleInfo) addArtifact(artifact HardHatArtifact) {
	abi, err := json.Marshal(artifact.Abi)
	if err != nil {
		log.Fatal(err)
	}
	m.contractNames = append(m.contractNames, artifact.ContractName)
	m.abis = append(m.abis, string(abi))
	m.bytecodes = append(m.bytecodes, artifact.Bytecode)
}

func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("bad path")
	}
	root := filepath.Dir(filename)
	parent := filepath.Dir(root)
	filePaths, err := filepath.Glob(filepath.Join(parent, "nitro-contracts", "build", "contracts", "src", "*", "*.sol", "*.json"))
	if err != nil {
		log.Fatal(err)
	}

	patterns := []string{
		filepath.Join(parent, "token-bridge-contracts", "build", "contracts", "contracts", "tokenbridge", "*", "*.sol", "*.json"),
		filepath.Join(parent, "token-bridge-contracts", "build", "contracts", "contracts", "tokenbridge", "*", "*", "*.sol", "*.json"),
	}

	var filePathsTokenBridgeContract []string

	// 각 패턴에 대해 filepath.Glob을 호출하고 결과를 합침
	for _, pattern := range patterns {
		m, err := filepath.Glob(pattern)
		if err != nil {
			fmt.Printf("Error matching pattern %s: %v\n", pattern, err)
			continue
		}
		filePathsTokenBridgeContract = append(filePathsTokenBridgeContract, m...)
	}

	filePaths = append(filePaths, filePathsTokenBridgeContract...)

	modules := make(map[string]*moduleInfo) // 모든 abis, bytecodes, contractNames 수집

	for _, path := range filePaths {
		if strings.Contains(path, ".dbg.json") {
			continue
		}
		// path의 가장 마지막 경로 분리
		dir, file := filepath.Split(path)

		// dir[:len(dir)-1] : Path의 맨 끝 "/" 제거
		// path의 가장 마지막 경로 분리
		dir, _ = filepath.Split(dir[:len(dir)-1])
		_, module := filepath.Split(dir[:len(dir)-1])
		module = strings.ReplaceAll(module, "-", "_")
		module += "gen"

		// file[:len(file)-5] : .json 제거
		name := file[:len(file)-5]

		// fmt.Printf("%d \n", len(file)-5)

		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatal("could not read", path, "for contract", name, err)
		}

		artifact := HardHatArtifact{} // abi
		if err := json.Unmarshal(data, &artifact); err != nil {
			log.Fatal("failed to parse contract", name, err)
		}

		modInfo := modules[module]
		if modInfo == nil {
			modInfo = &moduleInfo{}
			modules[module] = modInfo
		}
		modInfo.addArtifact(artifact)
	}

	var aggregate []map[string]interface{}

	for module, info := range modules {

		for _, _abi := range info.abis {
			var functions []map[string]interface{}
			json.Unmarshal([]byte(_abi), &functions)
			if err != nil {
				panic(err)
			}
			aggregate = append(aggregate, functions...)
		}

		code, err := bind.Bind(
			info.contractNames,
			info.abis,
			info.bytecodes,
			nil,
			module,
			bind.LangGo,
			nil,
			nil,
		)
		if err != nil {
			log.Fatal(err)
		}

		folder := filepath.Join(root, "go", module)

		err = os.MkdirAll(folder, 0o755)
		if err != nil {
			log.Fatal(err)
		}

		/*
			#nosec G306
			This file contains no private information so the permissions can be lenient
		*/
		err = os.WriteFile(filepath.Join(folder, module+".go"), []byte(code), 0o644) // 664
		if err != nil {
			log.Fatal(err)
		}
	}

	// 중복을 추적하기 위한 맵
	// 맨 처음 값만 들어옴
	typeTracker := make(map[string]bool)

	// 중복 제거를 위한 슬라이스
	var result []map[string]interface{}

	for _, entry := range aggregate {
		if entryType, ok := entry["type"].(string); ok {
			if entryType == "fallback" || entryType == "receive" {
				// 중복된 "fallback" 및 "receive" 항목은 무시
				if !typeTracker[entryType] {
					typeTracker[entryType] = true
					result = append(result, entry)
				}
			} else {
				// 다른 타입의 항목은 항상 추가
				result = append(result, entry)
			}
		}
	}

	folder := filepath.Join(root, "abi", "aggregateAbi.json")

	aggregateBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(folder, aggregateBytes, 0o644)

	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Println("successfully generated go abi files")
}
