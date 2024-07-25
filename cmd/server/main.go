package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/atsuyaourt/xyz-books/internal"
	db "github.com/atsuyaourt/xyz-books/internal/db/sqlc"
	docs "github.com/atsuyaourt/xyz-books/internal/docs/api"
	"github.com/atsuyaourt/xyz-books/internal/util"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"golang.org/x/sync/errgroup"
	_ "modernc.org/sqlite"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

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

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	_, err = os.Stat(config.DBSource)
	if os.IsNotExist(err) {
		f, _ := os.Create(config.DBSource)
		f.Close()
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %s", err)
	}

	dbSource := fmt.Sprintf("%s://%s?query", config.DBDriver, config.DBSource)
	err = util.DBMigrationUp(config.MigrationSrc, dbSource)
	if err != nil {
		log.Fatalf("migration error: %s", err)
	}

	store := db.NewStore(conn)

	g, ctx := errgroup.WithContext(ctx)
	runGinServer(ctx, g, config, store)

	err = g.Wait()
	if err != nil {
		log.Fatalf("error from wait group: %s", err)
	}
}

func runGinServer(ctx context.Context, g *errgroup.Group, config util.Config, store db.Store) {
	server, err := internal.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create server: %s", err)
	}

	server.Start(ctx, g)
}
