package entity

import (
	"time"
)

type UserId string
type BetId string

type BetCond struct {
	CurrWeight  float64   `json:"currWeight"`
	GoalWeight  float64   `json:"goalWeight"`
	Money       int       `json:"money"`
	DurationSec int       `json:"durationSec"` //User sends this
	Deadline    time.Time `json:"-"`           //Repo saves actual deadline
	LeftSec     int       `json:"leftSec"`     //Service returns this
}

const (
	NeedWinProof = iota
	NeedPayProof
)

type BetStatus struct {
	Opened       bool   `json:"opened"`
	State        int    `json:"state"`
	Uploaded     bool   `json:"uploaded"`
	AdminComment string `json:"adminComment"`
	AdminManaged bool   `json:"-"`
}

type BetInfo struct {
	Cond   BetCond
	Status BetStatus
}

type UserBet struct {
	User UserId
	Bet  BetId
}
type BetToReview struct {
	UserBet
	State int
}
