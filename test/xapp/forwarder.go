package main

import (
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
)

type Forwarder struct {
}

func (m Forwarder) Consume(mtype, subId int, payload []byte, meid *xapp.RMRMeid) (err error) {
	xapp.Logger.Debug("Message received - type=%d subId=%d meid=%v", mtype, subId, meid)

	// Store data and reply with the same message payload
	if xapp.Config.GetInt("test.store") != 0 {
		xapp.Sdl.Store("myKey", payload)
	}

	mid := xapp.Config.GetInt("test.mtype")
	if mid != 0 {
		mtype = mid
	} else {
		mtype = mtype + 1
	}

	sid := xapp.Config.GetInt("test.subId")
	if sid != 0 {
		subId = sid
	}

	if ok := xapp.Rmr.Send(mtype, subId, payload, nil); !ok {
		xapp.Logger.Info("Rmr.Send failed ...")
	}
	return
}

func forwarder() {
	xapp.Run(Forwarder{})
}
