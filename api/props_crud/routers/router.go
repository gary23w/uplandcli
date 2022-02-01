package routers

import (
	"eos_bot/api/props_crud/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/upland",

		beego.NSNamespace("/properties",
			beego.NSInclude(
				&controllers.PropertiesController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
