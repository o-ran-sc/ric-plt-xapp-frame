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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var _ = func() bool {
	testing.Init()
	return true
}()

type Consumer struct{}

func (m Consumer) Consume(params *RMRParams) (err error) {
	Sdl.Store("myKey", params.Payload)
	return nil
}

// Test cases
func TestMain(m *testing.M) {
	os.Setenv("SERVICE_RICXAPP_UEEC_HTTP_PORT", "tcp://localhost:8080")
	os.Setenv("SERVICE_RICXAPP_UEEC_RMR_PORT", "tcp://localhost:4561")
	go RunWithParams(Consumer{}, viper.GetBool("controls.waitForSdl"))
	time.Sleep(time.Duration(5) * time.Second)
	code := m.Run()
	os.Exit(code)
}

func TestGetHealthCheckRetursServiceUnavailableError(t *testing.T) {
	Logger.Info("CASE: TestGetHealthCheckRetursServiceUnavailableError")
	req, _ := http.NewRequest("GET", "/ric/v1/health/ready", nil)
	/*response :=*/ executeRequest(req, nil)

	//checkResponseCode(t, http.StatusServiceUnavailable, response.Code)
}

func TestGetHealthCheckReturnsSuccess(t *testing.T) {
	Logger.Info("CASE: TestGetHealthCheckReturnsSuccess")
	for Rmr.IsReady() == false {
		time.Sleep(time.Duration(2) * time.Second)
	}

	req, _ := http.NewRequest("GET", "/ric/v1/health/ready", nil)
	response := executeRequest(req, nil)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestInjectQuerySinglePath(t *testing.T) {
	Logger.Info("CASE: TestInjectQuerySinglePath")
	var handler = func(w http.ResponseWriter, r *http.Request) {
	}

	Resource.InjectQueryRoute("/ric/v1/user", handler, "GET", "foo", "bar")

	req, _ := http.NewRequest("GET", "/ric/v1/user?foo=bar", nil)
	response := executeRequest(req, nil)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestInjectQueryMultiplePaths(t *testing.T) {
	Logger.Info("CASE: TestInjectQueryMultiplePaths")
	var handler = func(w http.ResponseWriter, r *http.Request) {
	}

	Resource.InjectQueryRoute("/ric/v1/user", handler, "GET", "foo", "bar", "id", "mykey")

	req, _ := http.NewRequest("GET", "/ric/v1/user?foo=bar&id=mykey", nil)
	response := executeRequest(req, nil)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestInjectQueryFailures(t *testing.T) {
	Logger.Info("CASE: TestInjectQueryFailures")
	var handler = func(w http.ResponseWriter, r *http.Request) {
	}

	Resource.InjectQueryRoute("/ric/v1/user", handler, "GET", "foo", "bar", "id", "mykey")

	req, _ := http.NewRequest("GET", "/ric/v1/user?invalid=bar&no=mykey", nil)
	response := executeRequest(req, nil)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestMessagesReceivedSuccessfully(t *testing.T) {
	Logger.Info("CASE: TestMessagesReceivedSuccessfully")
	time.Sleep(time.Duration(5) * time.Second)
	for i := 0; i < 100; i++ {
		params := &RMRParams{}
		params.Mtype = 10004
		params.SubId = -1
		params.Payload = []byte{1, 2, 3, 4, 5, 6}
		params.Meid = &RMRMeid{PlmnID: "1234", EnbID: "7788", RanName: "RanName-1234"}
		params.Xid = "TestXID"

		if i%2 == 0 {
			Rmr.SendMsg(params)
		} else {
			Rmr.SendWithRetry(params, false, 1)
		}
	}
	Rmr.RegisterMetrics()

	// Allow time to process the messages
	time.Sleep(time.Duration(5) * time.Second)

	waitForSdl := viper.GetBool("controls.waitForSdl")
	stats := getMetrics(t)
	if !strings.Contains(stats, "ricxapp_RMR_Transmitted 100") {
		t.Errorf("Error: ricxapp_RMR_Transmitted value incorrect: %v", stats)
	}

	if !strings.Contains(stats, "ricxapp_RMR_Received 100") {
		t.Errorf("Error: ricxapp_RMR_Received value incorrect: %v", stats)
	}

	if !strings.Contains(stats, "ricxapp_RMR_TransmitError 0") {
		t.Errorf("Error: ricxapp_RMR_TransmitError value incorrect")
	}

	if !strings.Contains(stats, "ricxapp_RMR_ReceiveError 0") {
		t.Errorf("Error: ricxapp_RMR_ReceiveError value incorrect")
	}

	if waitForSdl && !strings.Contains(stats, "ricxapp_SDL_Stored 100") {
		t.Errorf("Error: ricxapp_SDL_Stored value incorrect")
	}

	if waitForSdl && !strings.Contains(stats, "ricxapp_SDL_StoreError 0") {
		t.Errorf("Error: ricxapp_SDL_StoreError value incorrect")
	}
}

func TestMessagesReceivedSuccessfullyUsingWh(t *testing.T) {
	Logger.Info("CASE: TestMessagesReceivedSuccessfullyUsingWh")
	time.Sleep(time.Duration(5) * time.Second)
	whid := Rmr.Openwh("localhost:4560")
	time.Sleep(time.Duration(1) * time.Second)
	for i := 0; i < 100; i++ {
		params := &RMRParams{}
		params.Mtype = 10004
		params.SubId = -1
		params.Payload = []byte{1, 2, 3, 4, 5, 6}
		params.Meid = &RMRMeid{PlmnID: "1234", EnbID: "7788", RanName: "RanName-1234"}
		params.Xid = "TestXID"
		params.Whid = int(whid)

		if i == 0 {
			Logger.Info("%+v", params.String())
		}

		Rmr.SendMsg(params)
	}

	// Allow time to process the messages
	time.Sleep(time.Duration(5) * time.Second)

	waitForSdl := viper.GetBool("controls.waitForSdl")
	stats := getMetrics(t)
	if !strings.Contains(stats, "ricxapp_RMR_Transmitted 200") {
		t.Errorf("Error: ricxapp_RMR_Transmitted value incorrect: %v", stats)
	}

	if !strings.Contains(stats, "ricxapp_RMR_Received 200") {
		t.Errorf("Error: ricxapp_RMR_Received value incorrect: %v", stats)
	}

	if !strings.Contains(stats, "ricxapp_RMR_TransmitError 0") {
		t.Errorf("Error: ricxapp_RMR_TransmitError value incorrect")
	}

	if !strings.Contains(stats, "ricxapp_RMR_ReceiveError 0") {
		t.Errorf("Error: ricxapp_RMR_ReceiveError value incorrect")
	}

	if waitForSdl && !strings.Contains(stats, "ricxapp_SDL_Stored 200") {
		t.Errorf("Error: ricxapp_SDL_Stored value incorrect")
	}

	if waitForSdl && !strings.Contains(stats, "ricxapp_SDL_StoreError 0") {
		t.Errorf("Error: ricxapp_SDL_StoreError value incorrect")
	}
	Rmr.Closewh(int(whid))
}

func TestMessagesReceivedSuccessfullyUsingWhCall(t *testing.T) {
	Logger.Info("CASE: TestMessagesReceivedSuccessfullyUsingWhCall")
	time.Sleep(time.Duration(5) * time.Second)
	whid := Rmr.Openwh("localhost:4560")
	params := &RMRParams{}
	params.Payload = []byte("newrt|start\nnewrt|end\n")
	params.Whid = int(whid)
	params.Callid = 4
	params.Timeout = 1000
	Rmr.SendCallMsg(params)

	// Allow time to process the messages
	time.Sleep(time.Duration(2) * time.Second)

	waitForSdl := viper.GetBool("controls.waitForSdl")
	stats := getMetrics(t)
	if !strings.Contains(stats, "ricxapp_RMR_Transmitted 200") {
		t.Errorf("Error: ricxapp_RMR_Transmitted value incorrect: %v", stats)
	}

	if !strings.Contains(stats, "ricxapp_RMR_Received 201") {
		t.Errorf("Error: ricxapp_RMR_Received value incorrect: %v", stats)
	}

	if !strings.Contains(stats, "ricxapp_RMR_TransmitError 1") {
		t.Errorf("Error: ricxapp_RMR_TransmitError value incorrect")
	}

	if !strings.Contains(stats, "ricxapp_RMR_ReceiveError 0") {
		t.Errorf("Error: ricxapp_RMR_ReceiveError value incorrect")
	}

	if waitForSdl && !strings.Contains(stats, "ricxapp_SDL_Stored 201") {
		t.Errorf("Error: ricxapp_SDL_Stored value incorrect")
	}

	if waitForSdl && !strings.Contains(stats, "ricxapp_SDL_StoreError 0") {
		t.Errorf("Error: ricxapp_SDL_StoreError value incorrect")
	}
	Rmr.Closewh(int(whid))
}

func TestSubscribeChannels(t *testing.T) {
	Logger.Info("CASE: TestSubscribeChannels")
	if !viper.GetBool("controls.waitForSdl") {
		return
	}

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

	// Misc.
	Sdl.MStoreAndPublish([]string{"channel1"}, "event", "key1", "data1")
}

func TestGetRicMessageSuccess(t *testing.T) {
	Logger.Info("CASE: TestGetRicMessageSuccess")
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
	Logger.Info("CASE: TestGetRicMessageFails")
	ok := Rmr.IsRetryError(&RMRParams{status: 0})
	if ok {
		t.Errorf("Error: IsRetryError returned wrong value")
	}

	ok = Rmr.IsRetryError(&RMRParams{status: 10})
	if !ok {
		t.Errorf("Error: IsRetryError returned wrong value")
	}

	ok = Rmr.IsNoEndPointError(&RMRParams{status: 5})
	if ok {
		t.Errorf("Error: IsNoEndPointError returned wrong value")
	}

	ok = Rmr.IsNoEndPointError(&RMRParams{status: 2})
	if !ok {
		t.Errorf("Error: IsNoEndPointError returned wrong value")
	}
}

func TestIsErrorFunctions(t *testing.T) {
	Logger.Info("CASE: TestIsErrorFunctions")
	id, ok := Rmr.GetRicMessageId("RIC_SUB_REQ")
	if !ok || id != 12010 {
		t.Errorf("Error: GetRicMessageId failed: id=%d", id)
	}

	name := Rmr.GetRicMessageName(12010)
	if name != "RIC_SUB_REQ" {
		t.Errorf("Error: GetRicMessageName failed: name=%s", name)
	}
}

func TestAddConfigChangeListener(t *testing.T) {
	Logger.Info("CASE: AddConfigChangeListener")
	AddConfigChangeListener(func(f string) {})
}

func TestConfigAccess(t *testing.T) {
	Logger.Info("CASE: TestConfigAccess")

	assert.Equal(t, Config.GetString("name"), "xapp")
	assert.Equal(t, Config.GetInt("controls.logger.level"), 3)
	assert.Equal(t, Config.GetUint32("controls.logger.level"), uint32(3))
	assert.Equal(t, Config.GetBool("controls.waitForSdl"), false)
	Config.Get("controls")
	Config.GetStringSlice("messaging.ports")
	Config.GetStringMap("messaging.ports")
	Config.IsSet("messaging")
}

func TestPublishConfigChange(t *testing.T) {
	Logger.Info("CASE: TestPublishConfigChange")
	PublishConfigChange("testApp", "values")
	ReadConfig("testApp")
}

func TestNewSubscriber(t *testing.T) {
	Logger.Info("CASE: TestNewSubscriber")
	assert.NotNil(t, NewSubscriber("", 0), "NewSubscriber failed")
}

func TestNewRMRClient(t *testing.T) {
	c := map[string]interface{}{"protPort": "tcp:4560"}
	viper.Set("rmr", c)
	assert.NotNil(t, NewRMRClient(), "NewRMRClient failed")

	params := &RMRParams{}
	params.Mtype = 1234
	params.SubId = -1
	params.Payload = []byte{1, 2, 3, 4, 5, 6}
	Rmr.SendWithRetry(params, false, 1)
}

func TestInjectRoutePrefix(t *testing.T) {
	Logger.Info("CASE: TestInjectRoutePrefix")
	assert.NotNil(t, Resource.InjectRoutePrefix("test", nil), "InjectRoutePrefix failed")
}

func TestInjectStatusCb(t *testing.T) {
	Logger.Info("CASE: TestInjectStatusCb")

	var f = func() bool {
		return true
	}
	Resource.InjectStatusCb(f)
	Resource.CheckStatus()
}

func TestSdlInterfaces(t *testing.T) {
	Sdl.Read("myKey")
	Sdl.MRead([]string{"myKey"})
	Sdl.ReadAllKeys("myKey")
	Sdl.Store("myKey", "Values")
	Sdl.MStore("myKey", "Values")
	Sdl.RegisterMetrics()
	Sdl.UpdateStatCounter("Stored")

	// Misc.
	var NotificationCb = func(ch string, events ...string) {}
	Sdl.Subscribe(NotificationCb, "channel1")
	Sdl.MSubscribe(NotificationCb, "channel1", "channel2")
	Sdl.StoreAndPublish("channel1", "event", "key1", "data1")
	Sdl.MStoreAndPublish([]string{"channel1"}, "event", "key1", "data1")
}

func TestRnibInterfaces(t *testing.T) {
	Rnib.GetNodeb("test-gnb")
	Rnib.GetCellList("test-gnb")
	Rnib.GetListGnbIds()
	Rnib.GetListEnbIds()
	Rnib.GetCountGnbList()
	Rnib.GetCell("test-gnb", 0)
	Rnib.GetCell("test-gnb", 0)
	Rnib.GetCellById(0, "cell-1")

	// Misc.
	var NotificationCb = func(ch string, events ...string) {}
	Rnib.Subscribe(NotificationCb, "channel1")
	Rnib.StoreAndPublish("channel1", "event", "key1", "data1")
}

func TestLogger(t *testing.T) {
	Logger.Error("CASE: TestNewSubscriber")
	Logger.Warn("CASE: TestNewSubscriber")
	Logger.GetLevel()
}

func TestConfigHandler(t *testing.T) {
	Logger.Error("CASE: TestConfigHandler")
	req, _ := http.NewRequest("POST", "/ric/v1/cm/appname", bytes.NewBuffer([]byte{}))
	handleFunc := http.HandlerFunc(configHandler)
	executeRequest(req, handleFunc)
}

func TestappconfigHandler(t *testing.T) {
	Logger.Error("CASE: TestappconfigHandler")
	req, _ := http.NewRequest("POST", "/ric/v1/config", bytes.NewBuffer([]byte{}))
	handleFunc := http.HandlerFunc(appconfigHandler)
	executeRequest(req, handleFunc)
}

func TestConfigChange(t *testing.T) {
	Logger.Error("CASE: TestConfigChange: %s", os.Getenv("CFG_FILE"))

	input, err := ioutil.ReadFile(os.Getenv("CFG_FILE"))
	assert.Equal(t, err, nil)

	err = ioutil.WriteFile(os.Getenv("CFG_FILE"), input, 0644)
	assert.Equal(t, err, nil)
}

func TestRegisterXapp(t *testing.T) {
	Logger.Error("CASE: TestRegisterXapp")
	doRegister()
}

func TestDeregisterXapp(t *testing.T) {
	Logger.Error("CASE: TestDeregisterXapp")
	doDeregister()
}

func TestMisc(t *testing.T) {
	Logger.Info("CASE: TestMisc")
	var cb = func() {}
	IsReady()
	SetReadyCB(func(interface{}) {}, "")
	XappReadyCb("")
	SetShutdownCB(cb)
	XappShutdownCb()
	getService("ueec", SERVICE_HTTP)

	Logger.SetFormat(1)
	Logger.SetLevel(0)
	Logger.Error("...")
	Logger.Warn("...")
	Logger.Info("...")

	mb := Rmr.Allocate(100)
	Rmr.ReAllocate(mb, 200)

	NewMetrics("", "", Resource.router)
}

func TestTeardown(t *testing.T) {
	Logger.Info("CASE: TestTeardown")
	Sdl.Delete([]string{"myKey"})
	Sdl.Clear()
	Sdl.IsReady()
	Sdl.GetStat()
	Rnib.GetNodebByGlobalNbId(1, &RNIBGlobalNbId{})
	Rnib.SaveNodeb(&RNIBNbIdentity{}, &RNIBNodebInfo{})
	go Sdl.TestConnection()
	time.Sleep(time.Duration(2) * time.Second)
}

// Helper functions
func executeRequest(req *http.Request, handleR http.HandlerFunc) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	if handleR != nil {
		vars := map[string]string{"name": "myxapp"}
		req = mux.SetURLVars(req, vars)
		handleR.ServeHTTP(rr, req)
	} else {
		vars := map[string]string{"id": "1"}
		req = mux.SetURLVars(req, vars)
		Resource.router.ServeHTTP(rr, req)
	}
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func getMetrics(t *testing.T) string {
	req, _ := http.NewRequest("GET", "/ric/v1/metrics", nil)
	response := executeRequest(req, nil)

	return response.Body.String()
}
