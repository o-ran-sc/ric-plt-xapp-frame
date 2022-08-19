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
	"strconv"
	"strings"
)

//-----------------------------------------------------------------------------
//
//-----------------------------------------------------------------------------
type RmrEndpoint struct {
	Addr string // xapp addr
	Port uint16 // xapp port
}

func (endpoint RmrEndpoint) String() string {
	return endpoint.Addr + ":" + strconv.FormatUint(uint64(endpoint.Port), 10)
}

func (endpoint *RmrEndpoint) Equal(ep *RmrEndpoint) bool {
	if (endpoint.Addr == ep.Addr) &&
		(endpoint.Port == ep.Port) {
		return true
	}
	return false
}

func (endpoint *RmrEndpoint) GetAddr() string {
	return endpoint.Addr
}

func (endpoint *RmrEndpoint) GetPort() uint16 {
	return endpoint.Port
}

func (endpoint *RmrEndpoint) Set(src string) bool {
	if strings.Contains(src, ":") {
		lind := strings.LastIndexByte(src, ':')
		srcAddr := src[:lind]
		srcPort, err := strconv.ParseUint(src[lind+1:], 10, 16)
		if err == nil {
			endpoint.Addr = srcAddr
			endpoint.Port = uint16(srcPort)
			return true
		}
	}
	endpoint.Addr = ""
	endpoint.Port = 0
	return false
}

func NewRmrEndpoint(src string) *RmrEndpoint {
	ep := &RmrEndpoint{}
	if ep.Set(src) == false {
		return nil
	}
	return ep
}
