package glib

import (
	"encoding/json"
	"log"

	//	"../common"
)

type delegate struct {
	glib               *Glib
	GetBroadcastCalled int
}

func (d *delegate) NodeMeta(limit int) []byte {
	log.Printf("Delegate NodeMeta() is called")
	return []byte{}
}

func (d *delegate) NotifyMsg(buf []byte) {
	//Gossipers will recive the message others and update the global map
	log.Printf("Delegate NotifyMsg() is called %s", string(buf))
	var msg Msg
	err := json.Unmarshal(buf, &msg)
	if err != nil {
		log.Printf("Delegate NotifyMsg() unmarshall error %v", err)
		return
	}

	switch msg.Type {
	case "FrameWorkMsG":
		var FW FrameWorkMsG
		msg.Body = &FW
		err := json.Unmarshal(buf, &msg)
		if err != nil {
			log.Printf("Delegate NotifyMsg() unmarshall FrameWorkMsG error %v", err)
			return
		}

		//First check if the Daacenter entry is available otherwise remove it
		this_frmwrk, isvalid := AllFrameworks[msg.Name]
		if !isvalid {
			this_frmwrk = make(map[string]bool)
		}

		//Loop through the frameworks
		for _, n := range FW.FrameWorks {
			this_frmwrk[n] = false
		}

		FrmWrkLck.Lock()
		AllFrameworks[msg.Name] = this_frmwrk
		FrmWrkLck.Unlock()

	case "OOR":
		msg.Body = &OutOfResourceMsG{}
		err := json.Unmarshal(buf, &msg)
		if err != nil {
			log.Printf("Delegate NotifyMsg() unmarshall OOR error %v", err)
			return
		}

	case "DC":
		msg.Body = &DCMsG{}
		err := json.Unmarshal(buf, &msg)
		if err != nil {
			log.Printf("Delegate NotifyMsg() unmarshall Datacenter Mesage error %v", err)
			return
		}
	}

}

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	if d.GetBroadcastCalled == 1000 {
		log.Printf("Delegate GetBroadcasts() is called 1000 times")
		d.GetBroadcastCalled = 0
	}
	d.GetBroadcastCalled++
	return d.glib.BC.GetBroadcasts(overhead, limit)
}

func (d *delegate) LocalState(join bool) []byte {
	log.Printf("Delegate LocalState() is called")

	return []byte{}
}

func (d *delegate) MergeRemoteState(buf []byte, isJoin bool) {
	log.Printf("Delegate MergeRemoteState() is called")
}
