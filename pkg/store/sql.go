package store

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"strings"
)

type SQLX struct {
	Client *sqlx.DB
}

func New(dbSource string) (store SQLX, err error) {
	//dbSource := "postgres://postgres:postgres@localhost:5432/book_store?sslmode=disable"
	driverName := strings.ToLower(strings.Split(dbSource, "://")[0])
	store.Client, err = sqlx.Connect(driverName, dbSource)
	if err != nil {
		fmt.Println("inadsadsads")
		//log.Fatalf("Failed to connect to database: %v", err)
		panic(err)
		return
	}
	store.Client.SetMaxOpenConns(20)

	return
}
