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
	sdl "gerrit.o-ran-sc.org/r/ric-plt/sdlgo"
	"sync"
	"time"
)

// To be removed later
type SDLStatistics struct{}

var SDLCounterOpts = []CounterOpts{
	{Name: "Stored", Help: "The total number of stored SDL transactions"},
	{Name: "StoreError", Help: "The total number of SDL store errors"},
}

type SDLClient struct {
	db    *sdl.SdlInstance
	stat  map[string]Counter
	mux   sync.Mutex
	ready bool
}

type RNIBClient struct {
	db *sdl.SdlInstance
}

// NewSDLClient returns a new SDLClient.
func NewSDLClient(ns string) *SDLClient {
	return &SDLClient{
		db:    sdl.NewSdlInstance(ns, sdl.NewDatabase()),
		stat:  Metric.RegisterCounterGroup(SDLCounterOpts, "SDL"),
		ready: false,
	}
}

func (s *SDLClient) TestConnection() {
	// Test DB connection, and wait until ready!
	for {
		if _, err := s.db.GetAll(); err == nil {
			break
		}
		Logger.Warn("Database connection not ready, waiting ...")
		time.Sleep(time.Duration(5 * time.Second))
	}
	s.ready = true
	Logger.Info("Connection to database established!")
}

func (s *SDLClient) IsReady() bool {
	return s.ready
}

func (s *SDLClient) doSet(pairs ...interface{}) (err error) {
	err = s.db.Set(pairs)
	if err != nil {
		s.UpdateStatCounter("StoreError")
	} else {
		s.UpdateStatCounter("Stored")
	}
	return
}

func (s *SDLClient) Store(key string, value interface{}) (err error) {
	return s.doSet(key, value)
}

func (s *SDLClient) MStore(pairs ...interface{}) (err error) {
	return s.doSet(pairs)
}

func (s *SDLClient) Read(key string) (value map[string]interface{}, err error) {
	return s.db.Get([]string{key})
}

func (s *SDLClient) MRead(key []string) (value map[string]interface{}, err error) {
	return s.db.Get(key)
}

func (s *SDLClient) ReadAllKeys(key string) (value []string, err error) {
	return s.db.GetAll()
}

func (s *SDLClient) Subscribe(cb func(string, ...string), channel string) error {
	return s.db.SubscribeChannel(cb, channel)
}

func (s *SDLClient) MSubscribe(cb func(string, ...string), channels ...string) error {
	return s.db.SubscribeChannel(cb, channels...)
}

func (s *SDLClient) StoreAndPublish(channel string, event string, pairs ...interface{}) error {
	return s.db.SetAndPublish([]string{channel, event}, pairs...)
}

func (s *SDLClient) MStoreAndPublish(channelsAndEvents []string, pairs ...interface{}) error {
	return s.db.SetAndPublish(channelsAndEvents, pairs...)
}

func (s *SDLClient) Clear() {
	s.db.RemoveAll()
}

func (s *SDLClient) RegisterMetrics() {
	s.stat = Metric.RegisterCounterGroup(SDLCounterOpts, "SDL")
}

func (s *SDLClient) UpdateStatCounter(name string) {
	s.mux.Lock()
	s.stat[name].Inc()
	s.mux.Unlock()
}

func (c *SDLClient) GetStat() (t SDLStatistics) {
	return
}

// To be removed ...
func NewRNIBClient(ns string) *RNIBClient {
	return &RNIBClient{
		db: sdl.NewSdlInstance(ns, sdl.NewDatabase()),
	}
}

func (r *RNIBClient) GetgNBList() (values map[string]interface{}, err error) {
	keys, err := r.db.GetAll()
	if err == nil {
		values = make(map[string]interface{})
		for _, key := range keys {
			v, err := r.db.Get([]string{key})
			if err == nil {
				values[key] = v[key]
			}
		}
	}
	return
}

func (r *RNIBClient) GetNRCellList(key string) (value map[string]interface{}, err error) {
	return r.db.Get([]string{key})
}

func (r *RNIBClient) GetUE(key1, key2 string) (value map[string]interface{}, err error) {
	return r.db.Get([]string{key1 + key2})
}

func (r *RNIBClient) Store(key string, value interface{}) (err error) {
	return r.db.Set(key, value)
}

func (r *RNIBClient) Clear() {
	r.db.RemoveAll()
}
