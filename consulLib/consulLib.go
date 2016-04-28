package consulLib

import (
	"log"
	"sync"

	"github.com/hashicorp/consul/api"

	"../common"
	//"../policylib"
)

var (
	GlobalConsulDCChannel chan bool
	GlobalConsulDCInfo    *FederaConsulClient
	GlobalConsulMutex     *sync.Mutex
)

func init() {
	log.Println("initializing consulLib package")
	GlobalConsulDCChannel = make(chan bool)
	GlobalConsulMutex = new(sync.Mutex)

}

// global map struct for all the dc list
// this will be populated only for the leader
type globalDCMap struct {
	*sync.Mutex
	//here string is the node name for now.
	// TODO: if two nodes on a WAN have the same name?
	DCClientConnection map[string]*FederaConsulClient
	AvalibleDCInfo     map[string]*api.AgentMember
}

//consul client connection
type FederaConsulClient struct {
	*api.Client
	*api.KV
	Name     string //this will have the dc name
	IsLeader bool   //if this client is the leader
	DClist   *globalDCMap
}

type KVData struct {
	api.KVPairs
	*api.QueryMeta
}

const Prefix = "Fedra"

func Run(config *common.ConsulConfig, DCName string) {

	//addr, err := net.LookupHost("0.0.0.0")
	//log.Println(addr, err)

	newConsulClient, ok := NewFederaConsulClient(config.DCEndpoint, DCName, config.IsLeader)
	newConsulClient.UpdateAndSignalGlobalConsulInfo()

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
		go newConsulClient.PollAndUpdateKV(true)
	} else {
		// We dont need this else part.
		// Since all the gossiper need to read the local store for the policy update.
		// Leader dosent one thing in addition.

		log.Println("[INFO]:Starting the consul replicate Client")
		//go newConsulClient.GetDataFromLocalKVStore()
	}
	//We pass the reference to the local client KV store
	//go policylib.Run(Prefix, newConsulClient)
	log.Println("waiting")
	wait := make(chan bool)
	<-wait
}
