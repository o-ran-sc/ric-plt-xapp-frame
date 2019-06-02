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
