package glib

import (
	"encoding/json"
	"log"
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
	var msg GossipMsG
	err := json.Unmarshal(buf, &msg)
	if err != nil {
		log.Printf("Delegate NotifyMsg() unmarshall error %v", err)
		return
	}

	this_frmwrk, isvalid := AllFrameworks[msg.Name]
	if !isvalid {
		this_frmwrk = make(map[string]bool)
	}
	for _, n := range msg.FrameWorks {
		this_frmwrk[n] = false
	}
	FrmWrkLck.Lock()
	AllFrameworks[msg.Name] = this_frmwrk
	FrmWrkLck.Unlock()
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
