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
	"github.com/go-openapi/loads"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"time"

	apiclient "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientapi"
	apisub "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientapi/report"
	apimodel "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientmodel"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/models"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi/operations"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi/operations/report"
)

type SubscriptionReportHandler func(models.SubscriptionType, interface{}) (models.ReportResult, error)

type Subscriber struct {
	localAddr  string
	localPort  int
	remoteHost string
	remoteUrl  string
	remoteProt []string
	timeout    time.Duration
}

func NewSubscriber(host string, timo int) *Subscriber {
	if host == "" {
		host = "service-ricplt-submgr-http:8088"
	}

	if timo == 0 {
		timo = 20
	}

	return &Subscriber{
		remoteHost: host,
		remoteUrl:  "/ric/v1",
		remoteProt: []string{"http"},
		timeout:    time.Duration(timo) * time.Second,
		localAddr:  "0.0.0.0",
		localPort:  8088,
	}
}

// Interface for Subscription Manager to listen and receive subscription requests
func (r *Subscriber) Listen(handler SubscriptionReportHandler) error {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return err
	}

	api := operations.NewXappFrameworkAPI(swaggerSpec)
	api.ReportSubscribeReportHandler = report.SubscribeReportHandlerFunc(
		func(p report.SubscribeReportParams) middleware.Responder {
			if resp, err := handler(models.SubscriptionTypeReport, p.ReportParams); err == nil {
				return report.NewSubscribeReportCreated().WithPayload(resp)
			}
			return report.NewSubscribeReportInternalServerError()
		})

	server := restapi.NewServer(api)
	defer server.Shutdown()
	server.Host = r.localAddr
	server.Port = r.localPort

	Logger.Info("Serving subscriptions on %s:%d\n", server.Host, server.Port)
	if err := server.Serve(); err != nil {
		return err
	}
	return nil
}

// Interface for xApps to post subscriptions to Subscription Manager
func (r *Subscriber) SubscribeReport(s *apimodel.ReportParams) apimodel.ReportResult {
	params := apisub.NewSubscribeReportParamsWithTimeout(r.timeout).WithReportParams(s)
	for {
		result, err := r.CreateTransport().Report.SubscribeReport(params)
		if err == nil {
			Logger.Info("Subscription successful: payload=%v", result.Payload)
			return result.Payload
		}
		Logger.Error("SubscSubscriptionribe unsuccessful: %v", err)
		time.Sleep(time.Duration(5 * time.Second))
	}
	return apimodel.ReportResult{}
}

func (s *Subscriber) CreateTransport() *apiclient.RICSubscription {
	return apiclient.New(httptransport.New(s.remoteHost, s.remoteUrl, s.remoteProt), strfmt.Default)
}
