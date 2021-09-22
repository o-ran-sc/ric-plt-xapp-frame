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
	rnibcommon "gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	rnibentities "gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	rnibreader "gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	sdl "gerrit.o-ran-sc.org/r/ric-plt/sdlgo"
	rnibwriter "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/rnib"
	"sync"
	"time"
)

// To be removed later
type SDLStatistics struct{}

var SDLCounterOpts = []CounterOpts{
	{Name: "Stored", Help: "The total number of stored SDL transactions"},
	{Name: "StoreError", Help: "The total number of SDL store errors"},
}

type SDLStorage struct {
	db    *sdl.SyncStorage
	stat  map[string]Counter
	mux   sync.Mutex
	ready bool
}

//Deprecated: Will be removed in a future release, please use SDLStorage type
type SDLClient struct {
	db        *SDLStorage
	nameSpace string
}

// Alias
type RNIBNodeType = rnibentities.Node_Type
type RNIBGlobalNbId = rnibentities.GlobalNbId
type RNIBNodebInfo = rnibentities.NodebInfo
type RNIBIRNibError = error
type RNIBCells = rnibentities.Cells
type RNIBNbIdentity = rnibentities.NbIdentity
type RNIBCellType = rnibentities.Cell_Type
type RNIBCell = rnibentities.Cell
type RNIBEnb = rnibentities.Enb
type RNIBGnb = rnibentities.Gnb

const RNIBNodeENB = rnibentities.Node_ENB
const RNIBNodeGNB = rnibentities.Node_GNB

type RNIBServedCellInfo = rnibentities.ServedCellInfo
type RNIBNodebInfoEnb = rnibentities.NodebInfo_Enb
type RNIBNodebInfoGnb = rnibentities.NodebInfo_Gnb
type RNIBServedNRCell = rnibentities.ServedNRCell
type RNIBServedNRCellInformation = rnibentities.ServedNRCellInformation
type RNIBNrNeighbourInformation = rnibentities.NrNeighbourInformation

type RNIBClient struct {
	db     rnibcommon.ISdlSyncStorage
	reader rnibreader.RNibReader
	writer rnibwriter.RNibWriter
}

// NewSdlStorage returns a new instance of SDLStorage type.
func NewSdlStorage() *SDLStorage {
	return &SDLStorage{
		db:    sdl.NewSyncStorage(),
		stat:  Metric.RegisterCounterGroup(SDLCounterOpts, "SDL"),
		ready: false,
	}
}

func (s *SDLStorage) TestConnection(namespace string) {
	// Test DB connection, and wait until ready!
	for {
		if _, err := s.db.GetAll(namespace); err == nil {
			break
		}
		Logger.Warn("Database connection not ready, waiting ...")
		time.Sleep(time.Duration(5 * time.Second))
	}
	s.ready = true
	Logger.Info("Connection to database established!")
}

func (s *SDLStorage) IsReady() bool {
	return s.ready
}

func (s *SDLStorage) doSet(namespace string, pairs ...interface{}) (err error) {
	err = s.db.Set(namespace, pairs)
	if err != nil {
		s.UpdateStatCounter("StoreError")
	} else {
		s.UpdateStatCounter("Stored")
	}
	return
}

func (s *SDLStorage) Store(namespace string, key string, value interface{}) (err error) {
	return s.doSet(namespace, key, value)
}

func (s *SDLStorage) MStore(namespace string, pairs ...interface{}) (err error) {
	return s.doSet(namespace, pairs)
}

func (s *SDLStorage) Read(namespace string, key string) (value map[string]interface{}, err error) {
	return s.db.Get(namespace, []string{key})
}

func (s *SDLStorage) MRead(namespace string, key []string) (value map[string]interface{}, err error) {
	return s.db.Get(namespace, key)
}

func (s *SDLStorage) ReadAllKeys(namespace string) (value []string, err error) {
	return s.db.GetAll(namespace)
}

func (s *SDLStorage) Subscribe(namespace string, cb func(string, ...string), channel string) error {
	return s.db.SubscribeChannel(namespace, cb, channel)
}

func (s *SDLStorage) MSubscribe(namespace string, cb func(string, ...string), channels ...string) error {
	return s.db.SubscribeChannel(namespace, cb, channels...)
}

func (s *SDLStorage) StoreAndPublish(namespace string, channel string, event string, pairs ...interface{}) error {
	return s.db.SetAndPublish(namespace, []string{channel, event}, pairs...)
}

func (s *SDLStorage) MStoreAndPublish(namespace string, channelsAndEvents []string, pairs ...interface{}) error {
	return s.db.SetAndPublish(namespace, channelsAndEvents, pairs...)
}

func (s *SDLStorage) Delete(namespace string, keys []string) (err error) {
	return s.db.Remove(namespace, keys)
}

func (s *SDLStorage) Clear(namespace string) {
	s.db.RemoveAll(namespace)
}

func (s *SDLStorage) RegisterMetrics() {
	s.stat = Metric.RegisterCounterGroup(SDLCounterOpts, "SDL")
}

func (s *SDLStorage) UpdateStatCounter(name string) {
	s.mux.Lock()
	s.stat[name].Inc()
	s.mux.Unlock()
}

func (s *SDLStorage) GetStat() (t SDLStatistics) {
	return
}

//NewSDLClient returns a new SDLClient.
//Deprecated: Will be removed in a future release, please use NewSdlStorage
func NewSDLClient(ns string) *SDLClient {
	if ns == "" {
		ns = "sdl"
	}
	return &SDLClient{
		db:        NewSdlStorage(),
		nameSpace: ns,
	}
}

//Deprecated: Will be removed in a future release, please use the TestConnection Receiver function of the SDLStorage type.
func (s *SDLClient) TestConnection() {
	s.db.TestConnection(s.nameSpace)
}

func (s *SDLClient) IsReady() bool {
	return s.db.ready
}

//Deprecated: Will be removed in a future release, please use the Store Receiver function of the SDLStorage type.
func (s *SDLClient) Store(key string, value interface{}) (err error) {
	return s.db.Store(s.nameSpace, key, value)
}

//Deprecated: Will be removed in a future release, please use the MStore Receiver function of the SDLStorage type.
func (s *SDLClient) MStore(pairs ...interface{}) (err error) {
	return s.db.MStore(s.nameSpace, pairs)
}

//Deprecated: Will be removed in a future release, please use the Read Receiver function of the SDLStorage type.
func (s *SDLClient) Read(key string) (value map[string]interface{}, err error) {
	return s.db.Read(s.nameSpace, key)
}

//Deprecated: Will be removed in a future release, please use the MRead Receiver function of the SDLStorage type.
func (s *SDLClient) MRead(key []string) (value map[string]interface{}, err error) {
	return s.db.MRead(s.nameSpace, key)
}

//Deprecated: Will be removed in a future release, please use the ReadAllKeys Receiver function of the SDLStorage type.
func (s *SDLClient) ReadAllKeys(key string) (value []string, err error) {
	return s.db.ReadAllKeys(s.nameSpace)
}

//Deprecated: Will be removed in a future release, please use the Subscribe Receiver function of the SDLStorage type.
func (s *SDLClient) Subscribe(cb func(string, ...string), channel string) error {
	return s.db.Subscribe(s.nameSpace, cb, channel)
}

//Deprecated: Will be removed in a future release, please use the MSubscribe Receiver function of the SDLStorage type.
func (s *SDLClient) MSubscribe(cb func(string, ...string), channels ...string) error {
	return s.db.MSubscribe(s.nameSpace, cb, channels...)
}

//Deprecated: Will be removed in a future release, please use the StoreAndPublish Receiver function of the SDLStorage type.
func (s *SDLClient) StoreAndPublish(channel string, event string, pairs ...interface{}) error {
	return s.db.StoreAndPublish(s.nameSpace, channel, event, pairs...)
}

//Deprecated: Will be removed in a future release, please use the MStoreAndPublish Receiver function of the SDLStorage type.
func (s *SDLClient) MStoreAndPublish(channelsAndEvents []string, pairs ...interface{}) error {
	return s.db.MStoreAndPublish(s.nameSpace, channelsAndEvents, pairs...)
}

//Deprecated: Will be removed in a future release, please use the Delete Receiver function of the SDLStorage type.
func (s *SDLClient) Delete(keys []string) (err error) {
	return s.db.Delete(s.nameSpace, keys)
}

//Deprecated: Will be removed in a future release, please use the Clear Receiver function of the SDLStorage type.
func (s *SDLClient) Clear() {
	s.db.Clear(s.nameSpace)
}

//Deprecated: Will be removed in a future release, please use the RegisterMetrics Receiver function of the SDLStorage type.
func (s *SDLClient) RegisterMetrics() {
	s.db.RegisterMetrics()
}

//Deprecated: Will be removed in a future release, please use the UpdateStatCounter Receiver function of the SDLStorage type.
func (s *SDLClient) UpdateStatCounter(name string) {
	s.db.UpdateStatCounter(name)
}

//Deprecated: Will be removed in a future release, please use the GetStat Receiver function of the SDLStorage type.
func (c *SDLClient) GetStat() (t SDLStatistics) {
	return c.db.GetStat()
}

func GetNewRnibClient(sdlStorage rnibcommon.ISdlSyncStorage) *RNIBClient {
	return &RNIBClient{
		db:     sdlStorage,
		reader: rnibreader.GetNewRNibReader(sdlStorage),
		writer: rnibwriter.GetNewRNibWriter(sdlStorage),
	}
}

//Deprecated: Will be removed in a future release, please use GetNewRnibClient instead.
func NewRNIBClient() *RNIBClient {
	s := sdl.NewSyncStorage()
	return &RNIBClient{
		db:     s,
		reader: rnibreader.GetNewRNibReader(s),
		writer: rnibwriter.GetNewRNibWriter(s),
	}
}

func (r *RNIBClient) Subscribe(cb func(string, ...string), channel string) error {
	return r.db.SubscribeChannel(rnibcommon.GetRNibNamespace(), cb, channel)
}

func (r *RNIBClient) StoreAndPublish(channel string, event string, pairs ...interface{}) error {
	return r.db.SetAndPublish(rnibcommon.GetRNibNamespace(), []string{channel, event}, pairs...)
}

func (r *RNIBClient) GetNodeb(invName string) (*RNIBNodebInfo, RNIBIRNibError) {
	return r.reader.GetNodeb(invName)
}

func (r *RNIBClient) GetNodebByGlobalNbId(t RNIBNodeType, gid *RNIBGlobalNbId) (*RNIBNodebInfo, RNIBIRNibError) {
	return r.reader.GetNodebByGlobalNbId(t, gid)
}

func (r *RNIBClient) GetCellList(invName string) (*RNIBCells, RNIBIRNibError) {
	return r.reader.GetCellList(invName)
}

func (r *RNIBClient) GetListGnbIds() ([]*RNIBNbIdentity, RNIBIRNibError) {
	return r.reader.GetListGnbIds()
}

func (r *RNIBClient) GetListEnbIds() ([]*RNIBNbIdentity, RNIBIRNibError) {
	return r.reader.GetListEnbIds()
}

func (r *RNIBClient) GetCountGnbList() (int, RNIBIRNibError) {
	return r.reader.GetCountGnbList()
}

func (r *RNIBClient) GetCell(invName string, pci uint32) (*RNIBCell, RNIBIRNibError) {
	return r.reader.GetCell(invName, pci)
}

func (r *RNIBClient) GetCellById(cellType RNIBCellType, cellId string) (*RNIBCell, RNIBIRNibError) {
	return r.reader.GetCellById(cellType, cellId)
}

func (r *RNIBClient) SaveNodeb(nbIdentity *RNIBNbIdentity, entity *RNIBNodebInfo) RNIBIRNibError {
	return r.writer.SaveNodeb(nbIdentity, entity)
}
