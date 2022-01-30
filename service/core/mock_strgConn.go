package core

import (
	"github.com/ivangodev/spordieta/entity"
)

type mockStorageConnect struct {
	winProofs map[entity.UserBet]string
	payProofs map[entity.UserBet]string
}

func newMockStorageConnect() *mockStorageConnect {
	return &mockStorageConnect{
		map[entity.UserBet]string{},
		map[entity.UserBet]string{},
	}
}

func (c *mockStorageConnect) GetWinProofURL(userId entity.UserId, betId entity.BetId) (url string, err error) {
	return c.winProofs[entity.UserBet{userId, betId}], nil
}

func (c *mockStorageConnect) GetPayProofURL(userId entity.UserId, betId entity.BetId) (url string, err error) {
	return c.payProofs[entity.UserBet{userId, betId}], nil
}

func (c *mockStorageConnect) DeleteProofs(userId entity.UserId, betId entity.BetId) {
	delete(c.winProofs, entity.UserBet{userId, betId})
	delete(c.payProofs, entity.UserBet{userId, betId})
}

//Useful for tests:
func (c *mockStorageConnect) uploadWinProof(userId entity.UserId, betId entity.BetId) {
	c.winProofs[entity.UserBet{userId, betId}] = "never mind"
}

func (c *mockStorageConnect) uploadPayProof(userId entity.UserId, betId entity.BetId) {
	c.payProofs[entity.UserBet{userId, betId}] = "never mind"
}
