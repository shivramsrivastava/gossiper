package common

import (
	"fmt"
	"sync"
)

//Declare some modules that will eb common for both Anonymous and Gossiper modulesv
var (
	ToAnon toanon
)

type toanon struct {
	Ch  chan bool
	M   map[string]bool
	Lck sync.Mutex
}

//Declare somecommon types  that will be used accorss the goroutines

func init() {

	ToAnon.M = make(map[string]bool)
	ToAnon.Ch = make(chan bool)
	fmt.Printf("Initalizeing Common")
}
