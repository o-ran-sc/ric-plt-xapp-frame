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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/models"
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

func (r *Router) GetSymptomDataParams(w http.ResponseWriter, req *http.Request) SymptomDataParams {
	Logger.Info("GetSymptomDataParams ...")

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

	if metrics, err := r.GetLocalMetrics(GetPortData("http").Port); err == nil {
		if err := Util.WriteToFile(baseDir+"metrics.json", metrics); err != nil {
			Logger.Error("writeToFile failed for metrics.json: %v", err)
		}
	}

	if data != nil {
		if b, err := json.MarshalIndent(data, "", "  "); err == nil {
			Util.WriteToFile(baseDir+fileName, string(b))
		}
	}

	rtPath := os.Getenv("RMR_STASH_RT")
	if rtPath == "" {
		return baseDir
	}

	input, err := ioutil.ReadFile(rtPath)
	if err != nil {
		Logger.Error("ioutil.ReadFile failed: %v", err)
		return baseDir
	}

	Util.WriteToFile(baseDir+"rttable.txt", string(input))
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
	// Compress and reply with attachment
	tmpFile, err := ioutil.TempFile("", "symptom")
	if err != nil {
		r.SendSymptomDataError(w, req, "Failed to create a tmp file: "+err.Error())
		return
	}
	defer os.Remove(tmpFile.Name())

	var fileList []string
	fileList = Util.FetchFiles(baseDir, fileList)
	err = Util.ZipFiles(tmpFile, baseDir, fileList)
	if err != nil {
		r.SendSymptomDataError(w, req, "Failed to zip the files: "+err.Error())
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+zipFile)
	http.ServeFile(w, req, tmpFile.Name())
}

func (r *Router) SendSymptomDataError(w http.ResponseWriter, req *http.Request, message string) {
	w.Header().Set("Content-Disposition", "attachment; filename=error_status.txt")
	http.Error(w, message, http.StatusInternalServerError)
}

func (r *Router) GetLocalMetrics(port int) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/ric/v1/metrics", port))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	metrics, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(metrics), nil
}

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
