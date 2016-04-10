package glib

import (
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
	log.Printf("Delegate NotifyMsg() is called %s", string(buf))
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
