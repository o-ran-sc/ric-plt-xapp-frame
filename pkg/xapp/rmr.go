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
#include <sys/epoll.h>
#include <unistd.h>
#include <rmr/rmr.h>
#include <rmr/RIC_message_types.h>


void write_bytes_array(unsigned char *dst, void *data, int len) {
    memcpy((void *)dst, (void *)data, len);
}

int init_epoll(int rcv_fd) {
	struct	epoll_event epe;
	int epoll_fd = epoll_create1( 0 );
	epe.events = EPOLLIN;
	epe.data.fd = rcv_fd;
	epoll_ctl( epoll_fd, EPOLL_CTL_ADD, rcv_fd, &epe );
	return epoll_fd;
}

void close_epoll(int epoll_fd) {
	if(epoll_fd >= 0) {
		close(epoll_fd);
	}
}

int wait_epoll(int epoll_fd,int rcv_fd) {
	struct	epoll_event events[1];
	if( epoll_wait( epoll_fd, events, 1, -1 ) > 0 ) {
		if( events[0].data.fd == rcv_fd ) {
			return 1;
		}
	}
	return 0;
}

#cgo CFLAGS: -I../
#cgo LDFLAGS: -lrmr_si
*/
import "C"

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/spf13/viper"
)

var RMRCounterOpts = []CounterOpts{
	{Name: "Transmitted", Help: "The total number of transmited RMR messages"},
	{Name: "TransmitError", Help: "The total number of RMR transmission errors"},
	{Name: "TransmitRetry", Help: "The total number of transmit retries on failure"},
	{Name: "Received", Help: "The total number of received RMR messages"},
	{Name: "ReceiveError", Help: "The total number of RMR receive errors"},
	{Name: "SendWithRetryRetry", Help: "SendWithRetry service retries"},
}

var RMRGaugeOpts = []CounterOpts{
	{Name: "Enqueued", Help: "The total number of enqueued in RMR library"},
	{Name: "Dropped", Help: "The total number of dropped in RMR library"},
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

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------
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

func (params *RMRParams) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "params(Src=%s Mtype=%d SubId=%d Xid=%s Meid=%s Paylens=%d/%d Paymd5=%x)", params.Src, params.Mtype, params.SubId, params.Xid, params.Meid, params.PayloadLen, len(params.Payload), md5.Sum(params.Payload))
	return b.String()
}

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------
type RMRClientParams struct {
	StatDesc string
	RmrData  PortData
}

func (params *RMRClientParams) String() string {
	return fmt.Sprintf("ProtPort=%d MaxSize=%d ThreadType=%d StatDesc=%s LowLatency=%t FastAck=%t Policies=%v",
		params.RmrData.Port, params.RmrData.MaxSize, params.RmrData.ThreadType, params.StatDesc,
		params.RmrData.LowLatency, params.RmrData.FastAck, params.RmrData.Policies)
}

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------
func NewRMRClientWithParams(params *RMRClientParams) *RMRClient {
	p := C.CString(fmt.Sprintf("%d", params.RmrData.Port))
	m := C.int(params.RmrData.MaxSize)
	c := C.int(params.RmrData.ThreadType)
	defer C.free(unsafe.Pointer(p))
	ctx := C.rmr_init(p, m, c)
	if ctx == nil {
		Logger.Error("rmrClient: Initializing RMR context failed, bailing out!")
	}

	Logger.Info("new rmrClient with parameters: %s", params.String())

	if params.RmrData.LowLatency {
		C.rmr_set_low_latency(ctx)
	}
	if params.RmrData.FastAck {
		C.rmr_set_fack(ctx)
	}

	return &RMRClient{
		context:           ctx,
		consumers:         make([]MessageConsumer, 0),
		statc:             Metric.RegisterCounterGroup(RMRCounterOpts, params.StatDesc),
		statg:             Metric.RegisterGaugeGroup(RMRGaugeOpts, params.StatDesc),
		maxRetryOnFailure: params.RmrData.MaxRetryOnFailure,
	}
}

func NewRMRClient() *RMRClient {
	p := GetPortData("rmrdata")
	if p.Port == 0 || viper.IsSet("rmr.protPort") {
		// Old xApp descriptor used, fallback to rmr section
		fmt.Sscanf(viper.GetString("rmr.protPort"), "tcp:%d", &p.Port)
		p.MaxSize = viper.GetInt("rmr.maxSize")
		p.ThreadType = viper.GetInt("rmr.threadType")
		p.LowLatency = viper.GetBool("rmr.lowLatency")
		p.FastAck = viper.GetBool("rmr.fastAck")
		p.MaxRetryOnFailure = viper.GetInt("rmr.maxRetryOnFailure")
	}

	return NewRMRClientWithParams(
		&RMRClientParams{
			RmrData:  p,
			StatDesc: "RMR",
		})
}

func (m *RMRClient) Start(c MessageConsumer) {
	if c != nil {
		m.consumers = append(m.consumers, c)
	}

	var counter int = 0
	for {
		m.contextMux.Lock()
		m.ready = int(C.rmr_ready(m.context))
		m.contextMux.Unlock()
		if m.ready == 1 {
			Logger.Info("rmrClient: RMR is ready after %d seconds waiting...", counter)
			break
		}
		if counter%10 == 0 {
			Logger.Info("rmrClient: Waiting for RMR to be ready ...")
		}
		time.Sleep(1 * time.Second)
		counter++
	}

	if m.readyCb != nil {
		go m.readyCb(m.readyCbParams)
	}

	m.wg.Add(1)
	go func() {
		m.contextMux.Lock()
		rfd := C.rmr_get_rcvfd(m.context)
		m.contextMux.Unlock()
		efd := C.init_epoll(rfd)

		defer m.wg.Done()
		for {

			if int(C.wait_epoll(efd, rfd)) == 0 {
				continue
			}
			m.contextMux.Lock()
			rxBuffer := C.rmr_rcv_msg(m.context, nil)
			m.contextMux.Unlock()

			if rxBuffer == nil {
				m.LogMBufError("RecvMsg failed", rxBuffer)
				m.UpdateStatCounter("ReceiveError")
				continue
			}
			m.UpdateStatCounter("Received")
			m.parseMessage(rxBuffer)
		}
	}()

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		for {
			m.UpdateRmrStats()
			time.Sleep(1 * time.Second)
		}
	}()

	m.wg.Wait()
}

func (m *RMRClient) UpdateRmrStats() {
	param := (*C.rmr_rx_debug_t)(C.malloc(C.size_t(unsafe.Sizeof(C.rmr_rx_debug_t{}))))
	m.contextMux.Lock()
	C.rmr_get_rx_debug_info(m.context, param)
	m.contextMux.Unlock()
	m.mux.Lock()
	m.statg["Enqueued"].Set(float64(param.enqueue))
	m.statg["Dropped"].Set(float64(param.drop))
	m.mux.Unlock()
	C.free(unsafe.Pointer(param))
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

	/*
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
	*/
}

func (m *RMRClient) Allocate(size int) *C.rmr_mbuf_t {
	m.contextMux.Lock()
	defer m.contextMux.Unlock()
	outbuf := C.rmr_alloc_msg(m.context, C.int(size))
	if outbuf == nil {
		Logger.Error("rmrClient: Allocating message buffer failed!")
	}
	return outbuf
}

func (m *RMRClient) ReAllocate(inbuf *C.rmr_mbuf_t, size int) *C.rmr_mbuf_t {
	m.contextMux.Lock()
	defer m.contextMux.Unlock()
	outbuf := C.rmr_realloc_msg(inbuf, C.int(size))
	if outbuf == nil {
		Logger.Error("rmrClient: Allocating message buffer failed!")
	}
	return outbuf
}

func (m *RMRClient) Free(mbuf *C.rmr_mbuf_t) {
	if mbuf == nil {
		return
	}
	m.contextMux.Lock()
	defer m.contextMux.Unlock()
	C.rmr_free_msg(mbuf)
}

func (m *RMRClient) SendMsg(params *RMRParams) bool {
	return m.Send(params, false)
}

func (m *RMRClient) SendRts(params *RMRParams) bool {
	return m.Send(params, true)
}

func (m *RMRClient) SendWithRetry(params *RMRParams, isRts bool, to time.Duration) (err error) {
	status := m.Send(params, isRts)
	i := 0
	for ; i < int(to)*2 && status == false; i++ {
		status = m.Send(params, isRts)
		if status == false {
			m.UpdateStatCounter("SendWithRetryRetry")
			time.Sleep(500 * time.Millisecond)
		}
	}
	if status == false {
		err = fmt.Errorf("Failed with retries(%d) %s", i, params.String())
		if params.Mbuf != nil {
			m.Free(params.Mbuf)
			params.Mbuf = nil
		}
	}
	return
}

func (m *RMRClient) CopyBuffer(params *RMRParams) *C.rmr_mbuf_t {

	if params == nil {
		return nil
	}

	payLen := len(params.Payload)
	if params.PayloadLen != 0 {
		payLen = params.PayloadLen
	}

	txBuffer := params.Mbuf
	params.Mbuf = nil

	if txBuffer != nil {
		txBuffer = m.ReAllocate(txBuffer, payLen)
	} else {
		txBuffer = m.Allocate(payLen)
	}

	if txBuffer == nil {
		return nil
	}
	txBuffer.mtype = C.int(params.Mtype)
	txBuffer.sub_id = C.int(params.SubId)
	txBuffer.len = C.int(payLen)

	datap := C.CBytes(params.Payload)
	defer C.free(datap)

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
	txBuffer.state = 0

	// Just quick retry seems to help for K8s issue
	if m.maxRetryOnFailure == 0 {
		m.maxRetryOnFailure = 5
	}

	for j := 0; j <= m.maxRetryOnFailure; j++ {
		m.contextMux.Lock()
		if whid != 0 {
			txBuffer = C.rmr_wh_send_msg(m.context, C.rmr_whid_t(whid), txBuffer)
		} else {
			if isRts {
				txBuffer = C.rmr_rts_msg(m.context, txBuffer)
			} else {
				txBuffer = C.rmr_send_msg(m.context, txBuffer)
			}
		}
		m.contextMux.Unlock()
		if j+1 <= m.maxRetryOnFailure && txBuffer != nil && txBuffer.state == C.RMR_ERR_RETRY {
			m.UpdateStatCounter("TransmitRetry")
			continue
		}
		break
	}

	if txBuffer == nil {
		m.UpdateStatCounter("TransmitError")
		m.LogMBufError("SendBuf failed", txBuffer)
		return int(C.RMR_ERR_INITFAILED)
	}

	if txBuffer.state != C.RMR_OK {
		m.UpdateStatCounter("TransmitError")
		m.LogMBufError("SendBuf failed", txBuffer)
	} else {
		m.UpdateStatCounter("Transmitted")
	}
	defer m.Free(txBuffer)
	return int(txBuffer.state)

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

	m.contextMux.Lock()
	currBuffer = C.rmr_wh_call(m.context, C.int(params.Whid), txBuffer, C.int(params.Callid), C.int(params.Timeout))
	m.contextMux.Unlock()

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
	m.contextMux.Lock()
	defer m.contextMux.Unlock()
	endpoint := C.CString(target)
	return C.rmr_wh_open(m.context, endpoint)
}

func (m *RMRClient) Closewh(whid int) {
	m.Wh_close(C.rmr_whid_t(whid))
}

func (m *RMRClient) Wh_close(whid C.rmr_whid_t) {
	m.contextMux.Lock()
	defer m.contextMux.Unlock()
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
	m.statc[name].Inc()
	m.mux.Unlock()
}

func (m *RMRClient) RegisterMetrics() {
	m.statc = Metric.RegisterCounterGroup(RMRCounterOpts, "RMR")
	m.statg = Metric.RegisterGaugeGroup(RMRGaugeOpts, "RMR")
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
	if mbuf != nil {
		Logger.Debug(fmt.Sprintf("rmrClient: %s -> [tp=%v] %v - %s", text, mbuf.tp_state, mbuf.state, RMRErrors[int(mbuf.state)]))
		return int(mbuf.state)
	}
	Logger.Debug(fmt.Sprintf("rmrClient: %s -> mbuf nil", text))
	return 0
}
