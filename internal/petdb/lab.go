package petdb

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type Lab struct {
	mu     sync.Mutex
	file   string
	tables map[string][]map[string]any
}

func Open(dataDir string) (*Lab, error) {
	_ = os.MkdirAll(dataDir, 0o755)
	file := filepath.Join(dataDir, "malo-core.json")
	l := &Lab{file: file, tables: map[string][]map[string]any{}}
	if raw, err := os.ReadFile(file); err == nil {
		_ = json.Unmarshal(raw, &l.tables)
	}
	return l, nil
}

func (l *Lab) Close() {}

func (l *Lab) File() string { return l.file }

func (l *Lab) persist() {
	raw, _ := json.MarshalIndent(l.tables, "", "  ")
	_ = os.WriteFile(l.file, raw, 0o644)
}

func (l *Lab) Sync(dump map[string]any) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	t := map[string][]map[string]any{}
	put := func(name string, rows []map[string]any) { t[name] = rows }
	put("customers", project(asMaps(dump["customers"]), []string{"id", "name", "segment", "city", "opened_at"}, map[string]string{"openedAt": "opened_at"}))
	put("kyc", project(asMaps(dump["kyc"]), []string{"id", "customer_id", "status", "reason"}, map[string]string{"customer": "customer_id"}))
	put("merchants", project(asMaps(dump["merchants"]), []string{"id", "account_id", "name", "mcc", "city"}, map[string]string{"accountId": "account_id"}))
	put("accounts", project(asMaps(dump["accounts"]), []string{"id", "wallet_id", "iban", "name", "type", "normal", "currency", "ledger", "hold", "available"}, map[string]string{"walletId": "wallet_id"}))
	put("wallets", project(asMaps(dump["wallets"]), []string{"id", "customer_id", "account_id", "name", "currency", "status", "ledger", "hold", "available"}, map[string]string{"customer": "customer_id", "accountId": "account_id"}))
	put("holds", project(asMaps(dump["holds"]), []string{"id", "account_id", "wallet_id", "payment_id", "amount", "currency", "status", "source", "created_at"}, map[string]string{"accountId": "account_id", "walletId": "wallet_id", "paymentId": "payment_id", "createdAt": "created_at"}))
	jrn := []map[string]any{}
	for _, m := range asMaps(dump["journals"]) {
		jrn = append(jrn, map[string]any{
			"id": m["id"], "kind": m["kind"], "ref": m["ref"], "debit": num(m["debit"]), "credit": num(m["credit"]),
			"balanced": bool01(m["balanced"], true), "at": str(m["at"]), "note": m["note"],
			"four_eyes": bool01(m["fourEyes"], true), "night_shift": bool01(m["nightShift"], false),
		})
	}
	put("journals", jrn)
	legs := []map[string]any{}
	for _, m := range asMaps(dump["ledger"]) {
		legs = append(legs, map[string]any{
			"id": m["id"], "journal_id": m["journalId"], "account_id": m["accountId"], "wallet_id": m["walletId"],
			"dc": m["dc"], "amount": num(m["amount"]), "currency": m["currency"], "kind": m["kind"],
			"status": or(str(m["status"]), "POSTED"), "at": m["at"],
			"four_eyes": bool01(m["fourEyes"], true), "night_shift": bool01(m["nightShift"], false),
		})
	}
	put("ledger_legs", legs)
	trs := []map[string]any{}
	for _, m := range asMaps(dump["transfers"]) {
		key := m["idempotencyKey"]
		if str(key) == "" {
			key = nil
		}
		trs = append(trs, map[string]any{
			"id": m["id"], "from_wallet_id": m["fromWalletId"], "to_wallet_id": m["toWalletId"],
			"amount": num(m["amount"]), "currency": m["currency"], "status": m["status"], "journal_id": m["journalId"],
			"idempotency_key": key, "created_at": m["createdAt"],
		})
	}
	put("transfers", trs)
	pays := []map[string]any{}
	for _, m := range asMaps(dump["payments"]) {
		merch := str(m["toAccountId"])
		if str(m["merchantId"]) != "" {
			merch = str(m["merchantId"])
		}
		pays = append(pays, map[string]any{
			"id": m["id"], "kind": m["kind"], "status": m["status"], "amount": num(m["amount"]),
			"captured_amount": num(m["capturedAmount"]), "currency": m["currency"],
			"from_account_id": m["fromAccountId"], "merchant_account_id": merch, "wallet_id": m["walletId"],
			"rrn": m["rrn"], "hold_id": m["holdId"], "journal_id": m["journalId"], "created_at": m["createdAt"],
		})
	}
	put("payments", pays)
	put("cards", project(asMaps(dump["cards"]), []string{"id", "wallet_id", "last4", "scheme", "status", "created_at"}, map[string]string{"walletId": "wallet_id", "createdAt": "created_at"}))
	put("card_auths", project(asMaps(dump["auths"]), []string{"id", "payment_id", "wallet_id", "merchant_id", "amount", "status", "rrn", "created_at"}, map[string]string{"paymentId": "payment_id", "walletId": "wallet_id", "merchantId": "merchant_id", "createdAt": "created_at"}))
	put("incoming", project(asMaps(dump["incoming"]), []string{"id", "iban", "amount", "currency", "status", "journal_id", "created_at"}, map[string]string{"journalId": "journal_id", "createdAt": "created_at"}))
	put("outgoing", project(asMaps(dump["outgoing"]), []string{"id", "from_account_id", "iban", "amount", "currency", "status", "hold_id", "created_at"}, map[string]string{"fromAccountId": "from_account_id", "holdId": "hold_id", "createdAt": "created_at"}))
	put("suspense", project(asMaps(dump["suspense"]), []string{"id", "iban", "amount", "currency", "status", "journal_id", "created_at"}, map[string]string{"journalId": "journal_id", "createdAt": "created_at"}))
	fxs := []map[string]any{}
	for _, m := range asMaps(dump["fxDeals"]) {
		cover := m["coverJournalId"]
		if cover == nil {
			cover = m["uzsJournalId"]
		}
		fxs = append(fxs, map[string]any{
			"id": m["id"], "pair": m["pair"], "amount_usd": num(m["amountUsd"]), "amount_uzs": num(m["amountUzs"]),
			"status": m["status"], "cover_journal_id": cover, "uzs_journal_id": m["uzsJournalId"],
			"usd_journal_id": m["usdJournalId"], "created_at": m["createdAt"], "note": m["note"],
		})
	}
	put("fx_deals", fxs)
	put("tickets", project(asMaps(dump["tickets"]), []string{"id", "customer_id", "kind", "status", "priority", "title", "body", "opened_at", "sla_hours"}, map[string]string{"customer": "customer_id", "openedAt": "opened_at", "slaHours": "sla_hours"}))
	put("audit_log", project(asMaps(dump["audit"]), []string{"at", "method", "path", "status", "note"}, nil))
	ui := []map[string]any{}
	for _, m := range asMaps(dump["walletUi"]) {
		ui = append(ui, map[string]any{"wallet_id": m["walletId"], "shown_available": m["shownAvailable"]})
	}
	put("wallet_ui", ui)
	iw := []map[string]any{}
	for _, w := range t["wallets"] {
		iw = append(iw, map[string]any{"id": w["id"], "name": w["name"], "currency": w["currency"], "status": w["status"], "available": w["available"]})
	}
	put("intern_wallets", iw)
	ih := []map[string]any{}
	for _, h := range t["holds"] {
		ih = append(ih, map[string]any{"id": h["id"], "wallet_id": h["wallet_id"], "amount": h["amount"], "status": h["status"]})
	}
	put("intern_holds", ih)
	ij := []map[string]any{}
	for _, j := range t["journals"] {
		ij = append(ij, map[string]any{"id": j["id"], "kind": j["kind"], "at": j["at"], "balanced": j["balanced"]})
	}
	put("intern_journals", ij)
	it := []map[string]any{}
	for _, x := range t["tickets"] {
		it = append(it, map[string]any{"id": x["id"], "kind": x["kind"], "status": x["status"], "title": x["title"]})
	}
	put("intern_tickets", it)
	l.tables = t
	l.persist()
	return nil
}

type RunResult struct {
	OK      bool             `json:"ok"`
	Kind    string           `json:"kind"`
	Columns []string         `json:"columns,omitempty"`
	Rows    []map[string]any `json:"rows,omitempty"`
	Count   int              `json:"count"`
	Error   string           `json:"error,omitempty"`
	Hint    string           `json:"hint,omitempty"`
}

func (l *Lab) Run(sqlText string, intern bool, allowDML bool) RunResult {
	l.mu.Lock()
	defer l.mu.Unlock()
	kind, err := classifySQL(sqlText, intern, allowDML)
	if err != nil {
		return RunResult{OK: false, Kind: kind, Error: err.Error()}
	}
	sandbox := cloneTables(l.tables)
	if kind == "select" {
		cols, rows, err := execSelect(sqlText, sandbox)
		if err != nil {
			return RunResult{OK: false, Kind: kind, Error: err.Error(), Hint: "Таблицы: wallets, holds, journals, intern_wallets. JOIN и GROUP BY поддерживаются."}
		}
		if len(rows) > 200 {
			rows = rows[:200]
		}
		return RunResult{OK: true, Kind: kind, Columns: cols, Rows: rows, Count: len(rows)}
	}
	n, err := execDML(sqlText, sandbox)
	if err != nil {
		return RunResult{OK: false, Kind: kind, Error: err.Error()}
	}
	return RunResult{OK: true, Kind: kind, Count: n, Hint: "DML остался в песочнице запроса и не меняет живой банк."}
}

func (l *Lab) QueryGold(sqlText string) ([]map[string]any, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, rows, err := execSelect(sqlText, l.tables)
	return rows, err
}

func (l *Lab) Schema(intern bool) []map[string]any {
	if intern {
		return []map[string]any{
			{"name": "intern_wallets", "cols": "id, name, currency, status, available", "note": "Витрина кошельков."},
			{"name": "intern_holds", "cols": "id, wallet_id, amount, status", "note": "Холды."},
			{"name": "intern_journals", "cols": "id, kind, at, balanced", "note": "Журналы."},
			{"name": "intern_tickets", "cols": "id, kind, status, title", "note": "Тикеты линии."},
		}
	}
	return []map[string]any{
		{"name": "customers", "cols": "id, name, segment, city, opened_at"},
		{"name": "accounts", "cols": "id, iban, type, currency, ledger, hold, available"},
		{"name": "wallets", "cols": "id, customer_id, ledger, hold, available"},
		{"name": "intern_wallets", "cols": "id, name, currency, status, available", "note": "Урезанная витрина Intern-задач. Есть и в PRO."},
		{"name": "holds", "cols": "id, account_id, amount, status, created_at"},
		{"name": "journals", "cols": "id, kind, debit, credit, balanced, at, four_eyes, night_shift"},
		{"name": "ledger_legs", "cols": "id, journal_id, account_id, dc, amount, currency, kind, at"},
		{"name": "transfers", "cols": "id, from_wallet_id, to_wallet_id, amount, idempotency_key"},
		{"name": "payments", "cols": "id, status, amount, merchant_account_id, hold_id"},
		{"name": "merchants", "cols": "id, account_id, name, mcc"},
		{"name": "suspense", "cols": "id, iban, amount, status, created_at"},
		{"name": "fx_deals", "cols": "id, pair, amount_usd, status, cover_journal_id"},
		{"name": "tickets", "cols": "id, customer_id, kind, status, title, opened_at"},
		{"name": "wallet_ui", "cols": "wallet_id, shown_available"},
		{"name": "cards / card_auths / incoming / outgoing / kyc / audit_log", "cols": "учебные таблицы ядра"},
	}
}

var denyStmt = regexp.MustCompile(`(?i)\b(drop|truncate|alter|attach|detach|vacuum|pragma|reindex|grant|revoke)\b`)
var internTables = regexp.MustCompile(`(?i)(^|[^a-z_])(customers|accounts|wallets|holds|journals|ledger_legs|transfers|payments|cards|card_auths|incoming|outgoing|suspense|fx_deals|merchants|tickets|audit_log|wallet_ui|kyc)($|[^a-z_])`)

func classifySQL(sqlText string, intern, allowDML bool) (string, error) {
	raw := strings.TrimSpace(stripSQLComments(sqlText))
	if raw == "" {
		return "", fmt.Errorf("Пустой запрос.")
	}
	if denyStmt.MatchString(raw) {
		return "denied", fmt.Errorf("DROP / TRUNCATE / ALTER и служебные команды в учебном ядре запрещены.")
	}
	if strings.Contains(raw, ";") && strings.TrimSpace(strings.SplitN(raw, ";", 2)[1]) != "" {
		return "denied", fmt.Errorf("Один statement за раз.")
	}
	s := strings.TrimRight(raw, "; \t\n\r")
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("Пустой запрос.")
	}
	first := strings.ToLower(strings.Fields(s)[0])
	switch first {
	case "select", "with":
		if intern && internTables.MatchString(s) {
			return "select", fmt.Errorf("На Intern доступны только витрины intern_wallets, intern_holds, intern_journals, intern_tickets.")
		}
		return "select", nil
	case "insert", "update", "delete":
		if intern || !allowDML {
			return first, fmt.Errorf("Запись в SQL-лабе — в PRO. На Intern только SELECT по витринам.")
		}
		return first, nil
	default:
		return first, fmt.Errorf("Разрешены SELECT и ограниченный DML в песочнице. Не «%s».", first)
	}
}

func stripSQLComments(s string) string {
	lines := strings.Split(s, "\n")
	var b strings.Builder
	for _, ln := range lines {
		if i := strings.Index(ln, "--"); i >= 0 {
			ln = ln[:i]
		}
		b.WriteString(ln)
		b.WriteByte('\n')
	}
	out := b.String()
	for {
		i := strings.Index(out, "/*")
		if i < 0 {
			break
		}
		j := strings.Index(out[i:], "*/")
		if j < 0 {
			out = out[:i]
			break
		}
		out = out[:i] + " " + out[i+j+2:]
	}
	return out
}

func project(rows []map[string]any, cols []string, rename map[string]string) []map[string]any {
	out := []map[string]any{}
	for _, r := range rows {
		m := map[string]any{}
		src := map[string]any{}
		for k, v := range r {
			src[k] = v
			if rename != nil {
				if nk, ok := rename[k]; ok {
					src[nk] = v
				}
			}
		}
		for _, c := range cols {
			if v, ok := src[c]; ok {
				m[c] = v
			} else {
				m[c] = src[c]
			}
		}
		out = append(out, m)
	}
	return out
}

func bool01(v any, defTrue bool) int {
	switch t := v.(type) {
	case bool:
		if t {
			return 1
		}
		return 0
	case float64:
		if t != 0 {
			return 1
		}
		return 0
	case nil:
		if defTrue {
			return 1
		}
		return 0
	default:
		if defTrue {
			return 1
		}
		return 0
	}
}

func cloneTables(in map[string][]map[string]any) map[string][]map[string]any {
	out := map[string][]map[string]any{}
	for k, rows := range in {
		cp := make([]map[string]any, len(rows))
		for i, r := range rows {
			m := map[string]any{}
			for a, b := range r {
				m[a] = b
			}
			cp[i] = m
		}
		out[k] = cp
	}
	return out
}

func execDML(sqlText string, tables map[string][]map[string]any) (int, error) {
	s := strings.TrimSpace(strings.TrimSuffix(stripSQLComments(sqlText), ";"))
	low := strings.ToLower(s)
	if strings.HasPrefix(low, "delete from ") {
		rest := strings.TrimSpace(s[len("delete from "):])
		parts := strings.SplitN(rest, " where ", 2)
		name := strings.ToLower(strings.Fields(parts[0])[0])
		rows := tables[name]
		if rows == nil {
			return 0, fmt.Errorf("нет таблицы %s", name)
		}
		if len(parts) == 1 {
			n := len(rows)
			tables[name] = []map[string]any{}
			return n, nil
		}
		keep := []map[string]any{}
		n := 0
		for _, r := range rows {
			ok, err := evalBool(parts[1], r)
			if err != nil {
				return 0, err
			}
			if ok {
				n++
			} else {
				keep = append(keep, r)
			}
		}
		tables[name] = keep
		return n, nil
	}
	return 0, fmt.Errorf("в песочнице учебный DML: DELETE FROM ... WHERE. INSERT/UPDATE — следующей итерацией лабы.")
}

func asMaps(v any) []map[string]any {
	switch t := v.(type) {
	case []map[string]any:
		return t
	case []any:
		out := []map[string]any{}
		for _, x := range t {
			if m, ok := x.(map[string]any); ok {
				out = append(out, m)
			}
		}
		return out
	default:
		return nil
	}
}

func str(v any) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	default:
		return fmt.Sprint(t)
	}
}

func num(v any) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case float32:
		return float64(t)
	case int:
		return float64(t)
	case int64:
		return float64(t)
	case json.Number:
		f, _ := t.Float64()
		return f
	case string:
		f, _ := strconv.ParseFloat(t, 64)
		return f
	default:
		return 0
	}
}

func or(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
