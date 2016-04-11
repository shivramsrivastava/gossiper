package glib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	//"../common"
)

type GossipMsG struct {
	Name       string
	FrameWorks []string
}

var (
	AllFrameworks   map[string]map[string]bool //Map of maps for all the frameworks
	CommonFramework map[string]bool            //Common among all the framework
	FrmWrkLck       sync.Mutex                 //Lock for centralized framework
)

func init() {

	AllFrameworks = make(map[string]map[string]bool)
	CommonFramework = make(map[string]bool)
}

func GetListofFrameworks(G *Glib, MasterEP string) {

	resp, err := http.Get(fmt.Sprintf("http://%s/state-summary", MasterEP))

	if err != nil {

		log.Printf("Unable to reach the Master error = %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Unable to read the body error = %v", err)
		return
	}

	var data map[string]interface{}

	err = json.Unmarshal(body, &data)

	if err != nil {
		log.Printf("Json Unmarshall error = %v", err)
		return
	}

	frameworks_array, found := data["frameworks"].([]interface{})

	if found != true {
		log.Printf("Unable to find the framework array in the nested json object %v", data)
		return
	}

	this_frmwrk, exisit := AllFrameworks[G.Name]
	FrmWrkLck.Lock()
	if !exisit {
		AllFrameworks[G.Name] = make(map[string]bool)
		this_frmwrk = AllFrameworks[G.Name]
	}
	FrmWrkLck.Unlock()

	for _, frwrk_interface := range frameworks_array {

		frwrk, isvalid := frwrk_interface.(map[string]interface{})
		if !isvalid {
			log.Printf("Malformed json object recived from he master %v", frameworks_array)
			continue
		}
		id, isvalid := frwrk["id"].(string)
		if isvalid != true {
			log.Printf("Unabel to get framework id from the json = %v", frwrk)
			continue
		}
		log.Printf("framework ID %v", id)
		//common.ToAnon.Lck.Lock()
		//_, exisit := common.ToAnon.M[id]

		//common.ToAnon.M[id] = false
		this_frmwrk[id] = false

		//common.ToAnon.Lck.Unlock()

	}
	FrmWrkLck.Lock()
	AllFrameworks[G.Name] = this_frmwrk
	FrmWrkLck.Unlock()
	G.BC.QueueBroadcast(NewBroadcast(GossipFrameworks(G.Name)))
}

func CollectMasterData(G *Glib, MasterEP string) {

	for {
		select {
		case <-time.After(time.Second * 5):
			GetListofFrameworks(G, MasterEP)
		}
	}
}

func GossipFrameworks(Name string) string {

	var msg GossipMsG

	msg.Name = Name

	gmap, isvalid := AllFrameworks[Name]
	if isvalid {
		for name := range gmap {
			msg.FrameWorks = append(msg.FrameWorks, name)
		}
	}

	data, _ := json.Marshal(msg)

	return string(data)

}
