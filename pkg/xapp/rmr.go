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
#cgo LDFLAGS: -lrmr_si
*/
import "C"

import (
	"fmt"
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

var RMRErrors = map[int]string{
	C.RMR_OK:             "state is good",
	C.RMR_ERR_BADARG:     "argument passed to function was unusable",
	C.RMR_ERR_NOENDPT:    "send/call could not find an endpoint based on msg type",
	C.RMR_ERR_EMPTY:      "msg received had no payload; attempt to send an empty message",
	C.RMR_ERR_NOHDR:      "message didn't contain a valid header",
	C.RMR_ERR_SENDFAILED: "send failed; errno has nano reason",
	C.RMR_ERR_CALLFAILED: "unable to send call() message",
	C.RMR_ERR_NOWHOPEN:   "no wormholes are open",
	C.RMR_ERR_WHID:       "wormhole id was invalid",
	C.RMR_ERR_OVERFLOW:   "operation would have busted through a buffer/field size",
	C.RMR_ERR_RETRY:      "request (send/call/rts) failed, but caller should retry (EAGAIN for wrappers)",
	C.RMR_ERR_RCVFAILED:  "receive failed (hard error)",
	C.RMR_ERR_TIMEOUT:    "message processing call timed out",
	C.RMR_ERR_UNSET:      "the message hasn't been populated with a transport buffer",
	C.RMR_ERR_TRUNC:      "received message likely truncated",
	C.RMR_ERR_INITFAILED: "initialization of something (probably message) failed",
	C.RMR_ERR_NOTSUPP:    "the request is not supported, or RMr was not initialized for the request",
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
	Whid       int
	Callid     int
	Timeout    int
	status     int
}

func NewRMRClientWithParams(protPort string, maxSize int, numWorkers int, threadType int, statDesc string) *RMRClient {
	p := C.CString(protPort)
	m := C.int(maxSize)
	c := C.int(threadType)
	defer C.free(unsafe.Pointer(p))

	//ctx := C.rmr_init(p, m, C.int(0))
	//ctx := C.rmr_init(p, m, C.RMRFL_NOTHREAD)
	ctx := C.rmr_init(p, m, c)
	if ctx == nil {
		Logger.Error("rmrClient: Initializing RMR context failed, bailing out!")
	}

	return &RMRClient{
		protPort:   protPort,
		numWorkers: numWorkers,
		context:    ctx,
		consumers:  make([]MessageConsumer, 0),
		stat:       Metric.RegisterCounterGroup(RMRCounterOpts, statDesc),
	}
}

func NewRMRClient() *RMRClient {
	return NewRMRClientWithParams(viper.GetString("rmr.protPort"), viper.GetInt("rmr.maxSize"), viper.GetInt("rmr.numWorkers"), viper.GetInt("rmr.threadType"), "RMR")
}

func (m *RMRClient) Start(c MessageConsumer) {
	if c != nil {
		m.consumers = append(m.consumers, c)
	}

	var counter int = 0
	for {
		if m.ready = int(C.rmr_ready(m.context)); m.ready == 1 {
			Logger.Info("rmrClient: RMR is ready after %d seconds waiting...", counter)
			break
		}
		if counter%10 == 0 {
			Logger.Info("rmrClient: Waiting for RMR to be ready ...")
		}
		time.Sleep(1 * time.Second)
		counter++
	}
	m.wg.Add(m.numWorkers)

	if m.readyCb != nil {
		go m.readyCb(m.readyCbParams)
	}

	for w := 0; w < m.numWorkers; w++ {
		go m.Worker("worker-"+strconv.Itoa(w), 0)
	}
	m.Wait()
}

func (m *RMRClient) Worker(taskName string, msgSize int) {
	Logger.Info("rmrClient: '%s': receiving messages on [%s]", taskName, m.protPort)

	defer m.wg.Done()
	for {
		rxBuffer := C.rmr_rcv_msg(m.context, nil)
		if rxBuffer == nil {
			m.LogMBufError("RecvMsg failed", rxBuffer)
			m.UpdateStatCounter("ReceiveError")
			continue
		}
		m.UpdateStatCounter("Received")

		m.msgWg.Add(1)
		go m.parseMessage(rxBuffer)
		m.msgWg.Wait()
	}
}

func (m *RMRClient) parseMessage(rxBuffer *C.rmr_mbuf_t) {
	defer m.msgWg.Done()
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

	// Default case: a single consumer
	if len(m.consumers) == 1 && m.consumers[0] != nil {
		params.PayloadLen = int(rxBuffer.len)
		params.Payload = (*[1 << 30]byte)(unsafe.Pointer(rxBuffer.payload))[:params.PayloadLen:params.PayloadLen]
		err := m.consumers[0].Consume(params)
		if err != nil {
			Logger.Warn("rmrClient: Consumer returned error: %v", err)
		}
		return
	}

	// Special case for multiple consumers
	for _, c := range m.consumers {
		cptr := unsafe.Pointer(rxBuffer.payload)
		params.Payload = C.GoBytes(cptr, C.int(rxBuffer.len))
		params.PayloadLen = int(rxBuffer.len)
		params.Mtype = int(rxBuffer.mtype)
		params.SubId = int(rxBuffer.sub_id)

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
		return
	}
	C.rmr_free_msg(mbuf)
}

func (m *RMRClient) SendMsg(params *RMRParams) bool {
	return m.Send(params, false)
}

func (m *RMRClient) SendRts(params *RMRParams) bool {
	return m.Send(params, true)
}

func (m *RMRClient) CopyBuffer(params *RMRParams) *C.rmr_mbuf_t {
	txBuffer := params.Mbuf
	if txBuffer == nil {
		txBuffer = m.Allocate()
		if txBuffer == nil {
			return nil
		}
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
			b := make([]byte, int(C.RMR_MAX_XID))
			copy(b, []byte(params.Xid))
			C.rmr_bytes2xact(txBuffer, (*C.uchar)(unsafe.Pointer(&b[0])), C.int(len(b)))
		}
	}
	C.write_bytes_array(txBuffer.payload, datap, txBuffer.len)
	return txBuffer
}

func (m *RMRClient) Send(params *RMRParams, isRts bool) bool {

	txBuffer := m.CopyBuffer(params)
	if txBuffer == nil {
		return false
	}
	params.status = m.SendBuf(txBuffer, isRts, params.Whid)
	if params.status == int(C.RMR_OK) {
		return true
	}
	return false
}

func (m *RMRClient) SendBuf(txBuffer *C.rmr_mbuf_t, isRts bool, whid int) int {
	var (
		currBuffer  *C.rmr_mbuf_t
		counterName string = "Transmitted"
	)

	txBuffer.state = 0
	if whid != 0 {
		currBuffer = C.rmr_wh_send_msg(m.context, C.rmr_whid_t(whid), txBuffer)
	} else {
		if isRts {
			currBuffer = C.rmr_rts_msg(m.context, txBuffer)
		} else {
			currBuffer = C.rmr_send_msg(m.context, txBuffer)
		}
	}

	if currBuffer == nil {
		m.UpdateStatCounter("TransmitError")
		return m.LogMBufError("SendBuf failed", txBuffer)
	}

	// Just quick retry seems to help for K8s issue
	maxRetryOnFailure := viper.GetInt("rmr.maxRetryOnFailure")
	if maxRetryOnFailure == 0 {
		maxRetryOnFailure = 5
	}

	for j := 0; j < maxRetryOnFailure && currBuffer != nil && currBuffer.state == C.RMR_ERR_RETRY; j++ {
		if whid != 0 {
			currBuffer = C.rmr_wh_send_msg(m.context, C.rmr_whid_t(whid), txBuffer)
		} else {
			if isRts {
				currBuffer = C.rmr_rts_msg(m.context, txBuffer)
			} else {
				currBuffer = C.rmr_send_msg(m.context, txBuffer)
			}
		}
	}

	if currBuffer.state != C.RMR_OK {
		counterName = "TransmitError"
		m.LogMBufError("SendBuf failed", currBuffer)
	}

	m.UpdateStatCounter(counterName)
	defer m.Free(currBuffer)

	return int(currBuffer.state)
}

func (m *RMRClient) SendCallMsg(params *RMRParams) (int, string) {
	var (
		currBuffer  *C.rmr_mbuf_t
		counterName string = "Transmitted"
	)
	txBuffer := m.CopyBuffer(params)
	if txBuffer == nil {
		return C.RMR_ERR_INITFAILED, ""
	}

	txBuffer.state = 0

	currBuffer = C.rmr_wh_call(m.context, C.int(params.Whid), txBuffer, C.int(params.Callid), C.int(params.Timeout))

	if currBuffer == nil {
		m.UpdateStatCounter("TransmitError")
		return m.LogMBufError("SendBuf failed", txBuffer), ""
	}

	if currBuffer.state != C.RMR_OK {
		counterName = "TransmitError"
		m.LogMBufError("SendBuf failed", currBuffer)
	}

	m.UpdateStatCounter(counterName)
	defer m.Free(currBuffer)

	cptr := unsafe.Pointer(currBuffer.payload)
	payload := C.GoBytes(cptr, C.int(currBuffer.len))

	return int(currBuffer.state), string(payload)
}

func (m *RMRClient) Openwh(target string) C.rmr_whid_t {
	return m.Wh_open(target)
}

func (m *RMRClient) Wh_open(target string) C.rmr_whid_t {
	endpoint := C.CString(target)
	return C.rmr_wh_open(m.context, endpoint)
}

func (m *RMRClient) Closewh(whid int) {
	m.Wh_close(C.rmr_whid_t(whid))
}

func (m *RMRClient) Wh_close(whid C.rmr_whid_t) {
	C.rmr_wh_close(m.context, whid)
}

func (m *RMRClient) IsRetryError(params *RMRParams) bool {
	if params.status == int(C.RMR_ERR_RETRY) {
		return true
	}
	return false
}

func (m *RMRClient) IsNoEndPointError(params *RMRParams) bool {
	if params.status == int(C.RMR_ERR_NOENDPT) {
		return true
	}
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

func (m *RMRClient) LogMBufError(text string, mbuf *C.rmr_mbuf_t) int {
	Logger.Debug(fmt.Sprintf("rmrClient: %s -> [tp=%v] %v - %s", text, mbuf.tp_state, mbuf.state, RMRErrors[int(mbuf.state)]))
	return int(mbuf.state)
}

// To be removed ...
func (m *RMRClient) GetStat() (r RMRStatistics) {
	return
}
