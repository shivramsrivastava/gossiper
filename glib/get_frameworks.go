package glib

import (
	"encoding/json"
	"log"
	"net/http"

	"../common"
)

func GetListofFrameworks(MasterEP string) {

	resp, err := http.Get(fmt.Sprintf("http://%s/state-summary", MasterEP))

	if err != nil {

		log.Printf("Unable to reach the Master error = %v", err)
	}

}
