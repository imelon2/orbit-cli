package contractgen_test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/imelon2/orbit-cli/contractgen"
	"github.com/imelon2/orbit-cli/utils"
)

func Test_Path(t *testing.T) {
	path := utils.GetContractNetworkDir()

	t.Log(path, "\n")
	files, err := os.ReadDir(path)
	if err != nil {
		// return "", err
	}

	var networks []contractgen.NetworkInfo

	for _, file := range files {

		jsonFile := filepath.Join(path, file.Name())
		jsonData, err := os.Open(jsonFile)
		if err != nil {
			log.Fatalf("Failed to open JSON file: %s", err)
		}

		byteValue, err := io.ReadAll(jsonData)
		if err != nil {
			log.Fatalf("Failed to read JSON file: %s", err)
		}

		var networkInfo contractgen.NetworkInfo

		// JSON을 구조체로 언마셜링
		err = json.Unmarshal(byteValue, &networkInfo)
		if err != nil {
			log.Fatalf("Failed to unmarshal JSON: %s", err)
		}

		networks = append(networks, networkInfo)

		jsonOutput, err := json.MarshalIndent(networkInfo, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal JSON: %s", err)
		}

		// JSON 형식으로 출력
		fmt.Println(string(jsonOutput))
	}
}
