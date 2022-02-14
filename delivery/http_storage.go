package delivery

import (
	"github.com/gorilla/mux"
	"github.com/ivangodev/spordieta/entity"
	"github.com/ivangodev/spordieta/service/storage"
	"io/ioutil"
	"net/http"
)

type reqInfo struct {
	userId entity.UserId
	betId  entity.BetId
}

type HttpStorage struct {
	router  *mux.Router
	service *storage.Storage
}

func NewHttpStorage(router *mux.Router, serv *storage.Storage) *HttpStorage {
	return &HttpStorage{router, serv}
}

func getReqInfo(w http.ResponseWriter, r *http.Request) (reqInfo, error) {
	userId := entity.UserId(mux.Vars(r)["user_id"])
	betId := entity.BetId(mux.Vars(r)["bet_id"])
	return reqInfo{userId: userId, betId: betId}, nil
}

func (c *HttpStorage) RegisterEndpoints() {
	c.router.HandleFunc("/user/{user_id}/bet/{bet_id}/{proof_kind}",
		func(w http.ResponseWriter, r *http.Request) {
			info, err := getReqInfo(w, r)
			if err != nil {
				return
			}
			switch r.Method {
			case "PUT":
				b, err := ioutil.ReadAll(r.Body)
				if err != nil {
					return
				}
				c.service.UploadProof(info.userId, info.betId, b)
			case "GET":
				if c.service.Uploaded(info.userId, info.betId) {
					w.Write([]byte(r.URL.Path + "/data"))
				} else {
					w.Write([]byte(""))
				}
			case "DELETE":
				c.service.DeleteProofs(info.userId, info.betId)
			}
		})

	c.router.HandleFunc("/user/{user_id}/bet/{bet_id}/{proof_kind}/data",
		func(w http.ResponseWriter, r *http.Request) {
			info, err := getReqInfo(w, r)
			if err != nil {
				return
			}
			if r.Method == "GET" {
				data, err := c.service.GetProof(info.userId, info.betId)
				if err != nil {
					return
				}
				w.Write(data)
			}
		})

	http.Handle("/", c.router)
}
