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
	"net/http/httptest"
	"testing"
)

func TestGetHealthReadyCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/ric/v1/health/ready", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(readyHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	//expected := `{"ready": true}`
	//if rr.Body.String() != expected {
	//  t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	//}
}

func TestGetHealthAliveCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/ric/v1/health/alive", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(aliveHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	//expected := `{"alive": true}`
	//if rr.Body.String() != expected {
	//  t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	//}
}
