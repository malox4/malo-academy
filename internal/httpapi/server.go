package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/malox4/malo-academy/internal/access"
	"github.com/malox4/malo-academy/internal/content"
	"github.com/malox4/malo-academy/internal/interview"
	"github.com/malox4/malo-academy/internal/ledger"
	"github.com/malox4/malo-academy/internal/petdb"
	"github.com/malox4/malo-academy/internal/store"
)

const cookieName = "malo_sid"

type Server struct {
	Store  *store.Store
	Bank   *ledger.Bank
	Lab    *petdb.Lab
	Live   *interview.Hub
	WebDir string
	Secure bool
	Cat    *content.Catalog
}

func New(st *store.Store, bank *ledger.Bank, webDir string, cat *content.Catalog, lab *petdb.Lab) *Server {
	s := &Server{Store: st, Bank: bank, WebDir: webDir, Cat: cat, Lab: lab, Live: interview.NewHub(), Secure: os.Getenv("COOKIE_SECURE") == "1" || os.Getenv("NODE_ENV") == "production"}
	s.syncLab()
	return s
}

func (s *Server) syncLab() {
	if s.Lab == nil || s.Bank == nil {
		return
	}
	s.Bank.Lock()
	dump := s.Bank.TeachingDump()
	s.Bank.Unlock()
	_ = s.Lab.Sync(dump)
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", s.health)
	mux.HandleFunc("POST /api/auth/register", s.register)
	mux.HandleFunc("POST /api/auth/login", s.login)
	mux.HandleFunc("POST /api/auth/logout", s.logout)
	mux.HandleFunc("GET /api/auth/me", s.me)
	mux.HandleFunc("GET /api/me", s.me)
	mux.HandleFunc("GET /api/me/progress", s.getProgress)
	mux.HandleFunc("PUT /api/me/progress", s.putProgress)
	mux.HandleFunc("GET /api/progress", s.getProgress)
	mux.HandleFunc("PUT /api/progress", s.putProgress)
	mux.HandleFunc("GET /api/practice", s.listPractice)
	mux.HandleFunc("POST /api/practice", s.postPractice)
	mux.HandleFunc("GET /api/admin/users", s.adminUsers)
	mux.HandleFunc("PATCH /api/admin/users/{id}", s.adminPatchUser)
	mux.HandleFunc("GET /api/admin/practice", s.adminPractice)
	mux.HandleFunc("PATCH /api/admin/practice/{id}", s.adminReview)
	mux.HandleFunc("GET /api/catalog", s.catalog)
	mux.HandleFunc("GET /api/materials", s.materials)
	mux.HandleFunc("GET /api/materials/{id}", s.material)
	mux.HandleFunc("GET /api/lessons/{id}", s.lesson)
	mux.HandleFunc("GET /api/boards/{id}", s.board)
	mux.HandleFunc("GET /api/sql/tasks", s.sqlTasks)
	mux.HandleFunc("GET /api/sql/schema", s.sqlSchema)
	mux.HandleFunc("POST /api/sql/run", s.sqlRun)
	mux.HandleFunc("POST /api/sql/check", s.sqlCheck)
	mux.HandleFunc("GET /api/explorer", s.explorer)
	mux.HandleFunc("GET /api/interview/arc", s.interviewArc)
	mux.HandleFunc("POST /api/interview/live/start", s.interviewStart)
	mux.HandleFunc("GET /api/interview/live", s.interviewGet)
	mux.HandleFunc("POST /api/interview/live/answer", s.interviewAnswer)
	mux.HandleFunc("POST /api/interview/live/finish", s.interviewFinish)
	mux.HandleFunc("GET /api/me/assessment", s.getAssessment)
	s.mountV1(mux)
	mux.HandleFunc("/", s.spa)
	return cors(mux)
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		if w.Header().Get("Access-Control-Allow-Origin") == "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Idempotency-Key, X-ACS-Timeout")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(204)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(body)
}

func writeErr(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, map[string]any{"error": map[string]any{"code": code, "message": msg}})
}

func readJSON(r *http.Request) (map[string]any, error) {
	defer r.Body.Close()
	b, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(b))) == 0 {
		return map[string]any{}, nil
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return map[string]any{"__parseError": true}, nil
	}
	return m, nil
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"ok": true, "product": "analyst-hall", "version": ledger.Version, "academy": true, "store": s.Store.Mode()})
}

func (s *Server) user(r *http.Request) *access.User {
	c, err := r.Cookie(cookieName)
	if err != nil || c.Value == "" {
		return nil
	}
	u, _ := s.Store.GetSessionUser(r.Context(), store.HashToken(c.Value))
	return u
}

func (s *Server) authBody(u *access.User) map[string]any {
	body := access.SessionFields(u)
	body["user"] = access.PublicUser(u)
	body["entitlements"] = access.EntitlementsOf(u)
	return body
}

func (s *Server) setSession(w http.ResponseWriter, u *access.User) {
	tok, _ := store.NewToken()
	exp := time.Now().Add(30 * 24 * time.Hour)
	_ = s.Store.CreateSession(context.Background(), u.ID, store.HashToken(tok), exp)
	http.SetCookie(w, s.cookie(tok, int(30*24*time.Hour.Seconds())))
	writeJSON(w, 200, s.authBody(u))
}

func (s *Server) cookie(val string, maxAge int) *http.Cookie {
	c := &http.Cookie{Name: cookieName, Value: val, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, MaxAge: maxAge}
	if s.Secure {
		c.Secure = true
	}
	return c
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	body, _ := readJSON(r)
	email := strings.ToLower(strings.TrimSpace(str(body["email"])))
	name := strings.TrimSpace(str(body["name"]))
	pw := str(body["password"])
	if !strings.Contains(email, "@") || len(pw) < 8 || len(name) < 2 {
		writeErr(w, 400, "VALIDATION", "Имя, email и пароль от 8 символов.")
		return
	}
	if existing, _ := s.Store.GetUserByEmail(r.Context(), email); existing != nil {
		writeErr(w, 409, "EXISTS", "Этот email уже в зале.")
		return
	}
	hash, err := store.HashPassword(pw)
	if err != nil {
		writeErr(w, 500, "HASH", "Не удалось сохранить пароль.")
		return
	}
	until := access.NewTrialUntil(time.Now())
	u, err := s.Store.CreateUser(r.Context(), email, name, hash, "student", "free", &until)
	if err != nil {
		writeErr(w, 500, "STORE", err.Error())
		return
	}
	s.setSession(w, u)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	body, _ := readJSON(r)
	email := strings.ToLower(strings.TrimSpace(str(body["email"])))
	u, _ := s.Store.GetUserByEmail(r.Context(), email)
	if u == nil || !store.CheckPassword(str(body["password"]), u.PasswordHash) {
		writeErr(w, 401, "AUTH", "Неверный email или пароль.")
		return
	}
	s.setSession(w, u)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(cookieName); err == nil {
		_ = s.Store.DeleteSession(r.Context(), store.HashToken(c.Value))
	}
	http.SetCookie(w, s.cookie("", -1))
	writeJSON(w, 200, map[string]any{"ok": true})
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.authBody(s.user(r)))
}

func (s *Server) getProgress(w http.ResponseWriter, r *http.Request) {
	u := s.user(r)
	if u == nil {
		writeErr(w, 401, "UNAUTH", "Войдите в зал.")
		return
	}
	p, _ := s.Store.GetProgress(r.Context(), u.ID)
	if p == nil {
		p = []byte("{}")
	}
	writeJSON(w, 200, map[string]any{"payload": json.RawMessage(p)})
}

func (s *Server) putProgress(w http.ResponseWriter, r *http.Request) {
	u := s.user(r)
	if u == nil {
		writeErr(w, 401, "UNAUTH", "Войдите в зал.")
		return
	}
	body, _ := readJSON(r)
	var payload any = body["payload"]
	if payload == nil {
		payload = body
	}
	raw, _ := json.Marshal(payload)
	_ = s.Store.PutProgress(r.Context(), u.ID, raw)
	writeJSON(w, 200, map[string]any{"ok": true})
}

func (s *Server) listPractice(w http.ResponseWriter, r *http.Request) {
	u := s.user(r)
	if u == nil {
		writeErr(w, 401, "UNAUTH", "Войдите в зал.")
		return
	}
	items, _ := s.Store.ListPractice(r.Context(), u.ID)
	writeJSON(w, 200, map[string]any{"items": items})
}

func (s *Server) postPractice(w http.ResponseWriter, r *http.Request) {
	u := s.user(r)
	if u == nil {
		writeErr(w, 401, "UNAUTH", "Войдите в зал.")
		return
	}
	body, _ := readJSON(r)
	kind := or(str(body["kind"]), "module")
	target := str(body["targetId"])
	title := strings.TrimSpace(str(body["title"]))
	text := strings.TrimSpace(str(body["body"]))
	if len(text) < 40 || title == "" {
		writeErr(w, 400, "VALIDATION", "Работа: заголовок и текст от 40 символов.")
		return
	}
	if kind == "module" {
		grade := strings.Split(target, "-")[0]
		if !access.CanGrade(u, grade) {
			writeErr(w, 403, "PLAN", "Этот контур закрыт на вашем тарифе.")
			return
		}
	}
	rec, err := s.Store.AddPractice(r.Context(), u.ID, kind, target, title, text)
	if err != nil {
		writeErr(w, 500, "STORE", err.Error())
		return
	}
	writeJSON(w, 201, map[string]any{"item": rec})
}

func (s *Server) adminUsers(w http.ResponseWriter, r *http.Request) {
	u := s.user(r)
	if u == nil || u.Role != "admin" {
		writeErr(w, 403, "PLAN", "Этот контур закрыт на вашем тарифе.")
		return
	}
	items, _ := s.Store.ListUsers(r.Context())
	writeJSON(w, 200, map[string]any{"items": items})
}

func (s *Server) adminPatchUser(w http.ResponseWriter, r *http.Request) {
	u := s.user(r)
	if u == nil || u.Role != "admin" {
		writeErr(w, 403, "PLAN", "Этот контур закрыт на вашем тарифе.")
		return
	}
	body, _ := readJSON(r)
	var plan, role *string
	if p := str(body["plan"]); p == "free" || p == "pro" {
		plan = &p
	}
	if rr := str(body["role"]); rr == "student" || rr == "admin" {
		role = &rr
	}
	next, _ := s.Store.UpdateUser(r.Context(), r.PathValue("id"), plan, role)
	if next == nil {
		writeErr(w, 404, "NOT_FOUND", "Нет такого ученика.")
		return
	}
	if plan != nil && *plan == "free" {
		_ = s.Store.SetTrialUntil(r.Context(), next.ID, nil)
		next, _ = s.Store.GetUserByID(r.Context(), next.ID)
	}
	writeJSON(w, 200, s.authBody(next))
}

func (s *Server) adminPractice(w http.ResponseWriter, r *http.Request) {
	u := s.user(r)
	if u == nil || u.Role != "admin" {
		writeErr(w, 403, "PLAN", "Этот контур закрыт на вашем тарифе.")
		return
	}
	items, _ := s.Store.ListAllPractice(r.Context())
	writeJSON(w, 200, map[string]any{"items": items})
}

func (s *Server) adminReview(w http.ResponseWriter, r *http.Request) {
	u := s.user(r)
	if u == nil || u.Role != "admin" {
		writeErr(w, 403, "PLAN", "Этот контур закрыт на вашем тарифе.")
		return
	}
	body, _ := readJSON(r)
	st := "submitted"
	if str(body["status"]) == "reviewed" {
		st = "reviewed"
	}
	var score *int
	if v, ok := body["score"].(float64); ok {
		i := int(v)
		score = &i
	}
	note := str(body["note"])
	if note == "" {
		note = str(body["reviewer_note"])
	}
	rec, _ := s.Store.ReviewPractice(r.Context(), r.PathValue("id"), st, score, note)
	writeJSON(w, 200, map[string]any{"item": rec})
}

func (s *Server) catalog(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.Cat.Public())
}

func (s *Server) materials(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"topics": s.Cat.Materials})
}

func (s *Server) material(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	m, ok := s.Cat.Material(id)
	if !ok {
		writeErr(w, 404, "NOT_FOUND", "Статья не найдена.")
		return
	}
	writeJSON(w, 200, m)
}

func (s *Server) lesson(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	l, ok := s.Cat.Lesson(id)
	if !ok {
		writeErr(w, 404, "NOT_FOUND", "Урок не найден.")
		return
	}
	u := s.user(r)
	if !access.CanGrade(u, l.GradeID) {
		writeJSON(w, 403, map[string]any{
			"error":   map[string]any{"code": "PLAN", "message": "Этот грейд в PRO. Intern остаётся открытым."},
			"paywall": true,
			"trialEnded": access.TrialUntilISO(u) != nil && !access.IsPro(u),
			"lesson": map[string]any{"id": l.ID, "title": l.Title, "gradeId": l.GradeID},
		})
		return
	}
	writeJSON(w, 200, l)
}

func (s *Server) board(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	u := s.user(r)
	for _, b := range s.Cat.Boards {
		if str(b["id"]) == id {
			plan := str(b["plan"])
			if plan == "pro" && !access.CanStudio(u, id) {
				writeErr(w, 403, "PLAN", "Эта доска в PRO.")
				return
			}
			writeJSON(w, 200, b)
			return
		}
	}
	writeErr(w, 404, "NOT_FOUND", "Доска не найдена.")
}

func (s *Server) gateWrite(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
		return true
	}
	u := s.user(r)
	if u == nil {
		writeErr(w, 401, "UNAUTH", "Войдите в зал.")
		return false
	}
	if !access.CanPetWrite(u) {
		writeErr(w, 403, "PLAN", "Запись в учебный банк — в PRO. На бесплатном тарифе можно читать справочник и Intern.")
		return false
	}
	return true
}

func (s *Server) spa(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api") {
		writeErr(w, 404, "NO_ROUTE", "Нет "+r.Method+" "+r.URL.Path)
		return
	}
	if s.WebDir == "" {
		writeErr(w, 404, "NO_WEB", "web/dist не собран")
		return
	}
	rel := strings.TrimPrefix(r.URL.Path, "/")
	if rel == "" {
		rel = "index.html"
	}
	target := filepath.Join(s.WebDir, filepath.FromSlash(rel))
	if !strings.HasPrefix(target, filepath.Clean(s.WebDir)) {
		http.Error(w, "bad path", 400)
		return
	}
	if st, err := os.Stat(target); err == nil && !st.IsDir() {
		http.ServeFile(w, r, target)
		return
	}
	http.ServeFile(w, r, filepath.Join(s.WebDir, "index.html"))
}

func str(v any) string {
	s, _ := v.(string)
	return s
}

func or(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func header(r *http.Request, name string) string {
	return r.Header.Get(name)
}

func (s *Server) bankJSON(w http.ResponseWriter, r *http.Request, status int, body any, err error) {
	path := r.URL.Path
	key := header(r, "Idempotency-Key")
	if err != nil {
		var fe *ledger.FailError
		if errors.As(err, &fe) {
			s.Bank.Note(r.Method, path, fe.Status, fe.Code, key)
			out := map[string]any{"error": map[string]any{"code": fe.Code, "message": fe.Message}}
			for k, v := range fe.Extra {
				out[k] = v
			}
			st := fe.Status
			if s.Bank.ContractErrorInBody() && st >= 400 && st < 500 {
				out["ok"] = false
				st = 200
			}
			writeJSON(w, st, out)
			return
		}
		writeErr(w, 500, "CRASH", err.Error())
		return
	}
	s.Bank.Note(r.Method, path, status, "ok", key)
	writeJSON(w, status, body)
}
