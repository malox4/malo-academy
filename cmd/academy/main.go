package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/malox4/malo-academy/internal/content"
	"github.com/malox4/malo-academy/internal/httpapi"
	"github.com/malox4/malo-academy/internal/ledger"
	"github.com/malox4/malo-academy/internal/petdb"
	"github.com/malox4/malo-academy/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	root, _ := os.Getwd()
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = filepath.Join(root, "data")
	}
	webDir := os.Getenv("WEB_DIR")
	if webDir == "" {
		webDir = filepath.Join(root, "web", "dist")
	}

	ctx := context.Background()
	st, err := store.Open(ctx, os.Getenv("DATABASE_URL"), dataDir)
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()
	if err := seedAdmin(ctx, st); err != nil {
		log.Fatal(err)
	}
	if err := seedPro(ctx, st); err != nil {
		log.Fatal(err)
	}

	cat, err := content.Load()
	if err != nil {
		log.Fatal(err)
	}

	bank := ledger.Open(dataDir)
	lab, err := petdb.Open(dataDir)
	if err != nil {
		log.Fatal(err)
	}
	defer lab.Close()

	srv := httpapi.New(st, bank, webDir, cat, lab)

	addr := host + ":" + port
	fmt.Printf("Analyst Hall http://%s  (web %s, store %s)\n", addr, webDir, st.Mode())
	log.Fatal(http.ListenAndServe(addr, srv.Handler()))
}

func seedAdmin(ctx context.Context, st *store.Store) error {
	email := os.Getenv("ADMIN_EMAIL")
	if email == "" {
		email = "admin@malo.academy"
	}
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		password = "ChangeMe_Admin1!"
	}
	existing, err := st.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if existing != nil {
		if store.CheckPassword(password, existing.PasswordHash) {
			return nil
		}
		hash, err := store.HashPassword(password)
		if err != nil {
			return err
		}
		return st.SetPassword(ctx, existing.ID, hash)
	}
	hash, err := store.HashPassword(password)
	if err != nil {
		return err
	}
	_, err = st.CreateUser(ctx, email, "Хозяин зала", hash, "admin", "pro", nil)
	if err == nil {
		fmt.Println("Seeded admin", email)
	}
	return err
}

func seedPro(ctx context.Context, st *store.Store) error {
	email := os.Getenv("PRO_EMAIL")
	if email == "" {
		email = "pro@malo.academy"
	}
	password := os.Getenv("PRO_PASSWORD")
	if password == "" {
		password = "ChangeMe_Pro1!"
	}
	existing, err := st.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if existing != nil {
		plan := "pro"
		if existing.Plan != "pro" {
			if _, err := st.UpdateUser(ctx, existing.ID, &plan, nil); err != nil {
				return err
			}
		}
		if err := st.SetTrialUntil(ctx, existing.ID, nil); err != nil {
			return err
		}
		if store.CheckPassword(password, existing.PasswordHash) {
			return nil
		}
		hash, err := store.HashPassword(password)
		if err != nil {
			return err
		}
		return st.SetPassword(ctx, existing.ID, hash)
	}
	hash, err := store.HashPassword(password)
	if err != nil {
		return err
	}
	_, err = st.CreateUser(ctx, email, "Аналитик PRO", hash, "student", "pro", nil)
	if err == nil {
		fmt.Println("Seeded pro", email)
	}
	return err
}
