package prompt_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/imelon2/orbit-cli/common/path"
	"github.com/imelon2/orbit-cli/prompt"
	"github.com/spf13/viper"
)

func init() {
	configPath := path.GetConfigPath()
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
}

func Test_SelectChains(t *testing.T) {
	key, _ := prompt.SelectChains()
	fmt.Printf("Test_SelectChains : %s\n", key)
}
