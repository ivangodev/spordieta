package storage

import (
	"github.com/ivangodev/spordieta/entity"
	"os"
	"reflect"
	"testing"
)

func TestStorage(t *testing.T) {
	rootDir := os.Getenv("STORAGE_ROOT_DIR")
	if rootDir == "" {
		t.Fatalf("Empty root directory is not allowed")
	}

	s := NewStorage(rootDir)
	uid, bid := entity.UserId("0"), entity.BetId("0")

	if s.Uploaded(uid, bid) {
		t.Fatalf("Unexpected uploaded data when nothing's been uploaded yet")
	}

	proof, err := s.GetProof(uid, bid)
	if err == nil {
		t.Fatalf("Unexpected get of empty proof")
	}

	if err := s.DeleteProofs(uid, bid); err == nil {
		t.Fatalf("Unexpected delete of empty proof")
	}

	proof = []byte("blabla")
	err = s.UploadProof(uid, bid, proof)
	if err != nil {
		t.Fatalf("Failed upload proof: %s", err)
	}

	if !s.Uploaded(uid, bid) {
		t.Fatalf("No uploaded data when it's been uploaded already")
	}

	actualProof, err := s.GetProof(uid, bid)
	if err != nil {
		t.Fatalf("Failed to get proof: %s", err)
	}
	if !reflect.DeepEqual(proof, actualProof) {
		t.Fatalf("Unexpected content of proof: want %v VS actual %v",
			proof, actualProof)
	}

	if err := s.DeleteProofs(uid, bid); err != nil {
		t.Fatalf("Failed to delete proof: %s", err)
	}

	if s.Uploaded(uid, bid) {
		t.Fatalf("Unexpected uploaded data when nothing's been uploaded yet")
	}
}
