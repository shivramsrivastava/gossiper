package main

import (
	"flag"
	//	"fmt"
	//	"time"

	"./anonlib"
	"./httplib"
)

type GossiperConfig struct {
	MasterEndPoint string //MesosMaster's IP address
	ConfigType     string //What type of gossiper is this ? Join a federation or start a federation? Values = JOIN or NEW
	JoinEndPoint   string //If we are joining an already runnig federatoin the what is the EndPoint?
	LogFile        string //Name of the logfile
	HTTPPort       string //Defaults to 8080 if otherwise specify explicitly
	TCPPort        string //TCP port at which gossiper will bind and listen for anon module to connect to
}

func NewGossiperConfig() GossiperConfig {
	return GossiperConfig{
		MasterEndPoint: ":5050",
		ConfigType:     "NEW",
		JoinEndPoint:   "",
		LogFile:        "stderr",
		HTTPPort:       "8080",
		TCPPort:        "5555",
	}
}

//This is suppose to be a simple stub which call the go routines of other mouels and simply wait
//TODO: Evaluate something like https://github.com/tedsuo/ifrit for managing our goroutines
//  We will have four go routines started from main
// 1) TCP Server that talks to the Anonymous Moduels in Apache Mesos
// 2) HTTP Server that will expose few REST(json) endpoints to query gossiper
// 3) Member list goroutine that will provide the feedback about the federation.
// 4) A Goroutine that will communicate with Policy Engine

func main() {

	//Get the default Config populated just in case no config.json was supplied via comamnd line argument
	config := NewGossiperConfig()

	//Try to parse the config file
	_ = flag.String("config", "./config.json", "Supply the location of MrRedis configuration file")
	flag.Parse()

	//Start Anon TCP server module
	go anonlib.Run("", config.TCPPort)

	//start http server
	go httplib.Run(config.HTTPPort)

	//start gossiper module
	//go glib.Run()

	//start mesos master poller
	//go mesoslib.Run()

	//Start the Policy Engine module
	//PE.Run()

	/*
		x := 1
		for {
			fmt.Println(string(anonlib.GetNFrameIds(x)))
			time.Sleep(3 * time.Second)
			x++
		}*/

}
