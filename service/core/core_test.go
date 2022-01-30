package core

import (
	"github.com/ivangodev/spordieta/entity"
	"testing"
	"time"
)

func testCreate(t *testing.T, s *Core, durationSec int) (entity.UserId, entity.BetId) {
	//Create user
	userId := entity.UserId("TelegramId:42493039")
	err := s.CreateUser(userId)
	if err != nil {
		t.Fatalf("Failed to create user: %s", err)
	}

	//Create user with already used ID
	err = s.CreateUser(userId)
	if err == nil {
		t.Fatalf("Unexpected created user that already exists")
	}

	//Create bet
	cond := entity.BetCond{
		CurrWeight:  80,
		GoalWeight:  70,
		Money:       100,
		DurationSec: durationSec,
	}
	betId, err := s.CreateBet(userId, cond)
	if err != nil {
		t.Fatalf("Failed to create bet: %s", err)
	}

	//Create bet while another is still opened
	_, err = s.CreateBet(userId, cond)
	if err == nil {
		t.Fatalf("Unexpected created bet while another is still opened")
	}

	//Test if created bet is available
	betIdActual, err := s.GetOpenedBet(userId)
	if err != nil {
		t.Fatalf("Failed to get opened bet: %s", err)
	}
	if betId != *betIdActual {
		t.Fatalf("Unexpected ID of opened bet: want %v VS actual %v", betId,
			*betIdActual)
	}

	//Test bet info
	info, err := s.GetBetInfo(userId, betId)
	if err != nil {
		t.Fatalf("Failed to get bet info: %s", err)
	}
	want := entity.BetStatus{Opened: true, State: entity.NeedWinProof}
	if info.Status != want {
		t.Fatalf("Unexpected bet status: want %v VS actual %v", want, info.Status)
	}

	return userId, betId
}

func testPay(t *testing.T, s *Core, c *mockStorageConnect, userId entity.UserId, betId entity.BetId) {
	info, err := s.GetBetInfo(userId, betId)
	if err != nil {
		t.Fatalf("Failed to get bet info: %s", err)
	}
	if info.Status.State != entity.NeedPayProof {
		t.Fatalf("Unexpected bet info: must be NEED PAY state")
	}

	//User uploaded pay
	c.uploadPayProof(userId, betId)
	info, err = s.GetBetInfo(userId, betId)
	if err != nil {
		t.Fatalf("Failed to get bet info: %s", err)
	}
	if !info.Status.Uploaded {
		t.Fatalf("Unexpected bet info: must be uploaded")
	}

	//Admin didn't like it
	status := entity.BetStatus{
		Opened:       true,
		State:        entity.NeedPayProof,
		AdminComment: "bad photo quality",
	}
	err = s.PatchBetStatus(userId, betId, status)
	if err != nil {
		t.Fatalf("Failed to patch status %v", err)
	}
	info, err = s.GetBetInfo(userId, betId)
	if err != nil {
		t.Fatalf("Failed to get bet info: %s", err)
	}
	if info.Status.Uploaded {
		t.Fatalf("Unexpected bet info: must be unuploaded")
	}

	//User reuploaded pay
	c.uploadPayProof(userId, betId)

	//Admin liked it
	status.Opened = false
	err = s.PatchBetStatus(userId, betId, status)
	if err != nil {
		t.Fatalf("Failed to patch status %v", err)
	}
	info, err = s.GetBetInfo(userId, betId)
	if err != nil {
		t.Fatalf("Failed to get bet info: %s", err)
	}
	if info.Status.Opened {
		t.Fatalf("Unexpected bet info: must be closed")
	}
}

func testBadWinProof(t *testing.T, s *Core, c *mockStorageConnect, userId entity.UserId, betId entity.BetId) {
	info, err := s.GetBetInfo(userId, betId)
	if err != nil {
		t.Fatalf("Failed to get bet info: %s", err)
	}
	if info.Status.State != entity.NeedWinProof {
		t.Fatalf("Unexpected bet info: must be NEED WIN state")
	}

	//User uploaded pay
	c.uploadWinProof(userId, betId)
	info, err = s.GetBetInfo(userId, betId)
	if err != nil {
		t.Fatalf("Failed to get bet info: %s", err)
	}
	if !info.Status.Uploaded {
		t.Fatalf("Unexpected bet info: must be uploaded")
	}

	//Admin didn't like it and oblige user to pay
	status := entity.BetStatus{
		Opened:       true,
		State:        entity.NeedPayProof,
		Uploaded:     false,
		AdminComment: "You haven't lost weight",
	}
	err = s.PatchBetStatus(userId, betId, status)
	if err != nil {
		t.Fatalf("Failed to patch status %v", err)
	}
	info, err = s.GetBetInfo(userId, betId)
	if err != nil {
		t.Fatalf("Failed to get bet info: %s", err)
	}
	if info.Status.Uploaded {
		t.Fatalf("Unexpected bet info: must be unuploaded")
	}
}

func TestTimeoutExpries(t *testing.T) {
	r := newMockRepo()
	c := newMockStorageConnect()
	s := newCoreService(r, c)
	u, b := testCreate(t, s, 1)
	time.Sleep(2 * time.Second)
	testPay(t, s, c, u, b)
}

func TestBadWinProof(t *testing.T) {
	r := newMockRepo()
	c := newMockStorageConnect()
	s := newCoreService(r, c)
	u, b := testCreate(t, s, 1)
	testBadWinProof(t, s, c, u, b)
	testPay(t, s, c, u, b)
}
