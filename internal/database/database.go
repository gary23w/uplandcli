package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gary23w/uplandcli/internal/models"

	_ "github.com/lib/pq"
)

type EOSDatabaseManager struct {
	db *sql.DB
	creds UserCredentials
}

var (
	globalOnce    sync.Once
	globalManager *EOSDatabaseManager
	globalErr     error
)

func getManager() (*EOSDatabaseManager, error) {
	globalOnce.Do(func() {
		x := &EOSDatabaseManager{}
		x.SetCredentials()
		x.ConnectToDatabase()
		if globalErr != nil {
			return
		}
		x.CreateTables()
		if globalErr != nil {
			return
		}
		globalManager = x
	})
	return globalManager, globalErr
}

func (x *EOSDatabaseManager) SetCredentials() {
	x.creds = getPostgresCredentials()
}

func (x *EOSDatabaseManager) ConnectToDatabase() {
	psql := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", x.creds.Host, x.creds.Port, x.creds.User, x.creds.Password, x.creds.Database)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		globalErr = err
		return
	}

	// Connection pool tuning (safe defaults for CLI + long-running tail mode)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		globalErr = fmt.Errorf("database connection failed (run 'upld database --deploy'): %w", err)
		return
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
		globalErr = err
		return
	}
	_, err = x.db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS properties_prop_id_uidx ON properties(prop_id);`)
	if err != nil {
		globalErr = err
		return
	}
}

func EOSDatabaseMan(route string, model []models.DataPackageBLOCK) {
	x, err := getManager()
	if err != nil {
		log.Println("database unavailable:", err)
		return
	}
	// Refresh credentials since RowLoad is persisted to JSON on disk.
	x.SetCredentials()
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
	if x.db == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := x.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("db begin failed:", err)
		return
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO properties (type, prop_id, address, latitude, longitude, upx, fiat)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (prop_id) DO NOTHING
	`)
	if err != nil {
		log.Println("db prepare failed:", err)
		return
	}
	defer stmt.Close()

	for _, value := range properties {
		if value.Type == "NULL-DATA" {
			continue
		}
		_, err := stmt.ExecContext(ctx, value.Type, value.ID, value.Address, value.Lat, value.Long, value.UPX, value.FIAT)
		if err != nil {
			log.Println("db insert failed:", err)
			return
		}
		x.creds.RowLoad++
	}

	if err := tx.Commit(); err != nil {
		log.Println("db commit failed:", err)
		return
	}
	setLoadVar(x.creds)
}

func (x *EOSDatabaseManager) GetPropertiesFromDatabase() {
	if x.db == nil {
		return
	}
	rows, err := x.db.Query("SELECT prop_id, address, upx, fiat FROM properties")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var propID, address, upx, fiat string
		err := rows.Scan(&propID, &address, &upx, &fiat)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(propID, address, upx, fiat)
	}
}