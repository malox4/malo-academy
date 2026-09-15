package httpapi

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/malox4/malo-academy/internal/content"
	"github.com/malox4/malo-academy/internal/store"
)

func (s *Server) here(w http.ResponseWriter, r *http.Request) {
	u := s.user(r)
	if u == nil {
		writeErr(w, 401, "UNAUTH", "Войдите в зал.")
		return
	}
	body, _ := readJSON(r)
	path := str(body["path"])
	title := str(body["title"])
	if path == "" {
		writeJSON(w, 200, map[string]any{"ok": true})
		return
	}
	_ = s.Store.TouchPresence(r.Context(), u.ID, path, title)
	writeJSON(w, 200, map[string]any{"ok": true})
}

func (s *Server) adminUsers(w http.ResponseWriter, r *http.Request) {
	if adminOK(s, w, r) == nil {
		return
	}
	rows, err := s.Store.ListUsers(r.Context())
	if err != nil {
		writeErr(w, 500, "STORE", err.Error())
		return
	}
	progress, _ := s.Store.MapProgress(r.Context())
	assess, _ := s.Store.MapAssessments(r.Context())
	prac, _ := s.Store.PracticeCounts(r.Context())
	items := make([]map[string]any, 0, len(rows))
	online, silent, walking := 0, 0, 0
	for _, row := range rows {
		item := s.deskRow(row, progress[str(row["id"])], assess[str(row["id"])], prac[str(row["id"])])
		if item == nil {
			continue
		}
		items = append(items, item)
		if item["online"] == true {
			online++
		}
		if n, _ := item["lessonsDone"].(int); n == 0 {
			silent++
		} else {
			walking++
		}
	}
	writeJSON(w, 200, map[string]any{
		"items": items,
		"store": s.Store.Mode(),
		"stats": map[string]any{"all": len(items), "online": online, "silent": silent, "walking": walking},
	})
}

func (s *Server) deskRow(row map[string]any, raw json.RawMessage, assess json.RawMessage, prac store.PracStat) map[string]any {
	pub := publicFromListRow(row)
	if pub == nil {
		return nil
	}
	seen := or(str(row["last_seen"]), str(row["lastSeen"]))
	path := or(str(row["last_path"]), str(row["lastPath"]))
	title := or(str(row["last_title"]), str(row["lastTitle"]))
	online := false
	if t, err := time.Parse(time.RFC3339, seen); err == nil {
		online = time.Since(t.UTC()) < 15*time.Minute
	}
	track := summarizeProgress(raw, s.Cat)
	item := pub
	item["lastSeen"] = nil
	if seen != "" {
		item["lastSeen"] = seen
	}
	item["herePath"] = path
	item["here"] = placeLabel(path, title, s.Cat)
	item["online"] = online
	item["lessonsDone"] = track.done
	item["lessonsTotal"] = track.total
	item["pct"] = track.pct
	item["floor"] = track.floor
	item["floorTitle"] = track.floorTitle
	item["lastLesson"] = track.last
	item["recent"] = track.recent
	item["floors"] = track.floors
	item["interviews"] = track.interviews
	item["practiceOpen"] = prac.Open
	item["practiceTotal"] = prac.Total
	item["assessLevel"] = assessLevel(assess)
	return item
}

type progTrack struct {
	done       int
	total      int
	pct        int
	floor      string
	floorTitle string
	last       map[string]any
	recent     []map[string]any
	floors     []map[string]any
	interviews int
}

func summarizeProgress(raw json.RawMessage, cat *content.Catalog) progTrack {
	out := progTrack{recent: []map[string]any{}, floors: []map[string]any{}}
	done := map[string]int64{}
	interviews := 0
	if len(raw) > 0 {
		var payload map[string]any
		if json.Unmarshal(raw, &payload) == nil {
			if lessons, ok := payload["lessons"].(map[string]any); ok {
				for id, v := range lessons {
					m, _ := v.(map[string]any)
					if m == nil {
						continue
					}
					if m["done"] == true {
						done[id] = millis(m["at"])
					}
				}
			}
			if iv, ok := payload["interview"].(map[string]any); ok {
				interviews = len(iv)
			}
		}
	}
	out.done = len(done)
	out.interviews = interviews
	if cat == nil {
		out.total = out.done
		if out.total > 0 {
			out.pct = 100
		}
		return out
	}
	type rec struct {
		id, title, grade, gradeTitle string
		at                           int64
	}
	var recs []rec
	for _, g := range cat.Grades {
		gDone, gTotal := 0, 0
		for _, lv := range g.Levels {
			for _, id := range lv.ModuleIDs {
				gTotal++
				l, _ := cat.Lessons[id]
				title := l.Title
				if title == "" {
					title = id
				}
				if at, ok := done[id]; ok {
					gDone++
					recs = append(recs, rec{id: id, title: title, grade: g.ID, gradeTitle: g.Title, at: at})
				}
			}
		}
		pct := 0
		if gTotal > 0 {
			pct = gDone * 100 / gTotal
		}
		out.floors = append(out.floors, map[string]any{"id": g.ID, "title": g.Title, "done": gDone, "total": gTotal, "pct": pct})
		out.total += gTotal
	}
	if out.total > 0 {
		out.pct = out.done * 100 / out.total
	}
	sort.Slice(recs, func(i, j int) bool { return recs[i].at > recs[j].at })
	if len(recs) > 0 {
		out.last = map[string]any{"id": recs[0].id, "title": recs[0].title, "at": recs[0].at}
		out.floor = recs[0].grade
		out.floorTitle = recs[0].gradeTitle
	}
	limit := 6
	if len(recs) < limit {
		limit = len(recs)
	}
	for i := 0; i < limit; i++ {
		out.recent = append(out.recent, map[string]any{"id": recs[i].id, "title": recs[i].title, "at": recs[i].at})
	}
	if out.floor == "" && len(cat.Grades) > 0 {
		out.floor = cat.Grades[0].ID
		out.floorTitle = cat.Grades[0].Title
	}
	return out
}

func placeLabel(path, title string, cat *content.Catalog) string {
	if title != "" {
		return title
	}
	p := path
	if i := strings.Index(p, "?"); i >= 0 {
		p = p[:i]
	}
	switch {
	case p == "" || p == "/":
		return ""
	case p == "/hall":
		return "Зал"
	case p == "/pet":
		return "Пет"
	case strings.HasPrefix(p, "/pet"):
		return "Пет"
	case p == "/practice":
		return "Практика"
	case strings.HasPrefix(p, "/practice"):
		return "Практика"
	case p == "/interview":
		return "Собес"
	case p == "/live" || strings.HasPrefix(p, "/live"):
		return "Live собес"
	case p == "/materials":
		return "Материалы"
	case strings.HasPrefix(p, "/materials/"):
		return "Материал"
	case strings.HasPrefix(p, "/boards"):
		return "Доска"
	case p == "/profile":
		return "Профиль"
	case p == "/admin":
		return "Журнал"
	case strings.HasPrefix(p, "/lesson/"):
		id := strings.TrimPrefix(p, "/lesson/")
		if cat != nil {
			if l, ok := cat.Lesson(id); ok {
				return l.Title
			}
		}
		return "Урок"
	case strings.HasPrefix(p, "/level/"):
		return "Этаж"
	default:
		return p
	}
}

func assessLevel(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var m map[string]any
	if json.Unmarshal(raw, &m) != nil {
		return ""
	}
	return str(m["level"])
}

func millis(v any) int64 {
	switch n := v.(type) {
	case float64:
		return int64(n)
	case int64:
		return n
	case json.Number:
		i, _ := n.Int64()
		return i
	default:
		return 0
	}
}
