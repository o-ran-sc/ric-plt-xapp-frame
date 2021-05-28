/*
==================================================================================
  Copyright (c) 2019 Nokia
==================================================================================
*/

package xapp

import (
	"fmt"
	"testing"
	"time"

	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientmodel"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/models"
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
	subscriptionParams := clientmodel.SubscriptionParams{
		SubscriptionID: "",
		Meid:           &meid,
		RANFunctionID:  &funId,
		ClientEndpoint: &clientEndpoint,
		SubscriptionDetails: clientmodel.SubscriptionDetailsList{
			&clientmodel.SubscriptionDetail{
				XappEventInstanceID: &eventInstanceId,
				EventTriggers: &clientmodel.EventTriggerDefinition{
					OctetString: "1234",
				},
				ActionToBeSetupList: clientmodel.ActionsToBeSetup{
					&clientmodel.ActionToBeSetup{
						ActionID:   &actionId,
						ActionType: &actionType,
						ActionDefinition: &clientmodel.ActionDefinition{
							OctetString: "5678",
						},
						SubsequentAction: &clientmodel.SubsequentAction{
							SubsequentActionType: &subsequestActioType,
							TimeToWait:           &timeToWait,
						},
					},
				},
			},
		},
	}

	Subscription.SetResponseCB(func(resp *clientmodel.SubscriptionResponse) {
		assert.Equal(t, len(resp.SubscriptionInstances), 1)
		assert.Equal(t, *resp.SubscriptionInstances[0].XappEventInstanceID, int64(11))
		assert.Equal(t, *resp.SubscriptionInstances[0].E2EventInstanceID, int64(22))
	})

	_, err := Subscription.Subscribe(&subscriptionParams)
	assert.Equal(t, err, nil)
}

func TestSubscriptionWithClientProvidedIdHandling(t *testing.T) {
	subscriptionParams := clientmodel.SubscriptionParams{
		SubscriptionID: "myxapp",
		Meid:           &meid,
		RANFunctionID:  &funId,
		ClientEndpoint: &clientEndpoint,
		SubscriptionDetails: clientmodel.SubscriptionDetailsList{
			&clientmodel.SubscriptionDetail{
				XappEventInstanceID: &eventInstanceId,
				EventTriggers: &clientmodel.EventTriggerDefinition{
					OctetString: "1234",
				},
				ActionToBeSetupList: clientmodel.ActionsToBeSetup{
					&clientmodel.ActionToBeSetup{
						ActionID:   &actionId,
						ActionType: &actionType,
						ActionDefinition: &clientmodel.ActionDefinition{
							OctetString: "5678",
						},
						SubsequentAction: &clientmodel.SubsequentAction{
							SubsequentActionType: &subsequestActioType,
							TimeToWait:           &timeToWait,
						},
					},
				},
			},
		},
	}

	Subscription.SetResponseCB(func(resp *clientmodel.SubscriptionResponse) {
		assert.Equal(t, len(resp.SubscriptionInstances), 1)
		assert.Equal(t, *resp.SubscriptionInstances[0].XappEventInstanceID, int64(11))
		assert.Equal(t, *resp.SubscriptionInstances[0].E2EventInstanceID, int64(22))
	})

	_, err := Subscription.Subscribe(&subscriptionParams)
	assert.Equal(t, err, nil)
}

func TestSubscriptionDeleteHandling(t *testing.T) {
	err := Subscription.Unsubscribe(subscriptionId)
	fmt.Println(err)
	assert.Equal(t, err, nil)
}

// Helper functions
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

func subscriptionHandler(params interface{}) (*models.SubscriptionResponse, error) {
	p := params.(*models.SubscriptionParams)

	assert.Equal(suite, meid, *p.Meid)
	assert.Equal(suite, funId, *p.RANFunctionID)
	assert.Equal(suite, clientEndpoint.Host, p.ClientEndpoint.Host)
	assert.Equal(suite, clientEndpoint.HTTPPort, p.ClientEndpoint.HTTPPort)
	assert.Equal(suite, clientEndpoint.RMRPort, p.ClientEndpoint.RMRPort)

	assert.Equal(suite, xappEventInstanceId, *p.SubscriptionDetails[0].XappEventInstanceID)
	assert.Equal(suite, "1234", p.SubscriptionDetails[0].EventTriggers.OctetString)
	assert.Equal(suite, actionId, *p.SubscriptionDetails[0].ActionToBeSetupList[0].ActionID)
	assert.Equal(suite, actionType, *p.SubscriptionDetails[0].ActionToBeSetupList[0].ActionType)

	assert.Equal(suite, subsequestActioType, *p.SubscriptionDetails[0].ActionToBeSetupList[0].SubsequentAction.SubsequentActionType)
	assert.Equal(suite, timeToWait, *p.SubscriptionDetails[0].ActionToBeSetupList[0].SubsequentAction.TimeToWait)
	assert.Equal(suite, "5678", p.SubscriptionDetails[0].ActionToBeSetupList[0].ActionDefinition.OctetString)

	// Generate a unique subscriptionId
	subscriptionId = fmt.Sprintf("%s-%s", meid, clientEndpoint.Host)

	// Process subscriptions on the background
	go processSubscriptions(subscriptionId)

	// and send response immediately
	return &models.SubscriptionResponse{
		SubscriptionID: &subscriptionId,
	}, nil
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

func deleteHandler(ep string) error {
	assert.Equal(suite, subscriptionId, ep)
	return nil
}
