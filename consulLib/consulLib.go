package consulLib

import (
	"log"

	"../common"
)

/*func init() {

}*/

//InitConsul should have the consulDCAddr with port
/*
func InitConsul(isLeader bool, consulDCAddr string) {

}*/

//
func Run(config *common.ConsulConfig, DCName string) {

	//addr, err := net.LookupHost("0.0.0.0")
	//log.Println(addr, err)

	newConsulClient, ok := NewFederaConsulClient(config.DCEndpoint, DCName, config.IsLeader)

	if !ok {
		log.Println("Cannot start the Consul Client")
		return
	}

	if config.IsLeader {

		//read from the KV store and send to other dc list
		//get the list of DC's
		//get the KV store with the configured prefix
		log.Println("[INFO]:Starting the consul replicate server")
		go newConsulClient.PopulatetheGlobalDCMap()
		go newConsulClient.PollAndUpdateKV(false)

	} else {

		log.Println("[INFO]:Starting the consul replicate Client")
		go newConsulClient.GetDataFromLocalKVStore()

	}
	log.Println("waiting")
	wait := make(chan bool)
	<-wait
}
