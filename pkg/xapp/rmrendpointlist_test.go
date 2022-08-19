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
	"fmt"
	"testing"
)

func TestRmrEndpoint1(t *testing.T) {
	addr := "127.0.0.1"
	port := uint16(8080)
	str := fmt.Sprintf("%s:%d", addr, port)
	Logger.Info("CASE: TestRmrEndpoint1 %s", str)
	ep := NewRmrEndpoint(str)
	if ep == nil || ep.Addr != addr || ep.Port != port {
		t.Errorf("NewRmrEndpoint: %s failed", str)
	}
}

func TestRmrEndpoint2(t *testing.T) {
	addr := "[2001:2003:fb69:ea00:c894:288b:4582:b5c/64]"
	port := uint16(8080)
	str := fmt.Sprintf("%s:%d", addr, port)
	Logger.Info("CASE: TestRmrEndpoint2 %s", str)
	ep := NewRmrEndpoint(str)
	if ep == nil || ep.Addr != addr || ep.Port != port {
		t.Errorf("NewRmrEndpoint: %s failed", str)
	}
}

func TestRmrEndpoint3(t *testing.T) {
	addr := "127.0.0.1"
	str := fmt.Sprintf("%s:port", addr)
	Logger.Info("CASE: TestRmrEndpoint3 %s", str)
	ep := NewRmrEndpoint(str)
	if ep != nil {
		t.Errorf("NewRmrEndpoint: %s successful while should fail", str)
	}
}

func TestRmrEndpoint4(t *testing.T) {
	addr := "127.0.0.1"
	str := fmt.Sprintf("%s:port", addr)
	Logger.Info("CASE: TestRmrEndpoint4 %s", str)
	ep := &RmrEndpoint{}
	if ep.Set(str) == true {
		t.Errorf("NewRmrEndpoint: Set %s successful while should fail", str)
	}
	if ep.Addr != "" || ep.Port != 0 {
		t.Errorf("NewRmrEndpoint: Values %s successful while should fail", str)
	}
}

func TestRmrEndpointList(t *testing.T) {
	Logger.Info("CASE: TestRmrEndpointList")

	epl := NewRmrEndpointList()

	// Simple add / has / delete
	if epl.AddEndpoint(NewRmrEndpoint("127.0.0.1:8080")) == false {
		t.Errorf("RmrEndpointList: 8080 add failed")
	}
	if epl.AddEndpoint(NewRmrEndpoint("127.0.0.1:8080")) == true {
		t.Errorf("RmrEndpointList: 8080 duplicate add success")
	}
	if epl.AddEndpoint(NewRmrEndpoint("127.0.0.1:8081")) == false {
		t.Errorf("RmrEndpointList: 8081 add failed")
	}
	if epl.HasEndpoint(NewRmrEndpoint("127.0.0.1:8081")) == false {
		t.Errorf("RmrEndpointList: 8081 has failed")
	}

	Logger.Info("%+v -- %+v -- %d", epl.String(), epl.StringList(), epl.Size())

	if epl.DelEndpoint(NewRmrEndpoint("127.0.0.1:8081")) == false {
		t.Errorf("RmrEndpointList: 8081 del failed")
	}
	if epl.HasEndpoint(NewRmrEndpoint("127.0.0.1:8081")) == true {
		t.Errorf("RmrEndpointList: 8081 has non existing success")
	}
	if epl.DelEndpoint(NewRmrEndpoint("127.0.0.1:8081")) == true {
		t.Errorf("RmrEndpointList: 8081 del non existing success")
	}
	if epl.DelEndpoint(NewRmrEndpoint("127.0.0.1:8080")) == false {
		t.Errorf("RmrEndpointList: 8080 del failed")
	}

	// list delete
	if epl.AddEndpoint(NewRmrEndpoint("127.0.0.1:8080")) == false {
		t.Errorf("RmrEndpointList: 8080 add failed")
	}
	if epl.AddEndpoint(NewRmrEndpoint("127.0.0.1:8081")) == false {
		t.Errorf("RmrEndpointList: 8081 add failed")
	}
	if epl.AddEndpoint(NewRmrEndpoint("127.0.0.1:8082")) == false {
		t.Errorf("RmrEndpointList: 8082 add failed")
	}

	epl2 := &RmrEndpointList{}
	if epl2.AddEndpoint(NewRmrEndpoint("127.0.0.1:9080")) == false {
		t.Errorf("RmrEndpointList: othlist add 9080 failed")
	}

	if epl.DelEndpoints(epl2) == true {
		t.Errorf("RmrEndpointList: delete list not existing successs")
	}

	if epl2.AddEndpoint(NewRmrEndpoint("127.0.0.1:8080")) == false {
		t.Errorf("RmrEndpointList: othlist add 8080 failed")
	}
	if epl.DelEndpoints(epl2) == false {
		t.Errorf("RmrEndpointList: delete list 8080,9080 failed")
	}

	if epl2.AddEndpoint(NewRmrEndpoint("127.0.0.1:8081")) == false {
		t.Errorf("RmrEndpointList: othlist add 8081 failed")
	}
	if epl2.AddEndpoint(NewRmrEndpoint("127.0.0.1:8082")) == false {
		t.Errorf("RmrEndpointList: othlist add 8082 failed")
	}

	if epl.DelEndpoints(epl2) == false {
		t.Errorf("RmrEndpointList: delete list 8080,8081,8082,9080 failed")
	}

}
