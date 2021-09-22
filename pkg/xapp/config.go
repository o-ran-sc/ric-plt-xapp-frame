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
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

//-----------------------------------------------------------------------------
//
//-----------------------------------------------------------------------------

type mtype struct {
	Name string
	Id   int
}

type Configurator struct {
}

type ConfigChangeCB func(filename string)

type SDLNotificationCB func(string, ...string)

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

	updateMTypes := func() {
		var mtypes []mtype
		viper.UnmarshalKey("messaging.mtypes", &mtypes)

		if len(mtypes) > 0 {
			l.Info("Config mtypes before RICMessageTypes:%d RicMessageTypeToName:%d", len(RICMessageTypes), len(RicMessageTypeToName))
			for _, v := range mtypes {
				nadd := false
				iadd := false
				if _, ok := RICMessageTypes[v.Name]; ok == false {
					nadd = true
				}
				if _, ok := RicMessageTypeToName[int(v.Id)]; ok == false {
					iadd = true
				}
				if iadd != nadd {
					l.Error("Config mtypes rmr.mtypes entry skipped due conflict with existing values %s(%t) %d(%t) ", v.Name, nadd, v.Id, iadd)
				} else if iadd {
					l.Info("Config mtypes rmr.mtypes entry added %s(%t) %d(%t) ", v.Name, nadd, v.Id, iadd)
					RICMessageTypes[v.Name] = int(v.Id)
					RicMessageTypeToName[int(v.Id)] = v.Name
				} else {
					l.Info("Config mtypes rmr.mtypes entry skipped %s(%t) %d(%t) ", v.Name, nadd, v.Id, iadd)
				}
			}
			l.Info("Config mtypes after RICMessageTypes:%d RicMessageTypeToName:%d", len(RICMessageTypes), len(RicMessageTypeToName))
		}
	}

	updateMTypes()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		l.Info("config file %s changed ", e.Name)

		updateMTypes()
		if viper.IsSet("controls.logger.level") {
			Logger.SetLevel(viper.GetInt("controls.logger.level"))
		} else {
			Logger.SetLevel(viper.GetInt("logger.level"))
		}

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

func PublishConfigChange(appName, eventJson string) error {
	channel := fmt.Sprintf("CM_UPDATE:%s", appName)
	if err := SdlStorage.StoreAndPublish(getCmSdlNs(), channel, eventJson, appName, eventJson); err != nil {
		Logger.Error("Sdl.Store failed: %v", err)
		return err
	}
	return nil
}

func ReadConfig(appName string) (map[string]interface{}, error) {
	return SdlStorage.Read(getCmSdlNs(), appName)
}

func GetPortData(pname string) (d PortData) {
	var getPolicies = func(policies []interface{}) (plist []int) {
		for _, p := range policies {
			plist = append(plist, int(p.(float64)))
		}
		return plist
	}

	if viper.IsSet("messaging") == false {
		if pname == "http" {
			d.Port = 8080
		}
		if pname == "rmrdata" {
			d.Port = 4560
		}
		return
	}

	for _, v := range viper.GetStringMap("messaging")["ports"].([]interface{}) {
		if n, ok := v.(map[string]interface{})["name"].(string); ok && n == pname {
			d.Name = n
			if p, _ := v.(map[string]interface{})["port"].(float64); ok {
				d.Port = int(p)
			}
			if m, _ := v.(map[string]interface{})["maxSize"].(float64); ok {
				d.MaxSize = int(m)
			}
			if m, _ := v.(map[string]interface{})["threadType"].(float64); ok {
				d.ThreadType = int(m)
			}
			if m, _ := v.(map[string]interface{})["lowLatency"].(bool); ok {
				d.LowLatency = bool(m)
			}
			if m, _ := v.(map[string]interface{})["fastAck"].(bool); ok {
				d.FastAck = bool(m)
			}
			if m, _ := v.(map[string]interface{})["maxRetryOnFailure"].(float64); ok {
				d.MaxRetryOnFailure = int(m)
			}
			if policies, ok := v.(map[string]interface{})["policies"]; ok {
				d.Policies = getPolicies(policies.([]interface{}))
			}
		}
	}
	return
}

func getCmSdlNs() string {
	return fmt.Sprintf("cm/%s", viper.GetString("name"))
}

func (*Configurator) SetSDLNotificationCB(appName string, sdlNotificationCb SDLNotificationCB) error {
	return SdlStorage.Subscribe(getCmSdlNs(), sdlNotificationCb, fmt.Sprintf("CM_UPDATE:%s", appName))
}

func (*Configurator) GetString(key string) string {
	return viper.GetString(key)
}

func (*Configurator) GetInt(key string) int {
	return viper.GetInt(key)
}

func (*Configurator) GetUint32(key string) uint32 {
	return viper.GetUint32(key)
}

func (*Configurator) GetBool(key string) bool {
	return viper.GetBool(key)
}

func (*Configurator) Get(key string) interface{} {
	return viper.Get(key)
}

func (*Configurator) GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func (*Configurator) GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

func (*Configurator) IsSet(key string) bool {
	return viper.IsSet(key)
}
