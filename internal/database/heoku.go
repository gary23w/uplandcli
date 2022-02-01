package database

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"eos_bot/internal/utils"
)

func checkHeroku() bool {
    out, err := exec.Command("heroku", "--version").Output()
    if err != nil {
        log.Println(err)
		log.Println("Heroku is not installed. See https://devcenter.heroku.com/articles/heroku-cli for more details.")
		return false
    }
    log.Printf("HEROKU VERSION: %v\n", string(out))
	return true
}

func loginHeroku() {
	if checkHeroku() == false {
		// force user to login to heroku?
		log.Fatalf("Please exit this app and login to heroku before proceeding any further. See https://devcenter.heroku.com/articles/heroku-cli for more details.")
	}
}

func buildPostgres(name string) {
	log.Println("Building heroku postgres...")

	cmd, err := exec.Command("heroku", "create", "--app", name).Output()
	if err != nil {
		log.Fatal(cmd)
		log.Fatal(err)
	}
	log.Printf("%s\n", cmd)

	ur := strings.Split(string(cmd), "|")
	url := ur[0]
	
	cmd2, err2 := exec.Command("heroku", "addons:create", "heroku-postgresql", "-a", name).Output()
	if err2 != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", cmd2)
    
	// get postgres name from output
	regg := regexp.MustCompile(`Created\s(.*)\sas`)
	urp := regg.FindStringSubmatch(string(cmd2))
	url2 := strings.Split(urp[0], " ")


	time.Sleep(10 * time.Second)

	cmd3, err3 := exec.Command("heroku", "config", "-a", name).Output()
	if err3 != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", cmd3)
    
	// split database config and collect login credentials
	postg_reg := regexp.MustCompile(`DATABASE_URL:\s(.*)`)
	post_rget := postg_reg.FindStringSubmatch(string(cmd3))
	user_reg := regexp.MustCompile(`\/\/(.*)@(.*):5432\/(.*)`)
	user := user_reg.FindStringSubmatch(string(cmd3))
	u := strings.Split(user[1], ":")
	// get url for API config
	post_url := post_rget[1] + "?sslmode=require"

	// set user credential struct
	user_cred := UserCredentials{
		Cur_Name: name,
		Url: 	url,
		PSQLurl: url2[1],
		User:     u[0],
		Password: u[1],
		Host:     user[2],
		Port:     "5432",
		Database: user[3],
		RowLoad: 0,
	}
	log.Println(user_cred)
	writeConfigFiles(user_cred, post_url)
}


func DestroyPostgres() {
	name := getPostgresCredentials().Cur_Name
	args := []string{"destroy", "--app", name, "--confirm", name}
	cmd, err := exec.Command("heroku", args...).Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", cmd)
}

func DeployHeroku(name string) {
	nameQuoted := fmt.Sprintf("%v", name)
    loginHeroku()
	var build string
	log.Println("Setup a postgres instance within heroku? (y/n)")
	fmt.Scanln(&build)
	if build == "y" {
		log.Println("Building postgres instance with name: " + nameQuoted)
		buildPostgres(nameQuoted)
	}
	log.Println("[*] Waiting for service to come online...")
	var bar utils.Bar
	bar.NewOption(0, 100)
	for i := 0; i <= 100; i++ {
		// sleep for half a second
		time.Sleep(time.Millisecond * 500)
		bar.Play(int64(i))
	}
	bar.Finish()
}