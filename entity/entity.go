package entity

import (
	"time"
)

type UserId string
type BetId string

type BetCond struct {
	CurrWeight  float64
	GoalWeight  float64
	Money       int
	DurationSec int       //User sends this
	Deadline    time.Time //Repo saves actual deadline
}

const (
	NeedWinProof = iota
	NeedPayProof
)

type BetStatus struct {
	Opened       bool
	State        int
	Uploaded     bool
	AdminComment string
	AdminManaged bool
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
