package storageconnect

import (
	"github.com/ivangodev/spordieta/entity"
)

type StorageConnectI interface {
	GetWinProofURL(userId entity.UserId, betId entity.BetId) (url string, err error)
	GetPayProofURL(userId entity.UserId, betId entity.BetId) (url string, err error)
	DeleteProofs(userId entity.UserId, betId entity.BetId)
}
