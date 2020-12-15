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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

func xappShutdownCb() {
	SendDeregistermsg()
	Logger.Info("Wait for xapp to get unregistered")
	time.Sleep(10 * time.Second)
}

func registerxapp() {
	var (
		retries int = 10
	)
	for retries > 0 {
		name, _ := os.Hostname()
		httpservicename := "SERVICE_RICXAPP_" + strings.ToUpper(name) + "_HTTP_PORT"
		httpendpoint := os.Getenv(strings.Replace(httpservicename, "-", "_", -1))
		urlString := strings.Split(httpendpoint, "//")
		// Added this check to make UT pass
		if urlString[0] == "" {
			return
		}
		resp, err := http.Get(fmt.Sprintf("http://%s/ric/v1/health/ready", urlString[1]))
		retries -= 1
		time.Sleep(5 * time.Second)
		if err != nil {
			Logger.Error("Error in health check: %v", err)
		}
		if err == nil {
			retries -= 10
			Logger.Info("Health Probe Success with resp.StatusCode is %v", resp.StatusCode)
			if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
				go SendRegistermsg()
			}
		} else {
			Logger.Info("Health Probe failed, retrying...")
		}
	}
}

func SendRegistermsg() {
	name, _ := os.Hostname()
	xappname := viper.GetString("name")
	xappversion := viper.GetString("version")
	appnamespace := os.Getenv("APP_NAMESPACE")
	if appnamespace == "" {
		appnamespace = "ricxapp"
	}
	httpservicename := "SERVICE_" + strings.ToUpper(appnamespace) + "_" + strings.ToUpper(name) + "_HTTP_PORT"
	rmrservicename := "SERVICE_" + strings.ToUpper(appnamespace) + "_" + strings.ToUpper(name) + "_RMR_PORT"
	httpendpointstr := os.Getenv(strings.Replace(httpservicename, "-", "_", -1))
	rmrendpointstr := os.Getenv(strings.Replace(rmrservicename, "-", "_", -1))
	httpendpoint := strings.Split(httpendpointstr, "//")
	rmrendpoint := strings.Split(rmrendpointstr, "//")
	if httpendpoint[0] == "" || rmrendpoint[0] == "" {
		return
	}

	pltnamespace := os.Getenv("PLT_NAMESPACE")
	if pltnamespace == "" {
		pltnamespace = "ricplt"
	}

	configpath := "/ric/v1/config"

	requestBody, err := json.Marshal(map[string]string{
		"appName":         name,
		"httpEndpoint":    httpendpoint[1],
		"rmrEndpoint":     rmrendpoint[1],
		"appInstanceName": xappname,
		"appVersion":      xappversion,
		"configPath":      configpath,
	})

	if err != nil {
		Logger.Info("Error while compiling request to appmgr: %v", err)
	} else {
		url := fmt.Sprintf("http://service-%v-appmgr-http.%v:8080/ric/v1/register", pltnamespace, pltnamespace)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
		Logger.Info(" Resp is %v", resp)
		if err != nil {
			Logger.Info("Error  compiling request to appmgr: %v", err)
		}
		Logger.Info("Registering request sent. Response received is :%v", resp)

		if resp != nil {
			body, err := ioutil.ReadAll(resp.Body)
			Logger.Info("Post body is %v", resp.Body)
			if err != nil {
				Logger.Info("rsp: Error  compiling request to appmgr: %v", string(body))
			}
			defer resp.Body.Close()
		}
	}
}

func SendDeregistermsg() {

	name, _ := os.Hostname()
	xappname := viper.GetString("name")

	appnamespace := os.Getenv("APP_NAMESPACE")
	if appnamespace == "" {
		appnamespace = "ricxapp"
	}
	pltnamespace := os.Getenv("PLT_NAMESPACE")
	if pltnamespace == "" {
		pltnamespace = "ricplt"
	}

	requestBody, err := json.Marshal(map[string]string{
		"appName":         name,
		"appInstanceName": xappname,
	})

	if err != nil {
		Logger.Info("Error while compiling request to appmgr: %v", err)
	} else {
		url := fmt.Sprintf("http://service-%v-appmgr-http.%v:8080/ric/v1/deregister", pltnamespace, pltnamespace)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
		Logger.Info(" Resp is %v", resp)
		if err != nil {
			Logger.Info("Error  compiling request to appmgr: %v", err)
		}
		Logger.Info("Deregistering request sent. Response received is :%v", resp)

		if resp != nil {
			body, err := ioutil.ReadAll(resp.Body)
			Logger.Info("Post body is %v", resp.Body)
			if err != nil {
				Logger.Info("rsp: Error  compiling request to appmgr: %v", string(body))
			}
			defer resp.Body.Close()
		}
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
	SetShutdownCB(xappShutdownCb)
	host := fmt.Sprintf(":%d", GetPortData("http").Port)
	go http.ListenAndServe(host, Resource.router)
	Logger.Info(fmt.Sprintf("Xapp started, listening on: %s", host))
	if sdlcheck {
		Sdl.TestConnection()
	}
	go registerxapp()
	Rmr.Start(c)
}

func Run(c MessageConsumer) {
	RunWithParams(c, true)
}
