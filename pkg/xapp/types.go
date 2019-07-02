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
	"sync"
	"unsafe"
)

// To be removed ...
type RMRStatistics struct{}

type RMRClient struct {
	context       unsafe.Pointer
	ready         int
	wg            sync.WaitGroup
	mux           sync.Mutex
	stat          map[string]Counter
	consumers     []MessageConsumer
	readyCb       ReadyCB
	readyCbParams interface{}
}

type RMRMeid struct {
	PlmnID string
	EnbID  string
}

type MessageConsumer interface {
	Consume(mtype int, sid int, payload []byte, meid *RMRMeid) error
}
