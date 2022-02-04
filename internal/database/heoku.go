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

type UserCredentials struct {
	Cur_Name string
	Url      string
	PSQLurl  string
	User     string
	Password string
	Host     string
	Port     string
	Database string
	RowLoad  int
	PostUrl string
}

type HerokuMain struct {
	Name string	
	Creds UserCredentials
}

func (x *HerokuMain) CheckHeroku() bool {
    out, err := exec.Command("heroku", "--version").Output()
    if err != nil {
        log.Println(err)
		log.Println("Heroku is not installed. See https://devcenter.heroku.com/articles/heroku-cli for more details.")
		return false
    }
    log.Printf("HEROKU VERSION: %v\n", string(out))
	return true
}

func (x *HerokuMain) LoginHeroku() {
	if x.CheckHeroku() == false {
		// force user to login to heroku?
		log.Fatalf("Please exit this app and login to heroku before proceeding any further. See https://devcenter.heroku.com/articles/heroku-cli for more details.")
	}
}

func (x *HerokuMain) DetailParser(cd1 []byte, cd2 []byte, cd3 []byte) {
	ur := strings.Split(string(cd1), "|")
	url := ur[0]

	// get postgres name from output
	regg := regexp.MustCompile(`Created\s(.*)\sas`)
	urp := regg.FindStringSubmatch(string(cd2))
	url2 := strings.Split(urp[0], " ")

	// split database config and collect login credentials
	postg_reg := regexp.MustCompile(`DATABASE_URL:\s(.*)`)
	post_rget := postg_reg.FindStringSubmatch(string(cd3))
	user_reg := regexp.MustCompile(`\/\/(.*)@(.*):5432\/(.*)`)
	user := user_reg.FindStringSubmatch(string(cd3))
	u := strings.Split(user[1], ":")
	// get url for API config
	post_url := post_rget[1] + "?sslmode=require"

	// set user credential struct
	x.Creds = UserCredentials{
		Cur_Name: x.Name,
		Url: 	url,
		PSQLurl: url2[1],
		User:     u[0],
		Password: u[1],
		Host:     user[2],
		Port:     "5432",
		Database: user[3],
		RowLoad: 0,
		PostUrl: post_url,
	}
	log.Println(x.Creds)
}

func (x *HerokuMain) BuildPostgres() {
	log.Println("Building heroku postgres...")
	cmd, err := exec.Command("heroku", "create", "--app", x.Name).Output()
	if err != nil {
		log.Fatal(cmd)
		log.Fatal(err)
	}
	log.Printf("%s\n", cmd)	
	cmd2, err2 := exec.Command("heroku", "addons:create", "heroku-postgresql", "-a", x.Name).Output()
	if err2 != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", cmd2)
	time.Sleep(10 * time.Second)
	cmd3, err3 := exec.Command("heroku", "config", "-a", x.Name).Output()
	if err3 != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", cmd3)
	x.DetailParser(cmd, cmd2, cmd3)
	writeConfigFiles(x.Creds, x.Creds.PostUrl)
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
	// quote name to work around exec.Command() bug?
	nameQuoted := fmt.Sprintf("%v", name)
	sys := HerokuMain{ Name: nameQuoted }
    sys.LoginHeroku()	
	log.Println("Building postgres instance with name: " + nameQuoted)
	sys.BuildPostgres()
	log.Println("[*] Waiting for service to come online...")
	var bar utils.Bar
	bar.NewOption(0, 100)
	for i := 0; i <= 100; i++ {
		time.Sleep(time.Millisecond * 800)
		bar.Play(int64(i))
	}
	bar.Finish()
}