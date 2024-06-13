package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/atsuyaourt/xyz-books/internal"
	db "github.com/atsuyaourt/xyz-books/internal/db/sqlc"
	docs "github.com/atsuyaourt/xyz-books/internal/docs/api"
	"github.com/atsuyaourt/xyz-books/internal/services"
	"github.com/atsuyaourt/xyz-books/internal/util"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

// XYZBooksAPI
//
//	@title			XYZ Books API
//	@version		1.0
//	@description	XYZ Books API
//	@contact.name	Emilio Gozo
//	@contact.email	emiliogozo@proton.me
func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %s", err)
	}
	docs.SwaggerInfo.BasePath = config.APIBasePath

	_, err = os.Stat(config.DBSource)
	if os.IsNotExist(err) {
		f, _ := os.Create(config.DBSource)
		f.Close()
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %s", err)
	}

	dbSource := fmt.Sprintf("sqlite3://%s?query", config.DBSource)
	err = util.DBMigrationUp(config.MigrationSrc, dbSource)
	if err != nil {
		log.Fatalf("migration error: %s", err)
	}

	store := db.NewStore(conn)

	runGinServer(config, store)

	interval := 1 * time.Hour // Adjust the interval as needed
	ticker := time.NewTicker(interval)
	shutdownChan := make(chan struct{})

	var wg sync.WaitGroup

	isbnService := services.NewISBNService(config.HTTPServerAddress, config.OutputPath)

	wg.Add(1)
	go func() {
		defer wg.Done()
		isbnService.Run()
	}()

	for {
		select {
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				isbnService.Run()
			}()
		case <-shutdownChan:
			ticker.Stop()
			close(shutdownChan)
			wg.Wait()
			return
		}
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := internal.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create server: %s", err)
	}

	server.Start(config.HTTPServerAddress)
}
