package config

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 环境
var _enviro string

func SetEnviro(enviro string) {
	_enviro = enviro
}

type ViperStruct interface {
	FileType() string
}

func Yaml(conf ViperStruct) error {
	return settingViper(&conf, VIPER_YAML)
}

func JSON(conf ViperStruct) error {
	return settingViper(&conf, VIPER_JSON)
}

func DOTENV(conf ViperStruct) error {
	return settingViper(&conf, VIPER_DOTENV)
}

func settingViper(conf *ViperStruct, type_ string) error {

	var name string

	if _enviro == "" {
		name = "app.conf"
	} else {
		name = fmt.Sprintf("app.conf.%v", _enviro)
	}

	v := viper.New()
	v.SetConfigName(name)
	v.SetConfigType(type_)
	v.AddConfigPath("../resource")

	err := v.ReadInConfig()
	if err != nil {
		log.Printf("Fatal error %v file: %s \n", v.ConfigFileUsed(), err)
		return err
	}

	if err := v.Unmarshal(conf); err != nil {
		log.Printf("unmarshal conf failed, err:%s \n", err)
		return err
	}

	autoModified(conf, v)
	return nil
}

func autoModified(conf *ViperStruct, v *viper.Viper) {

	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {

		if err := v.Unmarshal(conf); err != nil {
			log.Printf("unmarshal conf failed, err:%s \n", err)
		} else {
			log.Printf("%v changed", e.Name)

		}

	})
}
