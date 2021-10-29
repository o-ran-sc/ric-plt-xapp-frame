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
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
)

type CM struct {
	appName string
	cmData  map[string]interface{}
	cmCB    CMNotificationCB
}

const (
	CM_DATA_KEY string = "controls"
)

type CMNotificationCB func(string, ...string)

// Creates a CM instance
func NewCM(name string, cb CMNotificationCB) *CM {
	c := &CM{appName: name, cmCB: cb}

	if err := viper.UnmarshalKey(CM_DATA_KEY, &c.cmData); err != nil {
		Logger.Warn("Unable to read initial CM controls, %v", err)
		return nil
	}

	return c
}

// Reads CM data from database
func (c *CM) Read() map[string]interface{} {
	m, err := SdlStorage.Read(c.getCMNs(), CM_DATA_KEY)
	if err != nil {
		Logger.Error("Unable to read CM data: %v", err)
		return nil
	}
	Logger.Info("CM data: %+v", m[CM_DATA_KEY])

	return m
}

// Stores/persists initial CM data to database
func (c *CM) Store() error {
	err := SdlStorage.Store(c.getCMNs(), CM_DATA_KEY, fmt.Sprint(c.cmData))
	if err != nil {
		Logger.Error("Unable to store CM data: %v", err)
		return err
	}

	if err := SdlStorage.Subscribe(c.getCMNs(), c.cmCB, c.getCMChannel()); err != nil {
		Logger.Warn("Unable to subscribe for CM changes, %v", err)
	}

	return err
}

// Updates CM data in the database
func (c *CM) Update(cmJsonStr string, publish bool) error {
	cmJsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(cmJsonStr), &cmJsonMap)
	if err != nil {
		Logger.Error("Unable to unmarshal CM data: %v", err)
		return err
	}
	Logger.Info("cmJsonMap: %+v", cmJsonMap)

	ns := c.getCMNs()
	event := fmt.Sprint(cmJsonMap)
	if !publish {
		err = SdlStorage.Store(ns, CM_DATA_KEY, event)
		if err != nil {
			Logger.Error("Sdl.Store failed: %v", err)
			return err
		}
		return err
	}

	if err := SdlStorage.StoreAndPublish(ns, c.getCMChannel(), event, CM_DATA_KEY, event); err != nil {
		Logger.Error("Sdl.Store failed: %v", err)
		return err
	}

	return err
}

// Sets the CM notification callback
func (c *CM) SetCB(cb CMNotificationCB) {
	c.cmCB = cb
}

// Builds CM namaspace
func (c *CM) getCMNs() string {
	return fmt.Sprintf("O1:cmData/%s", c.appName)
}

// Builds CM notification channel
func (c *CM) getCMChannel() string {
	return fmt.Sprintf("O1:cmNotif/%s", c.appName)
}
