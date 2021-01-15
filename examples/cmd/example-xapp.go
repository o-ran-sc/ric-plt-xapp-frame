/*
==================================================================================
  Copyright (c) 2020 AT&T Intellectual Property.
  Copyright (c) 2020 Nokia

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

   This source code is part of the near-RT RIC (RAN Intelligent Controller)
   platform project (RICP).
==================================================================================
*/
package main

import (
	"gerrit.o-ran-sc.org/r/ric-plt/alarm-go.git/alarm"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientmodel"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
	"net/http"
)

// This could be defined in types.go
type ExampleXapp struct {
	stats                 map[string]xapp.Counter
	rmrReady              bool
	waitForSdl            bool
	subscriptionInstances []*clientmodel.SubscriptionInstance
	subscriptionId        *string
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

func (e *ExampleXapp) Subscribe() {
	// Setup response callback to handle subscription response from SubMgr
	xapp.Subscription.SetResponseCB(func(resp *clientmodel.SubscriptionResponse) {
		if *e.subscriptionId == *resp.SubscriptionID {
			for _, s := range resp.SubscriptionInstances {
				e.subscriptionInstances = append(e.subscriptionInstances, s)
			}
		}
	})

	// Fill subscription parameters: type=REPORT
	ranName := "en-gnb:369-11105-aaaaa8"
	functionId := int64(1)
	clientEndpoint := "localhost:4560"

	reportParams := clientmodel.ReportParams{
		Meid:           ranName,
		RANFunctionID:  &functionId,
		ClientEndpoint: &clientEndpoint,
		EventTriggers: clientmodel.EventTriggerList{
			&clientmodel.EventTrigger{
				InterfaceDirection: int64(0),
				ProcedureCode:      int64(27),
				TypeOfMessage:      int64(1),
			},
		},
	}

	// Now subscribe ...
	if resp, err := xapp.Subscription.SubscribeReport(&reportParams); err == nil {
		e.subscriptionId = resp.SubscriptionID
		xapp.Logger.Info("Subscriptions for RanName='%s' [subscriptionId=%s] done!", ranName, *resp.SubscriptionID)
		return
	}

	// Subscription failed, raise alarm (only for demo purpose!)
	if err := xapp.Alarm.Raise(8006, alarm.SeverityCritical, ranName, "subscriptionFailed"); err != nil {
		xapp.Logger.Info("Raising alarm failed with error: %v", err)
	}
}

func (e *ExampleXapp) Consume(msg *xapp.RMRParams) (err error) {
	id := xapp.Rmr.GetRicMessageName(msg.Mtype)

	xapp.Logger.Info("Message received: name=%s meid=%s subId=%d txid=%s len=%d", id, msg.Meid.RanName, msg.SubId, msg.Xid, msg.PayloadLen)

	switch id {
	case "RIC_INDICATION":
		e.handleRICIndication(msg.Meid.RanName, msg)
	case "RIC_EXAMPLE_MESSAGE":
		e.handleRICExampleMessage(msg.Meid.RanName, msg)
	default:
		xapp.Logger.Info("Unknown Message Type '%d', discarding", msg.Mtype)
	}

	defer func() {
		xapp.Rmr.Free(msg.Mbuf)
		msg.Mbuf = nil
	}()
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
	xapp.Logger.SetMdc("example-xapp", "0.1.3")

	// Register various callback functions for application management
	xapp.SetReadyCB(func(d interface{}) { e.rmrReady = true }, true)
	xapp.AddConfigChangeListener(e.ConfigChangeHandler)
	xapp.Resource.InjectStatusCb(e.StatusCB)

	// Inject own REST handler for testing purpose
	xapp.Resource.InjectRoute("/ric/v1/testing", e.TestRestHandler, "POST")

	xapp.RunWithParams(e, e.waitForSdl)
}

func GetMetricsOpts() []xapp.CounterOpts {
	return []xapp.CounterOpts{
		{Name: "RICIndicationsRx", Help: "The total number of RIC inidcation events received"},
		{Name: "RICExampleMessageRx", Help: "The total number of RIC example messages received"},
	}
}

func NewExampleXapp(rmrReady bool) *ExampleXapp {
	metrics := GetMetricsOpts()
	return &ExampleXapp{
		stats:      xapp.Metric.RegisterCounterGroup(metrics, "ExampleXapp"),
		rmrReady:   rmrReady,
		waitForSdl: xapp.Config.GetBool("db.waitForSdl"),
	}
}

func main() {
	NewExampleXapp(true).Run()
}
