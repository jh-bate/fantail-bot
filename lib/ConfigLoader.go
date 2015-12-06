package lib

import "github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/jinzhu/configor"

func ConfigLoader(config interface{}, name string) {
	const config_path = "./config/"
	configor.Load(&config, config_path+name)
}
