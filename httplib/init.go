package httplib

import (
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/v1/STATUS/", &MainController{}, "get:StatusAll")
	beego.Router("/healthz/", &MainController{}, "get:Healthz")
}
