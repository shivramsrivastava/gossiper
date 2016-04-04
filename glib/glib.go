package glib

import (
	"fmt"
	"log"

	ml "github.com/hashicorp/memberlist"
)

type ML_interface interface {
	Create(conf *ml.Config) (*ml.Memberlist, error)
	Join(existing []string) (int, error)
	Leave(timeout time.Duration) error
	Members() []*ml.Node
}

type Gossiplib struct {
	list       *ml.Memberlist //Main gossiper library
	Configtype string         //Type of the configfile default or LAN etc.,
	config     *ml.Config     //They type of config that has been used for this gossiper
	Name       string         //Name of this cluster or gosiper
	Zone       string         //Zone an optional value
	M          int            //Number of members those are in the federation
	Msg        []byes         //Message or Payload that will be passed around among gossipers
	New        bool           //Is this a new cluster or joining an existing cluster
	ToJoin     []string       //List of cluster this gossiper can join

}

func NewGossiplib() *Gossiplib {
}

func (G *Gossiplib) Init() {
}

func (G *Gossiplib) Join() {
}

func (G *Gossiplib) Leave() {
}

func (G *Gossiplib) SendMsg() {
}

func (G *Gossiplib) RecvMsg() {
}
