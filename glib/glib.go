package glib

import (
	Q "container/list"
	"fmt"
	"log"
	"time"

	ml "github.com/hashicorp/memberlist"
)

type ML_interface interface {
	Join([]string) (int, error)
	Leave(time.Duration) error
	Members() []*ml.Node
}

//Glib is short for GossipLibrary
//This is written on top of Hashicorps Memberlist library
//Hashi corps memberlist is well tested and is used in consula nd nomad for gossip (SWIM) among different notes in the cluster
//we want to use a such a well eshtablished library for Mesos Federation
type Glib struct {
	list       ML_interface             //Main gossiper library
	Configtype string                   //Type of the configfile default or LAN etc.,
	config     *ml.Config               //They type of config that has been used for this gossiper
	BC         *ml.TransmitLimitedQueue //Broadcast queue
	Name       string                   //Name of this cluster or gosiper
	Zone       string                   //Zone an optional value
	BindPort   int                      //The port at which gossiper will bind to
	M          int                      //Number of members those are in the federation
	Msg        []byte                   //Message or Payload that will be passed around among gossipers
	New        bool                     //Is this a new cluster or joining an existing cluster
	ToJoin     []string                 //List of cluster this gossiper can join
	ToQ        *Q.List                  //Messages to be broadcasted to other gossipers
	FromQ      *Q.List                  //Messages recived from other gossipers

}

func NewGlib(name string, myport int, zone string, new bool, ToJoin []string) *Glib {

	return &Glib{Name: name, BindPort: myport, Zone: zone, New: new}
}

func (G *Glib) Init() error {

	var err error

	G.ToQ = Q.New()
	G.FromQ = Q.New()

	G.config = ml.DefaultLocalConfig()

	G.config.BindPort = G.BindPort
	G.config.Name = G.Name
	G.config.Delegate = &delegate{glib: G}

	G.list, err = ml.Create(G.config)

	G.BC = &ml.TransmitLimitedQueue{
		NumNodes: func() int {
			return len(G.list.Members())
		},
	}

	if err != nil {
		return err
	}
	return nil
}

func (G *Glib) Join(others []string) error {

	_, err := G.list.Join(others)

	if err != nil {

		log.Printf("Join Error %v", err)
		return err
	}

	return nil
}

func (G *Glib) Leave() error {

	err := G.list.Leave(0 * time.Second)

	if err != nil {
		log.Printf("Gossiper Leave Error %v", err)

	}
	return nil
}

func (G *Glib) SendMsg() {
	//Send Message to all the gossiper if a new framework is added to this master

}

func (G *Glib) RecvMsg() {

	//Wait on a channel and once you recive a message from others addit to the map
}

func Run(name string, myport int, isnew bool, others []string) {

	var wait chan struct{}

	wait = make(chan struct{})

	//Create and Initalize the gossiper libraray

	//NewGlib(name string, zone string, new bool, ToJoin []string) *Glib {
	g := NewGlib(name, myport, "", isnew, others)

	err := g.Init()

	if err != nil {
		log.Fatalf("Error Initializeing Gossiper LIbrary %v", err)
	}

	//Join with the other masters
	err = g.Join(others)

	if err != nil {

		log.Fatalf("Error unable to join other gossipers %v", err)
	}

	//Start the goroutine that will examine the recived q and process
	//go ExamineFramework()

	//wait
	g.BC.QueueBroadcast(NewBroadcast(fmt.Sprintf("Initial %s:%v", g.Name, time.Now())))
	for {

		select {

		case <-time.After(time.Second * 10):
			log.Printf("A Second elapsed for %s", g.Name)
			g.BC.QueueBroadcast(NewBroadcast(fmt.Sprintf("%s:%v", g.Name, time.Now())))

		case <-wait:
			log.Printf("End of glib library, terminating")
			break
		}
	}
}
