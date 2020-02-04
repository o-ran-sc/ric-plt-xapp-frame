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
	"time"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	httptransport "github.com/go-openapi/runtime/client"

	apiclient "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientapi"
	apimodel "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientmodel"
	apisub "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientapi/subscription"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/models"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi/operations"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi/operations/subscription"
)

type SubscriptionHandler func(*subscription.SubscribeParams) (*models.SubscriptionResult, error)

type Subscriber struct {
	localAddr 	string
	localPort 	int
	remoteHost 	string
	remoteUrl	string
	remoteProt  []string
	timeout 	time.Duration
}

func NewSubscriber(host, timo int) *Subscriber {
	if host == "" {
		host = "service-ricplt-submgr-http:8088"
	}

	if timo == 0 {
		timo = 20
	}

	return &Subscriber{
		remoteHost: host, 
		remoteUrl: "/ric/v1", 
		remoteProt: []string{"http"}, 
		timeout: time.Duration(timo) * time.Second,
		localAddr: "0.0.0.0",
		localPort: 8088,
	}
}

func (r *Subscriber) Listen(handler SubscriptionHandler) error {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return err
	}
	
	api := operations.NewXappFrameworkAPI(swaggerSpec)
	api.SubscriptionSubscribeHandler = subscription.SubscribeHandlerFunc(
		func(params subscription.SubscribeParams) middleware.Responder {
			if result, err := handler(&params); err == nil {
				return subscription.NewSubscribeCreated().WithPayload(result)
			}
			return subscription.NewSubscribeInternalServerError()
		})

	server := restapi.NewServer(api)
	defer server.Shutdown()
	server.Host = r.localAddr
	server.Port = r.localPort

	Logger.Info("Subscription serving on %s:%d\n", server.Host, server.Port)
	if err := server.Serve(); err != nil {
		return err
	}
	return nil
}

func (r *Subscriber) Subscribe(s *apimodel.Subscription) *apimodel.SubscriptionResult {
	params := apisub.NewSubscribeParamsWithTimeout(r.timeout).WithSubscriptionParams(s)
	for {
		result, err := r.CreateTransport().Subscription.Subscribe(params)
		if err == nil {
			Logger.Info("Subscription successful: payload=%v", result.Payload)
			return result.Payload
		}
		Logger.Error("SubscSubscriptionribe unsuccessful: %v", err)
		time.Sleep(time.Duration(5 * time.Second))
	}
	return &apimodel.SubscriptionResult{}
}

func (s *Subscriber) CreateTransport() *apiclient.RICSubscription {
	return apiclient.New(httptransport.New(s.remoteHost, s.remoteUrl, s.remoteProt), strfmt.Default)
}