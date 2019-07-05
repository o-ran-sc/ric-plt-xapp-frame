//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package writer

/*
ISdlInstance integrates (wraps) the functionality that sdlgo library provides
*/
type ISdlInstance interface {
	SubscribeChannel(cb func(string, ...string), channels ...string) error
	UnsubscribeChannel(channels ...string) error
	Close() error
	SetAndPublish(channelsAndEvents []string, pairs ...interface{}) error
	Set(pairs ...interface{}) error
	Get(keys []string) (map[string]interface{}, error)
	SetIfAndPublish(channelsAndEvents []string, key string, oldData, newData interface{}) (bool, error)
	SetIf(key string, oldData, newData interface{}) (bool, error)
	SetIfNotExistsAndPublish(channelsAndEvents []string, key string, data interface{}) (bool, error)
	SetIfNotExists(key string, data interface{}) (bool, error)
	RemoveAndPublish(channelsAndEvents []string, keys []string) error
	Remove(keys []string) error
	RemoveIfAndPublish(channelsAndEvents []string, key string, data interface{}) (bool, error)
	RemoveIf(key string, data interface{}) (bool, error)
	GetAll() ([]string, error)
	RemoveAll() error
	RemoveAllAndPublish(channelsAndEvents []string) error
	AddMember(group string, member ...interface{}) error
	RemoveMember(group string, member ...interface{}) error
	RemoveGroup(group string) error
	GetMembers(group string) ([]string, error)
	IsMember(group string, member interface{}) (bool, error)
	GroupSize(group string) (int64, error)
}
