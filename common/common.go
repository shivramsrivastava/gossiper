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

type toanon struct {
	Ch  chan bool
	M   map[string]bool
	Lck sync.Mutex
}

//Declare somecommon types that will be used accorss the goroutines
var (
	ToAnon      toanon
	ALLDCs      alldcs
	ThisDCName  string
	ThisEP      string
	ThisCity    string
	ThisCountry string
)

func init() {

	ToAnon.M = make(map[string]bool)
	ToAnon.Ch = make(chan bool)
	ALLDCs.List = make(map[string]*DC)
	fmt.Printf("Initalizeing Common")

}

//global consul config
type ConsulConfig struct {
	IsLeader    bool
	DCEndpoint  string
	StorePreFix string
	DCName      string
}
