package storage

import (
	"errors"
	"github.com/ivangodev/spordieta/entity"
	"io/ioutil"
	"os"
)

type Storage struct {
	rootDir string
}

func NewStorage(rootDir string) *Storage {
	return &Storage{rootDir}
}

func (c *Storage) fname(userId entity.UserId, betId entity.BetId) string {
	return c.rootDir + "/" + string(userId) + "#" + string(betId)
}

func (c *Storage) Uploaded(userId entity.UserId, betId entity.BetId) bool {
	_, e := os.Stat(c.fname(userId, betId))
	return !errors.Is(e, os.ErrNotExist)
}

func (c *Storage) DeleteProofs(userId entity.UserId, betId entity.BetId) error {
	return os.Remove(c.fname(userId, betId))
}

func (c *Storage) UploadProof(userId entity.UserId, betId entity.BetId, body []byte) error {
	if c.Uploaded(userId, betId) {
		err := c.DeleteProofs(userId, betId)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(c.fname(userId, betId))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	return nil
}

func (c *Storage) GetProof(userId entity.UserId, betId entity.BetId) ([]byte, error) {
	data, err := ioutil.ReadFile(c.fname(userId, betId))
	if err != nil {
		return nil, err
	}
	return data, nil
}
