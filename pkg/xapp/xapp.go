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
	Util          *Utils
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

func XappShutdownCb() {
	if err := doDeregister(); err != nil {
		Logger.Info("xApp deregistration failed: %v, terminating ungracefully!", err)
	} else {
		Logger.Info("xApp deregistration successfull!")
	}

	if shutdownCb != nil {
		shutdownCb()
	}
}

func registerXapp() {
	for {
		time.Sleep(5 * time.Second)
		if !IsHealthProbeReady() {
			Logger.Info("Application='%s' is not ready yet, waiting ...", viper.GetString("name"))
			continue
		}

		Logger.Debug("Application='%s' is now up and ready, continue with registration ...", viper.GetString("name"))
		if err := doRegister(); err == nil {
			Logger.Info("Registration done, proceeding with startup ...")
			break
		}
	}
}

func getService(host, service string) string {
	appnamespace := os.Getenv("APP_NAMESPACE")
	if appnamespace == "" {
		appnamespace = DEFAULT_XAPP_NS
	}

	svc := fmt.Sprintf(service, strings.ToUpper(appnamespace), strings.ToUpper(host))
	url := strings.Split(os.Getenv(strings.Replace(svc, "-", "_", -1)), "//")
	if len(url) > 1 {
		return url[1]
	}
	return ""
}

func getPltNamespace(envName, defVal string) string {
	pltnamespace := os.Getenv("PLT_NAMESPACE")
	if pltnamespace == "" {
		pltnamespace = defVal
	}

	return pltnamespace
}

func doPost(pltNs, url string, msg []byte, status int) error {
	resp, err := http.Post(fmt.Sprintf(url, pltNs, pltNs), "application/json", bytes.NewBuffer(msg))
	if err != nil || resp == nil || resp.StatusCode != status {
		Logger.Info("http.Post to '%s' failed with error: %v", fmt.Sprintf(url, pltNs, pltNs), err)
		return err
	}
	Logger.Info("Post to '%s' done, status:%v", fmt.Sprintf(url, pltNs, pltNs), resp.Status)

	return err
}

func doRegister() error {
	host, _ := os.Hostname()
	xappname := viper.GetString("name")
	xappversion := viper.GetString("version")
	pltNs := getPltNamespace("PLT_NAMESPACE", DEFAULT_PLT_NS)

	httpEp, rmrEp := getService(host, SERVICE_HTTP), getService(host, SERVICE_RMR)
	if httpEp == "" || rmrEp == "" {
		Logger.Warn("Couldn't resolve service endpoints: httpEp=%s rmrEp=%s", httpEp, rmrEp)
		return nil
	}

	requestBody, err := json.Marshal(map[string]string{
		"appName":         host,
		"httpEndpoint":    httpEp,
		"rmrEndpoint":     rmrEp,
		"appInstanceName": xappname,
		"appVersion":      xappversion,
		"configPath":      CONFIG_PATH,
	})

	if err != nil {
		Logger.Error("json.Marshal failed with error: %v", err)
		return err
	}

	return doPost(pltNs, REGISTER_PATH, requestBody, http.StatusCreated)
}

func doDeregister() error {
	if !IsHealthProbeReady() {
		return nil
	}

	name, _ := os.Hostname()
	xappname := viper.GetString("name")
	pltNs := getPltNamespace("PLT_NAMESPACE", DEFAULT_PLT_NS)

	requestBody, err := json.Marshal(map[string]string{
		"appName":         name,
		"appInstanceName": xappname,
	})

	if err != nil {
		Logger.Error("json.Marshal failed with error: %v", err)
		return err
	}

	return doPost(pltNs, DEREGISTER_PATH, requestBody, http.StatusNoContent)
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
						XappShutdownCb()
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

	if !viper.IsSet("controls.logger.noFormat") || !viper.GetBool("controls.logger.noFormat") {
		Logger.SetFormat(0)
	}

	Resource = NewRouter()
	Config = Configurator{}
	Metric = NewMetrics(viper.GetString("metrics.url"), viper.GetString("metrics.namespace"), Resource.router)
	Subscription = NewSubscriber(viper.GetString("controls.subscription.host"), viper.GetInt("controls.subscription.timeout"))
	Sdl = NewSDLClient(viper.GetString("controls.db.namespace"))
	Rnib = NewRNIBClient()
	Util = NewUtils()

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
	go registerXapp()

	Rmr.Start(c)
}

func Run(c MessageConsumer) {
	RunWithParams(c, true)
}
