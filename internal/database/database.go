package database

import (
	"database/sql"
	"eos_bot/internal/models"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var creds UserCredentials = getPostgresCredentials()

func connectToDatabase() *sql.DB {
	psql := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", creds.Host, creds.Port, creds.User, creds.Password, creds.Database)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		log.Fatal(err)
	}
    err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func CreateTables() {
	db := connectToDatabase()
	defer db.Close()
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS properties (
		id SERIAL PRIMARY KEY,
		type VARCHAR(255) NOT NULL,
		prop_id VARCHAR(255) NOT NULL,
		address VARCHAR(255) NOT NULL,
		latitude VARCHAR(255) NOT NULL,
		longitude VARCHAR(255) NOT NULL,
		upx VARCHAR(255) NOT NULL,
		fiat VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`)
	if err != nil {
		log.Fatal(err)
	}
}

func AddPropertiesToDatabase(properties []models.DataPackageBLOCK) {
	usr := getPostgresCredentials()
	db := connectToDatabase()
	defer db.Close()
	CreateTables() // if table doesn't exist, create it
	if usr.RowLoad < 10000 {
		for _, value := range properties {
			_, err := db.Exec(fmt.Sprintf("INSERT INTO properties (type, prop_id, address, latitude, longitude, upx, fiat) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s') ON CONFLICT (ID) DO NOTHING", value.Type, value.ID, value.Address, value.Lat, value.Long, value.UPX, value.FIAT))
			if err != nil {
				log.Fatal(err)
			}
			usr.RowLoad = usr.RowLoad + 1
		}
	} else {
		log.Println("Row load limit reached, DATABASE OFFLINE...")
	}
	fmt.Println("[*] Rows loaded: ", usr.RowLoad)
	setLoadVar(usr)
}

func GetPropertiesFromDatabase() {
	db := connectToDatabase()
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