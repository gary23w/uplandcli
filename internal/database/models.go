package database

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type UserCredentials struct {
	Url      string
	PSQLurl  string
	User     string
	Password string
	Host     string
	Port     string
	Database string
	RowLoad  int
}

func writeConfigFiles(userJson UserCredentials, userConf string) {
	// write user credentials to database json file
	json_file, err := json.MarshalIndent(userJson, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("./conf/database.json", json_file, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// set API url
	myConf := `
	appname = props_crud
	httpport = 1337
	runmode = dev
	autorender = false
	copyrequestbody = true
	EnableDocs = true
	sqlconn = ` + string(userConf) + `
	`

	// write string to config file
	err = ioutil.WriteFile("./conf/app.conf", []byte(myConf), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// get postgres credentials from json file
func getPostgresCredentials() UserCredentials {
	var userCredentials UserCredentials
	jsonFile, err := os.Open("conf/database.json")
	if err != nil {
		log.Fatalln("Couldn't open the json file", err)
		return userCredentials
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &userCredentials)
	return userCredentials
}

func setLoadVar(rowCount UserCredentials) {
	json_file, err := json.MarshalIndent(rowCount, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("./conf/database.json", json_file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
