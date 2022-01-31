package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"github.com/ivangodev/spordieta/service/core"
	"github.com/ivangodev/spordieta/delivery"
)

const (
	corePort    = ":8080"
	storagePort = ":8081"
)

func domainedCore(url string) string {
	return "http://localhost" + corePort + url
}

func domainedStorage(url string) string {
	return "http://localhost" + storagePort + url
}

func startEndpoints() {
	routerStorage := mux.NewRouter()
	storage := core.NewMockStorage(routerStorage)
	storage.RegisterEndpoints()
	go func() {
		err := http.ListenAndServe(storagePort, routerStorage)
		if err != nil {
			panic(err)
		}
	}()

	routerCore := mux.NewRouter()
	r := core.NewMockRepo()
	c := core.NewStrgConn(domainedStorage(""))
	s := core.NewCoreService(r, c)
	d := delivery.NewHttpDeliv(routerCore, s)
	d.RegisterEndpoints()
	go func() {
		err := http.ListenAndServe(corePort, routerCore)
		if err != nil {
			panic(err)
		}
	}()
}

func main() {
	startEndpoints()
	fmt.Println("Core and Storage started")
	var block chan int
	<-block
}
