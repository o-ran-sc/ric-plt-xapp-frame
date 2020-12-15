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
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/models"
)

const (
	ReadyURL     = "/ric/v1/health/ready"
	AliveURL     = "/ric/v1/health/alive"
	ConfigURL    = "/ric/v1/cm/{name}"
	AppConfigURL = "/ric/v1/config"
)

type StatusCb func() bool

type Router struct {
	router *mux.Router
	cbMap  []StatusCb
}

func NewRouter() *Router {
	r := &Router{
		router: mux.NewRouter().StrictSlash(true),
		cbMap:  make([]StatusCb, 0),
	}

	// Inject default routes for health probes
	r.InjectRoute(ReadyURL, readyHandler, "GET")
	r.InjectRoute(AliveURL, aliveHandler, "GET")
	r.InjectRoute(ConfigURL, configHandler, "POST")
	r.InjectRoute(AppConfigURL, appconfigHandler, "GET")

	return r
}

func (r *Router) serviceChecker(inner http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		Logger.Info("restapi: method=%s url=%s", req.Method, req.URL.RequestURI())
		if req.URL.RequestURI() == AliveURL || r.CheckStatus() {
			inner.ServeHTTP(w, req)
		} else {
			respondWithJSON(w, http.StatusServiceUnavailable, nil)
		}
	})
}

func (r *Router) InjectRoute(url string, handler http.HandlerFunc, method string) *mux.Route {
	return r.router.Path(url).HandlerFunc(r.serviceChecker(handler)).Methods(method)
}

func (r *Router) InjectQueryRoute(url string, h http.HandlerFunc, m string, q ...string) *mux.Route {
	return r.router.Path(url).HandlerFunc(r.serviceChecker(h)).Methods(m).Queries(q...)
}

func (r *Router) InjectRoutePrefix(prefix string, handler http.HandlerFunc) *mux.Route {
	return r.router.PathPrefix(prefix).HandlerFunc(r.serviceChecker(handler))
}

func (r *Router) InjectStatusCb(f StatusCb) {
	r.cbMap = append(r.cbMap, f)
}

func (r *Router) CheckStatus() (status bool) {
	if len(r.cbMap) == 0 {
		return true
	}

	for _, f := range r.cbMap {
		status = f()
	}
	return
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, nil)
}

func aliveHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, nil)
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	xappName := mux.Vars(r)["name"]
	if xappName == "" || r.Body == nil {
		respondWithJSON(w, http.StatusBadRequest, nil)
		return
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Logger.Error("ioutil.ReadAll failed: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, nil)
		return
	}

	if err := PublishConfigChange(xappName, string(body)); err != nil {
		respondWithJSON(w, http.StatusInternalServerError, nil)
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		response, _ := json.Marshal(payload)
		w.Write(response)
	}
}

func appconfigHandler(w http.ResponseWriter, r *http.Request) {

	Logger.Info("Inside appconfigHandler")

	var appconfig models.XappConfigList
	var metadata models.ConfigMetadata
	var xappconfig models.XAppConfig
	name := viper.GetString("name")
	configtype := "json"
	metadata.XappName = &name
	metadata.ConfigType = &configtype

	configFile, err := os.Open("/opt/ric/config/config-file.json")
	if err != nil {
		Logger.Error("Cannot open config file: %v", err)
		respondWithJSON(w, http.StatusInternalServerError, nil)
		// return nil,errors.New("Could Not parse the config file")
	}

	body, err := ioutil.ReadAll(configFile)

	defer configFile.Close()

	xappconfig.Metadata = &metadata
	xappconfig.Config = string(body)

	appconfig = append(appconfig, &xappconfig)

	respondWithJSON(w, http.StatusOK, appconfig)

	//return appconfig,nil
}
