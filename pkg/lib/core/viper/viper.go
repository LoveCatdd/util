package viper

import (
	"fmt"
	"sync"

	"github.com/LoveCatdd/util/pkg/lib/core/log"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 环境
var (
	_enviro string
	once    sync.Once
)

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

func settingViper(conf *ViperStruct, _type string) error {

	var name string

	if _enviro == "" {
		name = "app.conf"
	} else {
		name = fmt.Sprintf("app.conf.%v", _enviro)
	}

	v := viper.New()
	v.SetConfigName(name)
	v.SetConfigType(_type)
	v.AddConfigPath("../resource")

	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error %v file: %s \n", v.ConfigFileUsed(), err)
		return err
	}

	if err := v.Unmarshal(conf); err != nil {
		log.Errorf("unmarshal conf failed, err:%s \n", err)
		return err
	}

	autoModified(conf, v)
	return nil
}

func autoModified(conf *ViperStruct, v *viper.Viper) {

	once.Do(func() {
		v.WatchConfig()

		v.OnConfigChange(func(e fsnotify.Event) {

			if err := v.Unmarshal(conf); err != nil {
				log.Errorf("unmarshal conf failed, err:%s \n", err)
			}
		})
	})
}
