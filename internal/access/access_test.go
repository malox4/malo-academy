package access

import (
	"testing"
	"time"
)

func TestIsProTrialAndPaid(t *testing.T) {
	if IsPro(nil) {
		t.Fatal("nil user must not be pro")
	}
	free := &User{Plan: "free"}
	if IsPro(free) || TrialActive(free) {
		t.Fatal("free without trial is not pro")
	}
	future := time.Now().Add(2 * time.Hour)
	trial := &User{Plan: "free", TrialUntil: &future}
	if !TrialActive(trial) || !IsPro(trial) {
		t.Fatal("active trial must grant pro")
	}
	if TrialRemainingSeconds(trial) <= 0 {
		t.Fatal("trial remaining must be positive")
	}
	past := time.Now().Add(-time.Hour)
	expired := &User{Plan: "free", TrialUntil: &past}
	if TrialActive(expired) || IsPro(expired) {
		t.Fatal("expired trial must fall back to intern")
	}
	paid := &User{Plan: "pro"}
	if !IsPro(paid) || TrialActive(paid) {
		t.Fatal("paid pro is pro without trial")
	}
	admin := &User{Plan: "free", Role: "admin"}
	if !IsPro(admin) {
		t.Fatal("admin is pro")
	}
}

func TestGrades(t *testing.T) {
	free := &User{Plan: "free"}
	if !CanGrade(free, "intern") || CanGrade(free, "junior") {
		t.Fatal("intern free, junior locked")
	}
	future := time.Now().Add(time.Hour)
	trial := &User{Plan: "free", TrialUntil: &future}
	if !CanGrade(trial, "junior") || !CanPetWrite(trial) {
		t.Fatal("trial opens junior and pet write")
	}
}
