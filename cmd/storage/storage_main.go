package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ivangodev/spordieta/delivery"
	"github.com/ivangodev/spordieta/service/storage"
	"net/http"
	"os"
)

func main() {
	rootDir := os.Getenv("STORAGE_ROOT_DIR")
	if rootDir == "" {
		panic("Unspecified root directory is not allowed")
	}

	storagePort := os.Getenv("STORAGE_PORT")
	if storagePort == "" {
		panic("Unspecified storage port is not allowed")
	}
	storagePort = ":" + storagePort

	s := storage.NewStorage(rootDir)
	r := mux.NewRouter()
	d := delivery.NewHttpStorage(r, s)
	d.RegisterEndpoints()
	fmt.Println("Starting storage")
	panic(http.ListenAndServe(storagePort, r))
}
