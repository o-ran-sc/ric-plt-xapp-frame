/*
==================================================================================
  Copyright (c) 2019 Nokia
==================================================================================
*/

package xapp

import (
	apimodel "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientmodel"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var suite *testing.T

var meid = "gnb123456"
var funId = int64(1)
var clientEndpoint = "localhost:4561"
var direction = int64(0)
var procedureCode = int64(27)
var typeOfMessage = int64(1)

var reportParams = apimodel.ReportParams{
	Meid: meid,
	RANFunctionID: &funId,
	ClientEndpoint: &clientEndpoint,
	EventTriggers: apimodel.EventTriggerList{
		&apimodel.EventTrigger{
			InterfaceDirection: direction,
			ProcedureCode:      procedureCode,
			TypeOfMessage:      typeOfMessage,
		},
	},
}

var policyParams = apimodel.PolicyParams{
	Meid: &meid,
	RANFunctionID: &funId,
	ClientEndpoint: &clientEndpoint,
	EventTriggers: apimodel.EventTriggerList{
		&apimodel.EventTrigger{
			InterfaceDirection: direction,
			ProcedureCode:      procedureCode,
			TypeOfMessage:      typeOfMessage,
		},
	},
	PolicyActionDefinitions: &apimodel.PolicyActionDefinition{},
}

func subscriptionHandler(stype models.SubscriptionType, params interface{}) (models.SubscriptionResponse, error) {
	switch stype {
	case models.SubscriptionTypeReport:
		p := params.(*models.ReportParams)
		assert.Equal(suite, meid, p.Meid)
		assert.Equal(suite, funId, *p.RANFunctionID)
		assert.Equal(suite, clientEndpoint, *p.ClientEndpoint)
		assert.Equal(suite, direction, p.EventTriggers[0].InterfaceDirection)
		assert.Equal(suite, procedureCode, p.EventTriggers[0].ProcedureCode)
		assert.Equal(suite, typeOfMessage, p.EventTriggers[0].TypeOfMessage)
	case models.SubscriptionTypePolicy:
		p := params.(*models.PolicyParams)
		assert.Equal(suite, clientEndpoint, *p.ClientEndpoint)
	}

	subId := "xapp-11"
	reqId := int64(11)
	instanceId := int64(22)
	return models.SubscriptionResponse{
		&models.SubscriptionResponseItem{SubscriptionID: &subId, RequestorID: &reqId, InstanceID: &instanceId}, 
		&models.SubscriptionResponseItem{SubscriptionID: &subId, RequestorID: &reqId, InstanceID: &instanceId},
	}, nil
}

func queryHandler() (models.SubscriptionList, error) {
	resp := models.SubscriptionList{
		&models.SubscriptionData{
			SubscriptionID: 11,
			Meid:           "Test-Gnb",
			Endpoint:       []string{"127.0.0.1:4056"},
		},
	}

	return resp, nil
}

func deleteHandler(ep string) error {
	assert.Equal(suite, clientEndpoint, ep)

	return nil
}

func TestSetup(t *testing.T) {
	suite = t

	// Start the server to simulate SubManager
	go Subscription.Listen(subscriptionHandler, queryHandler, deleteHandler)
	time.Sleep(time.Duration(2) * time.Second)
}

func TestSubscriptionQueryHandling(t *testing.T) {
	resp, err := Subscription.QuerySubscriptions()

	assert.Equal(t, err, nil)
	assert.Equal(t, resp[0].SubscriptionID, int64(11))
	assert.Equal(t, resp[0].Meid, "Test-Gnb")
	assert.Equal(t, resp[0].Endpoint, []string{"127.0.0.1:4056"})
}

func TestSubscriptionReportHandling(t *testing.T) {
	result, err := Subscription.SubscribeReport(&reportParams)

	assert.Equal(t, err, nil)
	assert.Equal(t, len(result), 2)
	assert.Equal(t, *result[0].RequestorID, int64(11))
	assert.Equal(t, *result[0].InstanceID, int64(22))
	assert.Equal(t, *result[1].RequestorID, int64(11))
	assert.Equal(t, *result[1].InstanceID, int64(22))
}

func TestSubscriptionPolicytHandling(t *testing.T) {
	result, err := Subscription.SubscribePolicy(&policyParams)

	assert.Equal(t, err, nil)
	assert.Equal(t, len(result), 2)
	assert.Equal(t, *result[0].RequestorID, int64(11))
	assert.Equal(t, *result[0].InstanceID, int64(22))
}

func TestSubscriptionDeleteHandling(t *testing.T) {
	err := Subscription.UnSubscribe(clientEndpoint)

	assert.Equal(t, err, nil)
}