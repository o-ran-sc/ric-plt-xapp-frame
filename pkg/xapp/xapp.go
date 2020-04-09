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
	"os"
	"io/ioutil"
)

type ReadyCB func(interface{})

var (
	// XApp is an application instance
	Rmr           *RMRClient
	Sdl           *SDLClient
	Rnib          *RNIBClient
	Resource      *Router
	Metric        *Metrics
	Logger        *Log
	Config        Configurator
	Subscription  *Subscriber
	Alarm         *AlarmClient
	readyCb       ReadyCB
	readyCbParams interface{}
	alarmRTFile   string = "/tmp/alarm.rt"
	alarmRT       string = "newrt|start\nrte|13111|service-ricplt-alarmadapter-rmr.ricplt:4560\nnewrt|end\n"
)

func IsReady() bool {
	return Rmr != nil && Rmr.IsReady() && Sdl != nil && Sdl.IsReady()
}

func SetReadyCB(cb ReadyCB, params interface{}) {
	readyCb = cb
	readyCbParams = params
}

func xappReadyCb(params interface{}) {
	// Setup static RT for alarm system
	if err := ioutil.WriteFile(alarmRTFile, []byte(alarmRT), 0644); err != nil {
		Logger.Error("ioutil.WriteFile failed with error: %v", err)
	}
	os.Setenv( "RMR_SEED_RT", alarmRTFile )
	os.Setenv( "RMR_RTG_SVC", "-1" )
	
	Subscription = NewSubscriber(viper.GetString("subscription.host"), viper.GetInt("subscription.timeout"))

	if readyCb != nil {
		readyCb(readyCbParams)
	}
}

func init() {
	// Load xapp configuration
	Logger = LoadConfig()

	Logger.SetLevel(viper.GetInt("logger.level"))
	Resource = NewRouter()
	Config = Configurator{}
	Metric = NewMetrics(viper.GetString("metrics.url"), viper.GetString("metrics.namespace"), Resource.router)
	Alarm = NewAlarmClient(viper.GetString("alarm.MOId"), viper.GetString("alarm.APPId"))

	if viper.IsSet("db.namespaces") {
		namespaces := viper.GetStringSlice("db.namespaces")
		if len(namespaces) > 0 && namespaces[0] != "" {
			Sdl = NewSDLClient(viper.GetStringSlice("db.namespaces")[0])
		}
		if len(namespaces) > 1 && namespaces[1] != "" {
			Rnib = NewRNIBClient(viper.GetStringSlice("db.namespaces")[1])
		}
	} else {
		Sdl = NewSDLClient(viper.GetString("db.namespace"))
	}
}

func RunWithParams(c MessageConsumer, sdlcheck bool) {
	Rmr = NewRMRClient()
	Rmr.SetReadyCB(xappReadyCb, nil)
	go http.ListenAndServe(viper.GetString("local.host"), Resource.router)
	Logger.Info(fmt.Sprintf("Xapp started, listening on: %s", viper.GetString("local.host")))
	if sdlcheck {
		Sdl.TestConnection()
	}
	Rmr.Start(c)
}

func Run(c MessageConsumer) {
	RunWithParams(c, true)
}
