package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"./anonlib"
	"./common"
	"./consulLib"
	"./glib"
	"./httplib"
)

type GossiperConfig struct {
	Name           string //Name of this gossiepr
	MasterEndPoint string //MesosMaster's IP address
	ConfigType     string //What type of gossiper is this ? Join a federation or start a federation? Values = JOIN or NEW
	JoinEndPoint   string //If we are joining an already runnig federatoin the what is the EndPoint?
	LogFile        string //Name of the logfile
	HTTPPort       string //Defaults to 8080 if otherwise specify explicitly
	TCPPort        string //TCP port at which gossiper will bind and listen for anon module to connect to
	GPort          int    //Port at which gossper should start
	ConsulConfig   common.ConsulConfig
}

func NewGossiperConfig() GossiperConfig {
	return GossiperConfig{
		MasterEndPoint: ":5050",
		ConfigType:     "NEW",
		JoinEndPoint:   "",
		LogFile:        "stderr",
		HTTPPort:       "8080",
		TCPPort:        "5555",
		GPort:          4444,
	}
}

//This is suppose to be a simple stub which call the go routines of other mouels and simply wait
//TODO: Evaluate something like https://github.com/tedsuo/ifrit for managing our goroutines
//  We will have four go routines started from main
// 1) TCP Server that talks to the Anonymous Moduels in Apache Mesos
// 2) HTTP Server that will expose few REST(json) endpoints to query gossiper
// 3) Member list goroutine that will provide the feedback about the federation.
// 4) A Goroutine that will communicate with Policy Engine

func ProcessConfFile(filename string, conf *GossiperConfig) {

	file_content, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatalf("Unable to read the config file %v", err)
	}

	err = json.Unmarshal(file_content, conf)

	if err != nil {
		log.Fatalf("unable to unmarshall the config file not a valid json err=%", err)
	}
}

func main() {

	log.Printf("The code just started")

	//Get the default Config populated just in case no config.json was supplied via comamnd line argument
	config := NewGossiperConfig()

	//Try to parse the config file
	conffile := flag.String("config", "./config.json", "Supply the location of MrRedis configuration file")
	flag.Parse()
	ProcessConfFile(*conffile, &config)

	//Start Anon TCP server module
	go anonlib.Run(config.MasterEndPoint, config.TCPPort)

	//start http server
	go httplib.Run(config.HTTPPort)

	//start gossiper module
	var isnew bool
	var others string
	if config.ConfigType == "New" {
		isnew = true
		others = fmt.Sprintf("localhost:%d", config.GPort)
	} else {
		isnew = false
		others = config.JoinEndPoint
	}

	go glib.Run(config.Name, config.GPort, isnew, []string{others}, config.MasterEndPoint)

	val, _ := json.Marshal(&common.ALLDCs)

	log.Println("Marshalled output:", string(val))

	//start mesos master poller
	//go mesoslib.Run()

	//start consul client
	go consulLib.Run(&config.ConsulConfig, config.Name)

	//Start the Policy Engine module
	//PE.Run()

	/*
		x := 1
		for {
			fmt.Println(string(anonlib.GetNFrameIds(x)))
			time.Sleep(3 * time.Second)
			x++
		}*/

	//wait for ever
	wait := make(chan struct{})
	<-wait
}
