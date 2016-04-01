package httplib

import (
	"log"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) StatusAll() {

	this.Ctx.WriteString("OK")

}

func (this *MainController) Healthz() {
	this.Ctx.WriteString("Healthy")
}

func Run(config string) {

	log.Printf("Starting the HTTP server at port %s", config)

	beego.Run(":" + config)

}
