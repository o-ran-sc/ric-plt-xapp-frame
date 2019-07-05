/*
==================================================================================
  Copyright (c) 2019 AT&T Intellectual Property.
  Copyright (c) 2019 Nokia

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
==================================================================================
*/

package xapp

import (
	"flag"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type Configurator struct {
}

type ConfigChangeCB func(filename string)

var ConfigChangeListeners []ConfigChangeCB

func parseCmd() string {
	var fileName *string
	fileName = flag.String("f", os.Getenv("CFG_FILE"), "Specify the configuration file.")
	flag.Parse()

	return *fileName
}

func LoadConfig() (l *Log) {
	l = NewLogger(filepath.Base(os.Args[0]))
	viper.SetConfigFile(parseCmd())

	if err := viper.ReadInConfig(); err != nil {
		l.Error("Reading config file failed: %v", err.Error())
	}
	l.Info("Using config file: %s", viper.ConfigFileUsed())

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		l.Info("config file %s changed ", e.Name)

		Logger.SetLevel(viper.GetInt("logger.level"))
		if len(ConfigChangeListeners) > 0 {
			for _, f := range ConfigChangeListeners {
				go f(e.Name)
			}
		}
	})

	return
}

func AddConfigChangeListener(f ConfigChangeCB) {
	if ConfigChangeListeners == nil {
		ConfigChangeListeners = make([]ConfigChangeCB, 0)
	}
	ConfigChangeListeners = append(ConfigChangeListeners, f)
}

func (*Configurator) GetString(key string) string {
	return viper.GetString(key)
}

func (*Configurator) GetInt(key string) int {
	return viper.GetInt(key)
}

func (*Configurator) GetBool(key string) bool {
	return viper.GetBool(key)
}
