package psql

import (
	"database/sql"
	"fmt"
	"github.com/ivangodev/spordieta/entity"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"
)

var db *sql.DB

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	//Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func testEmptyTables(r *PsqlRepo, t *testing.T) {
	uid := entity.UserId(1)
	_, err := r.GetOpenedBet(uid)
	if err == nil {
		t.Fatalf("Unexpected non-empty result of get opened bet over empty tables")
	}

	bets, err := r.GetOpenedBets()
	if err != nil {
		t.Fatalf("Failed to get opened bets over empty tables: %s", err)
	}
	if len(bets) != 0 {
		t.Fatalf("Unexpected non-empty result of get opened bets over empty tables")
	}

	bid := entity.BetId(1)
	_, err = r.GetBetInfo(uid, bid)
	if err == nil {
		t.Fatalf("Unexpected non-empty result of get bet info over empty tables")
	}
}

func createUserAndBet(r *PsqlRepo, t *testing.T, uid entity.UserId) entity.BetId {
	err := r.CreateUser(uid)
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}

	info := entity.BetInfo{
		entity.BetCond{},
		entity.BetStatus{Opened: true},
	}
	bid, err := r.CreateBet(uid, info)
	if err != nil {
		t.Fatalf("Failed to create bet: %s", err)
	}

	bidp, err := r.GetOpenedBet(uid)
	if err != nil {
		t.Fatalf("Failed to get opened bet: %s", err)
	}
	if *bidp != bid {
		t.Fatalf("Unexpected opened bet: want %v VS actual %v", bid, *bidp)
	}

	infoActual, err := r.GetBetInfo(uid, bid)
	if err != nil {
		t.Fatalf("Failed to get bet info: %s", err)
	}
	if info != infoActual {
		t.Fatalf("Unexpected bet info: want %v VS actual %v", info, infoActual)
	}

	return bid
}

func TestPsql(t *testing.T) {
	var err error
	r := NewPsql(db)
	if err != nil {
		t.Fatalf("Failed to create repository: %s", err)
	}

	testEmptyTables(r, t)

	uid1, uid2 := entity.UserId(1), entity.UserId(2)
	bid1 := createUserAndBet(r, t, uid1)
	bid2 := createUserAndBet(r, t, uid2)

	_, err = r.GetBetInfo(uid1, "1337")
	if err == nil {
		t.Fatalf("Unexpected non-empty result of bet info with wrong bet id")
	}

	bets, err := r.GetOpenedBets()
	if err != nil {
		t.Fatalf("Failed to get opened bets: %s", err)
	}
	want := []entity.BetToReview{
		{entity.UserBet{uid1, bid1},
			entity.NeedWinProof},
		{entity.UserBet{uid2, bid2},
			entity.NeedWinProof}}
	match := make([]bool, 2)
	if len(bets) != 2 {
		t.Fatalf("Unexpected bets: want %v VS actual %v", want, bets)
	}
	for _, b1 := range bets {
		for i, b2 := range want {
			if b1 == b2 {
				match[i] = true
			}
		}
	}
	if !match[0] || !match[1] {
		t.Fatalf("Unexpected bets: want %v VS actual %v", want, bets)
	}

	status := entity.BetStatus{Opened: false}
	info := entity.BetInfo{entity.BetCond{}, status}
	err = r.PatchBetStatus(uid1, bid1, status)
	if err != nil {
		t.Fatalf("Failed to patch bet status: %s", err)
	}

	infoActual, err := r.GetBetInfo(uid1, bid1)
	if err != nil {
		t.Fatalf("Failed to get bet info: %s", err)
	}
	if info != infoActual {
		t.Fatalf("Unexpected bet info: want %v VS actual %v", info, infoActual)
	}
}
