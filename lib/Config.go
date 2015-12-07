package lib

import (
	"os"
	"path/filepath"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/jinzhu/configor"
)

func LoadConfig(config interface{}, name string) {

	// we need to get the absolute path and also test it because of the
	// difference of how this code runs in different environments
	const config_path = "./config/"
	absPath, _ := filepath.Abs(config_path + name)
	_, err := os.Open(absPath)

	if err != nil {
		const config_lib_path = "./lib/config/"
		absPath, _ = filepath.Abs(config_lib_path + name)
	}

	configor.Load(&config, absPath)
}
