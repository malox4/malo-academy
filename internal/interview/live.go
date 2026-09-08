package interview

import (
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Answer struct {
	QuestionID string `json:"questionId"`
	PickID     string `json:"pickId,omitempty"`
	Text       string `json:"text,omitempty"`
	Score      int    `json:"score"`
	Why        string `json:"why"`
	Follow     bool   `json:"follow"`
	At         string `json:"at"`
}

type Session struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Mode        string    `json:"mode"`
	StartedAt   time.Time `json:"startedAt"`
	EndsAt      time.Time `json:"endsAt"`
	Index       int       `json:"index"`
	IDs         []string  `json:"ids"`
	Answers     []Answer  `json:"answers"`
	Status      string    `json:"status"`
	FollowOpen  bool      `json:"followOpen"`
	LastWhy     string    `json:"lastWhy,omitempty"`
	LastGood    bool      `json:"lastGood"`
}

type Assessment struct {
	Level      string         `json:"level"`
	Confidence float64        `json:"confidence"`
	Score      int            `json:"score"`
	Max        int            `json:"max"`
	Hards      map[string]int `json:"hards"`
	Softs      map[string]int `json:"softs"`
	Notes      []string       `json:"notes"`
	At         string         `json:"at"`
	Practice   int            `json:"practiceFold"`
	SQL        int            `json:"sqlFold"`
}

type Evidence struct {
	PracticeAvg int
	PracticeN   int
	SQLSolved   int
	SQLTotal    int
}

type Hub struct {
	mu   sync.Mutex
	live map[string]*Session
}

func NewHub() *Hub { return &Hub{live: map[string]*Session{}} }

func (h *Hub) Start(userID, mode string) *Session {
	h.mu.Lock()
	defer h.mu.Unlock()
	ids := LiveIDs()
	mins := LiveMinutes
	if mode != "live" {
		mode = "path"
		ids = OrderedIDs()
		mins = 24 * 60
	}
	s := &Session{
		ID: uuid.NewString(), UserID: userID, Mode: mode,
		StartedAt: time.Now().UTC(), EndsAt: time.Now().UTC().Add(time.Duration(mins) * time.Minute),
		Index: 0, IDs: ids, Status: "running",
	}
	h.live[userID] = s
	return s
}

func (h *Hub) Get(userID string) *Session {
	h.mu.Lock()
	defer h.mu.Unlock()
	s := h.live[userID]
	if s == nil {
		return nil
	}
	cp := *s
	return &cp
}

func (h *Hub) Public(s *Session) map[string]any {
	if s == nil {
		return map[string]any{"session": nil}
	}
	timedOut := s.Mode == "live" && time.Now().UTC().After(s.EndsAt) && s.Status == "running"
	done := s.Status == "done" || s.Index >= len(s.IDs) || timedOut
	var q *Question
	if !done && s.Index < len(s.IDs) {
		q = QuestionByID(s.IDs[s.Index])
	}
	remain := int(time.Until(s.EndsAt).Seconds())
	if remain < 0 {
		remain = 0
	}
	out := map[string]any{
		"session": map[string]any{
			"id": s.ID, "mode": s.Mode, "status": s.Status, "index": s.Index,
			"total": len(s.IDs), "ids": s.IDs, "startedAt": s.StartedAt.Format(time.RFC3339),
			"endsAt": s.EndsAt.Format(time.RFC3339), "remainSec": remain,
			"followOpen": s.FollowOpen, "lastWhy": s.LastWhy, "lastGood": s.LastGood,
			"answered": len(s.Answers),
		},
		"question": q,
		"done":     done,
	}
	if q != nil && s.Index < len(s.IDs) {
		out["nextId"] = s.IDs[s.Index]
		if s.Index+1 < len(s.IDs) {
			out["peekNextId"] = s.IDs[s.Index+1]
		}
	}
	return out
}

func (h *Hub) Answer(userID, pickID, text string, advance bool) (*Session, *Question, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	s := h.live[userID]
	if s == nil || s.Status != "running" {
		return s, nil, errNo
	}
	if s.Mode == "live" && time.Now().UTC().After(s.EndsAt) {
		s.Status = "done"
		return s, nil, nil
	}
	if s.Index < 0 || s.Index >= len(s.IDs) {
		s.Status = "done"
		return s, nil, nil
	}
	qid := s.IDs[s.Index]
	q := QuestionByID(qid)
	if q == nil {
		s.Index++
		if s.Index >= len(s.IDs) {
			s.Status = "done"
		}
		return s, nil, nil
	}
	if pickID == "" && strings.TrimSpace(text) == "" {
		if advance && !s.FollowOpen && hasMainAnswer(s, qid) {
			s.Index++
			if s.Index >= len(s.IDs) {
				s.Status = "done"
			}
			return s, q, nil
		}
		return s, q, errEmpty
	}
	score, why := 0, ""
	if pickID != "" {
		for _, r := range q.Replies {
			if r.ID == pickID {
				if r.Good {
					score = 2
				}
				why = r.Why
				break
			}
		}
	}
	if strings.TrimSpace(text) != "" {
		ts, tw := ScoreText(*q, text)
		if pickID == "" || ts >= score {
			score, why = ts, tw
		}
	}
	s.LastWhy = why
	s.LastGood = score >= 2
	s.Answers = append(s.Answers, Answer{
		QuestionID: qid, PickID: pickID, Text: text, Score: score, Why: why,
		Follow: s.FollowOpen, At: time.Now().UTC().Format(time.RFC3339),
	})
	if !s.FollowOpen && score < 2 && q.Follow != "" && s.Mode == "live" {
		s.FollowOpen = true
		return s, q, nil
	}
	s.FollowOpen = false
	s.Index++
	if s.Index >= len(s.IDs) {
		s.Status = "done"
	}
	return s, q, nil
}

func hasMainAnswer(s *Session, qid string) bool {
	for _, a := range s.Answers {
		if a.QuestionID == qid && !a.Follow {
			return true
		}
	}
	return false
}

func (h *Hub) Finish(userID string, ev Evidence) (*Session, Assessment) {
	h.mu.Lock()
	defer h.mu.Unlock()
	s := h.live[userID]
	if s == nil {
		return nil, Assessment{}
	}
	s.Status = "done"
	a := Evaluate(s, ev)
	return s, a
}

type simpleErr string

func (e simpleErr) Error() string { return string(e) }

var errNo = simpleErr("Нет живой сессии. Начните Live собес.")
var errEmpty = simpleErr("Скажите ответ или выберите ход.")

func Evaluate(s *Session, ev Evidence) Assessment {
	hards := map[string]int{"sql": 0, "process": 0, "bank": 0, "integrations": 0, "bpmn": 0, "data": 0}
	softs := map[string]int{"clarifying": 0, "stakeholders": 0, "escalation": 0, "clarity": 0}
	hardN := map[string]int{}
	softN := map[string]int{}
	score, max := 0, 0
	seen := map[string]bool{}
	for _, a := range s.Answers {
		if a.Follow {
			score += a.Score
			max += 2
			continue
		}
		if seen[a.QuestionID] {
			continue
		}
		seen[a.QuestionID] = true
		score += a.Score
		max += 2
		q := QuestionByID(a.QuestionID)
		if q == nil {
			continue
		}
		pct := 0
		if a.Score >= 2 {
			pct = 100
		} else if a.Score == 1 {
			pct = 55
		} else {
			pct = 20
		}
		for _, sk := range q.Skills {
			if _, ok := hards[sk]; ok {
				hards[sk] += pct
				hardN[sk]++
			}
			if _, ok := softs[sk]; ok {
				softs[sk] += pct
				softN[sk]++
			}
			if sk == "bank" {
				hards["integrations"] += pct / 2
				hardN["integrations"]++
			}
		}
	}
	avg := func(sum, n int) int {
		if n == 0 {
			return 35
		}
		return sum / n
	}
	for k := range hards {
		hards[k] = avg(hards[k], hardN[k])
	}
	for k := range softs {
		softs[k] = avg(softs[k], softN[k])
	}
	notes := []string{}
	if ev.SQLSolved > 0 && ev.SQLTotal > 0 {
		fold := 40 + int(60*float64(ev.SQLSolved)/float64(ev.SQLTotal))
		hards["sql"] = (hards["sql"]*2 + fold) / 3
		hards["data"] = (hards["data"]*2 + fold) / 3
		notes = append(notes, "В уровень вошли решённые SQL-задачи учебного банка.")
	}
	if ev.PracticeN > 0 {
		p := ev.PracticeAvg
		if p <= 5 {
			p = p * 20
		}
		hards["process"] = (hards["process"]*2 + p) / 3
		softs["clarity"] = (softs["clarity"]*2 + p) / 3
		notes = append(notes, "Практика (лабы и работы) сдвинула процесс и ясность.")
	}
	level, conf := LevelOf(score, max)
	if ev.SQLSolved >= 4 && level == "intern" {
		level = "junior"
		conf = 0.5
		notes = append(notes, "SQL-лаба подняла нижнюю оценку с intern.")
	}
	if max < 6 {
		conf = conf * 0.7
		notes = append(notes, "Мало ответов — уверенность снижена.")
	}
	return Assessment{
		Level: level, Confidence: conf, Score: score, Max: max,
		Hards: hards, Softs: softs, Notes: notes,
		At: time.Now().UTC().Format(time.RFC3339),
		Practice: ev.PracticeAvg, SQL: ev.SQLSolved,
	}
}
