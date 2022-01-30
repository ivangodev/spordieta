package core

import (
	"fmt"
	"github.com/ivangodev/spordieta/entity"
)

type mockRepo struct {
	users   map[entity.UserId]struct{}
	betsCnt int
	userBet map[entity.UserId]entity.BetId
	betInfo map[entity.BetId]entity.BetInfo
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		users:   map[entity.UserId]struct{}{},
		userBet: map[entity.UserId]entity.BetId{},
		betInfo: map[entity.BetId]entity.BetInfo{},
	}
}

func (s *mockRepo) validateUser(id entity.UserId) error {
	if _, exist := s.users[id]; !exist {
		return fmt.Errorf("User with ID %v does not exist", id)
	}
	return nil
}

func (s *mockRepo) validateBet(id entity.BetId) error {
	if _, exist := s.betInfo[id]; !exist {
		return fmt.Errorf("Bet with ID %v does not exist", id)
	}
	return nil
}

func (s *mockRepo) CreateUser(id entity.UserId) error {
	err := s.validateUser(id)
	if err == nil {
		return fmt.Errorf("User with ID %v already exists", id)
	}
	s.users[id] = struct{}{}

	return nil
}

func (s *mockRepo) CreateBet(id entity.UserId, info entity.BetInfo) (betId entity.BetId, err error) {
	err = s.validateUser(id)
	if err != nil {
		return
	}

	betId = entity.BetId(s.betsCnt)
	s.userBet[id] = betId
	s.betInfo[betId] = info
	s.betsCnt++
	return
}

func (s *mockRepo) GetOpenedBet(id entity.UserId) (*entity.BetId, error) {
	err := s.validateUser(id)
	if err != nil {
		return nil, err
	}
	betId := s.userBet[id]
	info := s.betInfo[betId]
	if info.Status.Opened {
		return &betId, nil
	}
	return nil, nil
}

func (s *mockRepo) GetOpenedBets() (res []entity.BetToReview, err error) {
	for userId, betId := range s.userBet {
		info := s.betInfo[betId]
		if info.Status.Opened {
			res = append(res, entity.BetToReview{
				entity.UserBet{userId, betId},
				info.Status.State,
			})
		}
	}
	return
}

func (s *mockRepo) GetBetInfo(userId entity.UserId, betId entity.BetId) (info entity.BetInfo, err error) {
	err = s.validateUser(userId)
	if err != nil {
		return
	}
	err = s.validateBet(betId)
	if err != nil {
		return
	}

	return s.betInfo[betId], nil
}

func (s *mockRepo) PatchBetStatus(userId entity.UserId, betId entity.BetId, status entity.BetStatus) error {
	err := s.validateUser(userId)
	if err != nil {
		return err
	}
	err = s.validateBet(betId)
	if err != nil {
		return err
	}

	s.betInfo[betId] = entity.BetInfo{s.betInfo[betId].Cond, status}
	return nil
}
