package delivery

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/ivangodev/spordieta/service/core"
	"io/ioutil"
	"net/http"
	"testing"
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
	d := NewHttpDeliv(routerCore, s)
	d.RegisterEndpoints()
	go func() {
		err := http.ListenAndServe(corePort, routerCore)
		if err != nil {
			panic(err)
		}
	}()
}

func TestEndpoints(t *testing.T) {
	startEndpoints()
	client := &http.Client{}

	method := "POST"
	url := domainedCore("/user/0")
	t.Logf("Method: %s\nURL: %s\n", method, url)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	t.Logf("Response status:\n%s\n", resp.Status)
	t.Logf("Response body:\n%s\n\n\n", string(respBody))

	method = "POST"
	url = domainedCore("/user/0/bet")
	body := bytes.NewBuffer([]byte(`
	{
		"currWeight":  80,
		"goalWeight":  75,
		"money":	   100,
		"durationSec": 60
	}
	`))
	t.Logf("Method: %s\nURL: %s\nBody: %s\n", method, url, body)
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	t.Logf("Response status:\n%s\n", resp.Status)
	t.Logf("Response body:\n%s\n\n\n", string(respBody))

	method = "GET"
	url = domainedCore("/user/0/bet/opened")
	t.Logf("Method: %s\nURL: %s\n", method, url)
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	t.Logf("Response status:\n%s\n", resp.Status)
	t.Logf("Response body:\n%s\n\n\n", string(respBody))

	method = "GET"
	url = domainedCore("/user/0/bet/0")
	t.Logf("Method: %s\nURL: %s\n", method, url)
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	t.Logf("Response status:\n%s\n", resp.Status)
	t.Logf("Response body:\n%s\n\n\n", string(respBody))

	method = "PUT"
	url = domainedStorage("/user/0/bet/0/winproof")
	body = bytes.NewBuffer([]byte(`I lost weight!`))
	t.Logf("Method: %s\nURL: %s\nBody: %s\n", method, url, body)
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	t.Logf("Response status:\n%s\n", resp.Status)
	t.Logf("Response body:\n%s\n\n\n", string(respBody))

	method = "GET"
	url = domainedCore("/toreview")
	t.Logf("Method: %s\nURL: %s\n", method, url)
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	t.Logf("Response status:\n%s\n", resp.Status)
	t.Logf("Response body:\n%s\n\n\n", string(respBody))

	method = "GET"
	url = domainedStorage("/user/0/bet/0/winproof")
	t.Logf("Method: %s\nURL: %s\n", method, url)
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	t.Logf("Response status:\n%s\n", resp.Status)
	t.Logf("Response body:\n%s\n\n\n", string(respBody))

	method = "GET"
	url = domainedStorage("/user/0/bet/0/winproof/data")
	t.Logf("Method: %s\nURL: %s\n", method, url)
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	t.Logf("Response status:\n%s\n", resp.Status)
	t.Logf("Response body:\n%s\n\n\n", string(respBody))

	method = "PATCH"
	url = domainedCore("/user/0/bet/0")
	body = bytes.NewBuffer([]byte(`
	{
		"opened": true,
		"state": 1,
		"uploaded": false,
		"adminComment": "You failed to lose weight. You must pay."
	}
	`))
	t.Logf("Method: %s\nURL: %s\nBody: %s\n", method, url, body)
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	t.Logf("Response status:\n%s\n", resp.Status)
	t.Logf("Response body:\n%s\n\n\n", string(respBody))

	method = "PUT"
	url = domainedStorage("/user/0/bet/0/payproof")
	body = bytes.NewBuffer([]byte(`I payed!`))
	t.Logf("Method: %s\nURL: %s\nBody: %s\n", method, url, body)
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	t.Logf("Response status:\n%s\n", resp.Status)
	t.Logf("Response body:\n%s\n\n\n", string(respBody))

	method = "PATCH"
	url = domainedCore("/user/0/bet/0")
	body = bytes.NewBuffer([]byte(`
	{
		"opened": false
	}
	`))
	t.Logf("Method: %s\nURL: %s\nBody: %s\n", method, url, body)
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed request: %s", err)
	}
	t.Logf("Response status:\n%s\n", resp.Status)
	t.Logf("Response body:\n%s\n\n\n", string(respBody))
}
