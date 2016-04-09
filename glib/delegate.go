package glib

import (
	"log"
)

type delegate struct {
}

func (d *delegate) NodeMeta(limit int) []byte {
	log.Printf("Delegate NodeMeta() is called")
	return []byte{}
}

func (d *delegate) NotifyMsg(buf []byte) {
	log.Printf("Delegate NotifyMsg() is called")
}

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	log.Printf("Delegate GetBroadcasts() is called")
	return [][]byte{}
}

func (d *delegate) LocalState(join bool) []byte {
	log.Printf("Delegate LocalState() is called")

	return []byte{}
}

func (d *delegate) MergeRemoteState(buf []byte, isJoin bool) {
	log.Printf("Delegate MergeRemoteState() is called")
}
