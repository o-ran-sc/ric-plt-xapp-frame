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
package main

import "gitlabe1.ext.net.nokia.com/ric_dev/nokia-xapps/xapp/pkg/xapp"

type MessageCounter struct {
}

func (m MessageCounter) Consume(mtype, len int, payload []byte) (err error) {
	xapp.Logger.Debug("Message received - type=%d len=%d", mtype, len)

	xapp.Sdl.Store("myKey", payload)
	xapp.Rmr.Send(10005, len, payload)
	return nil
}

func main() {
	xapp.Run(MessageCounter{})
}
