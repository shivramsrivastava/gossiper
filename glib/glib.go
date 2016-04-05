package glib

import (
	//	"fmt"
	//	"log"
	"time"

	ml "github.com/hashicorp/memberlist"
)

type ML_interface interface {
	Create(*ml.Config) (*ml.Memberlist, error)
	Join([]string) (int, error)
	Leave(time.Duration) error
	Members() []*ml.Node
}

//Glib is short for GossipLibrary
//This is written on top of Hashicorps Memberlist library
//Hashi corps memberlist is well tested and is used in consula nd nomad for gossip (SWIM) among different notes in the cluster
//we want to use a such a well eshtablished library for Mesos Federation
type Glib struct {
	list       ML_interface //Main gossiper library
	Configtype string       //Type of the configfile default or LAN etc.,
	config     *ml.Config   //They type of config that has been used for this gossiper
	Name       string       //Name of this cluster or gosiper
	Zone       string       //Zone an optional value
	M          int          //Number of members those are in the federation
	Msg        []byte       //Message or Payload that will be passed around among gossipers
	New        bool         //Is this a new cluster or joining an existing cluster
	ToJoin     []string     //List of cluster this gossiper can join

}

func NewGlib(name string, zone string, new bool, ToJoin []string) *Glib {

	return &Glib{Name: name, Zone: zone, New: new}
}

func (G *Glib) Init() error {

	var err error

	G.list, err = ml.Create(ml.DefaultLocalConfig())

	if err != nil {
		return err
	}
	return nil
}

func (G *Glib) Join() {
}

func (G *Glib) Leave() {
}

func (G *Glib) SendMsg() {
}

func (G *Glib) RecvMsg() {
}
