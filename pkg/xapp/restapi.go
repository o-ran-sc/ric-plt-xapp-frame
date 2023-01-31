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
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/models"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"github.com/spf13/viper"
)

const (
	ReadyURL     = "/ric/v1/health/ready"
	AliveURL     = "/ric/v1/health/alive"
	ConfigURL    = "/ric/v1/cm/{name}"
	AppConfigURL = "/ric/v1/config"
)

var (
	healthReady bool
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
		Logger.Debug("restapi: method=%s url=%s", req.Method, req.URL.RequestURI())
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

func (r *Router) GetSymptomDataParams(w http.ResponseWriter, req *http.Request) SymptomDataParams {
	params := SymptomDataParams{}
	queryParams := req.URL.Query()

	Logger.Info("GetSymptomDataParams: %+v", queryParams)

	for p := range queryParams {
		if p == "timeout" {
			fmt.Sscanf(p, "%d", &params.Timeout)
		}
		if p == "fromtime" {
			fmt.Sscanf(p, "%d", &params.FromTime)
		}
		if p == "totime" {
			fmt.Sscanf(p, "%d", &params.ToTime)
		}
	}
	return params
}

func (r *Router) CollectDefaultSymptomData(fileName string, data interface{}) string {
	baseDir := Config.GetString("controls.symptomdata.baseDir")
	if baseDir == "" {
		baseDir = "/tmp/xapp/"
	}

	if err := Util.CreateDir(baseDir); err != nil {
		Logger.Error("CreateDir failed: %v", err)
		return ""
	}

	//
	// Collect some general information into one file
	//
	var lines []string

	// uptime
	d := XappUpTime()
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	lines = append(lines, fmt.Sprintf("uptime: %02d:%02d:%02d", h, m, s))

	Util.WriteToFile(baseDir+"information.txt", strings.Join(lines, "\n"))

	//
	// Collect metrics
	//
	if metrics, err := r.GetLocalMetrics(); err == nil {
		if err := Util.WriteToFile(baseDir+"metrics.json", metrics); err != nil {
			Logger.Error("writeToFile failed for metrics.json: %v", err)
		}
	}

	//
	// Collect currently used config file
	//
	cfile := viper.ConfigFileUsed()
	input, err := ioutil.ReadFile(cfile)
	if err == nil {
		Util.WriteToFile(baseDir+path.Base(cfile), string(input))
	} else {
		Logger.Error("ioutil.ReadFile failed: %v", err)
	}

	//
	// Collect environment
	//
	Util.WriteToFile(baseDir+"environment.txt", strings.Join(os.Environ(), "\n"))

	//
	// Collect RMR rt table
	//
	rtPath := os.Getenv("RMR_STASH_RT")
	if rtPath != "" {
		input, err = ioutil.ReadFile(rtPath)
		if err == nil {
			Util.WriteToFile(baseDir+"rttable.txt", string(input))
		} else {
			Logger.Error("ioutil.ReadFile failed: %v", err)
		}
	}

	//
	// Put data that was provided as argument
	//
	if data != nil {
		if b, err := json.MarshalIndent(data, "", "  "); err == nil {
			Util.WriteToFile(baseDir+fileName, string(b))
		}
	}

	return baseDir
}

func (r *Router) SendSymptomDataJson(w http.ResponseWriter, req *http.Request, data interface{}, n string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename="+n)
	w.WriteHeader(http.StatusOK)
	if data != nil {
		response, _ := json.MarshalIndent(data, "", " ")
		w.Write(response)
	}
}

func (r *Router) SendSymptomDataFile(w http.ResponseWriter, req *http.Request, baseDir, zipFile string) {

	var fileList []string
	fileList = Util.FetchFiles(baseDir, fileList)
	tmpFileName, err := Util.ZipFilesToTmpFile(baseDir, "symptom", fileList)
	if err != nil {
		r.SendSymptomDataError(w, req, err.Error())
		return
	}
	defer os.Remove(tmpFileName)

	w.Header().Set("Content-Disposition", "attachment; filename="+zipFile)
	http.ServeFile(w, req, tmpFileName)
}

func (r *Router) SendSymptomDataError(w http.ResponseWriter, req *http.Request, message string) {
	w.Header().Set("Content-Disposition", "attachment; filename=error_status.txt")
	http.Error(w, message, http.StatusInternalServerError)
}

func (r *Router) GetLocalMetrics() (string, error) {
	buf := &bytes.Buffer{}
	enc := expfmt.NewEncoder(buf, expfmt.FmtText)
	vals, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return fmt.Sprintf("#metrics get error: %s\n", err.Error()), fmt.Errorf("Could get local metrics %w", err)
	}
	for _, val := range vals {
		err = enc.Encode(val)
		if err != nil {
			buf.WriteString(fmt.Sprintf("#metrics enc err:%s\n", err.Error()))
		}
	}
	return string(buf.Bytes()), nil
}

//Resource.InjectRoute(url, metricsHandler, "GET")
//func metricsHandler(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "text/plain")
//	w.WriteHeader(http.StatusOK)
//	metrics, _ := Resource.GetLocalMetrics()
//	w.Write([]byte(metrics))
//}

func IsHealthProbeReady() bool {
	return healthReady
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	healthReady = true
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

	// Read config-files
	cfiles := []string{viper.ConfigFileUsed(), "/opt/ric/config/config-file.json"}

	var err error
	var configFile *os.File
	for _, cfile := range cfiles {
		configFile, err = os.Open(cfile)
		if err != nil {
			configFile = nil
			Logger.Error("Cannot open config file: %s err: %v", cfile, err)
		}
	}
	if err != nil || configFile == nil {
		Logger.Error("Cannot open any of listed config files: %v", cfiles)
		respondWithJSON(w, http.StatusInternalServerError, nil)
		return
	}

	body, err := ioutil.ReadAll(configFile)
	defer configFile.Close()

	xappconfig.Metadata = &metadata
	xappconfig.Config = string(body)

	appconfig = append(appconfig, &xappconfig)

	respondWithJSON(w, http.StatusOK, appconfig)
}
