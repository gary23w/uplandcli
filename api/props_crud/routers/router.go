package routers

import (
	"eos_bot/api/props_crud/controllers"

	beego "eos_bot/api/props_crud/web"
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
