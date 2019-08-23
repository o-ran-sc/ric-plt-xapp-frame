package main

import (
	"gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
)

type Forwarder struct {
}

func (m Forwarder) Consume(params *xapp.RMRParams) (err error) {
	xapp.Logger.Debug("Message received - type=%d txid=%s subId=%d meid=%s",
		params.Mtype, params.Xid, params.SubId, params.Meid.RanName)

	// Store data and reply with the same message payload
	if xapp.Config.GetInt("test.store") != 0 {
		xapp.Sdl.Store("myKey", params.Payload)
	}

	mid := xapp.Config.GetInt("test.mtype")
	if mid != 0 {
		params.Mtype = mid
	} else {
		params.Mtype = params.Mtype + 1
	}

	sid := xapp.Config.GetInt("test.subId")
	if sid != 0 {
		params.SubId = sid
	}

	if ok := xapp.Rmr.SendRts(params); !ok {
		xapp.Logger.Info("Rmr.Send failed ...")
	}
	return
}

func forwarder() {
	xapp.Run(Forwarder{})
}
