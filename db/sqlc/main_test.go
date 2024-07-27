package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/chowchow-dev/bankie/util"
	_ "github.com/golang/mock/gomock"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot read config:", err)
	}

	testDb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect database")
	}
	testQueries = New(testDb)

	os.Exit(m.Run())
}
