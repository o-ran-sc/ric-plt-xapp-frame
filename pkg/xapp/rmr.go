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

/*
#include <time.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <rmr/rmr.h>
#include <rmr/RIC_message_types.h>

void write_bytes_array(unsigned char *dst, void *data, int len) {
    memcpy((void *)dst, (void *)data, len);
}

#cgo CFLAGS: -I../
#cgo LDFLAGS: -lrmr_nng -lnng
*/
import "C"

import (
	"github.com/spf13/viper"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

var RMRCounterOpts = []CounterOpts{
	{Name: "Transmitted", Help: "The total number of transmited RMR messages"},
	{Name: "Received", Help: "The total number of received RMR messages"},
	{Name: "TransmitError", Help: "The total number of RMR transmission errors"},
	{Name: "ReceiveError", Help: "The total number of RMR receive errors"},
}

// To be removed ...
type RMRStatistics struct{}

type RMRClient struct {
	context   unsafe.Pointer
	ready     int
	wg        sync.WaitGroup
	mux       sync.Mutex
	stat      map[string]Counter
	consumers []MessageConsumer
}

type MessageConsumer interface {
	Consume(mtype int, sid int, len int, payload []byte) error
}

func NewRMRClient() *RMRClient {
	r := &RMRClient{}
	r.consumers = make([]MessageConsumer, 0)

	p := C.CString(viper.GetString("rmr.protPort"))
	m := C.int(viper.GetInt("rmr.maxSize"))
	defer C.free(unsafe.Pointer(p))

	r.context = C.rmr_init(p, m, C.int(0))
	if r.context == nil {
		Logger.Fatal("rmrClient: Initializing RMR context failed, bailing out!")
	}

	return r
}

func (m *RMRClient) Start(c MessageConsumer) {
	m.RegisterMetrics()

	for {
		Logger.Info("rmrClient: Waiting for RMR to be ready ...")

		if m.ready = int(C.rmr_ready(m.context)); m.ready == 1 {
			break
		}
		time.Sleep(10 * time.Second)
	}
	m.wg.Add(viper.GetInt("rmr.numWorkers"))

	if c != nil {
		m.consumers = append(m.consumers, c)
	}

	for w := 0; w < viper.GetInt("rmr.numWorkers"); w++ {
		go m.Worker("worker-"+strconv.Itoa(w), 0)
	}

	m.Wait()
}

func (m *RMRClient) Worker(taskName string, msgSize int) {
	p := viper.GetString("rmr.protPort")
	Logger.Info("rmrClient: '%s': receiving messages on [%s]", taskName, p)

	defer m.wg.Done()
	for {
		rxBuffer := C.rmr_rcv_msg(m.context, nil)
		if rxBuffer == nil {
			m.UpdateStatCounter("ReceiveError")
			continue
		}
		m.UpdateStatCounter("Received")

		go m.parseMessage(rxBuffer)
	}
}

func (m *RMRClient) parseMessage(rxBuffer *C.rmr_mbuf_t) {
	if len(m.consumers) == 0 {
		Logger.Info("rmrClient: No message handlers defined, message discarded!")
		return
	}

	for _, c := range m.consumers {
		cptr := unsafe.Pointer(rxBuffer.payload)
		payload := C.GoBytes(cptr, C.int(rxBuffer.len))

		err := c.Consume(int(rxBuffer.mtype), int(rxBuffer.sub_id), int(rxBuffer.len), payload)
		if err != nil {
			Logger.Warn("rmrClient: Consumer returned error: %v", err)
		}
	}
}

func (m *RMRClient) Allocate() *C.rmr_mbuf_t {
	buf := C.rmr_alloc_msg(m.context, 0)
	if buf == nil {
		Logger.Fatal("rmrClient: Allocating message buffer failed!")
	}

	return buf
}

func (m *RMRClient) Send(mtype int, sid int, len int, payload []byte) bool {
	buf := m.Allocate()

	buf.mtype = C.int(mtype)
	buf.sub_id = C.int(sid)
	buf.len = C.int(len)
	datap := C.CBytes(payload)
	defer C.free(datap)

	C.write_bytes_array(buf.payload, datap, C.int(len))

	return m.SendBuf(buf)
}

func (m *RMRClient) SendBuf(txBuffer *C.rmr_mbuf_t) bool {
	for i := 0; i < 10; i++ {
		txBuffer.state = 0
		txBuffer := C.rmr_send_msg(m.context, txBuffer)
		if txBuffer == nil {
			break
		} else if txBuffer.state != C.RMR_OK {
			if txBuffer.state != C.RMR_ERR_RETRY {
				time.Sleep(100 * time.Microsecond)
				m.UpdateStatCounter("TransmitError")
			}
			for j := 0; j < 100 && txBuffer.state == C.RMR_ERR_RETRY; j++ {
				txBuffer = C.rmr_send_msg(m.context, txBuffer)
			}
		}

		if txBuffer.state == C.RMR_OK {
			m.UpdateStatCounter("Transmitted")
			return true
		}
	}
	m.UpdateStatCounter("TransmitError")
	return false
}

func (m *RMRClient) UpdateStatCounter(name string) {
	m.mux.Lock()
	m.stat[name].Inc()
	m.mux.Unlock()
}

func (m *RMRClient) RegisterMetrics() {
	m.stat = Metric.RegisterCounterGroup(RMRCounterOpts, "RMR")
}

// To be removed ...
func (m *RMRClient) GetStat() (r RMRStatistics) {
	return
}

func (m *RMRClient) Wait() {
	m.wg.Wait()
}

func (m *RMRClient) IsReady() bool {
	return m.ready != 0
}

func (m *RMRClient) GetRicMessageId(mid string) int {
	return RICMessageTypes[mid]
}
