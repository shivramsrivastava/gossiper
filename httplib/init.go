package httplib

import (
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/v1/BOOTSTRAP/", &MainController{}, "get:BootStrap")
	beego.Router("/v1/STATUS/", &MainController{}, "get:StatusAll")
	beego.Router("/healthz/", &MainController{}, "get:Healthz")
	beego.Router("/v1/LATENCY/", &MainController{}, "get:LatencyAll")
}
