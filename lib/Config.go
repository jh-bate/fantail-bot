package lib

import (
	"log"
	"path/filepath"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/jinzhu/configor"
)

func LoadConfig(config interface{}, name string) {

	const config_path = "./config/"

	absPath, err := filepath.Abs(config_path + name)
	if err != nil {
		log.Println("error getting config path", err.Error())
		return
	}
	configor.Load(&config, absPath)
}
