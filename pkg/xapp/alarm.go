/*
==================================================================================
  Copyright (c) 2020 AT&T Intellectual Property.
  Copyright (c) 2020 Nokia

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
	"os"

	"gerrit.o-ran-sc.org/r/ric-plt/alarm-go.git/alarm"
)

type AlarmClient struct {
	alarmer *alarm.RICAlarm
}

// NewAlarmClient returns a new AlarmClient.
func NewAlarmClient(moId, appId string) *AlarmClient {
	if moId == "" {
		moId = "RIC"
		if hostname, err := os.Hostname(); err == nil {
			moId = hostname
		}
	}

	if appId == "" {
		appId = os.Args[0]
	}

	alarmInstance, err := alarm.InitAlarm(moId, appId)
	if err == nil {
		return &AlarmClient{
			alarmer: alarmInstance,
		}
	}
	return nil
}

func (c *AlarmClient) Raise(sp int, severity alarm.Severity, additionalInfo, identifyingInfo string) error {
	alarmData := c.alarmer.NewAlarm(sp, severity, additionalInfo, identifyingInfo)
	return c.alarmer.Raise(alarmData)
}

func (c *AlarmClient) Clear(sp int, severity alarm.Severity, additionalInfo, identifyingInfo string) error {
	alarmData := c.alarmer.NewAlarm(sp, severity, additionalInfo, identifyingInfo)
	return c.alarmer.Clear(alarmData)
}

func (c *AlarmClient) Reraise(sp int, severity alarm.Severity, additionalInfo, identifyingInfo string) error {
	alarmData := c.alarmer.NewAlarm(sp, severity, additionalInfo, identifyingInfo)
	return c.alarmer.Reraise(alarmData)
}

func (c *AlarmClient) ClearAll() error {
	return c.alarmer.ClearAll()
}
