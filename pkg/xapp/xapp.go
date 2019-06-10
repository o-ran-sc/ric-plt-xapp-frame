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
	"fmt"
	"github.com/spf13/viper"
	"net/http"
)

type ReadyCB func()

var (
	// XApp is an application instance
	Rmr      *RMRClient
	Sdl      *SDLClient
	Rnib     *RNIBClient
	Resource *Router
	Metric   *Metrics
	Logger   *Log
	Config   Configurator
)

func IsReady() bool {
	return Rmr.IsReady() && Sdl.IsReady()
}

func SetReadyCB(cb ReadyCB) {
	Rmr.SetReadyCB(cb)
}

func init() {
	// Load xapp configuration
	Logger = LoadConfig()

	Logger.SetLevel(viper.GetInt("logger.level"))
	Resource = NewRouter()
	Config = Configurator{}
	Metric = NewMetrics(viper.GetString("metrics.url"), viper.GetString("metrics.namespace"), Resource.router)
	Rmr = NewRMRClient()

	if viper.IsSet("db.namespaces") {
		namespaces := viper.GetStringSlice("db.namespaces")
		if namespaces[0] != "" {
			Sdl = NewSDLClient(viper.GetStringSlice("db.namespaces")[0])
		}
		if namespaces[1] != "" {
			Rnib = NewRNIBClient(viper.GetStringSlice("db.namespaces")[1])
		}
	} else {
		Sdl = NewSDLClient(viper.GetString("db.namespace"))
	}
}

func Run(c MessageConsumer) {
	go http.ListenAndServe(viper.GetString("local.host"), Resource.router)

	Logger.Info(fmt.Sprintf("Xapp started, listening on: %s", viper.GetString("local.host")))

	Sdl.TestConnection()
	Rmr.Start(c)
}
