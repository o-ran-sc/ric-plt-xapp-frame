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
	"strings"
	"time"
	"unsafe"
)

var RMRCounterOpts = []CounterOpts{
	{Name: "Transmitted", Help: "The total number of transmited RMR messages"},
	{Name: "Received", Help: "The total number of received RMR messages"},
	{Name: "TransmitError", Help: "The total number of RMR transmission errors"},
	{Name: "ReceiveError", Help: "The total number of RMR receive errors"},
}

type RMRParams struct {
	Mtype      int
	Payload    []byte
	PayloadLen int
	Meid       *RMRMeid
	Xid        string
	SubId      int
	Src        string
	Mbuf       *C.rmr_mbuf_t
}

func NewRMRClient() *RMRClient {
	p := C.CString(viper.GetString("rmr.protPort"))
	m := C.int(viper.GetInt("rmr.maxSize"))
	defer C.free(unsafe.Pointer(p))

	ctx := C.rmr_init(p, m, C.int(0))
	if ctx == nil {
		Logger.Error("rmrClient: Initializing RMR context failed, bailing out!")
	}

	return &RMRClient{
		context:   ctx,
		consumers: make([]MessageConsumer, 0),
		stat:      Metric.RegisterCounterGroup(RMRCounterOpts, "RMR"),
	}
}

func (m *RMRClient) Start(c MessageConsumer) {
	if c != nil {
		m.consumers = append(m.consumers, c)
	}

	for {
		Logger.Info("rmrClient: Waiting for RMR to be ready ...")

		if m.ready = int(C.rmr_ready(m.context)); m.ready == 1 {
			break
		}
		time.Sleep(10 * time.Second)
	}
	m.wg.Add(viper.GetInt("rmr.numWorkers"))

	if m.readyCb != nil {
		go m.readyCb(m.readyCbParams)
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

	params := &RMRParams{}
	params.Mbuf = rxBuffer
	params.Mtype = int(rxBuffer.mtype)
	params.SubId = int(rxBuffer.sub_id)
	params.Meid = &RMRMeid{}

	meidBuf := make([]byte, int(C.RMR_MAX_MEID))
	if meidCstr := C.rmr_get_meid(rxBuffer, (*C.uchar)(unsafe.Pointer(&meidBuf[0]))); meidCstr != nil {
		params.Meid.RanName = strings.TrimRight(string(meidBuf), "\000")
	}

	xidBuf := make([]byte, int(C.RMR_MAX_XID))
	if xidCstr := C.rmr_get_xact(rxBuffer, (*C.uchar)(unsafe.Pointer(&xidBuf[0]))); xidCstr != nil {
		params.Xid = strings.TrimRight(string(xidBuf[0:32]), "\000")
	}

	srcBuf := make([]byte, int(C.RMR_MAX_SRC))
	if srcStr := C.rmr_get_src(rxBuffer, (*C.uchar)(unsafe.Pointer(&srcBuf[0]))); srcStr != nil {
		params.Src = strings.TrimRight(string(srcBuf[0:64]), "\000")
	}

	for _, c := range m.consumers {
		cptr := unsafe.Pointer(rxBuffer.payload)
		params.Payload = C.GoBytes(cptr, C.int(rxBuffer.len))
		params.PayloadLen = int(rxBuffer.len)

		err := c.Consume(params)
		if err != nil {
			Logger.Warn("rmrClient: Consumer returned error: %v", err)
		}
	}
}

func (m *RMRClient) Allocate() *C.rmr_mbuf_t {
	buf := C.rmr_alloc_msg(m.context, 0)
	if buf == nil {
		Logger.Error("rmrClient: Allocating message buffer failed!")
	}
	return buf
}

func (m *RMRClient) Free(mbuf *C.rmr_mbuf_t) {
	if mbuf == nil {
		Logger.Error("rmrClient: Can't free mbuffer, given nil pointer")
		return
	}
	C.rmr_free_msg(mbuf)
}

func (m *RMRClient) SendMsg(params *RMRParams) bool {
	return m.SendBuffer(params, false)
}

func (m *RMRClient) SendRts(params *RMRParams) bool {
	return m.SendBuffer(params, true)
}

func (m *RMRClient) SendBuffer(params *RMRParams, isRts bool) bool {
	defer m.Free(params.Mbuf)
	for i := 0; i < 10; i++ {
		errCode := m.Send(params, isRts)
		if errCode == C.RMR_OK {
			m.UpdateStatCounter("Transmitted")
			return true
		}
		if errCode != C.RMR_ERR_RETRY {
			Logger.Error("rmrClient: rmr_send returned hard error - %d", errCode)
			break
		}

	}
	m.UpdateStatCounter("TransmitError")
	return false
}

func (m *RMRClient) Send(params *RMRParams, isRts bool) C.int {
	txBuffer := params.Mbuf
	if txBuffer == nil {
		txBuffer = m.Allocate()
	}

	txBuffer.mtype = C.int(params.Mtype)
	txBuffer.sub_id = C.int(params.SubId)
	txBuffer.len = C.int(len(params.Payload))
	if params.PayloadLen != 0 {
		txBuffer.len = C.int(params.PayloadLen)
	}
	datap := C.CBytes(params.Payload)
	defer C.free(datap)

	if params != nil {
		if params.Meid != nil {
			b := make([]byte, int(C.RMR_MAX_MEID))
			copy(b, []byte(params.Meid.RanName))
			C.rmr_bytes2meid(txBuffer, (*C.uchar)(unsafe.Pointer(&b[0])), C.int(len(b)))
		}
		xidLen := len(params.Xid)
		if xidLen > 0 && xidLen <= C.RMR_MAX_XID {
			b := make([]byte, int(C.RMR_MAX_MEID))
			copy(b, []byte(params.Xid))
			C.rmr_bytes2xact(txBuffer, (*C.uchar)(unsafe.Pointer(&b[0])), C.int(len(b)))
		}
	}
	C.write_bytes_array(txBuffer.payload, datap, txBuffer.len)

	txBuffer.state = 0
	currBuffer := txBuffer
	if isRts {
		currBuffer = C.rmr_rts_msg(m.context, txBuffer)
	} else {
		currBuffer = C.rmr_send_msg(m.context, txBuffer)
	}

	if currBuffer != nil {
		return currBuffer.state
	}
	return -1
}

func (m *RMRClient) UpdateStatCounter(name string) {
	m.mux.Lock()
	m.stat[name].Inc()
	m.mux.Unlock()
}

func (m *RMRClient) RegisterMetrics() {
	m.stat = Metric.RegisterCounterGroup(RMRCounterOpts, "RMR")
}

func (m *RMRClient) Wait() {
	m.wg.Wait()
}

func (m *RMRClient) IsReady() bool {
	return m.ready != 0
}

func (m *RMRClient) SetReadyCB(cb ReadyCB, params interface{}) {
	m.readyCb = cb
	m.readyCbParams = params
}

func (m *RMRClient) GetRicMessageId(name string) (int, bool) {
	id, ok := RICMessageTypes[name]
	return id, ok
}

func (m *RMRClient) GetRicMessageName(id int) (s string) {
	for k, v := range RICMessageTypes {
		if id == v {
			return k
		}
	}
	return
}

// To be removed ...
func (m *RMRClient) GetStat() (r RMRStatistics) {
	return
}
