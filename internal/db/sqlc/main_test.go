package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/atsuyaourt/xyz-books/internal/util"
)

var (
	testConfig util.Config
	testDBUrl  string
	testStore  Store
)

func TestMain(m *testing.M) {
	tempFile, _ := os.CreateTemp("", "test.db")
	cleanup := func() {
		fmt.Println("Cleaning up...")
		tempFile.Close()
		os.Remove(tempFile.Name())
	}

	fmt.Println(tempFile.Name())

	testConfig = util.Config{
		MigrationSrc: "../../db/migrations",
		DBDriver:     "sqlite",
		DBSource:     tempFile.Name(),
	}

	testDBUrl = fmt.Sprintf("%s://%s?query", testConfig.DBDriver, testConfig.DBSource)

	testDB, err := sql.Open(testConfig.DBDriver, testConfig.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer testDB.Close()

	testStore = NewStore(testDB)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Received OS interrupt - exiting.")
		cleanup()
		os.Exit(0)
	}()

	exitVal := m.Run()
	cleanup()
	os.Exit(exitVal)
}
