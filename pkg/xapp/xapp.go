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
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

type ReadyCB func(interface{})
type ShutdownCB func()

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
	shutdownCb    ShutdownCB
	shutdownFlag  int32
	shutdownCnt   int32
)

func IsReady() bool {
	return Rmr != nil && Rmr.IsReady() && Sdl != nil && Sdl.IsReady()
}

func SetReadyCB(cb ReadyCB, params interface{}) {
	readyCb = cb
	readyCbParams = params
}

func XappReadyCb(params interface{}) {
	Alarm = NewAlarmClient(viper.GetString("moId"), viper.GetString("name"))
	if readyCb != nil {
		readyCb(readyCbParams)
	}
}

func SetShutdownCB(cb ShutdownCB) {
	shutdownCb = cb
}

func InstallSignalHandler() {
	//
	// Signal handlers to really exit program.
	// shutdownCb can hang until application has
	// made all needed gracefull shutdown actions
	// hardcoded limit for shutdown is 20 seconds
	//
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	//signal handler function
	go func() {
		for range interrupt {
			if atomic.CompareAndSwapInt32(&shutdownFlag, 0, 1) {
				// close function
				go func() {
					timeout := int(20)
					sentry := make(chan struct{})
					defer close(sentry)

					// close callback
					go func() {
						if shutdownCb != nil {
							shutdownCb()
						}
						sentry <- struct{}{}
					}()
					select {
					case <-time.After(time.Duration(timeout) * time.Second):
						Logger.Info("xapp-frame shutdown callback took more than %d seconds", timeout)
					case <-sentry:
						Logger.Info("xapp-frame shutdown callback handled within %d seconds", timeout)
					}
					os.Exit(0)
				}()
			} else {
				newCnt := atomic.AddInt32(&shutdownCnt, 1)
				Logger.Info("xapp-frame shutdown already ongoing. Forced exit counter %d/%d ", newCnt, 5)
				if newCnt >= 5 {
					Logger.Info("xapp-frame shutdown forced exit")
					os.Exit(0)
				}
				continue
			}

		}
	}()
}

func init() {
	// Load xapp configuration
	Logger = LoadConfig()

	if viper.IsSet("controls.logger.level") {
		Logger.SetLevel(viper.GetInt("controls.logger.level"))
	} else {
		Logger.SetLevel(viper.GetInt("logger.level"))
	}
	Logger.SetFormat(0)

	Resource = NewRouter()
	Config = Configurator{}
	Metric = NewMetrics(viper.GetString("metrics.url"), viper.GetString("metrics.namespace"), Resource.router)
	Subscription = NewSubscriber(viper.GetString("subscription.host"), viper.GetInt("subscription.timeout"))
	Sdl = NewSDLClient(viper.GetString("controls.db.namespace"))
	Rnib = NewRNIBClient()

	InstallSignalHandler()
}

func RunWithParams(c MessageConsumer, sdlcheck bool) {
	Rmr = NewRMRClient()
	Rmr.SetReadyCB(XappReadyCb, nil)

	host := fmt.Sprintf(":%d", GetPortData("http").Port)
	go http.ListenAndServe(host, Resource.router)
	Logger.Info(fmt.Sprintf("Xapp started, listening on: %s", host))
	if sdlcheck {
		Sdl.TestConnection()
	}
	Rmr.Start(c)
}

func Run(c MessageConsumer) {
	RunWithParams(c, true)
}
