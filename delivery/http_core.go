package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ivangodev/spordieta/entity"
	"github.com/ivangodev/spordieta/service/core"
	"net/http"
)

type HttpCore struct {
	router  *mux.Router
	service *core.Core
}

func NewHttpCore(router *mux.Router, serv *core.Core) *HttpCore {
	return &HttpCore{router, serv}
}

func writeJson(w http.ResponseWriter, js []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func errToJson(e error) []byte {
	return []byte(`{"error": "` + e.Error() + `"}`)
}

func writeError(w http.ResponseWriter, e error) {
	w.WriteHeader(http.StatusInternalServerError)
	writeJson(w, errToJson(e))
}

func (d *HttpCore) createUser(w http.ResponseWriter, r *http.Request) {
	userId := entity.UserId(mux.Vars(r)["user_id"])
	err := d.service.CreateUser(userId)
	if err != nil {
		writeError(w, fmt.Errorf("Failed to create user: %w", err))
		return
	}
}

func (d *HttpCore) createBet(w http.ResponseWriter, r *http.Request) {
	userId := entity.UserId(mux.Vars(r)["user_id"])

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var cond entity.BetCond
	err := dec.Decode(&cond)
	if err != nil {
		writeError(w, err)
		return
	}

	betId, err := d.service.CreateBet(userId, cond)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJson(w, []byte([]byte(`{"bet_id": "`+betId+`"}`)))
}

func (d *HttpCore) getOpenedBet(w http.ResponseWriter, r *http.Request) {
	userId := entity.UserId(mux.Vars(r)["user_id"])
	betId, err := d.service.GetOpenedBet(userId)
	if err != nil {
		writeError(w, err)
		return
	}

	if betId != nil {
		writeJson(w, []byte([]byte(`{"bet_id": "`+*betId+`"}`)))
	} else {
		writeJson(w, []byte([]byte(`{"bet_id": ""}`)))
	}
}

func (d *HttpCore) getBetInfo(w http.ResponseWriter, r *http.Request) {
	userId := entity.UserId(mux.Vars(r)["user_id"])
	betId := entity.BetId(mux.Vars(r)["bet_id"])
	info, err := d.service.GetBetInfo(userId, betId)
	if err != nil {
		writeError(w, err)
		return
	}

	resp, err := json.Marshal(info)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJson(w, resp)
}

func (d *HttpCore) patchBetStatus(w http.ResponseWriter, r *http.Request) {
	userId := entity.UserId(mux.Vars(r)["user_id"])
	betId := entity.BetId(mux.Vars(r)["bet_id"])

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var status entity.BetStatus
	err := dec.Decode(&status)
	if err != nil {
		writeError(w, err)
		return
	}

	err = d.service.PatchBetStatus(userId, betId, status)
	if err != nil {
		writeError(w, err)
		return
	}
}

func (d *HttpCore) getBetToReview(w http.ResponseWriter, r *http.Request) {
	bets, err := d.service.GetBetsToReview()
	if err != nil {
		writeError(w, err)
		return
	}

	var resp string
	for _, b := range bets {
		userId, betId := b.User, b.Bet
		var url string
		switch b.State {
		case entity.NeedWinProof:
			url = fmt.Sprintf("/user/%v/bet/%v/winproof", userId, betId)
		case entity.NeedPayProof:
			url = fmt.Sprintf("/user/%v/bet/%v/payproof", userId, betId)
		default:
			writeError(w, fmt.Errorf("Unknown state %d", b.State))
			return
		}
		resp += fmt.Sprintf("\"%s\",", url)
	}

	writeJson(w, []byte("["+resp+"]"))
}

func (d *HttpCore) RegisterEndpoints() {
	d.router.HandleFunc("/user/{user_id}",
		func(w http.ResponseWriter, r *http.Request) {
			d.createUser(w, r)
		}).Methods("POST")
	d.router.HandleFunc("/user/{user_id}/bet",
		func(w http.ResponseWriter, r *http.Request) {
			d.createBet(w, r)
		}).Methods("POST")
	d.router.HandleFunc("/user/{user_id}/bet/opened",
		func(w http.ResponseWriter, r *http.Request) {
			d.getOpenedBet(w, r)
		}).Methods("GET")
	d.router.HandleFunc("/user/{user_id}/bet/{bet_id}",
		func(w http.ResponseWriter, r *http.Request) {
			d.getBetInfo(w, r)
		}).Methods("GET")
	d.router.HandleFunc("/user/{user_id}/bet/{bet_id}",
		func(w http.ResponseWriter, r *http.Request) {
			d.patchBetStatus(w, r)
		}).Methods("PATCH")
	d.router.HandleFunc("/toreview",
		func(w http.ResponseWriter, r *http.Request) {
			d.getBetToReview(w, r)
		}).Methods("GET")
}
