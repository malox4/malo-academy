package httpapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"testing"
	"time"
)

func TestTelegramInitData(t *testing.T) {
	token := "123456:ABC"
	user := `{"id":42,"first_name":"Анна"}`
	auth := fmt.Sprintf("%d", time.Now().Unix())
	vals := url.Values{"auth_date": {auth}, "user": {user}}
	check := "auth_date=" + auth + "\nuser=" + user
	mac := hmac.New(sha256.New, telegramSecret(token))
	mac.Write([]byte(check))
	vals.Set("hash", hex.EncodeToString(mac.Sum(nil)))
	got, err := telegramUserFromInit(vals.Encode(), token)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 42 || got.First != "Анна" {
		t.Fatalf("%+v", got)
	}
	if _, err := telegramUserFromInit(vals.Encode(), "other"); err == nil {
		t.Fatal("bad token must fail")
	}
}
