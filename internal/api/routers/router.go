package routers

import (
	"github.com/gary23w/uplandcli/internal/api/controllers"

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
