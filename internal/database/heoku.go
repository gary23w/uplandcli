package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"eos_bot/internal/utils"
)

// appname = heroes_crud
// httpport = 1337
// runmode = dev
// autorender = false
// copyrequestbody = true
// EnableDocs = true
// sqlconn =

func checkHeroku() bool {
	//get heroku version
    out, err := exec.Command("heroku", "--version").Output()
    if err != nil {
        fmt.Println(err)
		fmt.Println("Heroku is not installed. See https://devcenter.heroku.com/articles/heroku-cli for more details.")
		return false
    }
    fmt.Printf("HEROKU VERSION: %v\n", string(out))
	return true
}

func loginHeroku() {
	if checkHeroku() == false {
		// force user to login to heroku?
		log.Fatalf("Please exit this app and login to heroku before proceeding any further. See https://devcenter.heroku.com/articles/heroku-cli for more details.")
	}
}

func buildPostgres() {
	cmd, err := exec.Command("heroku", "create", "--app", "upland-cli").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", cmd)

	// get app url fronm output
	ur := strings.Split(string(cmd), "|")
	url := ur[0]
	
	cmd2, err2 := exec.Command("heroku", "addons:create", "heroku-postgresql", "-a", "upland-cli").Output()
	if err2 != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", cmd2)
    
	// get postgres name from output
	regg := regexp.MustCompile(`Created\s(.*)\sas`)
	urp := regg.FindStringSubmatch(string(cmd2))
	url2 := strings.Split(urp[0], " ")


	time.Sleep(10 * time.Second)
	cmd3, err3 := exec.Command("heroku", "config", "-a", "upland-cli").Output()
	if err3 != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", cmd3)
    
	// split database config and collect login credentials
	user_reg := regexp.MustCompile(`\/\/(.*)@(.*):5432\/(.*)`)
	user := user_reg.FindStringSubmatch(string(cmd3))
	u := strings.Split(user[1], ":")

	// set user credential struct
	user_cred := UserCredentials{
		Url: 	url,
		PSQLurl: url2[1],
		User:     u[0],
		Password: u[1],
		Host:     user[2],
		Port:     "5432",
		Database: user[3],
	}
	fmt.Println(user_cred)

	// write user credentials to database json file
	json_file, err := json.MarshalIndent(user_cred, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("./conf/database.json", json_file, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// set API url
	myConf := `
	appname = heroes_crud
	httpport = 1337
	runmode = dev
	autorender = false
	copyrequestbody = true
	EnableDocs = true
	sqlconn = ` + string(cmd3) + `
	`

	// write string to config file
	err = ioutil.WriteFile("./conf/app.conf", []byte(myConf), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func DestroyPostgres() {
	cmd, err := exec.Command("heroku", "destroy", "--app", "upland-cli", "--confirm", "upland-cli").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", cmd)
}

func DeployHeroku() {
    loginHeroku()
	var build string
	fmt.Println("Do you really want to build a postgresql database on heroku? (y/n)")
	fmt.Scanln(&build)
	if build == "y" {
		buildPostgres()
	}
	//loading bar
	fmt.Println("[*] Waiting for service to come online...")
	var bar utils.Bar
	bar.NewOption(0, 100)
	for i := 0; i <= 100; i++ {
		// sleep for 1 second
		time.Sleep(1 * time.Second)
		bar.Play(int64(i))
	}
	bar.Finish()
}