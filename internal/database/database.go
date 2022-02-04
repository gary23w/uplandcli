package database

import (
	"database/sql"
	"eos_bot/internal/models"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type EOSDatabaseManager struct {
	db *sql.DB
	creds UserCredentials
}

func (x *EOSDatabaseManager) SetCredentials() {
	x.creds = getPostgresCredentials()
}

func (x *EOSDatabaseManager) ConnectToDatabase() {
	psql := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", x.creds.Host, x.creds.Port, x.creds.User, x.creds.Password, x.creds.Database)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		log.Fatal(err)
	}
    err = db.Ping()
	if err != nil {
		log.Fatal("Database connection failed. Please run 'upld database --deploy' for more details.")
	}
	
	x.db = db
}

func(x *EOSDatabaseManager) CreateTables() {
	_, err := x.db.Exec(`
	CREATE TABLE IF NOT EXISTS properties (
		id SERIAL PRIMARY KEY,
		type VARCHAR(255) NOT NULL,
		prop_id VARCHAR(255) NOT NULL,
		address VARCHAR(255) NOT NULL,
		latitude VARCHAR(255) NOT NULL,
		longitude VARCHAR(255) NOT NULL,
		upx VARCHAR(255) NOT NULL,
		fiat VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func EOSDatabaseMan(route string, model []models.DataPackageBLOCK) {
	log.Println("[*] DATABASEMANAGER: ", route)
	var x EOSDatabaseManager	
	x.SetCredentials()
	x.ConnectToDatabase()
	x.CreateTables()
	usr := x.creds
	if usr.RowLoad >= 10000 {
			log.Println("Row load limit reached, DATABASE RESETTING...")		
			newDb := writeNewName()
			log.Printf("New database name: upland-cli-%s \n", newDb)
			DeployHeroku(newDb)
	} 
	switch route {
		case "add":
			x.AddPropertyList(model)
			break
		case "remove":
			log.Println("TODO: add removal of properties")
			break
		case "update":
			log.Println("TODO: add update of properties")
			break
		default:
			log.Println("Invalid route in database manager")
	}
}

func (x *EOSDatabaseManager) AddPropertyList (properties []models.DataPackageBLOCK) {
	db := x.db
	defer db.Close()
	for _, value := range properties {
		_, err := db.Exec(fmt.Sprintf("INSERT INTO properties (type, prop_id, address, latitude, longitude, upx, fiat) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s') ON CONFLICT (ID) DO NOTHING", value.Type, value.ID, value.Address, value.Lat, value.Long, value.UPX, value.FIAT))
		if err != nil {
			log.Fatal(err)
		}
		x.creds.RowLoad = x.creds.RowLoad + 1
	}
	setLoadVar(x.creds)
}

func (x *EOSDatabaseManager) GetPropertiesFromDatabase() {
	db := x.db
	defer db.Close()
	rows, err := db.Query("SELECT * FROM properties")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var value string
		err := rows.Scan(&name, &value)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(name, value)
	}
}