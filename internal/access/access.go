package access

import "time"

const TrialDuration = 72 * time.Hour

type User struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	Name         string     `json:"name"`
	Role         string     `json:"role"`
	Plan         string     `json:"plan"`
	TrialUntil   *time.Time `json:"trialUntil,omitempty"`
	CreatedAt    string     `json:"createdAt"`
	PasswordHash string     `json:"-"`
}

type Entitlements struct {
	Plan      string   `json:"plan"`
	Grades    any      `json:"grades"`
	Labs      any      `json:"labs"`
	Quests    any      `json:"quests"`
	Studio    any      `json:"studio"`
	PetWrite  bool     `json:"petWrite"`
	Exam      bool     `json:"exam"`
	Interview bool     `json:"interview"`
}

func NewTrialUntil(now time.Time) time.Time {
	return now.UTC().Add(TrialDuration)
}

func TrialActive(u *User) bool {
	if u == nil || u.TrialUntil == nil || u.TrialUntil.IsZero() {
		return false
	}
	return time.Now().Before(u.TrialUntil.UTC())
}

func TrialRemainingSeconds(u *User) int64 {
	if !TrialActive(u) {
		return 0
	}
	sec := int64(time.Until(u.TrialUntil.UTC()).Seconds())
	if sec < 0 {
		return 0
	}
	return sec
}

func IsPaidPro(u *User) bool {
	if u == nil {
		return false
	}
	return u.Role == "admin" || u.Plan == "pro"
}

func IsPro(u *User) bool {
	return IsPaidPro(u) || TrialActive(u)
}

func TrialUntilISO(u *User) any {
	if u == nil || u.TrialUntil == nil || u.TrialUntil.IsZero() {
		return nil
	}
	return u.TrialUntil.UTC().Format(time.RFC3339)
}

func PublicUser(u *User) map[string]any {
	if u == nil {
		return nil
	}
	return map[string]any{
		"id":             u.ID,
		"email":          u.Email,
		"name":           u.Name,
		"role":           u.Role,
		"plan":           u.Plan,
		"createdAt":      u.CreatedAt,
		"trialUntil":     TrialUntilISO(u),
		"trialActive":    TrialActive(u),
		"trialRemaining": TrialRemainingSeconds(u),
		"isPro":          IsPro(u),
	}
}

func SessionFields(u *User) map[string]any {
	if u == nil {
		return map[string]any{
			"plan":           "guest",
			"trialUntil":     nil,
			"trialActive":    false,
			"trialRemaining": int64(0),
			"isPro":          false,
		}
	}
	return map[string]any{
		"plan":           u.Plan,
		"trialUntil":     TrialUntilISO(u),
		"trialActive":    TrialActive(u),
		"trialRemaining": TrialRemainingSeconds(u),
		"isPro":          IsPro(u),
	}
}

func EntitlementsOf(u *User) Entitlements {
	if u == nil {
		return Entitlements{
			Plan: "guest", Grades: []string{}, Labs: []string{}, Quests: []string{},
			Studio: []string{}, PetWrite: false, Exam: false, Interview: false,
		}
	}
	if IsPro(u) {
		plan := "pro"
		if TrialActive(u) && !IsPaidPro(u) {
			plan = "pro"
		}
		return Entitlements{
			Plan:      plan,
			Grades:    []string{"intern", "junior", "middle", "senior"},
			Labs:      "*",
			Quests:    "*",
			Studio:    "*",
			PetWrite:  true,
			Exam:      true,
			Interview: true,
		}
	}
	return Entitlements{
		Plan:      "free",
		Grades:    []string{"intern"},
		Labs:      []string{"lab-bank-ledger", "lab-intern-1", "lab-intern-2", "lab-intern-3"},
		Quests:    []string{"wallet", "bank", "q-wallet", "q-hold", "q-iban"},
		Studio:    []string{"hold-capture", "t-account", "hall-intro"},
		PetWrite:  false,
		Exam:      false,
		Interview: true,
	}
}

func CanGrade(u *User, grade string) bool {
	if IsPro(u) {
		return true
	}
	e := EntitlementsOf(u)
	ids, ok := e.Grades.([]string)
	if !ok {
		return false
	}
	for _, g := range ids {
		if g == grade {
			return true
		}
	}
	return false
}

func CanPetWrite(u *User) bool {
	return IsPro(u)
}

func CanStudio(u *User, id string) bool {
	if IsPro(u) {
		return true
	}
	e := EntitlementsOf(u)
	ids, ok := e.Studio.([]string)
	if !ok {
		return true
	}
	for _, s := range ids {
		if s == id {
			return true
		}
	}
	return false
}
