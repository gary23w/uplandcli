package database

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
)

func writeNewName() string {
	var b strings.Builder
	b.Grow(7)
	for i := 0; i < 7; i++ {
		b.WriteByte(byte(65 + rand.Intn(25)))
	}
	return fmt.Sprintf("upl%s", strings.ToLower(b.String()))
}

func writeConfigFiles(userJson UserCredentials, userConf string) {
	// write user credentials to database json file
	log.Println("Writing user credentials to database json file")
	json_file, err := json.MarshalIndent(userJson, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("./conf/database.json", json_file, 0644)
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
	log.Println("Writing config file for API")
	err = os.WriteFile("./conf/app.conf", []byte(myConf), 0644)
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
	dec := json.NewDecoder(io.LimitReader(jsonFile, 2*1024*1024))
	_ = dec.Decode(&userCredentials)
	return userCredentials
}

func setLoadVar(rowCount UserCredentials) {
	json_file, err := json.MarshalIndent(rowCount, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("./conf/database.json", json_file, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
