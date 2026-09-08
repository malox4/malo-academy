package store

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/malox4/malo-academy/internal/access"
)

//go:embed schema.sql
var schemaSQL string

type Store struct {
	mode string
	pool *pgxpool.Pool
	mu   sync.Mutex
	file string
	mem  memory
}

type memory struct {
	Users    []memUser    `json:"users"`
	Sessions []memSession `json:"sessions"`
	Progress []memProg    `json:"progress"`
	Practice []memPractice `json:"practice"`
	Assess   []memAssess   `json:"assessments"`
}

type memAssess struct {
	UserID  string          `json:"user_id"`
	Payload json.RawMessage `json:"payload"`
}

type memUser struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
	Plan         string `json:"plan"`
	CreatedAt    string `json:"created_at"`
	TrialUntil   string `json:"trial_until,omitempty"`
}

type memSession struct {
	TokenHash string `json:"token_hash"`
	UserID    string `json:"user_id"`
	ExpiresAt string `json:"expires_at"`
}

type memProg struct {
	UserID  string          `json:"user_id"`
	Payload json.RawMessage `json:"payload"`
}

type memPractice struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	Kind         string  `json:"kind"`
	TargetID     string  `json:"target_id"`
	Title        string  `json:"title"`
	Body         string  `json:"body"`
	Status       string  `json:"status"`
	Score        *int    `json:"score"`
	ReviewerNote *string `json:"reviewer_note"`
	CreatedAt    string  `json:"created_at"`
	Email        string  `json:"email,omitempty"`
	Author       string  `json:"author,omitempty"`
}

func Open(ctx context.Context, databaseURL, dataDir string) (*Store, error) {
	_ = os.MkdirAll(dataDir, 0o755)
	s := &Store{file: filepath.Join(dataDir, "academy-users.json")}
	if databaseURL != "" {
		pool, err := pgxpool.New(ctx, databaseURL)
		if err != nil {
			return nil, err
		}
		for _, stmt := range strings.Split(schemaSQL, ";") {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := pool.Exec(ctx, stmt); err != nil {
				pool.Close()
				return nil, err
			}
		}
		s.pool = pool
		s.mode = "pg"
		fmt.Println("Academy store: postgres")
		return s, nil
	}
	s.mode = "file"
	s.loadFile()
	fmt.Println("Academy store: JSON file (set DATABASE_URL for Postgres)")
	return s, nil
}

func (s *Store) Mode() string { return s.mode }

func (s *Store) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}

func (s *Store) loadFile() {
	b, err := os.ReadFile(s.file)
	if err != nil {
		s.mem = memory{}
		return
	}
	_ = json.Unmarshal(b, &s.mem)
}

func (s *Store) saveFile() {
	_ = os.MkdirAll(filepath.Dir(s.file), 0o755)
	b, _ := json.MarshalIndent(s.mem, "", "  ")
	_ = os.WriteFile(s.file, b, 0o644)
}

func HashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPassword(pw, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func NewToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func toUser(u memUser) *access.User {
	usr := &access.User{
		ID: u.ID, Email: u.Email, Name: u.Name, Role: u.Role, Plan: u.Plan,
		CreatedAt: u.CreatedAt, PasswordHash: u.PasswordHash,
	}
	if u.TrialUntil != "" {
		if t, err := time.Parse(time.RFC3339, u.TrialUntil); err == nil {
			tt := t.UTC()
			usr.TrialUntil = &tt
		}
	}
	return usr
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*access.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if s.pool != nil {
		row := s.pool.QueryRow(ctx, `SELECT id, email, name, password_hash, role, plan, created_at, trial_until FROM users WHERE email=$1`, email)
		return scanUser(row)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range s.mem.Users {
		if u.Email == email {
			return toUser(u), nil
		}
	}
	return nil, nil
}

func (s *Store) GetUserByID(ctx context.Context, id string) (*access.User, error) {
	if s.pool != nil {
		row := s.pool.QueryRow(ctx, `SELECT id, email, name, password_hash, role, plan, created_at, trial_until FROM users WHERE id=$1`, id)
		return scanUser(row)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range s.mem.Users {
		if u.ID == id {
			return toUser(u), nil
		}
	}
	return nil, nil
}

type scanR interface {
	Scan(dest ...any) error
}

func scanUser(row scanR) (*access.User, error) {
	var u access.User
	var created time.Time
	var trial sql.NullTime
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.Role, &u.Plan, &created, &trial)
	if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	u.CreatedAt = created.UTC().Format(time.RFC3339)
	if trial.Valid {
		t := trial.Time.UTC()
		u.TrialUntil = &t
	}
	return &u, nil
}

func (s *Store) CreateUser(ctx context.Context, email, name, passwordHash, role, plan string, trialUntil *time.Time) (*access.User, error) {
	if role == "" {
		role = "student"
	}
	if plan == "" {
		plan = "free"
	}
	trialStr := ""
	var trial any
	if trialUntil != nil && !trialUntil.IsZero() {
		t := trialUntil.UTC()
		trialStr = t.Format(time.RFC3339)
		trial = t
	}
	u := memUser{
		ID: uuid.NewString(), Email: strings.ToLower(email), Name: name,
		PasswordHash: passwordHash, Role: role, Plan: plan, CreatedAt: time.Now().UTC().Format(time.RFC3339),
		TrialUntil: trialStr,
	}
	if s.pool != nil {
		_, err := s.pool.Exec(ctx, `INSERT INTO users (id,email,name,password_hash,role,plan,created_at,trial_until) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
			u.ID, u.Email, u.Name, u.PasswordHash, u.Role, u.Plan, u.CreatedAt, trial)
		if err != nil {
			return nil, err
		}
		return toUser(u), nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mem.Users = append(s.mem.Users, u)
	s.saveFile()
	return toUser(u), nil
}

func (s *Store) SetPassword(ctx context.Context, id, passwordHash string) error {
	if s.pool != nil {
		_, err := s.pool.Exec(ctx, `UPDATE users SET password_hash=$1 WHERE id=$2`, passwordHash, id)
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, u := range s.mem.Users {
		if u.ID == id {
			s.mem.Users[i].PasswordHash = passwordHash
			s.saveFile()
			return nil
		}
	}
	return nil
}

func (s *Store) UpdateUser(ctx context.Context, id string, plan, role *string) (*access.User, error) {
	if s.pool != nil {
		u, err := s.GetUserByID(ctx, id)
		if err != nil || u == nil {
			return u, err
		}
		p, r := u.Plan, u.Role
		if plan != nil {
			p = *plan
		}
		if role != nil {
			r = *role
		}
		_, err = s.pool.Exec(ctx, `UPDATE users SET plan=$1, role=$2 WHERE id=$3`, p, r, id)
		if err != nil {
			return nil, err
		}
		return s.GetUserByID(ctx, id)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, u := range s.mem.Users {
		if u.ID == id {
			if plan != nil {
				s.mem.Users[i].Plan = *plan
			}
			if role != nil {
				s.mem.Users[i].Role = *role
			}
			s.saveFile()
			return toUser(s.mem.Users[i]), nil
		}
	}
	return nil, nil
}

func (s *Store) ListUsers(ctx context.Context) ([]map[string]any, error) {
	out := []map[string]any{}
	if s.pool != nil {
		rows, err := s.pool.Query(ctx, `SELECT id,email,name,role,plan,created_at,trial_until FROM users ORDER BY created_at DESC`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var id, email, name, role, plan string
			var created time.Time
			var trial sql.NullTime
			if err := rows.Scan(&id, &email, &name, &role, &plan, &created, &trial); err != nil {
				return nil, err
			}
			row := map[string]any{"id": id, "email": email, "name": name, "role": role, "plan": plan, "created_at": created.UTC().Format(time.RFC3339), "trial_until": nil}
			if trial.Valid {
				row["trial_until"] = trial.Time.UTC().Format(time.RFC3339)
			}
			out = append(out, row)
		}
		return out, rows.Err()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := len(s.mem.Users) - 1; i >= 0; i-- {
		u := s.mem.Users[i]
		row := map[string]any{"id": u.ID, "email": u.Email, "name": u.Name, "role": u.Role, "plan": u.Plan, "created_at": u.CreatedAt, "trial_until": nil}
		if u.TrialUntil != "" {
			row["trial_until"] = u.TrialUntil
		}
		out = append(out, row)
	}
	return out, nil
}

func (s *Store) SetTrialUntil(ctx context.Context, id string, until *time.Time) error {
	var trial any
	trialStr := ""
	if until != nil && !until.IsZero() {
		t := until.UTC()
		trial = t
		trialStr = t.Format(time.RFC3339)
	}
	if s.pool != nil {
		_, err := s.pool.Exec(ctx, `UPDATE users SET trial_until=$1 WHERE id=$2`, trial, id)
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, u := range s.mem.Users {
		if u.ID == id {
			s.mem.Users[i].TrialUntil = trialStr
			s.saveFile()
			return nil
		}
	}
	return nil
}

func (s *Store) CreateSession(ctx context.Context, userID, tokenHash string, expires time.Time) error {
	if s.pool != nil {
		_, err := s.pool.Exec(ctx, `INSERT INTO sessions (token_hash,user_id,expires_at) VALUES ($1,$2,$3)`, tokenHash, userID, expires)
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mem.Sessions = append(s.mem.Sessions, memSession{TokenHash: tokenHash, UserID: userID, ExpiresAt: expires.UTC().Format(time.RFC3339)})
	s.saveFile()
	return nil
}

func (s *Store) GetSessionUser(ctx context.Context, tokenHash string) (*access.User, error) {
	if s.pool != nil {
		var uid string
		err := s.pool.QueryRow(ctx, `SELECT user_id FROM sessions WHERE token_hash=$1 AND expires_at>now()`, tokenHash).Scan(&uid)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		return s.GetUserByID(ctx, uid)
	}
	s.mu.Lock()
	now := time.Now()
	var uid string
	for _, sess := range s.mem.Sessions {
		if sess.TokenHash != tokenHash {
			continue
		}
		exp, _ := time.Parse(time.RFC3339, sess.ExpiresAt)
		if exp.After(now) {
			uid = sess.UserID
		}
		break
	}
	s.mu.Unlock()
	if uid == "" {
		return nil, nil
	}
	return s.GetUserByID(ctx, uid)
}

func (s *Store) DeleteSession(ctx context.Context, tokenHash string) error {
	if s.pool != nil {
		_, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE token_hash=$1`, tokenHash)
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	next := s.mem.Sessions[:0]
	for _, sess := range s.mem.Sessions {
		if sess.TokenHash != tokenHash {
			next = append(next, sess)
		}
	}
	s.mem.Sessions = next
	s.saveFile()
	return nil
}

func (s *Store) GetProgress(ctx context.Context, userID string) (json.RawMessage, error) {
	if s.pool != nil {
		var payload []byte
		err := s.pool.QueryRow(ctx, `SELECT payload FROM progress WHERE user_id=$1`, userID).Scan(&payload)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return payload, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, p := range s.mem.Progress {
		if p.UserID == userID {
			return p.Payload, nil
		}
	}
	return nil, nil
}

func (s *Store) PutProgress(ctx context.Context, userID string, payload json.RawMessage) error {
	if payload == nil {
		payload = []byte("{}")
	}
	if s.pool != nil {
		_, err := s.pool.Exec(ctx, `INSERT INTO progress (user_id, payload, updated_at) VALUES ($1,$2::jsonb, now())
			ON CONFLICT (user_id) DO UPDATE SET payload=EXCLUDED.payload, updated_at=now()`, userID, payload)
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, p := range s.mem.Progress {
		if p.UserID == userID {
			s.mem.Progress[i].Payload = payload
			s.saveFile()
			return nil
		}
	}
	s.mem.Progress = append(s.mem.Progress, memProg{UserID: userID, Payload: payload})
	s.saveFile()
	return nil
}

func (s *Store) AddPractice(ctx context.Context, userID, kind, targetID, title, body string) (map[string]any, error) {
	rec := memPractice{
		ID: uuid.NewString(), UserID: userID, Kind: kind, TargetID: targetID, Title: title, Body: body,
		Status: "submitted", CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if s.pool != nil {
		_, err := s.pool.Exec(ctx, `INSERT INTO practice_submissions
			(id,user_id,kind,target_id,title,body,status,score,reviewer_note,created_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
			rec.ID, rec.UserID, rec.Kind, rec.TargetID, rec.Title, rec.Body, rec.Status, rec.Score, rec.ReviewerNote, rec.CreatedAt)
		if err != nil {
			return nil, err
		}
		return practiceMap(rec), nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mem.Practice = append([]memPractice{rec}, s.mem.Practice...)
	s.saveFile()
	return practiceMap(rec), nil
}

func practiceMap(p memPractice) map[string]any {
	return map[string]any{
		"id": p.ID, "user_id": p.UserID, "kind": p.Kind, "target_id": p.TargetID,
		"title": p.Title, "body": p.Body, "status": p.Status, "score": p.Score,
		"reviewer_note": p.ReviewerNote, "created_at": p.CreatedAt, "email": p.Email, "author": p.Author,
	}
}

func (s *Store) ListPractice(ctx context.Context, userID string) ([]map[string]any, error) {
	out := []map[string]any{}
	if s.pool != nil {
		rows, err := s.pool.Query(ctx, `SELECT id,user_id,kind,target_id,title,body,status,score,reviewer_note,created_at
			FROM practice_submissions WHERE user_id=$1 ORDER BY created_at DESC LIMIT 80`, userID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			p, err := scanPractice(rows)
			if err != nil {
				return nil, err
			}
			out = append(out, practiceMap(p))
		}
		return out, rows.Err()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for _, p := range s.mem.Practice {
		if p.UserID != userID {
			continue
		}
		out = append(out, practiceMap(p))
		n++
		if n >= 80 {
			break
		}
	}
	return out, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanPractice(row rowScanner) (memPractice, error) {
	var p memPractice
	var score sql.NullInt64
	var note sql.NullString
	var created time.Time
	err := row.Scan(&p.ID, &p.UserID, &p.Kind, &p.TargetID, &p.Title, &p.Body, &p.Status, &score, &note, &created)
	if err != nil {
		return p, err
	}
	if score.Valid {
		v := int(score.Int64)
		p.Score = &v
	}
	if note.Valid {
		p.ReviewerNote = &note.String
	}
	p.CreatedAt = created.UTC().Format(time.RFC3339)
	return p, nil
}

func (s *Store) ListAllPractice(ctx context.Context) ([]map[string]any, error) {
	out := []map[string]any{}
	if s.pool != nil {
		rows, err := s.pool.Query(ctx, `SELECT p.id,p.user_id,p.kind,p.target_id,p.title,p.body,p.status,p.score,p.reviewer_note,p.created_at,u.email,u.name
			FROM practice_submissions p JOIN users u ON u.id=p.user_id ORDER BY p.created_at DESC LIMIT 200`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var p memPractice
			var score sql.NullInt64
			var note sql.NullString
			var created time.Time
			if err := rows.Scan(&p.ID, &p.UserID, &p.Kind, &p.TargetID, &p.Title, &p.Body, &p.Status, &score, &note, &created, &p.Email, &p.Author); err != nil {
				return nil, err
			}
			if score.Valid {
				v := int(score.Int64)
				p.Score = &v
			}
			if note.Valid {
				p.ReviewerNote = &note.String
			}
			p.CreatedAt = created.UTC().Format(time.RFC3339)
			out = append(out, practiceMap(p))
		}
		return out, rows.Err()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	limit := 200
	for i, p := range s.mem.Practice {
		if i >= limit {
			break
		}
		for _, u := range s.mem.Users {
			if u.ID == p.UserID {
				p.Email = u.Email
				p.Author = u.Name
				break
			}
		}
		out = append(out, practiceMap(p))
	}
	return out, nil
}

func (s *Store) ReviewPractice(ctx context.Context, id, status string, score *int, note string) (map[string]any, error) {
	if s.pool != nil {
		_, err := s.pool.Exec(ctx, `UPDATE practice_submissions SET status=$1, score=$2, reviewer_note=$3 WHERE id=$4`, status, score, note, id)
		if err != nil {
			return nil, err
		}
		row := s.pool.QueryRow(ctx, `SELECT id,user_id,kind,target_id,title,body,status,score,reviewer_note,created_at FROM practice_submissions WHERE id=$1`, id)
		p, err := scanPractice(row)
		if err != nil {
			return nil, nil
		}
		return practiceMap(p), nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, p := range s.mem.Practice {
		if p.ID == id {
			s.mem.Practice[i].Status = status
			s.mem.Practice[i].Score = score
			s.mem.Practice[i].ReviewerNote = &note
			s.saveFile()
			return practiceMap(s.mem.Practice[i]), nil
		}
	}
	return nil, nil
}

func (s *Store) GetAssessment(ctx context.Context, userID string) (json.RawMessage, error) {
	if s.pool != nil {
		var payload []byte
		err := s.pool.QueryRow(ctx, `SELECT payload FROM assessments WHERE user_id=$1`, userID).Scan(&payload)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return payload, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, a := range s.mem.Assess {
		if a.UserID == userID {
			return a.Payload, nil
		}
	}
	return nil, nil
}

func (s *Store) PutAssessment(ctx context.Context, userID string, payload json.RawMessage) error {
	if payload == nil {
		payload = []byte("{}")
	}
	if s.pool != nil {
		_, err := s.pool.Exec(ctx, `INSERT INTO assessments (user_id, payload, updated_at) VALUES ($1,$2::jsonb, now())
			ON CONFLICT (user_id) DO UPDATE SET payload=EXCLUDED.payload, updated_at=now()`, userID, payload)
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, a := range s.mem.Assess {
		if a.UserID == userID {
			s.mem.Assess[i].Payload = payload
			s.saveFile()
			return nil
		}
	}
	s.mem.Assess = append(s.mem.Assess, memAssess{UserID: userID, Payload: payload})
	s.saveFile()
	return nil
}
