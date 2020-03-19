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
	"encoding/json"
	"fmt"
	"github.com/go-openapi/loads"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"io/ioutil"
	"net/http"
	"time"

	apiclient "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientapi"
	apicontrol "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientapi/control"
	apipolicy "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientapi/policy"
	apireport "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientapi/report"
	apimodel "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientmodel"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/models"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi/operations"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi/operations/control"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi/operations/policy"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi/operations/query"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi/operations/report"
)

type SubscriptionHandler func(models.SubscriptionType, interface{}) (models.SubscriptionResult, error)
type SubscriptionQueryHandler func() (models.SubscriptionList, error)

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

// Server interface: listen and receive subscription requests
func (r *Subscriber) Listen(subHandler SubscriptionHandler, queryHandler SubscriptionQueryHandler) error {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return err
	}

	api := operations.NewXappFrameworkAPI(swaggerSpec)

	// Subscription: query
	api.QueryGetAllSubscriptionsHandler = query.GetAllSubscriptionsHandlerFunc(
		func(p query.GetAllSubscriptionsParams) middleware.Responder {
			if resp, err := queryHandler(); err == nil {
				return query.NewGetAllSubscriptionsOK().WithPayload(resp)
			}
			return query.NewGetAllSubscriptionsInternalServerError()
		})

	// SubscriptionType: Report
	api.ReportSubscribeReportHandler = report.SubscribeReportHandlerFunc(
		func(p report.SubscribeReportParams) middleware.Responder {
			if resp, err := subHandler(models.SubscriptionTypeReport, p.ReportParams); err == nil {
				return report.NewSubscribeReportCreated().WithPayload(resp)
			}
			return report.NewSubscribeReportInternalServerError()
		})

	// SubscriptionType: Control
	api.ControlSubscribeControlHandler = control.SubscribeControlHandlerFunc(
		func(p control.SubscribeControlParams) middleware.Responder {
			if resp, err := subHandler(models.SubscriptionTypeControl, p.ControlParams); err == nil {
				return control.NewSubscribeControlCreated().WithPayload(resp)
			}
			return control.NewSubscribeControlInternalServerError()
		})

	// SubscriptionType: policy
	api.PolicySubscribePolicyHandler = policy.SubscribePolicyHandlerFunc(
		func(p policy.SubscribePolicyParams) middleware.Responder {
			if resp, err := subHandler(models.SubscriptionTypePolicy, p.PolicyParams); err == nil {
				return policy.NewSubscribePolicyCreated().WithPayload(resp)
			}
			return policy.NewSubscribePolicyInternalServerError()
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

// Subscription interface for xApp: REPORT
func (r *Subscriber) SubscribeReport(p *apimodel.ReportParams) (apimodel.SubscriptionResult, error) {
	params := apireport.NewSubscribeReportParamsWithTimeout(r.timeout).WithReportParams(p)
	result, err := r.CreateTransport().Report.SubscribeReport(params)
	if err != nil {
		return apimodel.SubscriptionResult{}, err
	}

	return result.Payload, err
}

// Subscription interface for xApp: CONTROL
func (r *Subscriber) SubscribeControl(p *apimodel.ControlParams) (apimodel.SubscriptionResult, error) {
	params := apicontrol.NewSubscribeControlParamsWithTimeout(r.timeout).WithControlParams(p)
	result, err := r.CreateTransport().Control.SubscribeControl(params)
	if err != nil {
		return apimodel.SubscriptionResult{}, err
	}

	return result.Payload, err
}

// Subscription interface for xApp: POLICY
func (r *Subscriber) SubscribePolicy(p *apimodel.PolicyParams) (apimodel.SubscriptionResult, error) {
	params := apipolicy.NewSubscribePolicyParamsWithTimeout(r.timeout).WithPolicyParams(p)
	result, err := r.CreateTransport().Policy.SubscribePolicy(params)
	if err != nil {
		return apimodel.SubscriptionResult{}, err
	}

	return result.Payload, err
}

// Subscription interface for xApp: QUERY
func (r *Subscriber) QuerySubscriptions() (models.SubscriptionList, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/%s/subscriptions", r.remoteHost, r.remoteUrl))
	if err != nil {
		return models.SubscriptionList{}, err
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.SubscriptionList{}, err
	}

	subscriptions := models.SubscriptionList{}
	err = json.Unmarshal([]byte(string(contents)), &subscriptions)
	if err != nil {
		return models.SubscriptionList{}, err
	}

	return subscriptions, nil
}

func (r *Subscriber) CreateTransport() *apiclient.RICSubscription {
	return apiclient.New(httptransport.New(r.remoteHost, r.remoteUrl, r.remoteProt), strfmt.Default)
}