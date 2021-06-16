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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUtils(t *testing.T) {
	utils := NewUtils()
	assert.NotNil(t, utils, "NewUtils failed")

	assert.Equal(t, utils.FileExists("abcd"), false)
	assert.Equal(t, utils.CreateDir("/tmp/abcd"), nil)
	assert.Equal(t, utils.WriteToFile("/tmp/abcd/file.txt", "data"), nil)
	utils.FetchFiles("./", []string{"go.mod"})
	utils.FetchFiles("./", []string{"go.mod"})

	utils.DeleteFile("/tmp/abcd")
}

func TestSymptomdata(t *testing.T) {
	assert.Equal(t, Resource.CollectDefaultSymptomData("abcd.tgz", "data"), "/tmp/xapp/")
}

func TestSymptomdataCollection(t *testing.T) {
	var handler = func(w http.ResponseWriter, r *http.Request) {
		Resource.SendSymptomDataJson(w, r, "data", "aaaa")
		Resource.SendSymptomDataFile(w, r, "./config", "symptomdata.gz")
	}

	Resource.InjectQueryRoute("/ric/v1/user", handler, "GET", "foo", "bar", "id", "mykey")

	req, _ := http.NewRequest("GET", "/ric/v1/user?foo=bar&id=mykey", nil)
	resp := executeRequest(req, nil)
	checkResponseCode(t, http.StatusOK, resp.Code)
}

func TestSymptomdataCollectionError(t *testing.T) {
	var handler = func(w http.ResponseWriter, r *http.Request) {
		Resource.SendSymptomDataError(w, r, "Error text")
	}

	Resource.InjectQueryRoute("/ric/v1/user", handler, "GET", "foo", "bar", "id", "mykey")

	req, _ := http.NewRequest("GET", "/ric/v1/user?foo=bar&id=mykey", nil)
	resp := executeRequest(req, nil)
	checkResponseCode(t, http.StatusOK, resp.Code)
}

func TestGetSymptomDataParams(t *testing.T) {
	var handler = func(w http.ResponseWriter, r *http.Request) {
		Resource.GetSymptomDataParams(w, r)
	}

	Resource.InjectQueryRoute("/ric/v1/user", handler, "GET", "foo", "bar", "id", "mykey")

	req, _ := http.NewRequest("GET", "/ric/v1/user?foo=bar&id=mykey", nil)
	resp := executeRequest(req, nil)
	checkResponseCode(t, http.StatusOK, resp.Code)
}

func TestAppconfigHandler(t *testing.T) {
	var handler = func(w http.ResponseWriter, r *http.Request) {
		appconfigHandler(w, r)
	}

	Resource.InjectQueryRoute("/ric/v1/user", handler, "GET", "foo", "bar", "id", "mykey")

	req, _ := http.NewRequest("GET", "/ric/v1/user?foo=bar&id=mykey", nil)
	resp := executeRequest(req, nil)
	checkResponseCode(t, http.StatusOK, resp.Code)
}
