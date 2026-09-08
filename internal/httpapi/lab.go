package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/malox4/malo-academy/internal/access"
	"github.com/malox4/malo-academy/internal/interview"
	"github.com/malox4/malo-academy/internal/petdb"
)

func (s *Server) needUser(w http.ResponseWriter, r *http.Request) *access.User {
	u := s.user(r)
	if u == nil {
		writeErr(w, 401, "UNAUTH", "Войдите в зал.")
		return nil
	}
	return u
}

func (s *Server) sqlTasks(w http.ResponseWriter, r *http.Request) {
	u := s.needUser(w, r)
	if u == nil {
		return
	}
	s.syncLab()
	show := map[string]bool{}
	p, _ := s.Store.GetProgress(r.Context(), u.ID)
	if p != nil {
		var payload map[string]any
		_ = json.Unmarshal(p, &payload)
		if sql, ok := payload["sqlLab"].(map[string]any); ok {
			if solved, ok := sql["solved"].(map[string]any); ok {
				for id := range solved {
					show[id] = true
				}
			}
		}
	}
	items := []map[string]any{}
	for _, row := range petdb.PublicTasks(show) {
		if !access.IsPro(u) && !truthy(row["internSafe"]) {
			row["locked"] = true
			delete(row, "hint")
		}
		items = append(items, row)
	}
	writeJSON(w, 200, map[string]any{"items": items, "intern": !access.IsPro(u)})
}

func (s *Server) sqlSchema(w http.ResponseWriter, r *http.Request) {
	u := s.needUser(w, r)
	if u == nil {
		return
	}
	writeJSON(w, 200, map[string]any{"tables": s.Lab.Schema(!access.IsPro(u)), "intern": !access.IsPro(u)})
}

func (s *Server) sqlRun(w http.ResponseWriter, r *http.Request) {
	u := s.needUser(w, r)
	if u == nil {
		return
	}
	body, _ := readJSON(r)
	sqlText := strings.TrimSpace(str(body["sql"]))
	s.syncLab()
	pro := access.IsPro(u)
	res := s.Lab.Run(sqlText, !pro, pro)
	status := 200
	if !res.OK {
		status = 400
	}
	writeJSON(w, status, res)
}

func (s *Server) sqlCheck(w http.ResponseWriter, r *http.Request) {
	u := s.needUser(w, r)
	if u == nil {
		return
	}
	body, _ := readJSON(r)
	task := petdb.TaskByID(str(body["taskId"]))
	if task == nil {
		writeErr(w, 404, "NOT_FOUND", "Нет такой задачи.")
		return
	}
	if !access.IsPro(u) && !task.InternSafe {
		writeErr(w, 403, "PLAN", "Эта задача SQL — в PRO и в триале.")
		return
	}
	sqlText := strings.TrimSpace(str(body["sql"]))
	s.syncLab()
	pro := access.IsPro(u)
	res := s.Lab.Run(sqlText, !pro, false)
	if !res.OK {
		writeJSON(w, 400, res)
		return
	}
	gold, err := s.Lab.QueryGold(task.Gold)
	if err != nil {
		writeErr(w, 500, "GOLD", err.Error())
		return
	}
	ok, why := petdb.MatchGold(res.Rows, gold)
	out := map[string]any{"ok": ok, "why": why, "result": res, "taskId": task.ID}
	if ok {
		out["gold"] = task.Gold
	s.markSQL(r.Context(), u.ID, task.ID)
	}
	writeJSON(w, 200, out)
}

func (s *Server) markSQL(ctx context.Context, userID, taskID string) {
	p, _ := s.Store.GetProgress(ctx, userID)
	payload := map[string]any{}
	if p != nil {
		_ = json.Unmarshal(p, &payload)
	}
	sqlLab, _ := payload["sqlLab"].(map[string]any)
	if sqlLab == nil {
		sqlLab = map[string]any{}
	}
	solved, _ := sqlLab["solved"].(map[string]any)
	if solved == nil {
		solved = map[string]any{}
	}
	solved[taskID] = map[string]any{"ok": true, "at": time.Now().UTC().Format(time.RFC3339)}
	sqlLab["solved"] = solved
	payload["sqlLab"] = sqlLab
	raw, _ := json.Marshal(payload)
	_ = s.Store.PutProgress(ctx, userID, raw)
}

func (s *Server) explorer(w http.ResponseWriter, r *http.Request) {
	u := s.user(r)
	if u == nil {
		writeErr(w, 401, "UNAUTH", "Войдите в зал.")
		return
	}
	items := petdb.ExplorerCatalog()
	writeJSON(w, 200, map[string]any{"items": items, "count": len(items), "groups": []string{"Счета", "Холды", "Журнал", "P2P", "IBAN", "FX", "Карты", "Сверка", "KYC", "Контракт", "Тикеты", "Магазин"}})
}

func (s *Server) interviewArc(w http.ResponseWriter, r *http.Request) {
	if s.needUser(w, r) == nil {
		return
	}
	writeJSON(w, 200, map[string]any{"stages": interview.Arc(), "ids": interview.OrderedIDs(), "liveIds": interview.LiveIDs()})
}

func (s *Server) interviewStart(w http.ResponseWriter, r *http.Request) {
	u := s.needUser(w, r)
	if u == nil {
		return
	}
	body, _ := readJSON(r)
	mode := or(str(body["mode"]), "live")
	sess := s.Live.Start(u.ID, mode)
	writeJSON(w, 200, s.Live.Public(sess))
}

func (s *Server) interviewGet(w http.ResponseWriter, r *http.Request) {
	u := s.needUser(w, r)
	if u == nil {
		return
	}
	writeJSON(w, 200, s.Live.Public(s.Live.Get(u.ID)))
}

func (s *Server) interviewAnswer(w http.ResponseWriter, r *http.Request) {
	u := s.needUser(w, r)
	if u == nil {
		return
	}
	body, _ := readJSON(r)
	advance, _ := body["advance"].(bool)
	sess, _, err := s.Live.Answer(u.ID, str(body["pickId"]), str(body["text"]), advance)
	if err != nil {
		writeErr(w, 400, "ANSWER", err.Error())
		return
	}
	out := s.Live.Public(sess)
	out["assessment"] = nil
	writeJSON(w, 200, out)
}

func (s *Server) interviewFinish(w http.ResponseWriter, r *http.Request) {
	u := s.needUser(w, r)
	if u == nil {
		return
	}
	ev := s.evidence(r.Context(), u.ID)
	sess, a := s.Live.Finish(u.ID, ev)
	raw, _ := json.Marshal(a)
	_ = s.Store.PutAssessment(r.Context(), u.ID, raw)
	p, _ := s.Store.GetProgress(r.Context(), u.ID)
	payload := map[string]any{}
	if p != nil {
		_ = json.Unmarshal(p, &payload)
	}
	payload["assessment"] = a
	pr, _ := json.Marshal(payload)
	_ = s.Store.PutProgress(r.Context(), u.ID, pr)
	writeJSON(w, 200, map[string]any{"session": s.Live.Public(sess)["session"], "assessment": a})
}

func (s *Server) getAssessment(w http.ResponseWriter, r *http.Request) {
	u := s.needUser(w, r)
	if u == nil {
		return
	}
	raw, _ := s.Store.GetAssessment(r.Context(), u.ID)
	if raw == nil {
		writeJSON(w, 200, map[string]any{"assessment": nil})
		return
	}
	writeJSON(w, 200, map[string]any{"assessment": json.RawMessage(raw)})
}

func (s *Server) evidence(ctx context.Context, userID string) interview.Evidence {
	ev := interview.Evidence{SQLTotal: len(petdb.Tasks())}
	items, _ := s.Store.ListPractice(ctx, userID)
	sum, n := 0, 0
	for _, it := range items {
		switch sc := it["score"].(type) {
		case *int:
			if sc != nil {
				sum += *sc
				n++
			}
		case float64:
			sum += int(sc)
			n++
		case int:
			sum += sc
			n++
		}
	}
	ev.PracticeN = n
	if n > 0 {
		ev.PracticeAvg = sum / n
	}
	p, _ := s.Store.GetProgress(ctx, userID)
	if p != nil {
		var payload map[string]any
		_ = json.Unmarshal(p, &payload)
		if sql, ok := payload["sqlLab"].(map[string]any); ok {
			if solved, ok := sql["solved"].(map[string]any); ok {
				ev.SQLSolved = len(solved)
			}
		}
	}
	return ev
}

func truthy(v any) bool {
	b, _ := v.(bool)
	return b
}
