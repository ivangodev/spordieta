package repository

import (
	"github.com/ivangodev/spordieta/entity"
)

type RepositoryI interface {
	CreateUser(id entity.UserId) error
	CreateBet(id entity.UserId, info entity.BetInfo) (betId entity.BetId, err error)
	GetOpenedBet(id entity.UserId) (*entity.BetId, error)
	GetOpenedBets() (res []entity.BetToReview, err error)
	GetBetInfo(userId entity.UserId, betId entity.BetId) (info entity.BetInfo, err error)
	PatchBetStatus(userId entity.UserId, betId entity.BetId, info entity.BetStatus) error
}
