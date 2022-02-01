package props_crud

import (
	_ "eos_bot/api/props_crud/routers"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	_ "github.com/lib/pq"
)

func StartAPI() {
	sqlconn, _ := beego.AppConfig.String("sqlconn")
	orm.RegisterDataBase("default", "postgres", sqlconn)
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		//beego.BConfig.WebConfig.StaticDir["/docs"] = "metrics"
	}
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: []string{"*"},
        AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
        AllowHeaders: []string{"Origin", "x-requested-with",
            "content-type",
            "accept",
            "origin",
            "authorization",
            "x-csrftoken"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,

    }))
    beego.Run()

}

