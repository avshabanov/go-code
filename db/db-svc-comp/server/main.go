package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/avshabanov/go-code/db/db-svc-comp/server/logic"
)

var (
	address = flag.String("server-address", "localhost:9090", "The server address in the format of host:port")
	dbPath  = flag.String("db-path", "/tmp/identity-service.db", "Path to identity service database")
	dbType  = flag.String("db-type", "sqlite", "Type of the database to test")
)

func main() {
	flag.Parse()

	var err error

	dao, err := logic.NewSqliteDao()
	if err != nil {
		log.Fatalf("cannot create dao: %v", err)
		return
	}
	defer dao.Close()

	dao.QueryOrders("1", "", 20)

	svr := &http.Server{
		Addr:         *address,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	if err = svr.ListenAndServe(); err != nil {
		log.Fatalf("cannot start the server: %v", err)
		return
	}
}
