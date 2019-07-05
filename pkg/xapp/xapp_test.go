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
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

type Consumer struct {
}

func (m Consumer) Consume(mtype, sid int, payload []byte, meid *RMRMeid) (err error) {
	Sdl.Store("myKey", payload)
	return nil
}

// Test cases
func TestMain(m *testing.M) {
	// Just run on the background (for coverage)
	go Run(Consumer{})

	code := m.Run()
	os.Exit(code)
}

func TestGetHealthCheckRetursServiceUnavailableError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/ric/v1/health/ready", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusServiceUnavailable, response.Code)
}

func TestGetHealthCheckReturnsSuccess(t *testing.T) {
	// Wait until RMR is up-and-running
	for Rmr.IsReady() == false {
		time.Sleep(time.Duration(2) * time.Second)
	}

	req, _ := http.NewRequest("GET", "/ric/v1/health/ready", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestInjectQuerySinglePath(t *testing.T) {
	var handler = func(w http.ResponseWriter, r *http.Request) {
	}

	Resource.InjectQueryRoute("/ric/v1/user", handler, "GET", "foo", "bar")

	req, _ := http.NewRequest("GET", "/ric/v1/user?foo=bar", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestInjectQueryMultiplePaths(t *testing.T) {
	var handler = func(w http.ResponseWriter, r *http.Request) {
	}

	Resource.InjectQueryRoute("/ric/v1/user", handler, "GET", "foo", "bar", "id", "mykey")

	req, _ := http.NewRequest("GET", "/ric/v1/user?foo=bar&id=mykey", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestInjectQueryFailures(t *testing.T) {
	var handler = func(w http.ResponseWriter, r *http.Request) {
	}

	Resource.InjectQueryRoute("/ric/v1/user", handler, "GET", "foo", "bar", "id", "mykey")

	req, _ := http.NewRequest("GET", "/ric/v1/user?invalid=bar&no=mykey", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestMessagesReceivedSuccessfully(t *testing.T) {
	for i := 0; i < 100; i++ {
		Rmr.Send(10004, 1111, []byte{1, 2, 3, 4, 5, 6}, &RMRMeid{PlmnID: "1234", EnbID: "7788"})
	}

	// Allow time to process the messages
	time.Sleep(time.Duration(2) * time.Second)

	stats := getMetrics(t)
	if !strings.Contains(stats, "ricxapp_RMR_Transmitted 100") {
		t.Errorf("Error: ricxapp_RMR_Transmitted value incorrect")
	}

	if !strings.Contains(stats, "ricxapp_RMR_Received 100") {
		t.Errorf("Error: ricxapp_RMR_Received value incorrect")
	}

	if !strings.Contains(stats, "ricxapp_RMR_TransmitError 0") {
		t.Errorf("Error: ricxapp_RMR_TransmitError value incorrect")
	}

	if !strings.Contains(stats, "ricxapp_RMR_ReceiveError 0") {
		t.Errorf("Error: ricxapp_RMR_ReceiveError value incorrect")
	}

	if !strings.Contains(stats, "ricxapp_SDL_Stored 100") {
		t.Errorf("Error: ricxapp_SDL_Stored value incorrect")
	}

	if !strings.Contains(stats, "ricxapp_SDL_StoreError 0") {
		t.Errorf("Error: ricxapp_SDL_StoreError value incorrect")
	}
}

func TestSubscribeChannels(t *testing.T) {
	var NotificationCb = func(ch string, events ...string) {
		if ch != "channel1" {
			t.Errorf("Error: Callback function called with incorrect params")
		}
	}

	if err := Sdl.Subscribe(NotificationCb, "channel1"); err != nil {
		t.Errorf("Error: Subscribe failed: %v", err)
	}
	time.Sleep(time.Duration(2) * time.Second)

	if err := Sdl.StoreAndPublish("channel1", "event", "key1", "data1"); err != nil {
		t.Errorf("Error: Publish failed: %v", err)
	}
}

func TestGetRicMessageSuccess(t *testing.T) {
	id, ok := Rmr.GetRicMessageId("RIC_SUB_REQ")
	if !ok || id != 12010 {
		t.Errorf("Error: GetRicMessageId failed: id=%d", id)
	}

	name := Rmr.GetRicMessageName(12010)
	if name != "RIC_SUB_REQ" {
		t.Errorf("Error: GetRicMessageName failed: name=%s", name)
	}
}

func TestGetRicMessageFails(t *testing.T) {
	id, ok := Rmr.GetRicMessageId("INVALID")
	if ok {
		t.Errorf("Error: GetRicMessageId returned invalid value id=%d", id)
	}

	name := Rmr.GetRicMessageName(123456)
	if name != "" {
		t.Errorf("Error: GetRicMessageName returned invalid value: name=%s", name)
	}
}

func TestTeardown(t *testing.T) {
	Sdl.Clear()
}

// Helper functions
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	vars := map[string]string{"id": "1"}
	req = mux.SetURLVars(req, vars)
	Resource.router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func getMetrics(t *testing.T) string {
	req, _ := http.NewRequest("GET", "/ric/v1/metrics", nil)
	response := executeRequest(req)

	return response.Body.String()
}
