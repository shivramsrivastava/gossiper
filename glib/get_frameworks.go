package glib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//	"sync"
	"time"

	"../common"
)

func GetListofFrameworks(MasterEP string) {

	resp, err := http.Get(fmt.Sprintf("http://%s/state-summary", MasterEP))
	defer resp.Body.Close()

	if err != nil {

		log.Printf("Unable to reach the Master error = %v", err)
		return
	}

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

	var isNewFrwrk bool

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
		common.ToAnon.Lck.Lock()
		_, exisit := common.ToAnon.M[id]
		if !exisit {
			log.Printf("framework ID %v is new", id)
			common.ToAnon.M[id] = true
			if isNewFrwrk == false {
				isNewFrwrk = true
			}

		}
		common.ToAnon.Lck.Unlock()

	}

	if isNewFrwrk {
		common.ToAnon.Ch <- true
	}

}

func CollectMasterData(MasterEP string) {

	for {
		select {
		case <-time.After(time.Second * 5):
			GetListofFrameworks(MasterEP)
		}
	}
}
