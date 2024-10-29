package path_test

import (
	"fmt"
	"testing"

	"github.com/imelon2/orbit-cli/common/path"
)

func Test_GetConfigPath(t *testing.T) {
	config := path.GetConfigPath()
	fmt.Printf("Config.yml Path : %s\n", config)
}

func Test_GetAbiPath(t *testing.T) {
	path := path.GetAbiPath()
	fmt.Printf("aggregateAbi.json Path : %s\n", path)
}
