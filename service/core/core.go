package core

import (
	"fmt"
	"github.com/ivangodev/spordieta/entity"
	"github.com/ivangodev/spordieta/repository"
	"log"
	"time"
)

type Core struct {
	repo     repository.RepositoryI
	strgConn StorageConnectI
}

func NewCoreService(repo repository.RepositoryI, strgConn StorageConnectI) *Core {
	return &Core{repo, strgConn}
}

func (s *Core) CreateUser(id entity.UserId) error {
	return s.repo.CreateUser(id)
}

func (s *Core) validateBetCond(cond entity.BetCond) error {
	return nil
}

func (s *Core) CreateBet(id entity.UserId, cond entity.BetCond) (betId entity.BetId, err error) {
	b, err := s.GetOpenedBet(id)
	if err != nil {
		return
	}
	if b != nil {
		err = fmt.Errorf("There is opened bet already: %v", betId)
		return
	}

	err = s.validateBetCond(cond)
	if err != nil {
		return
	}

	cond.Deadline = time.Now().Add(time.Duration(cond.DurationSec) * time.Second)
	defaultStatus := entity.BetStatus{
		Opened:       true,
		State:        entity.NeedWinProof,
		AdminManaged: false,
	}
	info := entity.BetInfo{Cond: cond, Status: defaultStatus}
	return s.repo.CreateBet(id, info)
}

func (s *Core) GetOpenedBet(id entity.UserId) (*entity.BetId, error) {
	return s.repo.GetOpenedBet(id)
}

func (s *Core) GetBetInfo(userId entity.UserId, betId entity.BetId) (info entity.BetInfo, err error) {
	info, err = s.repo.GetBetInfo(userId, betId)
	if err != nil {
		return
	}

	st := info.Status
	if !st.Opened {
		return
	}

	switch st.State {
	case entity.NeedWinProof:
		var url string
		url, err = s.strgConn.GetWinProofURL(userId, betId)
		if err != nil {
			return
		}
		info.Status.Uploaded = url != ""

		if !st.AdminManaged && !st.Uploaded && time.Now().After(info.Cond.Deadline) {
			info.Status.State = entity.NeedPayProof
			err = s.PatchBetStatus(userId, betId, info.Status)
			if err != nil {
				return
			}
		}
	case entity.NeedPayProof:
		var url string
		url, err = s.strgConn.GetPayProofURL(userId, betId)
		if err != nil {
			return
		}
		info.Status.Uploaded = url != ""
	default:
		err = fmt.Errorf("Unknokn bet state %v", st.State)
		return
	}

	return
}

func (s *Core) PatchBetStatus(userId entity.UserId, betId entity.BetId, status entity.BetStatus) error {
	status.AdminManaged = true
	err := s.repo.PatchBetStatus(userId, betId, status)
	if err != nil {
		return err
	}
	if !status.Opened || !status.Uploaded {
		if s.strgConn.DeleteProofs(userId, betId) != err {
			log.Printf("Failed to delete proofs (%v/%v): %s", userId, betId, err)
		}
	}
	return nil
}

func (s *Core) GetBetsToReview() (res []entity.BetToReview, err error) {
	usersBets, err := s.repo.GetOpenedBets()
	if err != nil {
		return
	}
	if usersBets == nil {
		usersBets = make([]entity.BetToReview, 0)
	}

	for _, userBet := range usersBets {
		userId := userBet.User
		betId := userBet.Bet

		var info entity.BetInfo
		info, err = s.GetBetInfo(userId, betId)
		if err != nil {
			return
		}

		st := info.Status
		if st.Uploaded {
			res = append(res, entity.BetToReview{
				entity.UserBet{userId, betId},
				st.State})
		}
	}

	return
}
