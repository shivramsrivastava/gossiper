package main

import (
	"fmt"
	"sync"
	"time"

	"./anonlib"
	"./config"
)

func main() {

	var mylocagrp sync.WaitGroup

	fmt.Println(config.BindAddr, config.BindPort)

	mylocagrp.Add(2)

	x := 1
	for {
		fmt.Println(string(anonlib.GetNFrameIds(x)))
		time.Sleep(3 * time.Second)
		x++
	}

}
