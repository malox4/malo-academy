package httpapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/malox4/malo-academy/internal/access"
	"github.com/malox4/malo-academy/internal/store"
)

func (s *Server) telegramAuth(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if token == "" {
		writeErr(w, 503, "TELEGRAM", "Mini App вход не настроен. Задайте TELEGRAM_BOT_TOKEN.")
		return
	}
	body, _ := readJSON(r)
	initData := str(body["initData"])
	if initData == "" {
		writeErr(w, 400, "VALIDATION", "Нет initData.")
		return
	}
	user, err := telegramUserFromInit(initData, token)
	if err != nil {
		writeErr(w, 401, "AUTH", "Telegram не подтвердил вход.")
		return
	}
	existing, err := s.Store.GetUserByTelegramID(r.Context(), user.ID)
	if err != nil {
		writeErr(w, 500, "STORE", err.Error())
		return
	}
	if existing != nil {
		s.setSession(w, existing)
		return
	}
	name := strings.TrimSpace(user.First + " " + user.Last)
	if name == "" {
		name = user.Username
	}
	if name == "" {
		name = "Telegram"
	}
	email := fmt.Sprintf("tg%d@telegram.user", user.ID)
	pw, _ := store.NewToken()
	hash, err := store.HashPassword(pw)
	if err != nil {
		writeErr(w, 500, "HASH", "Не удалось сохранить пароль.")
		return
	}
	until := access.NewTrialUntil(time.Now())
	created, err := s.Store.CreateUser(r.Context(), email, name, hash, "student", "free", &until)
	if err != nil {
		if reused, _ := s.Store.GetUserByEmail(r.Context(), email); reused != nil {
			_ = s.Store.BindTelegram(r.Context(), reused.ID, user.ID)
			s.setSession(w, reused)
			return
		}
		writeErr(w, 500, "STORE", err.Error())
		return
	}
	_ = s.Store.BindTelegram(r.Context(), created.ID, user.ID)
	s.setSession(w, created)
}

type tgUser struct {
	ID       int64  `json:"id"`
	First    string `json:"first_name"`
	Last     string `json:"last_name"`
	Username string `json:"username"`
}

func telegramUserFromInit(initData, botToken string) (tgUser, error) {
	var zero tgUser
	vals, err := url.ParseQuery(initData)
	if err != nil {
		return zero, err
	}
	hash := vals.Get("hash")
	if hash == "" {
		return zero, fmt.Errorf("no hash")
	}
	vals.Del("hash")
	keys := make([]string, 0, len(vals))
	for k := range vals {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for i, k := range keys {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(vals.Get(k))
	}
	mac := hmac.New(sha256.New, telegramSecret(botToken))
	mac.Write([]byte(b.String()))
	if !hmac.Equal([]byte(hex.EncodeToString(mac.Sum(nil))), []byte(hash)) {
		return zero, fmt.Errorf("bad hash")
	}
	auth, _ := strconv.ParseInt(vals.Get("auth_date"), 10, 64)
	if auth == 0 || time.Since(time.Unix(auth, 0)) > 24*time.Hour {
		return zero, fmt.Errorf("stale")
	}
	var u tgUser
	if err := json.Unmarshal([]byte(vals.Get("user")), &u); err != nil || u.ID == 0 {
		return zero, fmt.Errorf("no user")
	}
	return u, nil
}

func telegramSecret(botToken string) []byte {
	mac := hmac.New(sha256.New, []byte("WebAppData"))
	mac.Write([]byte(botToken))
	return mac.Sum(nil)
}
