package consulLib

import (
	"log"
	"runtime"
	"strings"
	"sync"

	"time"

	"github.com/hashicorp/consul/api"
)

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

func NewglobalDCMap() *globalDCMap {
	return &globalDCMap{Mutex: new(sync.Mutex),
		DCClientConnection: make(map[string]*FederaConsulClient),
		AvalibleDCInfo:     make(map[string]*api.AgentMember),
	}
}

//this function to create the KV client object for the any DC's address passed.
func NewFederaConsulClient(DcEndPoint string, DCName string, isLeader bool) (*FederaConsulClient, bool) {

	var dcInfoMap *globalDCMap

	log.Println("[INFO]NewFederaConsulClient: called")

	defaultConfig := api.DefaultConfig()
	if isLeader == true {
		//populate the DC maps
		dcInfoMap = NewglobalDCMap()
	}
	defaultConfig.Datacenter = DCName
	defaultConfig.Address = DcEndPoint

	log.Println("[INFO]NewFederaConsulClient:", defaultConfig)

	dcClient, err := api.NewClient(defaultConfig)

	if err != nil {
		log.Println("GetDefaultFederaConsulClient cannot connect to the default ", err)
		return nil, false
	}

	return &FederaConsulClient{Client: dcClient, Name: DCName, IsLeader: isLeader, DClist: dcInfoMap, KV: dcClient.KV()}, true
}

func (this *FederaConsulClient) CheckAndUpdateDCInfo(newAgentList []*api.AgentMember, localDcName string) {

	//newAgentList returns a pointer so it can be assigned?
	for index, dcAgents := range newAgentList {
		if strings.HasSuffix(dcAgents.Name, localDcName) != true {
			if _, ok := this.DClist.AvalibleDCInfo[dcAgents.Name]; !ok {
				//memebrs not found
				//get the connection
				this.GetNewClientConn(newAgentList[index])
			}
		} else {
			log.Println("Skipping local dc store", dcAgents.Name, localDcName)
		}
	}
}

func (this *FederaConsulClient) GetNewClientConn(newAgentList *api.AgentMember) {

	log.Println("Data centre name ", newAgentList.Name)
	//need to parse the node name out of DC name
	newAgentList.Name = strings.Split(newAgentList.Name, ".")[1]
	log.Println("after split Data centre name ", newAgentList.Name)
	newDcClient, ok := NewFederaConsulClient(newAgentList.Addr+":"+"8500", newAgentList.Name, false)
	if !ok {
		log.Println("Error in getting client connection to ", newAgentList.Name, " Failed to update the DC map ")
		return
	}
	this.DClist.Lock()
	this.DClist.AvalibleDCInfo[newAgentList.Name] = newAgentList
	this.DClist.DCClientConnection[newAgentList.Name] = newDcClient
	this.DClist.Unlock()
}

//type NodeName string

//gloablDCmap will hold the list of all the dc excluding the one on which this replicate is running

//this should
func (this *FederaConsulClient) PopulatetheGlobalDCMap() bool {

	log.Println("[INFO] PopulatetheGlobalDCMap: Called")

	clientAgent := this.Agent()
	//get wan members from agent

	localDCInfo, err := clientAgent.Self()

	if err != nil {
		log.Println("Error retrieving the local Client info", err)
		return false
	}

	localDCName := localDCInfo["Config"]["Datacenter"].(string)
	log.Println("[INFO] Localagent name", localDCName)
	for {

		dcMembers, err := clientAgent.Members(true)
		if err != nil {
			log.Println("err list from members", err)
			return false
		}

		for range dcMembers {
			//log.Println("List from all the dc's", val)
		}

		log.Println("CheckAndUpdateDCInfo called")
		this.CheckAndUpdateDCInfo(dcMembers, localDCName)
		<-time.After(10 * time.Second)
		runtime.Gosched()
	}

	return true
}
