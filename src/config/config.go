package config

import (
	"fmt"
	"os"
	"strconv"
)

type anonConfig struct {
	BindAddr string
	BindPort int
}

var (
	BindAddr string
	BindPort int
)

func init() {
	fmt.Println("config init called")
	BindAddr = os.Getenv("anonBindAddr")
	if len(BindAddr) == 0 {
		//set default bind addr
		BindAddr = "127.0.0.1"
	}
	tempBindPort := os.Getenv("anonBindPort")
	if len(tempBindPort) > 0 {

		tempPort, err := strconv.ParseInt(tempBindPort, 10, 10)
		if err != nil {
			panic("Cannot convert the Port")
		}

		BindPort = int(tempPort)
	} else {
		BindPort = 5050
	}

}
