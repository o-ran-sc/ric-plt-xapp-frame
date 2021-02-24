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
	"github.com/stretchr/testify/assert"
	"testing"

	"gerrit.o-ran-sc.org/r/ric-plt/alarm-go.git/alarm"
)

func TestNewAlarmClient(t *testing.T) {
	Logger.Info("TestNewAlarmClient")

	a := NewAlarmClient("", "")
	assert.NotNil(t, a, "NewAlarmClient failed")
}

func TestAlarmRaise(t *testing.T) {
	Logger.Info("TestAlarmRaise")

	a := NewAlarmClient("", "")
	assert.NotNil(t, a, "NewAlarmClient failed")

	a.Raise(1234, alarm.SeverityCritical, "Some App data", "eth 0 1")
}

func TestAlarmClear(t *testing.T) {
	Logger.Info("TestAlarmClear")

	a := NewAlarmClient("", "")
	assert.NotNil(t, a, "NewAlarmClient failed")

	a.Clear(1234, alarm.SeverityCritical, "Some App data", "eth 0 1")
}

func TestAlarmReraise(t *testing.T) {
	Logger.Info("TestAlarmReraise")

	a := NewAlarmClient("", "")
	assert.NotNil(t, a, "NewAlarmClient failed")

	a.Reraise(1234, alarm.SeverityCritical, "Some App data", "eth 0 1")
}

func TestAlarmClearall(t *testing.T) {
	Logger.Info("TestAlarmClearall")

	a := NewAlarmClient("", "")
	assert.NotNil(t, a, "NewAlarmClient failed")

	a.ClearAll()
}
