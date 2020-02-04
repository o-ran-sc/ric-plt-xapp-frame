/*
==================================================================================
  Copyright (c) 2019 Nokia
==================================================================================
*/

package xapp

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/models"
    apimodel "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/clientmodel"
)

var suite *testing.T

var requestorId = int64(0x4EEC)
var direction = int64(0)
var procedureCode = int64(27)
var typeOfMessage = int64(1)

var reportParams = apimodel.ReportParams{
	RequestorID: &requestorId,
	Interfaces: apimodel.EventTriggerList{
		&apimodel.EventTrigger{
			InterfaceDirection: &direction,
			ProcedureCode: &procedureCode,				
			TypeOfMessage: &typeOfMessage,
		},
	},
}

func subscriptionHandler(stype models.SubscriptionType, params interface{}) (models.SubscriptionResult, error) {
	switch stype {
	case models.SubscriptionTypeReport:
		p := params.(*models.ReportParams)
		assert.Equal(suite, requestorId, *p.RequestorID)
	case models.SubscriptionTypeControl:
	case models.SubscriptionTypePolicy:
	}
	
	return models.SubscriptionResult{11, 22, 33}, nil
}

func TestSubscriptionHandling(t *testing.T) {
	suite = t

	// Start the server to simulate SubManager
	go Subscription.Listen(subscriptionHandler)
	time.Sleep(time.Duration(2) * time.Second)

	// Subscribe X2AP events (action type: report)
	result := Subscription.SubscribeReport(&reportParams)

	assert.Equal(t, len(result), 3, "Should be equal!")
	assert.Equal(t, result[0], int64(11), "Should be equal!")
	assert.Equal(t, result[1], int64(22), "Should be equal!")
	assert.Equal(t, result[2], int64(33), "Should be equal!")
}