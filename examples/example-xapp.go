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
package main

import (
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"net/http"
)

// This could be defined in types.go
type ExampleXapp struct {
	msgChan  chan *xapp.RMRParams
	stats    map[string]xapp.Counter
	appReady bool
	rmrReady bool
}

func (e *ExampleXapp) handleRICIndication(ranName string, r *xapp.RMRParams) {
	// Just update metrics and store RMR message payload to SDL
	e.stats["E2APIndicationsRx"].Inc()

	xapp.Sdl.Store("myKey", r.Payload)
}

func (e *ExampleXapp) handleRICExampleMessage(ranName string, r *xapp.RMRParams) {
	// Just update metrics and echo the message back (update the message type)
	e.stats["RICExampleMessageRx"].Inc()

	r.Mtype = r.Mtype + 1
	if ok := xapp.Rmr.SendMsg(r); !ok {
		xapp.Logger.Info("Rmr.SendMsg failed ...")
	}
}

func (e *ExampleXapp) messageLoop() {
	for {
		msg := <-e.msgChan
		id := xapp.Rmr.GetRicMessageName(msg.Mtype)
		defer xapp.Rmr.Free(msg.Mbuf)

		xapp.Logger.Info("Message received: name=%s meid=%s subId=%d txid=%s len=%d", id, msg.Meid.RanName, msg.SubId, msg.Xid, msg.PayloadLen)

		switch id {
		case "RIC_INDICATION":
			e.handleRICIndication(msg.Meid.RanName, msg)
		case "RIC_EXAMPLE_MESSAGE":
			e.handleRICExampleMessage(msg.Meid.RanName, msg)
		default:
			xapp.Logger.Info("Unknown Message Type '%d', discarding", msg.Mtype)
		}
	}
}

func (e *ExampleXapp) Consume(rp *xapp.RMRParams) (err error) {
	e.msgChan <- rp
	return
}

func (u *ExampleXapp) TestRestHandler(w http.ResponseWriter, r *http.Request) {
	xapp.Logger.Info("TestRestHandler called!")
}

func (u *ExampleXapp) ConfigChangeHandler(f string) {
	xapp.Logger.Info("Config file changed, do something meaningful!")
}

func (u *ExampleXapp) StatusCB() bool {
	xapp.Logger.Info("Status callback called, do something meaningful!")
	return true
}

func (e *ExampleXapp) Run() {
	// Set MDC (read: name visible in the logs)
	xapp.Logger.SetMdc("example-xapp", "0.1.2")

	// Register various callback functions for application management
	xapp.SetReadyCB(func(d interface{}) { e.rmrReady = true }, true)
	xapp.AddConfigChangeListener(e.ConfigChangeHandler)
	xapp.Resource.InjectStatusCb(e.StatusCB)

	// Inject own REST handler for testing purpose
	xapp.Resource.InjectRoute("/ric/v1/testing", e.TestRestHandler, "POST")

	go e.messageLoop()
	xapp.Run(e)
}

func GetMetricsOpts() []xapp.CounterOpts {
	return []xapp.CounterOpts{
		{Name: "RICIndicationsRx", Help: "The total number of RIC inidcation events received"},
		{Name: "RICExampleMessageRx", Help: "The total number of RIC example messages received"},
	}
}

func NewExampleXapp(appReady, rmrReady bool) *ExampleXapp {
	metrics := GetMetricsOpts()
	return &ExampleXapp{
		stats:    xapp.Metric.RegisterCounterGroup(metrics, "ExampleXapp"),
		rmrReady: rmrReady,
		appReady: appReady,
	}
}

func main() {
	NewExampleXapp(true, false).Run()
}
