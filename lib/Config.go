package lib

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/jinzhu/configor"
)

func LoadConfig(config interface{}, name string) {

	const config_path = "./config/"
	const lib_config_path = "./lib/config/"

	absPath, _ := filepath.Abs(config_path + name)
	_, err := os.Open(absPath)

	if err != nil {
		absPath, _ = filepath.Abs(lib_config_path + name)
		log.Println("QandA path ", absPath)
	}

	configor.Load(&config, absPath)
}
