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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-openapi/loads"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	apiclient "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientapi"
	apicommon "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientapi/common"
	apipolicy "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientapi/policy"
	apireport "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientapi/report"
	apimodel "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientmodel"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/models"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi/operations"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi/operations/common"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi/operations/policy"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi/operations/query"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi/operations/report"
)

type SubscriptionHandler func(models.SubscriptionType, interface{}) (*models.SubscriptionResponse, error)
type SubscriptionQueryHandler func() (models.SubscriptionList, error)
type SubscriptionDeleteHandler func(string) error
type SubscriptionResponseCallback func(models.SubscriptionResponse)

type Subscriber struct {
	localAddr  string
	localPort  int
	remoteHost string
	remoteUrl  string
	remoteProt []string
	timeout    time.Duration
	clientUrl  string
	clientCB   SubscriptionResponseCallback
}

func NewSubscriber(host string, timo int) *Subscriber {
	if host == "" {
		host = "service-ricplt-submgr-http:8088"
	}

	if timo == 0 {
		timo = 20
	}

	r := &Subscriber{
		remoteHost: host,
		remoteUrl:  "/ric/v1",
		remoteProt: []string{"http"},
		timeout:    time.Duration(timo) * time.Second,
		localAddr:  "0.0.0.0",
		localPort:  8088,
		clientUrl:  "/ric/v1/subscriptions/response",
	}
	Resource.InjectRoute(r.clientUrl, r.ResponseHandler, "POST")

	return r
}

// Server interface: listen and receive subscription requests
func (r *Subscriber) Listen(createSubscription SubscriptionHandler, getSubscription SubscriptionQueryHandler, delSubscription SubscriptionDeleteHandler) error {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return err
	}

	api := operations.NewXappFrameworkAPI(swaggerSpec)

	// Subscription: Query
	api.QueryGetAllSubscriptionsHandler = query.GetAllSubscriptionsHandlerFunc(
		func(p query.GetAllSubscriptionsParams) middleware.Responder {
			if resp, err := getSubscription(); err == nil {
				return query.NewGetAllSubscriptionsOK().WithPayload(resp)
			}
			return query.NewGetAllSubscriptionsInternalServerError()
		})

	// SubscriptionType: Report
	api.ReportSubscribeReportHandler = report.SubscribeReportHandlerFunc(
		func(p report.SubscribeReportParams) middleware.Responder {
			if _, err := createSubscription(models.SubscriptionTypeReport, p.ReportParams); err == nil {
				return report.NewSubscribeReportCreated()
			}
			return report.NewSubscribeReportInternalServerError()
		})

	// SubscriptionType: Policy
	api.PolicySubscribePolicyHandler = policy.SubscribePolicyHandlerFunc(
		func(p policy.SubscribePolicyParams) middleware.Responder {
			if _, err := createSubscription(models.SubscriptionTypePolicy, p.PolicyParams); err == nil {
				return policy.NewSubscribePolicyCreated()
			}
			return policy.NewSubscribePolicyInternalServerError()
		})

	// SubscriptionType: Delete
	api.CommonUnsubscribeHandler = common.UnsubscribeHandlerFunc(
		func(p common.UnsubscribeParams) middleware.Responder {
			if err := delSubscription(p.SubscriptionID); err == nil {
				return common.NewUnsubscribeNoContent()
			}
			return common.NewUnsubscribeInternalServerError()
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
func (r *Subscriber) SubscribeReport(p *apimodel.ReportParams) (*apimodel.SubscriptionResponse, error) {
	params := apireport.NewSubscribeReportParamsWithTimeout(r.timeout).WithReportParams(p)
	result, err := r.CreateTransport().Report.SubscribeReport(params)
	if err != nil {
		return &apimodel.SubscriptionResponse{}, err
	}

	return result.Payload, err
}

// Subscription interface for xApp: POLICY
func (r *Subscriber) SubscribePolicy(p *apimodel.PolicyParams) (*apimodel.SubscriptionResponse, error) {
	params := apipolicy.NewSubscribePolicyParamsWithTimeout(r.timeout).WithPolicyParams(p)
	result, err := r.CreateTransport().Policy.SubscribePolicy(params)
	if err != nil {
		return &apimodel.SubscriptionResponse{}, err
	}

	return result.Payload, err
}

// Subscription interface for xApp: DELETE
func (r *Subscriber) UnSubscribe(subId string) error {
	params := apicommon.NewUnsubscribeParamsWithTimeout(r.timeout).WithSubscriptionID(subId)
	_, err := r.CreateTransport().Common.Unsubscribe(params)

	return err
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

func (r *Subscriber) ResponseHandler(w http.ResponseWriter, req *http.Request) {
	if req.Body != nil {
		var resp models.SubscriptionResponse
		if err := json.NewDecoder(req.Body).Decode(&resp); err == nil {
			if r.clientCB != nil {
				r.clientCB(resp)
			}
		}
		req.Body.Close()
	}
}

func (r *Subscriber) SetResponseCB(c SubscriptionResponseCallback) {
	r.clientCB = c
}

func (r *Subscriber) Notify(resp *models.SubscriptionResponse, clientEndpoint string) (err error) {
	respData, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	// Todo: delay, retry, etc. should be defined in config file
	port := strings.Split(viper.GetString("local.host"), ":")[1]
	clientUrl := fmt.Sprintf("http://%s:%s%s", clientEndpoint, port, r.clientUrl)
	for i := 0; i < 10; i++ {
		r, err := http.Post(clientUrl, "application/json", bytes.NewBuffer(respData))
		if err == nil && r.StatusCode == http.StatusOK {
			return err
		}
		Logger.Error("Post failed: %v", err)
		time.Sleep(time.Duration(5) * time.Second)
	}

	return err
}
