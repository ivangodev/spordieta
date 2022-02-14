package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ivangodev/spordieta/delivery"
	"github.com/ivangodev/spordieta/repository/psql"
	"github.com/ivangodev/spordieta/service/core"
	"net/http"
	"os"
)

func main() {
	corePort := os.Getenv("CORE_PORT")
	if corePort == "" {
		panic("Unspecified core port is not allowed")
	}
	corePort = ":" + corePort

	storageHostname := os.Getenv("STORAGE_HOSTNAME")
	if storageHostname == "" {
		panic("Unspecified storage hostname is not allowed")
	}

	storagePort := os.Getenv("STORAGE_PORT")
	if storagePort == "" {
		panic("Unspecified storage port is not allowed")
	}
	storagePort = ":" + storagePort

	host := os.Getenv("POSTGRES_HOSTNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	if host == "" || password == "" || dbname == "" {
		panic("Unspecified POSTGRES_HOSTNAME, _PASSWORD or _DB is not allowed")
	}

	router := mux.NewRouter()
	db, err := psql.OpenDB(host, password, dbname)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	repo := psql.NewPsql(db)
	cnct := core.NewStrgConn("http://" + storageHostname + storagePort)
	serv := core.NewCoreService(repo, cnct)
	deliv := delivery.NewHttpCore(router, serv)
	deliv.RegisterEndpoints()
	fmt.Println("Starting core")
	panic(http.ListenAndServe(corePort, router))
}
