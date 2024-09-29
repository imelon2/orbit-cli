package prompt_test

import (
	"log"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/imelon2/orbit-cli/utils"
	"github.com/spf13/viper"
)

func Test_SelectProvider(t *testing.T) {

	chains := viper.GetStringSlice("providers.sepolia")

	for _, item := range chains {
		t.Log("\n", item)
	}
}

func init() {
	var cfgFile string

	_, filename, _, _ := runtime.Caller(0)
	parsent := utils.GetParentRootDir(filename)

	cfgFile = filepath.Join(parsent, "config.yml")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetConfigFile(cfgFile)

	// 설정 파일 읽기
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

}
