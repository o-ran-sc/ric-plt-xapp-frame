/*
==================================================================================
  Copyright (c) 2019 Nokia
==================================================================================
*/

package xapp

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientmodel"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/models"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/restapi/operations/common"
	"github.com/stretchr/testify/assert"
)

var (
	suite *testing.T

	meid                = "gnb123456"
	xappEventInstanceId = int64(1)
	eventInstanceId     = int64(1)
	funId               = int64(1)
	actionId            = int64(1)
	actionType          = "report"
	subsequestActioType = "continue"
	timeToWait          = "w10ms"
	direction           = int64(0)
	procedureCode       = int64(27)
	typeOfMessage       = int64(1)
	subscriptionId      = ""
	hPort               = int64(8080)
	rPort               = int64(4560)
	clientEndpoint      = clientmodel.SubscriptionParamsClientEndpoint{Host: "localhost", HTTPPort: &hPort, RMRPort: &rPort}
)

// Test cases
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
	assert.Equal(t, resp[0].ClientEndpoint, []string{"127.0.0.1:4056"})
}

func TestSubscriptionHandling(t *testing.T) {
	subscriptionParams := GetSubscriptionparams()

	Subscription.SetResponseCB(func(resp *clientmodel.SubscriptionResponse) {
		assert.Equal(t, len(resp.SubscriptionInstances), 1)
		assert.Equal(t, *resp.SubscriptionInstances[0].XappEventInstanceID, int64(11))
		assert.Equal(t, *resp.SubscriptionInstances[0].E2EventInstanceID, int64(22))
	})

	_, err := Subscription.Subscribe(subscriptionParams)
	assert.Equal(t, err, nil)
}

func TestSubscriptionWithClientProvidedIdHandling(t *testing.T) {
	subscriptionParams := GetSubscriptionparams()
	subscriptionParams.SubscriptionID = "myxapp"
    
	Subscription.SetResponseCB(func(resp *clientmodel.SubscriptionResponse) {
		assert.Equal(t, len(resp.SubscriptionInstances), 1)
		assert.Equal(t, *resp.SubscriptionInstances[0].XappEventInstanceID, int64(11))
		assert.Equal(t, *resp.SubscriptionInstances[0].E2EventInstanceID, int64(22))
	})

	_, err := Subscription.Subscribe(subscriptionParams)
	assert.Equal(t, err, nil)
}

func TestBadRequestSubscriptionHandling(t *testing.T) {
	subscriptionParams := GetSubscriptionparams()
	subscriptionParams.SubscriptionID = "123_send_bad_request_response"

	Subscription.SetResponseCB(func(resp *clientmodel.SubscriptionResponse) {
		assert.Equal(t, len(resp.SubscriptionInstances), 1)
		assert.Equal(t, *resp.SubscriptionInstances[0].XappEventInstanceID, int64(11))
		assert.Equal(t, *resp.SubscriptionInstances[0].E2EventInstanceID, int64(22))
	})

	_, err := Subscription.Subscribe(subscriptionParams)
	assert.Equal(t, err.Error(), "[POST /subscriptions][400] subscribeBadRequest ")
	fmt.Println("Error:",err)
}

func TestInternalServerErrorSubscriptionHandling(t *testing.T) {
	subscriptionParams := GetSubscriptionparams()
	subscriptionParams.SubscriptionID = "123_send_internal_server_error_response"

	Subscription.SetResponseCB(func(resp *clientmodel.SubscriptionResponse) {
		assert.Equal(t, len(resp.SubscriptionInstances), 1)
		assert.Equal(t, *resp.SubscriptionInstances[0].XappEventInstanceID, int64(11))
		assert.Equal(t, *resp.SubscriptionInstances[0].E2EventInstanceID, int64(22))
	})

	_, err := Subscription.Subscribe(subscriptionParams)
	assert.Equal(t, err.Error(), "[POST /subscriptions][500] subscribeInternalServerError ")
	fmt.Println("Error:",err)
}

func GetSubscriptionparams() *clientmodel.SubscriptionParams {
	return &clientmodel.SubscriptionParams{
		SubscriptionID: "",
		Meid:           &meid,
		RANFunctionID:  &funId,
		ClientEndpoint: &clientEndpoint,
		SubscriptionDetails: clientmodel.SubscriptionDetailsList{
			&clientmodel.SubscriptionDetail{
				XappEventInstanceID: &eventInstanceId,
				EventTriggers:       clientmodel.EventTriggerDefinition{00, 0x11, 0x12, 0x13, 0x00, 0x21, 0x22, 0x24, 0x1B, 0x80},
				ActionToBeSetupList: clientmodel.ActionsToBeSetup{
					&clientmodel.ActionToBeSetup{
						ActionID:         &actionId,
						ActionType:       &actionType,
						ActionDefinition: clientmodel.ActionDefinition{5, 6, 7, 8},
						SubsequentAction: &clientmodel.SubsequentAction{
							SubsequentActionType: &subsequestActioType,
							TimeToWait:           &timeToWait,
						},
					},
				},
			},
		},
	}
}

func TestSuccessfulSubscriptionDeleteHandling(t *testing.T) {
	subscriptionId = "123_send_successful_response"
	err := Subscription.Unsubscribe(subscriptionId)
	assert.Equal(t, err, nil)
	fmt.Println("Error:",err)
}

func TestBadRequestSubscriptionDeleteHandling(t *testing.T) {
	subscriptionId = "123_send_bad_request_response"
	err := Subscription.Unsubscribe(subscriptionId)
	assert.NotEqual(t, err, nil)
	fmt.Println("Error:",err.Error())
	assert.Equal(t, err.Error(), "[DELETE /subscriptions/{subscriptionId}][400] unsubscribeBadRequest ")
}

func TestInternalServerErrorSubscriptionDeleteHandling(t *testing.T) {
	subscriptionId = "123_send_internal_server_error_response"
	err := Subscription.Unsubscribe(subscriptionId)
	assert.NotEqual(t, err, nil)
	fmt.Println("Error:",err.Error())
	assert.Equal(t, err.Error(), "[DELETE /subscriptions/{subscriptionId}][500] unsubscribeInternalServerError ")
}

func TestResponseHandler(t *testing.T) {
	Subscription.SetResponseCB(SubscriptionRespHandler)

	payload := []byte(`{"SubscriptionInstances":[{"tXappEventInstanceID": 1}]`)
	req, err := http.NewRequest("POST", "/ric/v1/subscriptions/response", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Subscription.ResponseHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	time.Sleep(time.Duration(2) * time.Second)
}

// Helper functions
func SubscriptionRespHandler(resp *clientmodel.SubscriptionResponse) {
}

func processSubscriptions(subscriptionId string) {
	// Generate requestorId, instanceId
	xappInstanceId := int64(11)
	e2InstanceId := int64(22)

	resp := &models.SubscriptionResponse{
		SubscriptionID: &subscriptionId,
		SubscriptionInstances: []*models.SubscriptionInstance{
			{
				XappEventInstanceID: &xappInstanceId,
				E2EventInstanceID:   &e2InstanceId,
			},
		},
	}

	// Notify the client: don't worry about errors ... Notify() will handle retries, etc.
	Subscription.Notify(resp, models.SubscriptionParamsClientEndpoint{Host: "localhost", HTTPPort: &hPort, RMRPort: &rPort})
}

func subscriptionHandler(params interface{}) (*models.SubscriptionResponse, int) {
	p := params.(*models.SubscriptionParams)

	assert.Equal(suite, meid, *p.Meid)
	assert.Equal(suite, funId, *p.RANFunctionID)
	assert.Equal(suite, clientEndpoint.Host, p.ClientEndpoint.Host)
	assert.Equal(suite, clientEndpoint.HTTPPort, p.ClientEndpoint.HTTPPort)
	assert.Equal(suite, clientEndpoint.RMRPort, p.ClientEndpoint.RMRPort)

	assert.Equal(suite, xappEventInstanceId, *p.SubscriptionDetails[0].XappEventInstanceID)
	et := []int64{00, 0x11, 0x12, 0x13, 0x00, 0x21, 0x22, 0x24, 0x1B, 0x80}
	assert.ElementsMatch(suite, et, p.SubscriptionDetails[0].EventTriggers)
	assert.Equal(suite, actionId, *p.SubscriptionDetails[0].ActionToBeSetupList[0].ActionID)
	assert.Equal(suite, actionType, *p.SubscriptionDetails[0].ActionToBeSetupList[0].ActionType)

	assert.Equal(suite, subsequestActioType, *p.SubscriptionDetails[0].ActionToBeSetupList[0].SubsequentAction.SubsequentActionType)
	assert.Equal(suite, timeToWait, *p.SubscriptionDetails[0].ActionToBeSetupList[0].SubsequentAction.TimeToWait)
	assert.ElementsMatch(suite, []int64{5, 6, 7, 8}, p.SubscriptionDetails[0].ActionToBeSetupList[0].ActionDefinition)

	// Generate a unique subscriptionId
	subscriptionId = fmt.Sprintf("%s-%s", meid, clientEndpoint.Host)

	if p.SubscriptionID == "123_send_bad_request_response" {
		// Simulate bad request case
		return &models.SubscriptionResponse{
		}, common.UnsubscribeBadRequestCode

	}

	if p.SubscriptionID == "123_send_internal_server_error_response" {
		// Simulate bad internal server error case
		return &models.SubscriptionResponse{
			}, common.UnsubscribeInternalServerErrorCode
	}
	// Process subscriptions on the background
	go processSubscriptions(subscriptionId)

	// and send response immediately
	return &models.SubscriptionResponse{
		SubscriptionID: &subscriptionId,
	}, common.SubscribeCreatedCode
}

func queryHandler() (models.SubscriptionList, error) {
	resp := models.SubscriptionList{
		&models.SubscriptionData{
			SubscriptionID: 11,
			Meid:           "Test-Gnb",
			ClientEndpoint: []string{"127.0.0.1:4056"},
		},
	}
	return resp, nil
}

func deleteHandler(ep string) int {
	assert.Equal(suite, subscriptionId, ep)
	if subscriptionId == "123_send_successful_response" {
		return common.UnsubscribeNoContentCode
	} else if subscriptionId == "123_send_bad_request_response" {
		return common.UnsubscribeBadRequestCode
	} else {
		return common.UnsubscribeInternalServerErrorCode
	}
}
