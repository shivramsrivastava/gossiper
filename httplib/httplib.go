package httplib

import (
	"encoding/json"
	"log"

	"github.com/astaxie/beego"

	"../common"
)

type BootStrapResponse struct {
	Name     string
	Country  string
	City     string
	EndPoint string
}

type StatusResponse struct {
	Name              string
	CPU, MEM, DISK    float64 //Total CPU MEM and DISK
	UCPU, UMEM, UDISK float64 //Used CPU MEM and DISK
	OutOfResource     bool
}

type MainController struct {
	beego.Controller
}

func (this *MainController) StatusAll() {
	var res StatusResponse

	common.ALLDCs.Lck.Lock()
	defer common.ALLDCs.Lck.Unlock()

	dc, available := common.ALLDCs.List[common.ThisDCName]

	if !available {
		this.Ctx.WriteString("DC information not available")
		log.Printf("DC information not available")
		return
	}

	res.CPU = dc.CPU
	res.MEM = dc.MEM
	res.DISK = dc.DISK
	res.UCPU = dc.Ucpu
	res.UMEM = dc.Umem
	res.UDISK = dc.Udisk
	res.OutOfResource = dc.OutOfResource

	resp_byte, err := json.Marshal(&res)

	if err != nil {

		log.Printf("Error Marshalling the response")
		this.Ctx.WriteString("Status Failed")
		return
	}

	this.Ctx.WriteString(string(resp_byte))
	log.Printf("HTTP Status %s", string(resp_byte))
}

func (this *MainController) BootStrap() {

	var resp []BootStrapResponse

	for _, v := range common.ALLDCs.List {
		var dc BootStrapResponse
		dc.Name = v.Name
		dc.Country = v.Country
		dc.City = v.City
		dc.EndPoint = v.Endpoint
		resp = append(resp, dc)
	}
	resp_byte, err := json.Marshal(&resp)

	if err == nil {
		log.Println("Something wrong bootstrap api failed")
		this.Ctx.WriteString("Boot Strap Failed")
		return
	}
	this.Ctx.WriteString(string(resp_byte))
}

func (this *MainController) Healthz() {
	this.Ctx.WriteString("Healthy")
}

func Run(config string) {

	log.Printf("Starting the HTTP server at port %s", config)

	beego.Run(":" + config)

}
