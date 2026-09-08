package ledger

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const Version = 4

type Contract struct {
	TransferIdempotencyRequired bool    `json:"transferIdempotencyRequired"`
	TransferDuplicateHTTP       int     `json:"transferDuplicateHttp"`
	ErrorInBody                 bool    `json:"errorInBody"`
	KycIsMaster                 bool    `json:"kycIsMaster"`
	AcsTimeoutStatus            string  `json:"acsTimeoutStatus"`
	OrderIdempotencyRequired    bool    `json:"orderIdempotencyRequired"`
	HoldThenCapture             bool    `json:"holdThenCapture"`
	RefundPostsReversal         bool    `json:"refundPostsReversal"`
	RequireBalancedJournal      bool    `json:"requireBalancedJournal"`
	FeeOnCaptureBps             float64 `json:"feeOnCaptureBps"`
	DailyP2pLimit               float64 `json:"dailyP2pLimit"`
	AllowPartialCapture         bool    `json:"allowPartialCapture"`
	FxSpreadBps                 float64 `json:"fxSpreadBps"`
	SalaryTrailerMustMatch      bool    `json:"salaryTrailerMustMatch"`
	SalaryFileOnce              bool    `json:"salaryFileOnce"`
	UnknownIbanToSuspense       bool    `json:"unknownIbanToSuspense"`
	BlockedCardCannotAuth       bool    `json:"blockedCardCannotAuth"`
}

type Account struct {
	ID       string  `json:"id"`
	WalletID any     `json:"walletId"`
	IBAN     any     `json:"iban"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Normal   string  `json:"normal"`
	Currency string  `json:"currency"`
	Ledger   float64 `json:"ledger"`
}

type Wallet struct {
	ID        string  `json:"id"`
	Customer  string  `json:"customer"`
	Name      string  `json:"name"`
	AccountID string  `json:"accountId"`
	Currency  string  `json:"currency"`
	Status    string  `json:"status"`
	Ledger    float64 `json:"ledger,omitempty"`
}

type State struct {
	Version   int              `json:"version"`
	Contract  Contract         `json:"contract"`
	Accounts  []Account        `json:"accounts"`
	Wallets   []Wallet         `json:"wallets"`
	KYC       []map[string]any `json:"kyc"`
	Plastic   []map[string]any `json:"plastic"`
	Transfers []map[string]any `json:"transfers"`
	Payments  []map[string]any `json:"payments"`
	Holds     []map[string]any `json:"holds"`
	Journals  []map[string]any `json:"journals"`
	Ledger    []map[string]any `json:"ledger"`
	Cards     []map[string]any `json:"cards"`
	Incoming  []map[string]any `json:"incoming"`
	Outgoing  []map[string]any `json:"outgoing"`
	Recon     []map[string]any `json:"recon"`
	FX        []map[string]any `json:"fx"`
	FXDeals   []map[string]any `json:"fxDeals"`
	Salary    []map[string]any `json:"salary"`
	Suspense  []map[string]any `json:"suspense"`
	Clearing  []map[string]any `json:"clearing"`
	Orders     []map[string]any `json:"orders"`
	Stock      []map[string]any `json:"stock"`
	Audit      []map[string]any `json:"audit"`
	Customers  []map[string]any `json:"customers"`
	Merchants  []map[string]any `json:"merchants"`
	Tickets    []map[string]any `json:"tickets"`
	WalletUI   []map[string]any `json:"walletUi"`
	Teaching   bool             `json:"teaching"`
	Seq        int              `json:"seq"`
}

type Bank struct {
	mu    sync.Mutex
	file  string
	state *State
}

type FailError struct {
	Status  int
	Code    string
	Message string
	Extra   map[string]any
}

func (e *FailError) Error() string { return e.Message }

func fail(status int, code, message string, extra ...map[string]any) *FailError {
	x := map[string]any{}
	if len(extra) > 0 && extra[0] != nil {
		x = extra[0]
	}
	return &FailError{Status: status, Code: code, Message: message, Extra: x}
}

func seed() *State {
	b := &Bank{state: seedBase()}
	b.plantTeaching()
	return b.state
}

func Open(dataDir string) *Bank {
	file := filepath.Join(dataDir, "malo-wallet.json")
	b := &Bank{file: file, state: seed()}
	raw, err := os.ReadFile(file)
	if err == nil {
		var st State
		if json.Unmarshal(raw, &st) == nil && st.Version >= Version && len(st.Accounts) > 0 && st.Teaching {
			b.state = &st
		}
	}
	b.save()
	return b
}

func (b *Bank) Lock()   { b.mu.Lock() }
func (b *Bank) Unlock() { b.mu.Unlock() }

func (b *Bank) Note(method, path string, status int, note, key string) {
	b.audit(method, path, status, note, key)
}

func (b *Bank) ContractErrorInBody() bool { return b.state.Contract.ErrorInBody }

func (b *Bank) Wallets() []map[string]any { return b.walletViews() }
func (b *Bank) AccountsTrial() map[string]any {
	tb := b.TrialBalance()
	return map[string]any{"items": tb["rows"], "trial": tb}
}
func (b *Bank) Holds() []map[string]any     { return b.state.Holds }
func (b *Bank) Payments() []map[string]any  { return b.state.Payments }
func (b *Bank) Journals() []map[string]any  { return b.state.Journals }
func (b *Bank) Legs() []map[string]any      { return b.state.Ledger }
func (b *Bank) Incoming() []map[string]any  { return b.state.Incoming }
func (b *Bank) Outgoing() []map[string]any  { return b.state.Outgoing }
func (b *Bank) CardsAuth() []map[string]any { return b.state.Cards }
func (b *Bank) Plastic() []map[string]any   { return b.state.Plastic }
func (b *Bank) KYCList() []map[string]any   { return b.state.KYC }
func (b *Bank) FXRates() map[string]any {
	return map[string]any{"items": b.state.FX, "spreadBps": b.state.Contract.FxSpreadBps}
}
func (b *Bank) FXDeals() []map[string]any   { return b.state.FXDeals }
func (b *Bank) Salary() []map[string]any    { return b.state.Salary }
func (b *Bank) Suspense() []map[string]any  { return b.state.Suspense }
func (b *Bank) Clearing() []map[string]any  { return b.state.Clearing }
func (b *Bank) Stock() []map[string]any     { return b.state.Stock }
func (b *Bank) Customers() []map[string]any { return b.state.Customers }
func (b *Bank) Merchants() []map[string]any { return b.state.Merchants }
func (b *Bank) Tickets() []map[string]any   { return b.state.Tickets }
func (b *Bank) AuditLog() []map[string]any  { return b.state.Audit }
func (b *Bank) WalletUI() []map[string]any  { return b.state.WalletUI }
func (b *Bank) GetHold(id string) (map[string]any, error) {
	h := find(b.state.Holds, id)
	if h == nil {
		return nil, fail(404, "NOT_FOUND", "Холд не найден")
	}
	return h, nil
}
func (b *Bank) GetTicket(id string) (map[string]any, error) {
	t := find(b.state.Tickets, id)
	if t == nil {
		return nil, fail(404, "NOT_FOUND", "Тикет не найден")
	}
	return t, nil
}
func (b *Bank) CreateTicket(body map[string]any) (int, any, error) {
	title := strings.TrimSpace(str(body["title"]))
	if title == "" {
		return 0, nil, fail(400, "TITLE", "Нужен заголовок тикета")
	}
	row := map[string]any{
		"id": b.nextID("tkt"), "customer": or(str(body["customer"]), "anna"),
		"kind": or(str(body["kind"]), "OPS"), "status": "OPEN", "title": title,
		"body": or(str(body["body"]), ""), "priority": or(str(body["priority"]), "P2"),
		"openedAt": now(), "slaHours": 24.0, "assignee": body["assignee"],
	}
	b.state.Tickets = append(b.state.Tickets, row)
	b.save()
	return 201, row, nil
}
func (b *Bank) CommentTicket(id string, body map[string]any) (int, any, error) {
	row := find(b.state.Tickets, id)
	if row == nil {
		return 0, nil, fail(404, "NOT_FOUND", "Тикет не найден")
	}
	text := strings.TrimSpace(str(body["text"]))
	if text == "" {
		return 0, nil, fail(400, "TEXT", "Пустой комментарий")
	}
	comments, _ := row["comments"].([]any)
	if comments == nil {
		if raw, ok := row["comments"].([]map[string]any); ok {
			for _, c := range raw {
				comments = append(comments, c)
			}
		}
	}
	c := map[string]any{"at": now(), "text": text, "author": or(str(body["author"]), "аналитик")}
	comments = append(comments, c)
	row["comments"] = comments
	if str(body["status"]) != "" {
		row["status"] = str(body["status"])
	}
	b.save()
	return 200, row, nil
}
func (b *Bank) LegsFilter(acc string) []map[string]any {
	if acc == "" {
		return b.state.Ledger
	}
	out := []map[string]any{}
	for _, l := range b.state.Ledger {
		if str(l["accountId"]) == acc {
			out = append(out, l)
		}
	}
	return out
}

func (b *Bank) save() {
	_ = os.MkdirAll(filepath.Dir(b.file), 0o755)
	raw, _ := json.MarshalIndent(b.state, "", "  ")
	_ = os.WriteFile(b.file, raw, 0o644)
}

func (b *Bank) Statement(id string) (map[string]any, error) { return b.statement(id) }
func (b *Bank) TAccount(id string) (map[string]any, error)  { return b.tAccount(id) }

func now() string { return time.Now().UTC().Format("2006-01-02T15:04:05.000Z") }

func (b *Bank) account(id string) *Account {
	for i := range b.state.Accounts {
		if b.state.Accounts[i].ID == id {
			return &b.state.Accounts[i]
		}
	}
	return nil
}

func (b *Bank) wallet(id string) *Wallet {
	for i := range b.state.Wallets {
		if b.state.Wallets[i].ID == id {
			return &b.state.Wallets[i]
		}
	}
	return nil
}

func liveLeg(l map[string]any) bool {
	st, _ := l["status"].(string)
	return st != "REVERSED"
}

func (b *Bank) openHoldSum(accountID string) float64 {
	sum := 0.0
	for _, h := range b.state.Holds {
		if str(h["accountId"]) == accountID && str(h["status"]) == "OPEN" {
			sum += num(h["amount"])
		}
	}
	return sum
}

func signedDelta(acc *Account, dc string, amount float64) float64 {
	if acc.Normal == "credit" {
		if dc == "C" {
			return amount
		}
		return -amount
	}
	if dc == "D" {
		return amount
	}
	return -amount
}

func (b *Bank) syncWallets() {
	for i := range b.state.Wallets {
		if acc := b.account(b.state.Wallets[i].AccountID); acc != nil {
			b.state.Wallets[i].Ledger = acc.Ledger
		}
	}
}

func (b *Bank) walletView(w Wallet) map[string]any {
	acc := b.account(w.AccountID)
	var kyc map[string]any
	for _, k := range b.state.KYC {
		if str(k["customer"]) == w.Customer {
			kyc = k
			break
		}
	}
	hold := 0.0
	ledger := 0.0
	iban := any(nil)
	accID := w.AccountID
	ccy := w.Currency
	if acc != nil {
		hold = b.openHoldSum(acc.ID)
		ledger = acc.Ledger
		iban = acc.IBAN
		accID = acc.ID
		ccy = acc.Currency
	}
	status := w.Status
	kycStatus := any(nil)
	if kyc != nil {
		kycStatus = kyc["status"]
		if b.state.Contract.KycIsMaster && str(kyc["status"]) != "APPROVED" {
			status = "PENDING_KYC"
		}
	}
	return map[string]any{
		"id": w.ID, "customer": w.Customer, "name": w.Name, "accountId": accID,
		"currency": ccy, "status": status, "iban": iban, "ledger": ledger, "hold": hold,
		"available": ledger - hold, "balance": ledger - hold, "kyc": kycStatus,
	}
}

func (b *Bank) ccyTrial(ccy string) map[string]any {
	var accs []Account
	for _, a := range b.state.Accounts {
		if a.Currency == ccy {
			accs = append(accs, a)
		}
	}
	debit, credit := 0.0, 0.0
	for _, l := range b.state.Ledger {
		if str(l["currency"]) == ccy && liveLeg(l) {
			if str(l["dc"]) == "D" {
				debit += num(l["amount"])
			} else if str(l["dc"]) == "C" {
				credit += num(l["amount"])
			}
		}
	}
	rows := []map[string]any{}
	for _, a := range accs {
		hold := b.openHoldSum(a.ID)
		rows = append(rows, map[string]any{
			"id": a.ID, "name": a.Name, "type": a.Type, "normal": a.Normal, "currency": a.Currency,
			"iban": a.IBAN, "ledger": a.Ledger, "hold": hold, "available": a.Ledger - hold,
		})
	}
	return map[string]any{"currency": ccy, "rows": rows, "debit": debit, "credit": credit, "balanced": debit == credit, "diff": debit - credit}
}

func (b *Bank) TrialBalance() map[string]any {
	seen := map[string]bool{}
	var ccys []string
	for _, a := range b.state.Accounts {
		if !seen[a.Currency] {
			seen[a.Currency] = true
			ccys = append(ccys, a.Currency)
		}
	}
	by := map[string]any{}
	rows := []map[string]any{}
	balanced := true
	for _, c := range ccys {
		t := b.ccyTrial(c)
		by[c] = t
		if r, ok := t["rows"].([]map[string]any); ok {
			rows = append(rows, r...)
		}
		if t["balanced"] != true {
			balanced = false
		}
	}
	uzs, _ := by["UZS"].(map[string]any)
	if uzs == nil {
		uzs = map[string]any{"debit": 0.0, "credit": 0.0, "diff": 0.0, "balanced": true, "rows": []any{}}
	}
	return map[string]any{
		"rows": rows, "debit": uzs["debit"], "credit": uzs["credit"], "diff": uzs["diff"],
		"byCurrency": by, "balanced": balanced,
	}
}

func (b *Bank) postJournal(kind, ref, note string, legs []map[string]any) (string, error) {
	byCcy := map[string]map[string]float64{}
	for _, leg := range legs {
		acc := b.account(str(leg["accountId"]))
		if acc == nil {
			return "", fail(409, "NO_ACCOUNT", "NO_ACCOUNT")
		}
		ccy := str(leg["currency"])
		if ccy == "" {
			ccy = acc.Currency
		}
		if byCcy[ccy] == nil {
			byCcy[ccy] = map[string]float64{"D": 0, "C": 0}
		}
		byCcy[ccy][str(leg["dc"])] += num(leg["amount"])
	}
	holes := false
	for _, v := range byCcy {
		if v["D"] != v["C"] {
			holes = true
		}
	}
	if holes && b.state.Contract.RequireBalancedJournal {
		return "", fail(409, "UNBALANCED", "UNBALANCED", map[string]any{"byCurrency": byCcy})
	}
	debit, credit := 0.0, 0.0
	for _, leg := range legs {
		if str(leg["dc"]) == "D" {
			debit += num(leg["amount"])
		} else {
			credit += num(leg["amount"])
		}
	}
	journalID := b.nextID("jrn")
	at := now()
	for _, leg := range legs {
		acc := b.account(str(leg["accountId"]))
		amount := num(leg["amount"])
		delta := signedDelta(acc, str(leg["dc"]), amount)
		acc.Ledger += delta
		n := note
		if leg["note"] != nil {
			n = str(leg["note"])
		}
		b.state.Ledger = append(b.state.Ledger, map[string]any{
			"id": b.nextID("ld"), "journalId": journalID, "accountId": acc.ID, "walletId": acc.WalletID,
			"dc": str(leg["dc"]), "amount": amount, "delta": delta, "currency": acc.Currency,
			"kind": kind, "ref": ref, "note": n, "status": "POSTED", "at": at,
		})
	}
	b.state.Journals = append(b.state.Journals, map[string]any{
		"id": journalID, "kind": kind, "ref": ref, "debit": debit, "credit": credit,
		"byCurrency": byCcy, "balanced": !holes, "at": at, "note": note,
	})
	b.syncWallets()
	return journalID, nil
}

func (b *Bank) statement(accountID string) (map[string]any, error) {
	acc := b.account(accountID)
	if acc == nil {
		return nil, fail(404, "NOT_FOUND", "Счёт не найден")
	}
	items := []map[string]any{}
	for _, l := range b.state.Ledger {
		if str(l["accountId"]) == accountID {
			items = append(items, l)
		}
	}
	sort.Slice(items, func(i, j int) bool { return str(items[i]["at"]) < str(items[j]["at"]) })
	posted := 0.0
	for _, l := range items {
		if liveLeg(l) {
			posted += num(l["delta"])
		}
	}
	opening := acc.Ledger - posted
	run := opening
	rows := []map[string]any{}
	for _, l := range items {
		if liveLeg(l) {
			run += num(l["delta"])
		}
		row := clone(l)
		row["balanceAfter"] = run
		rows = append(rows, row)
	}
	return map[string]any{
		"accountId": accountID, "name": acc.Name, "currency": acc.Currency, "iban": acc.IBAN,
		"opening": opening, "closing": acc.Ledger, "hold": b.openHoldSum(acc.ID), "items": rows,
	}, nil
}

func (b *Bank) tAccount(accountID string) (map[string]any, error) {
	acc := b.account(accountID)
	if acc == nil {
		return nil, fail(404, "NOT_FOUND", "Счёт не найден")
	}
	debit, credit := []map[string]any{}, []map[string]any{}
	dSum, cSum := 0.0, 0.0
	for _, l := range b.state.Ledger {
		if str(l["accountId"]) != accountID || !liveLeg(l) {
			continue
		}
		if str(l["dc"]) == "D" {
			debit = append(debit, l)
			dSum += num(l["amount"])
		} else {
			credit = append(credit, l)
			cSum += num(l["amount"])
		}
	}
	return map[string]any{
		"account": map[string]any{"id": acc.ID, "name": acc.Name, "type": acc.Type, "normal": acc.Normal, "currency": acc.Currency, "ledger": acc.Ledger, "iban": acc.IBAN},
		"debit": debit, "credit": credit, "debitSum": dSum, "creditSum": cSum, "hold": b.openHoldSum(acc.ID),
	}, nil
}

func str(v any) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	default:
		return fmt.Sprint(v)
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
	default:
		return 0
	}
}

func clone(m map[string]any) map[string]any {
	o := map[string]any{}
	for k, v := range m {
		o[k] = v
	}
	return o
}

func round(v float64) float64 { return math.Round(v) }

func (b *Bank) audit(method, path string, status int, note, key string) {
	b.state.Audit = append([]map[string]any{{
		"at": now(), "method": method, "path": path, "status": status, "note": note, "key": nilIfEmpty(key),
	}}, b.state.Audit...)
	if len(b.state.Audit) > 120 {
		b.state.Audit = b.state.Audit[:120]
	}
}

func nilIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
