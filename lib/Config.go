package lib

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func LoadConfig(config interface{}, name string) {

	const config_path = "config/"
	const config_lib_path = "lib/config/"

	absPath, _ := filepath.Abs(config_path + name)
	file, err := os.Open(absPath)

	if err != nil {
		log.Println("error getting config path", err.Error())
		absPath, _ = filepath.Abs(config_lib_path + name)
		file, err = os.Open(absPath)
	}

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Panic("error loading config ", err.Error())
	}

}
