package core

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ivangodev/spordieta/entity"
	"io/ioutil"
	"net/http"
)

type mockStorage struct {
	winProofs map[entity.UserBet]string
	payProofs map[entity.UserBet]string
}

type reqInfo struct {
	userId    entity.UserId
	betId     entity.BetId
	proofKind string
}

func newMockStorage() *mockStorage {
	return &mockStorage{
		map[entity.UserBet]string{},
		map[entity.UserBet]string{},
	}
}

func (c *mockStorage) uploaded(info reqInfo) bool {
	switch info.proofKind {
	case "winproof":
		_, ok := c.winProofs[entity.UserBet{info.userId, info.betId}]
		return ok
	case "payproof":
		_, ok := c.payProofs[entity.UserBet{info.userId, info.betId}]
		return ok
	}
	return false
}

func (c *mockStorage) deleteProofs(userId entity.UserId, betId entity.BetId) {
	delete(c.winProofs, entity.UserBet{userId, betId})
	delete(c.payProofs, entity.UserBet{userId, betId})
}

func (c *mockStorage) uploadProof(userId entity.UserId, betId entity.BetId, body []byte, proofKind string) {
	s := string(body)
	switch proofKind {
	case "winproof":
		c.winProofs[entity.UserBet{userId, betId}] = s
	case "payproof":
		c.payProofs[entity.UserBet{userId, betId}] = s
	}
}

func getReqInfo(w http.ResponseWriter, r *http.Request) (reqInfo, error) {
	userId := entity.UserId(mux.Vars(r)["user_id"])
	betId := entity.BetId(mux.Vars(r)["bet_id"])
	proofKind := mux.Vars(r)["proof_kind"]
	switch proofKind {
	case "winproof":
	case "payproof":
	default:
		w.Write([]byte(`{"error": "unknown proof kind"}`))
		return reqInfo{}, fmt.Errorf("Unknown proof kind")
	}
	return reqInfo{userId: userId, betId: betId, proofKind: proofKind}, nil
}

func (c *mockStorage) start() {
	r := mux.NewRouter()
	r.HandleFunc("/user/{user_id}/bet/{bet_id}/{proof_kind}",
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
				defer r.Body.Close()
				c.uploadProof(info.userId, info.betId, b, info.proofKind)
			case "GET":
				if c.uploaded(info) {
					w.Write([]byte(r.URL.Path + "/data"))
				} else {
					w.Write([]byte(""))
				}
			case "DELETE":
				c.deleteProofs(info.userId, info.betId)
			}
		})

	r.HandleFunc("/user/{user_id}/bet/{bet_id}/{proof_kind}/data",
		func(w http.ResponseWriter, r *http.Request) {
			info, err := getReqInfo(w, r)
			if err != nil {
				return
			}
			if r.Method == "GET" {
				var content string
				switch info.proofKind {
				case "winproof":
					content = c.winProofs[entity.UserBet{info.userId, info.betId}]
				case "payproof":
					content = c.payProofs[entity.UserBet{info.userId, info.betId}]
				}
				w.Write([]byte(content))
			}
		})

	http.Handle("/", r)
	go http.ListenAndServe(":8080", r)
}
