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
	"strings"
)

//-----------------------------------------------------------------------------
//
//-----------------------------------------------------------------------------
type RmrEndpointList struct {
	Endpoints []RmrEndpoint
}

func (eplist *RmrEndpointList) String() string {
	valuesText := eplist.StringList()
	return strings.Join(valuesText, ",")
}

func (eplist *RmrEndpointList) StringList() []string {
	tmpList := eplist.Endpoints
	valuesText := []string{}
	for i := range tmpList {
		valuesText = append(valuesText, tmpList[i].String())
	}
	return valuesText
}

func (eplist *RmrEndpointList) Size() int {
	return len(eplist.Endpoints)
}

func (eplist *RmrEndpointList) AddEndpoint(ep *RmrEndpoint) bool {
	for i := range eplist.Endpoints {
		if eplist.Endpoints[i].Equal(ep) {
			return false
		}
	}
	eplist.Endpoints = append(eplist.Endpoints, *ep)
	return true
}

func (eplist *RmrEndpointList) DelEndpoint(ep *RmrEndpoint) bool {
	for i := range eplist.Endpoints {
		if eplist.Endpoints[i].Equal(ep) {
			eplist.Endpoints[i] = eplist.Endpoints[len(eplist.Endpoints)-1]
			eplist.Endpoints[len(eplist.Endpoints)-1] = RmrEndpoint{"", 0}
			eplist.Endpoints = eplist.Endpoints[:len(eplist.Endpoints)-1]
			return true
		}
	}
	return false
}

func (eplist *RmrEndpointList) DelEndpoints(otheplist *RmrEndpointList) bool {
	var retval bool = false
	for i := range otheplist.Endpoints {
		if eplist.DelEndpoint(&otheplist.Endpoints[i]) {
			retval = true
		}
	}
	return retval
}

func (eplist *RmrEndpointList) HasEndpoint(ep *RmrEndpoint) bool {
	for i := range eplist.Endpoints {
		if eplist.Endpoints[i].Equal(ep) {
			return true
		}
	}
	return false
}

func NewRmrEndpointList() *RmrEndpointList {
	return &RmrEndpointList{}
}
