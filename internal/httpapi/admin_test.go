package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/malox4/malo-academy/internal/store"
)

func TestAdminGrantAccess(t *testing.T) {
	st, err := store.Open(context.Background(), "", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(st.Close)

	adminHash, _ := store.HashPassword("ChangeMe_Admin1!")
	studentHash, _ := store.HashPassword("Student_Pass1!")
	admin, err := st.CreateUser(context.Background(), "admin@malo.academy", "Хозяин", adminHash, "admin", "pro", nil)
	if err != nil {
		t.Fatal(err)
	}
	student, err := st.CreateUser(context.Background(), "a@b.co", "Аня", studentHash, "student", "free", nil)
	if err != nil {
		t.Fatal(err)
	}

	h := New(st, nil, "", nil, nil).Handler()

	login := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(`{"email":"admin@malo.academy","password":"ChangeMe_Admin1!"}`))
	login.Header.Set("Content-Type", "application/json")
	lw := httptest.NewRecorder()
	h.ServeHTTP(lw, login)
	if lw.Code != 200 {
		t.Fatalf("login %d %s", lw.Code, lw.Body.String())
	}
	cookie := lw.Result().Cookies()[0]

	list := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
	list.AddCookie(cookie)
	lr := httptest.NewRecorder()
	h.ServeHTTP(lr, list)
	if lr.Code != 200 {
		t.Fatalf("list %d %s", lr.Code, lr.Body.String())
	}
	var listed struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.Unmarshal(lr.Body.Bytes(), &listed); err != nil {
		t.Fatal(err)
	}
	if len(listed.Items) != 2 {
		t.Fatalf("want 2 users, got %d", len(listed.Items))
	}

	_ = st.PutProgress(context.Background(), student.ID, []byte(`{"lessons":{"intern-1-profession":{"done":true,"at":1700000000000}},"interview":{"hold":{"at":1}}}`))
	_ = st.TouchPresence(context.Background(), student.ID, "/lesson/intern-1-profession", "")
	list2 := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
	list2.AddCookie(cookie)
	lr2 := httptest.NewRecorder()
	h.ServeHTTP(lr2, list2)
	if err := json.Unmarshal(lr2.Body.Bytes(), &listed); err != nil {
		t.Fatal(err)
	}
	var anya map[string]any
	for _, it := range listed.Items {
		if it["email"] == "a@b.co" {
			anya = it
		}
	}
	if anya == nil {
		t.Fatal("missing student in desk")
	}
	if n, _ := anya["lessonsDone"].(float64); n < 1 {
		t.Fatalf("progress %+v", anya)
	}
	if str, _ := anya["herePath"].(string); str != "/lesson/intern-1-profession" {
		t.Fatalf("here %+v", anya)
	}

	patch := httptest.NewRequest(http.MethodPatch, "/api/admin/users/"+student.ID, bytes.NewBufferString(`{"trialDays":3}`))
	patch.Header.Set("Content-Type", "application/json")
	patch.AddCookie(cookie)
	pw := httptest.NewRecorder()
	h.ServeHTTP(pw, patch)
	if pw.Code != 200 {
		t.Fatalf("trial %d %s", pw.Code, pw.Body.String())
	}
	var trialBody struct {
		User map[string]any `json:"user"`
	}
	_ = json.Unmarshal(pw.Body.Bytes(), &trialBody)
	if trialBody.User["isPro"] != true || trialBody.User["trialActive"] != true {
		t.Fatalf("trial grant %+v", trialBody.User)
	}

	pro := httptest.NewRequest(http.MethodPatch, "/api/admin/users/"+student.ID, bytes.NewBufferString(`{"plan":"pro"}`))
	pro.Header.Set("Content-Type", "application/json")
	pro.AddCookie(cookie)
	pr := httptest.NewRecorder()
	h.ServeHTTP(pr, pro)
	var proBody struct {
		User map[string]any `json:"user"`
	}
	_ = json.Unmarshal(pr.Body.Bytes(), &proBody)
	if proBody.User["plan"] != "pro" || proBody.User["trialActive"] != false {
		t.Fatalf("pro grant %+v", proBody.User)
	}

	self := httptest.NewRequest(http.MethodPatch, "/api/admin/users/"+admin.ID, bytes.NewBufferString(`{"role":"student"}`))
	self.Header.Set("Content-Type", "application/json")
	self.AddCookie(cookie)
	sw := httptest.NewRecorder()
	h.ServeHTTP(sw, self)
	if sw.Code != 400 {
		t.Fatalf("self demote want 400 got %d %s", sw.Code, sw.Body.String())
	}

	stuLogin := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(`{"email":"a@b.co","password":"Student_Pass1!"}`))
	stuLogin.Header.Set("Content-Type", "application/json")
	slw := httptest.NewRecorder()
	h.ServeHTTP(slw, stuLogin)
	stuCookie := slw.Result().Cookies()[0]
	forbid := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
	forbid.AddCookie(stuCookie)
	fw := httptest.NewRecorder()
	h.ServeHTTP(fw, forbid)
	if fw.Code != 403 {
		t.Fatalf("student admin list want 403 got %d", fw.Code)
	}
}
