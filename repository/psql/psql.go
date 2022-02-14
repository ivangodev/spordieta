package psql

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"github.com/ivangodev/spordieta/entity"
	_ "github.com/lib/pq"
	"strconv"
	"time"
)

type PsqlRepo struct {
	Db *sql.DB
}

func OpenDB(host, password, dbname string) (*sql.DB, error) {
	port := 5432
	user := "postgres"

	var db *sql.DB
	err := fmt.Errorf("Error")
	for err != nil {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"dbname=%s password=%s sslmode=disable",
			host, port, user, dbname, password)
		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			fmt.Println("Failed to connect to db. Try again")
			continue
		}
		if err = db.Ping(); err != nil {
			time.Sleep(100 * time.Millisecond)
			fmt.Println("Failed to connect to db. Try again")
		}
	}

	return db, nil
}

func NewPsql(db *sql.DB) *PsqlRepo {
	q := `
		CREATE TABLE IF NOT EXISTS users (
			user_id text PRIMARY KEY
		);
		CREATE TABLE IF NOT EXISTS bets (
			user_id text REFERENCES users(user_id),
			bet_id text,
			bet_info bytea,
			UNIQUE (user_id, bet_id)
		);
	`
	r := &PsqlRepo{db}
	_, err := r.Db.Exec(q)
	if err != nil {
		panic(err)
	}
	return r
}

var ErrInvalidUser = fmt.Errorf("User does not exist")
var ErrInvalidBet = fmt.Errorf("Bet does not exist")

func (r *PsqlRepo) verifyUserBet(uid *entity.UserId, bid *entity.BetId) error {
	q := `SELECT user_id FROM users WHERE user_id = $1`
	rows, err := r.Db.Query(q, *uid)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return ErrInvalidUser
	}

	if bid != nil {
		q := `SELECT bet_id FROM bets WHERE user_id = $1 AND bet_id = $2`
		rows, err := r.Db.Query(q, *uid, *bid)
		if err != nil {
			return err
		}
		if !rows.Next() {
			return ErrInvalidBet
		}
	}

	return nil
}

func (r *PsqlRepo) CreateUser(id entity.UserId) error {
	q := `INSERT INTO users VALUES ($1)`
	_, err := r.Db.Exec(q, id)
	return err
}

func (r *PsqlRepo) CreateBet(id entity.UserId, info entity.BetInfo) (betId entity.BetId, err error) {
	err = r.verifyUserBet(&id, nil)
	if err != nil {
		return
	}
	q := `SELECT MAX(bet_id) FROM bets WHERE user_id = $1`
	rows, err := r.Db.Query(q, id)
	if err != nil {
		return
	}

	if rows.Next() {
		var bid sql.NullString
		err = rows.Scan(&bid)
		if err != nil {
			return
		}
		if !bid.Valid {
			bid.String = "-1"
		}
		b, _ := strconv.Atoi(bid.String)
		b++
		bid.String = strconv.Itoa(b)
		betId = entity.BetId(bid.String)
	} else {
		err = fmt.Errorf("Unexpected no rows")
		return
	}

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(info)
	if err != nil {
		return
	}
	q = `INSERT INTO bets VALUES ($1, $2, $3)`
	_, err = r.Db.Exec(q, id, betId, buf.Bytes())
	return
}

func (r *PsqlRepo) GetOpenedBet(id entity.UserId) (*entity.BetId, error) {
	err := r.verifyUserBet(&id, nil)
	if err != nil {
		return nil, err
	}
	q := `SELECT MAX(bet_id) FROM bets WHERE user_id = $1`
	rows, err := r.Db.Query(q, id)
	if err != nil {
		return nil, err
	}

	var betId entity.BetId
	if rows.Next() {
		var bid sql.NullString
		err = rows.Scan(&bid)
		if err != nil {
			return nil, err
		}
		if !bid.Valid {
			return nil, nil
		}
		betId = entity.BetId(bid.String)
	} else {
		err = fmt.Errorf("Unexpected no rows")
		return nil, err
	}

	info, err := r.GetBetInfo(id, betId)
	if err != nil {
		return nil, err
	}

	if info.Status.Opened {
		return &betId, nil
	}
	return nil, nil
}

func (r *PsqlRepo) GetOpenedBets() (res []entity.BetToReview, err error) {
	q := `SELECT user_id FROM users`
	var rows *sql.Rows
	rows, err = r.Db.Query(q)
	if err != nil {
		return
	}

	res = make([]entity.BetToReview, 0)
	for rows.Next() {
		var s sql.NullString
		err = rows.Scan(&s)
		if err != nil {
			return
		}
		if !s.Valid {
			continue
		}
		uid := entity.UserId(s.String)
		var bid *entity.BetId
		bid, err = r.GetOpenedBet(uid)
		if err != nil {
			return
		}
		if bid == nil {
			continue
		}
		var info entity.BetInfo
		info, err = r.GetBetInfo(uid, *bid)
		if err != nil {
			return
		}
		res = append(res, entity.BetToReview{entity.UserBet{uid, *bid},
			info.Status.State})
	}
	return
}

func (r *PsqlRepo) GetBetInfo(userId entity.UserId, betId entity.BetId) (info entity.BetInfo, err error) {
	err = r.verifyUserBet(&userId, &betId)
	if err != nil {
		return
	}
	q := `SELECT user_id, bet_info FROM bets WHERE user_id = $1 AND bet_id = $2`
	rows, err := r.Db.Query(q, userId, betId)
	if err != nil {
		return
	}

	var b []byte
	if rows.Next() {
		var uid sql.NullString
		err = rows.Scan(&uid, &b)
		if err != nil {
			return
		}
		if !uid.Valid {
			err = fmt.Errorf("No results for the request")
			return
		}
	} else {
		err = fmt.Errorf("Failed to get bets info")
		return
	}

	dec := gob.NewDecoder(bytes.NewReader(b))
	err = dec.Decode(&info)
	if err != nil {
		return
	}

	return
}

func (r *PsqlRepo) PatchBetStatus(userId entity.UserId, betId entity.BetId, status entity.BetStatus) error {
	info, err := r.GetBetInfo(userId, betId)
	if err != nil {
		return err
	}
	info.Status = status

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(info)
	if err != nil {
		return err
	}
	q := `UPDATE bets SET bet_info = $1 WHERE user_id = $2 AND bet_id = $3`
	_, err = r.Db.Exec(q, buf.Bytes(), userId, betId)
	return err
}
