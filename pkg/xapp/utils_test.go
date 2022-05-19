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
	"os"
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

	tmpFileName, err := utils.ZipFilesToTmpFile("/tmp/abcd", "symptom", []string{"/tmp/abcd/file.txt"})
	assert.Equal(t, err, nil)
	defer os.Remove(tmpFileName)

	assert.Equal(t, utils.CreateDir("/tmp/dcba"), nil)
	_, err = utils.UnZipFiles(tmpFileName, "/tmp/dcba")
	assert.Equal(t, err, nil)

	utils.DeleteFile("/tmp/abcd")
	utils.DeleteFile("/tmp/dcba")
}

func TestSymptomdata(t *testing.T) {
	os.Setenv("RMR_STASH_RT", "config/uta_rtg.rt.stash.inc")
	assert.Equal(t, Resource.CollectDefaultSymptomData("abcd.tgz", "data"), "/tmp/xapp/")
}

func TestSymptomdataCollection(t *testing.T) {
	var handler1 = func(w http.ResponseWriter, r *http.Request) {
		Resource.SendSymptomDataJson(w, r, "data", "aaaa")
		Resource.SendSymptomDataFile(w, r, "./config", "symptomdata.gz")
	}

	Resource.InjectQueryRoute("/ric/v1/user1", handler1, "GET", "foo", "bar", "id", "mykey")

	req, _ := http.NewRequest("GET", "/ric/v1/user1?foo=bar&id=mykey", nil)
	resp := executeRequest(req, nil)
	checkResponseCode(t, http.StatusOK, resp.Code)
}

func TestSymptomdataCollectionError(t *testing.T) {
	var handler2 = func(w http.ResponseWriter, r *http.Request) {
		Resource.SendSymptomDataError(w, r, "Error text")
	}

	Resource.InjectQueryRoute("/ric/v1/user2", handler2, "GET", "foo", "bar", "id", "mykey")

	req, _ := http.NewRequest("GET", "/ric/v1/user2?foo=bar&id=mykey", nil)
	resp := executeRequest(req, nil)
	checkResponseCode(t, http.StatusInternalServerError, resp.Code)
}

func TestGetSymptomDataParams(t *testing.T) {
	var handler3 = func(w http.ResponseWriter, r *http.Request) {
		Resource.GetSymptomDataParams(w, r)
	}

	Resource.InjectQueryRoute("/ric/v1/user3", handler3, "GET", "timeout", "10", "fromtime", "1", "totime", "2")

	req, _ := http.NewRequest("GET", "/ric/v1/user3?timeout=10&fromtime=1&totime=2", nil)
	resp := executeRequest(req, nil)
	checkResponseCode(t, http.StatusOK, resp.Code)
}

func TestAppconfigHandler(t *testing.T) {
	var handler4 = func(w http.ResponseWriter, r *http.Request) {
		appconfigHandler(w, r)
	}

	Resource.InjectQueryRoute("/ric/v1/user4", handler4, "GET", "foo", "bar", "id", "mykey")

	req, _ := http.NewRequest("GET", "/ric/v1/user4?foo=bar&id=mykey", nil)
	resp := executeRequest(req, nil)
	checkResponseCode(t, http.StatusInternalServerError, resp.Code)
}
