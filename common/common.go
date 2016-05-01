package common

import (
	"fmt"
	"sync"
	"time"
)

//Declare some structure that will eb common for both Anonymous and Gossiper modulesv
type DC struct {
	OutOfResource bool
	Name          string
	City          string
	Country       string
	Endpoint      string
	CPU           float64
	MEM           float64
	DISK          float64
	Ucpu          float64 //Remaining CPU
	Umem          float64 //Remaining Memory
	Udisk         float64 //Remaining Disk
	LastUpdate    time.Duration
}

type alldcs struct {
	Lck  sync.Mutex
	List map[string]*DC
}

type rttbwGossipers struct {
	Lck  sync.Mutex
	List map[string]int64
}

type toanon struct {
	Ch  chan bool
	M   map[string]bool
	Lck sync.Mutex
}

//global consul config
type ConsulConfig struct {
	IsLeader    bool
	DCEndpoint  string
	StorePreFix string
	DCName      string
}

//Declare somecommon types that will be used accorss the goroutines
var (
	ToAnon             toanon    //Structure Sending messages to FedComms module via TCP client
	ALLDCs             alldcs    //The data structure that stores all the Datacenter information
	ThisDCName         string    //This DataCenter's Name
	ThisEP             string    //Thsi Datacenter's Endpoint
	ThisCity           string    //This Datacenters City
	ThisCountry        string    //This Datacentes Country
	ResourceThresold   int       //Threshold value of any resource (CPU, MEM or Disk) after which we need to broadcast OOR
	TriggerPolicyCh    chan bool //Polcy Engine will listen in this Channel
	RttOfPeerGossipers rttbwGossipers
)

func init() {

	ToAnon.M = make(map[string]bool)
	ToAnon.Ch = make(chan bool)
	TriggerPolicyCh = make(chan bool)
	ALLDCs.List = make(map[string]*DC)
	ResourceThresold = 100
	RttOfPeerGossipers.List = make(map[string]int64)
	fmt.Printf("Initalizeing Common")

}

func SupressFrameWorks() {
	ToAnon.Lck.Lock()
	for k, _ := range ToAnon.M {
		ToAnon.M[k] = false
	}
	ToAnon.Lck.Unlock()

	ToAnon.Ch <- true
}
func UnSupressFrameWorks() {
	ToAnon.Lck.Lock()
	for k, _ := range ToAnon.M {
		ToAnon.M[k] = true
	}
	ToAnon.Lck.Unlock()

	ToAnon.Ch <- true
}
