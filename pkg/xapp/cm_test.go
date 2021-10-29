/*
==================================================================================
  Copyright (c) 2021 AT&T Intellectual Property.
  Copyright (c) 2021 Nokia

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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCM(t *testing.T) {
	Logger.Info("TestNewCM")

	assert.NotNil(t, NewCM("myApp", notificationCB), "NewCM failed")
}

func TestStoreCM(t *testing.T) {
	Logger.Info("TestStoreCM")

	cm := NewCM("myApp", notificationCB)
	assert.NotNil(t, cm, "NewCM failed")

	assert.Nil(t, cm.Store(), "StoreCM failed")
}

func TestReadCM(t *testing.T) {
	Logger.Info("TestReadCM")

	cm := NewCM("myApp", notificationCB)
	assert.NotNil(t, cm, "NewCM failed")

	cmData := cm.Read()
	assert.NotNil(t, cmData, "ReadCM failed")
}

func TestUpdateCM(t *testing.T) {
	Logger.Info("TestUpdateCM")

	cm := NewCM("myApp", notificationCB)
	assert.NotNil(t, cm, "NewCM failed")

	var jsonStr = `{"logger": {"level": 4,"noFormat": true}}`
	assert.Nil(t, cm.Update(jsonStr, false), "UpdateCM failed")

	cmData := cm.Read()
	assert.NotNil(t, cmData, "ReadCM failed")
}

func TestUpdateWithPublishCM(t *testing.T) {
	Logger.Info("TestUpdateCM")

	cm := NewCM("myApp", notificationCB)
	assert.NotNil(t, cm, "NewCM failed")

	var jsonStr = `{"logger": {"level": 4,"noFormat": true}}`
	assert.Nil(t, cm.Update(jsonStr, true), "UpdateCM failed")

	cmData := cm.Read()
	assert.NotNil(t, cmData, "ReadCM failed")
}

func notificationCB(ch string, events ...string) {
	Logger.Info("TESTCM notification received: channel=%+v events=%+v", ch, events[0])
}
