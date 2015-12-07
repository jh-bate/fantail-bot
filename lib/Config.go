package lib

import (
	"path/filepath"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/jinzhu/configor"
)

func LoadConfig(config interface{}, name string) {

	const config_path = "./config/"
	absPath, _ := filepath.Abs(config_path + name)
	configor.Load(&config, absPath)
}
