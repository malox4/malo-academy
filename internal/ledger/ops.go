package ledger

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func (b *Bank) nextID(prefix string) string {
	b.state.Seq++
	return prefix + "_" + strconv.FormatInt(int64(b.state.Seq), 36)
}

func (b *Bank) Snapshot() map[string]any {
	tb := b.TrialBalance()
	annaDebits := 0
	for _, l := range b.state.Ledger {
		if str(l["accountId"]) == "acc_anna" && str(l["kind"]) == "P2P" && str(l["dc"]) == "D" && num(l["amount"]) == 5000 {
			annaDebits++
		}
	}
	c := b.state.Contract
	balanced, _ := tb["balanced"].(bool)
	return map[string]any{
		"contract": c, "wallets": b.walletViews(), "accounts": tb["rows"], "trial": tb,
		"kyc": b.state.KYC, "plastic": b.state.Plastic, "transfers": b.state.Transfers,
		"payments": b.state.Payments, "holds": b.state.Holds, "journals": b.state.Journals,
		"ledger": b.state.Ledger, "cards": b.state.Cards, "incoming": b.state.Incoming,
		"outgoing": b.state.Outgoing, "recon": b.state.Recon, "fx": b.state.FX, "fxDeals": b.state.FXDeals,
		"salary": b.state.Salary, "suspense": b.state.Suspense, "clearing": b.state.Clearing,
		"orders": b.state.Orders, "stock": b.state.Stock, "audit": b.state.Audit,
		"customers": b.state.Customers, "merchants": b.state.Merchants, "tickets": b.state.Tickets,
		"walletUi": b.state.WalletUI, "teaching": b.state.Teaching,
		"missions": map[string]any{
			"reproducedDouble": annaDebits >= 2,
			"contractShipped":  c.TransferIdempotencyRequired && c.TransferDuplicateHTTP == 409 && !c.ErrorInBody,
			"kycHonest":        c.KycIsMaster,
			"acsHonest":        c.AcsTimeoutStatus == "UNKNOWN",
			"orderSafe":        c.OrderIdempotencyRequired,
			"capturedHold":     c.HoldThenCapture && somePay(b.state.Payments, func(p map[string]any) bool { return p["holdId"] != nil && in(str(p["status"]), "CAPTURED", "REFUNDED", "CHARGEBACK") }),
			"refundHole":       someLeg(b.state.Ledger, func(l map[string]any) bool { return str(l["kind"]) == "REFUND_HOLE" }) && !balanced,
			"refundClean":      somePay(b.state.Payments, func(p map[string]any) bool { return str(p["status"]) == "REFUNDED" && p["refundJournalId"] != nil }) && balanced,
			"incomingPosted":   some(b.state.Incoming, func(i map[string]any) bool { return str(i["status"]) == "POSTED" }),
			"reconMismatch":    reconBad(b.state.Recon),
			"fxPosted":         some(b.state.FXDeals, func(d map[string]any) bool { return str(d["status"]) == "POSTED" }),
			"salaryPosted":     some(b.state.Salary, func(s map[string]any) bool { return str(s["status"]) == "POSTED" }),
			"clearingSettled":  some(b.state.Clearing, func(c map[string]any) bool { return str(c["status"]) == "SETTLED" }),
			"partialCaptured": somePay(b.state.Payments, func(p map[string]any) bool {
				return num(p["capturedAmount"]) > 0 && num(p["capturedAmount"]) < num(p["amount"]) && in(str(p["status"]), "CAPTURED", "REFUNDED", "CHARGEBACK")
			}),
			"suspensePosted": some(b.state.Suspense, func(s map[string]any) bool { return in(str(s["status"]), "SUSPENSE", "ALLOCATED") }),
			"mt103Issued":    some(b.state.Outgoing, func(o map[string]any) bool { return str(o["status"]) == "POSTED" && o["mt103"] != nil }),
		},
	}
}

func (b *Bank) walletViews() []map[string]any {
	out := []map[string]any{}
	for _, w := range b.state.Wallets {
		out = append(out, b.walletView(w))
	}
	return out
}

func (b *Bank) Reset() map[string]any {
	b.state = seed()
	b.save()
	return map[string]any{"ok": true, "state": b.Snapshot()}
}

func (b *Bank) GetContract() Contract { return b.state.Contract }

func (b *Bank) PatchContract(body map[string]any) Contract {
	c := &b.state.Contract
	setBool := func(k string, dst *bool) {
		if v, ok := body[k]; ok {
			if b, ok := v.(bool); ok {
				*dst = b
			}
		}
	}
	setBool("transferIdempotencyRequired", &c.TransferIdempotencyRequired)
	setBool("errorInBody", &c.ErrorInBody)
	setBool("kycIsMaster", &c.KycIsMaster)
	setBool("orderIdempotencyRequired", &c.OrderIdempotencyRequired)
	setBool("holdThenCapture", &c.HoldThenCapture)
	setBool("refundPostsReversal", &c.RefundPostsReversal)
	setBool("requireBalancedJournal", &c.RequireBalancedJournal)
	setBool("allowPartialCapture", &c.AllowPartialCapture)
	setBool("salaryTrailerMustMatch", &c.SalaryTrailerMustMatch)
	setBool("salaryFileOnce", &c.SalaryFileOnce)
	setBool("unknownIbanToSuspense", &c.UnknownIbanToSuspense)
	setBool("blockedCardCannotAuth", &c.BlockedCardCannotAuth)
	if v, ok := body["transferDuplicateHttp"]; ok {
		c.TransferDuplicateHTTP = int(num(v))
	}
	if v, ok := body["acsTimeoutStatus"]; ok {
		c.AcsTimeoutStatus = str(v)
	}
	if v, ok := body["feeOnCaptureBps"]; ok {
		c.FeeOnCaptureBps = num(v)
	}
	if v, ok := body["dailyP2pLimit"]; ok {
		c.DailyP2pLimit = num(v)
	}
	if v, ok := body["fxSpreadBps"]; ok {
		c.FxSpreadBps = num(v)
	}
	if c.TransferDuplicateHTTP != 200 && c.TransferDuplicateHTTP != 409 {
		c.TransferDuplicateHTTP = 200
	}
	if c.AcsTimeoutStatus != "SUCCESS" && c.AcsTimeoutStatus != "UNKNOWN" {
		c.AcsTimeoutStatus = "SUCCESS"
	}
	b.save()
	return *c
}

func (b *Bank) CreateTransfer(key string, body map[string]any) (int, any, error) {
	if b.state.Contract.TransferIdempotencyRequired && key == "" {
		return 0, nil, fail(400, "IDEMPOTENCY_KEY_REQUIRED", "Idempotency-Key обязателен")
	}
	if key != "" {
		for _, t := range b.state.Transfers {
			if str(t["idempotencyKey"]) == key {
				code := b.state.Contract.TransferDuplicateHTTP
				if code == 409 {
					return 409, map[string]any{"error": map[string]any{"code": "DUPLICATE", "message": "Перевод с этим ключом уже есть"}, "existingId": t["id"]}, nil
				}
				return 200, t, nil
			}
		}
	}
	fromW, toW := b.wallet(str(body["fromWalletId"])), b.wallet(str(body["toWalletId"]))
	if fromW == nil || toW == nil {
		return 0, nil, fail(404, "WALLET", "from/to кошелёк не найден")
	}
	from, to := b.account(fromW.AccountID), b.account(toW.AccountID)
	if from.Currency != to.Currency {
		return 0, nil, fail(409, "CCY", "P2P только в одной валюте — иначе FX")
	}
	amount := num(body["amount"])
	if amount <= 0 {
		return 0, nil, fail(400, "AMOUNT", "amount > 0")
	}
	avail := from.Ledger - b.openHoldSum(from.ID)
	if avail < amount {
		return 0, nil, fail(409, "INSUFFICIENT", "Недостаточно available", map[string]any{"available": avail, "hold": b.openHoldSum(from.ID)})
	}
	used := b.dailyP2pOut(from.ID)
	if used+amount > b.state.Contract.DailyP2pLimit {
		return 0, nil, fail(409, "LIMIT", "Дневной лимит P2P", map[string]any{"used": used, "limit": b.state.Contract.DailyP2pLimit})
	}
	id := b.nextID("tr")
	jid, err := b.postJournal("P2P", id, "p2p transfer", []map[string]any{
		{"accountId": from.ID, "dc": "D", "amount": amount},
		{"accountId": to.ID, "dc": "C", "amount": amount},
	})
	if err != nil {
		return 0, nil, err
	}
	tx := map[string]any{
		"id": id, "fromWalletId": fromW.ID, "toWalletId": toW.ID, "fromAccountId": from.ID, "toAccountId": to.ID,
		"amount": amount, "currency": from.Currency, "status": "SUCCESS", "journalId": jid,
		"idempotencyKey": nilIfEmpty(key), "createdAt": now(),
	}
	b.state.Transfers = append(b.state.Transfers, tx)
	b.save()
	return 201, tx, nil
}

func (b *Bank) dailyP2pOut(accountID string) float64 {
	day := now()[:10]
	s := 0.0
	for _, t := range b.state.Transfers {
		if str(t["fromAccountId"]) == accountID && strings.HasPrefix(str(t["createdAt"]), day) && str(t["status"]) == "SUCCESS" {
			s += num(t["amount"])
		}
	}
	return s
}

func (b *Bank) IssueCard(body map[string]any) (int, any, error) {
	w := b.wallet(str(body["walletId"]))
	if w == nil {
		return 0, nil, fail(404, "WALLET", "Кошелёк не найден")
	}
	row := map[string]any{
		"id": b.nextID("crd"), "walletId": w.ID, "last4": fmt.Sprintf("%04d", 1000+rand.Intn(9000)),
		"scheme": or(str(body["scheme"]), "MIR"), "status": "ACTIVE", "createdAt": now(),
	}
	b.state.Plastic = append(b.state.Plastic, row)
	b.save()
	return 201, row, nil
}

func (b *Bank) BlockCard(id string) (int, any, error) {
	for _, row := range b.state.Plastic {
		if str(row["id"]) == id {
			row["status"] = "BLOCKED"
			row["blockedAt"] = now()
			b.save()
			return 200, row, nil
		}
	}
	return 0, nil, fail(404, "NOT_FOUND", "Карта не найдена")
}

func (b *Bank) AuthorizeCard(timeout bool, body map[string]any) (int, any, error) {
	w := b.wallet(str(body["walletId"]))
	if w == nil {
		return 0, nil, fail(404, "WALLET", "Кошелёк не найден")
	}
	if cid := str(body["cardId"]); cid != "" {
		var card map[string]any
		for _, c := range b.state.Plastic {
			if str(c["id"]) == cid {
				card = c
				break
			}
		}
		if card == nil {
			return 0, nil, fail(404, "CARD", "Карта не найдена")
		}
		if str(card["walletId"]) != w.ID {
			return 0, nil, fail(409, "CARD_WALLET", "Карта другого кошелька")
		}
		if str(card["status"]) == "BLOCKED" && b.state.Contract.BlockedCardCannotAuth {
			return 0, nil, fail(409, "CARD_BLOCKED", "Карта в блоке")
		}
	}
	acc := b.account(w.AccountID)
	merchantID := or(str(body["merchantAccountId"]), "acc_merchant")
	amount := num(body["amount"])
	if amount <= 0 {
		return 0, nil, fail(400, "AMOUNT", "amount > 0")
	}
	rrn, payID := b.nextID("rrn"), b.nextID("pay")
	if timeout {
		status := b.state.Contract.AcsTimeoutStatus
		row := map[string]any{"id": b.nextID("auth"), "paymentId": payID, "walletId": w.ID, "amount": amount, "rrn": rrn, "acs": "TIMEOUT", "status": status, "createdAt": now()}
		b.state.Cards = append(b.state.Cards, row)
		pay := map[string]any{"id": payID, "kind": "CARD", "status": status, "amount": amount, "currency": acc.Currency, "fromAccountId": acc.ID, "toAccountId": merchantID, "rrn": rrn, "holdId": nil, "createdAt": now()}
		if status == "SUCCESS" {
			jid, err := b.postJournal("CARD_FALSE_SUCCESS", payID, "ACS timeout but SUCCESS — баг контракта", []map[string]any{
				{"accountId": acc.ID, "dc": "D", "amount": amount},
				{"accountId": merchantID, "dc": "C", "amount": amount},
			})
			if err != nil {
				return 0, nil, err
			}
			pay["journalId"] = jid
			pay["status"] = "CAPTURED"
			pay["capturedAmount"] = amount
		}
		b.state.Payments = append(b.state.Payments, pay)
		b.save()
		out := clone(row)
		out["payment"] = pay
		return 200, out, nil
	}
	avail := acc.Ledger - b.openHoldSum(acc.ID)
	if avail < amount {
		return 0, nil, fail(409, "INSUFFICIENT", "Недостаточно available", map[string]any{"available": avail})
	}
	if b.state.Contract.HoldThenCapture {
		hold := map[string]any{"id": b.nextID("hld"), "accountId": acc.ID, "paymentId": payID, "amount": amount, "status": "OPEN", "createdAt": now()}
		b.state.Holds = append(b.state.Holds, hold)
		pay := map[string]any{
			"id": payID, "kind": "CARD", "status": "AUTHORIZED", "amount": amount, "currency": acc.Currency,
			"fromAccountId": acc.ID, "toAccountId": merchantID, "rrn": rrn, "holdId": hold["id"],
			"cardId": body["cardId"], "authCode": b.nextID("apv"), "createdAt": now(),
		}
		b.state.Payments = append(b.state.Payments, pay)
		b.state.Cards = append(b.state.Cards, map[string]any{"id": b.nextID("auth"), "paymentId": payID, "walletId": w.ID, "amount": amount, "rrn": rrn, "acs": "OK", "status": "AUTHORIZED", "createdAt": now()})
		b.save()
		return 201, map[string]any{"payment": pay, "hold": hold, "available": b.walletView(*w)["available"]}, nil
	}
	jid, err := b.postJournal("CARD_CAPTURE_DIRECT", payID, "нет холда — сразу capture", []map[string]any{
		{"accountId": acc.ID, "dc": "D", "amount": amount},
		{"accountId": merchantID, "dc": "C", "amount": amount},
	})
	if err != nil {
		return 0, nil, err
	}
	pay := map[string]any{
		"id": payID, "kind": "CARD", "status": "CAPTURED", "amount": amount, "capturedAmount": amount,
		"currency": acc.Currency, "fromAccountId": acc.ID, "toAccountId": merchantID, "rrn": rrn,
		"journalId": jid, "holdId": nil, "createdAt": now(),
	}
	b.state.Payments = append(b.state.Payments, pay)
	b.state.Cards = append(b.state.Cards, map[string]any{"id": b.nextID("auth"), "paymentId": payID, "walletId": w.ID, "amount": amount, "rrn": rrn, "acs": "OK", "status": "CAPTURED", "createdAt": now()})
	b.save()
	return 201, map[string]any{"payment": pay}, nil
}

func (b *Bank) CapturePayment(id string, body map[string]any) (int, any, error) {
	pay := find(b.state.Payments, id)
	if pay == nil {
		return 0, nil, fail(404, "NOT_FOUND", "Платёж не найден")
	}
	if str(pay["status"]) != "AUTHORIZED" {
		return 0, nil, fail(409, "STATE", "Нельзя capture из "+str(pay["status"]))
	}
	amount := num(pay["amount"])
	requested := num(body["amount"])
	if requested > 0 && requested != amount && b.state.Contract.AllowPartialCapture {
		if requested > amount {
			return 0, nil, fail(409, "AMOUNT", "partial > auth")
		}
		amount = requested
	}
	bps := b.state.Contract.FeeOnCaptureBps
	fee := round((amount * bps) / 10000)
	merchantNet := amount - fee
	legs := []map[string]any{
		{"accountId": pay["fromAccountId"], "dc": "D", "amount": amount},
		{"accountId": pay["toAccountId"], "dc": "C", "amount": merchantNet},
	}
	if fee > 0 {
		legs = append(legs, map[string]any{"accountId": "acc_fee", "dc": "C", "amount": fee, "note": "merchant discount"})
	} else {
		legs[1]["amount"] = amount
	}
	jid, err := b.postJournal("CARD_CAPTURE", str(pay["id"]), "RRN "+str(pay["rrn"]), legs)
	if err != nil {
		return 0, nil, err
	}
	pay["journalId"] = jid
	if hold := find(b.state.Holds, str(pay["holdId"])); hold != nil {
		hold["status"] = "CAPTURED"
		hold["capturedAmount"] = amount
		if amount < num(pay["amount"]) {
			hold["releasedAmount"] = num(pay["amount"]) - amount
		}
	}
	pay["status"] = "CAPTURED"
	pay["capturedAt"] = now()
	pay["capturedAmount"] = amount
	pay["fee"] = fee
	b.save()
	return 200, map[string]any{"payment": pay, "trial": b.TrialBalance()}, nil
}

func (b *Bank) CancelPayment(id string) (int, any, error) {
	pay := find(b.state.Payments, id)
	if pay == nil {
		return 0, nil, fail(404, "NOT_FOUND", "Платёж не найден")
	}
	if str(pay["status"]) != "AUTHORIZED" {
		return 0, nil, fail(409, "STATE", "Нельзя void из "+str(pay["status"]))
	}
	if hold := find(b.state.Holds, str(pay["holdId"])); hold != nil {
		hold["status"] = "RELEASED"
	}
	pay["status"] = "CANCELLED"
	b.save()
	acc := b.account(str(pay["fromAccountId"]))
	avail := acc.Ledger - b.openHoldSum(acc.ID)
	return 200, map[string]any{"payment": pay, "available": avail}, nil
}

func (b *Bank) RefundPayment(id string, body map[string]any) (int, any, error) {
	pay := find(b.state.Payments, id)
	if pay == nil {
		return 0, nil, fail(404, "NOT_FOUND", "Платёж не найден")
	}
	if str(pay["status"]) != "CAPTURED" {
		return 0, nil, fail(409, "STATE", "Refund только после capture")
	}
	amount := num(pay["capturedAmount"])
	if amount == 0 {
		amount = num(pay["amount"])
	}
	if v := num(body["amount"]); v > 0 {
		amount = v
	}
	if b.state.Contract.RefundPostsReversal {
		jid, err := b.postJournal("REFUND", str(pay["id"]), "реверс capture", []map[string]any{
			{"accountId": pay["toAccountId"], "dc": "D", "amount": amount},
			{"accountId": pay["fromAccountId"], "dc": "C", "amount": amount},
		})
		if err != nil {
			return 0, nil, err
		}
		pay["refundJournalId"] = jid
	} else {
		acc := b.account(str(pay["fromAccountId"]))
		acc.Ledger += amount
		b.state.Ledger = append(b.state.Ledger, map[string]any{
			"id": b.nextID("ld"), "journalId": nil, "accountId": acc.ID, "walletId": acc.WalletID,
			"dc": "C", "amount": amount, "delta": amount, "currency": acc.Currency, "kind": "REFUND_HOLE",
			"ref": pay["id"], "note": "возврат без DR мерчанта — дыра", "status": "POSTED", "at": now(),
		})
		b.syncWallets()
	}
	pay["status"] = "REFUNDED"
	pay["refundedAt"] = now()
	b.save()
	return 200, map[string]any{"payment": pay, "trial": b.TrialBalance()}, nil
}

func (b *Bank) ChargebackPayment(id string, body map[string]any) (int, any, error) {
	pay := find(b.state.Payments, id)
	if pay == nil {
		return 0, nil, fail(404, "NOT_FOUND", "Платёж не найден")
	}
	if !in(str(pay["status"]), "CAPTURED", "REFUNDED") {
		return 0, nil, fail(409, "STATE", "Chargeback после capture")
	}
	amount := num(body["amount"])
	if amount == 0 {
		amount = num(pay["capturedAmount"])
		if amount == 0 {
			amount = num(pay["amount"])
		}
	}
	jid, err := b.postJournal("CHARGEBACK", str(pay["id"]), "спор", []map[string]any{
		{"accountId": pay["toAccountId"], "dc": "D", "amount": amount},
		{"accountId": pay["fromAccountId"], "dc": "C", "amount": amount},
	})
	if err != nil {
		return 0, nil, err
	}
	pay["cbJournalId"] = jid
	if _, err := b.postJournal("CHARGEBACK_FEE", str(pay["id"]), "штраф эквайера", []map[string]any{
		{"accountId": "acc_cb", "dc": "D", "amount": 15000.0},
		{"accountId": "acc_nostro", "dc": "C", "amount": 15000.0},
	}); err != nil {
		return 0, nil, err
	}
	pay["status"] = "CHARGEBACK"
	b.save()
	return 200, map[string]any{"payment": pay, "trial": b.TrialBalance()}, nil
}

func (b *Bank) BankIncoming(body map[string]any) (int, any, error) {
	iban, amount := str(body["iban"]), num(body["amount"])
	if iban == "" || amount == 0 {
		return 0, nil, fail(400, "FIELDS", "iban и amount")
	}
	id := b.nextID("in")
	var acc *Account
	for i := range b.state.Accounts {
		if str(b.state.Accounts[i].IBAN) == iban {
			acc = &b.state.Accounts[i]
			break
		}
	}
	if acc == nil {
		if !b.state.Contract.UnknownIbanToSuspense {
			return 0, nil, fail(404, "IBAN", "IBAN не наш")
		}
		jid, err := b.postJournal("SUSPENSE_IN", id, or(str(body["paymentRef"]), "unknown iban"), []map[string]any{
			{"accountId": "acc_nostro", "dc": "D", "amount": amount},
			{"accountId": "acc_suspense", "dc": "C", "amount": amount},
		})
		if err != nil {
			return 0, nil, err
		}
		row := map[string]any{"id": id, "iban": iban, "amount": amount, "currency": "UZS", "paymentRef": body["paymentRef"], "journalId": jid, "status": "SUSPENSE", "createdAt": now()}
		b.state.Incoming = append(b.state.Incoming, row)
		sus := clone(row)
		sus["allocated"] = false
		b.state.Suspense = append(b.state.Suspense, sus)
		b.save()
		return 201, row, nil
	}
	jid, err := b.postJournal("IBAN_IN", id, or(str(body["paymentRef"]), "incoming"), []map[string]any{
		{"accountId": "acc_nostro", "dc": "D", "amount": amount},
		{"accountId": acc.ID, "dc": "C", "amount": amount},
	})
	if err != nil {
		return 0, nil, err
	}
	row := map[string]any{"id": id, "iban": acc.IBAN, "amount": amount, "currency": acc.Currency, "paymentRef": body["paymentRef"], "journalId": jid, "status": "POSTED", "createdAt": now()}
	b.state.Incoming = append(b.state.Incoming, row)
	b.save()
	return 201, row, nil
}

func (b *Bank) AllocateSuspense(id string, body map[string]any) (int, any, error) {
	row := find(b.state.Suspense, id)
	if row == nil {
		return 0, nil, fail(404, "NOT_FOUND", "Suspence не найден")
	}
	if row["allocated"] == true {
		return 0, nil, fail(409, "STATE", "Уже разнесён")
	}
	var acc *Account
	for i := range b.state.Accounts {
		if str(b.state.Accounts[i].IBAN) == str(body["iban"]) {
			acc = &b.state.Accounts[i]
			break
		}
	}
	if acc == nil {
		return 0, nil, fail(404, "IBAN", "IBAN не наш")
	}
	jid, err := b.postJournal("SUSPENSE_ALLOC", str(row["id"]), str(acc.IBAN), []map[string]any{
		{"accountId": "acc_suspense", "dc": "D", "amount": num(row["amount"])},
		{"accountId": acc.ID, "dc": "C", "amount": num(row["amount"])},
	})
	if err != nil {
		return 0, nil, err
	}
	row["allocJournalId"] = jid
	row["allocated"] = true
	row["allocatedIban"] = acc.IBAN
	row["status"] = "ALLOCATED"
	if inc := find(b.state.Incoming, id); inc != nil {
		inc["status"] = "ALLOCATED"
	}
	b.save()
	return 200, row, nil
}

func (b *Bank) BankOutgoing(body map[string]any) (int, any, error) {
	w := b.wallet(str(body["fromWalletId"]))
	if w == nil {
		return 0, nil, fail(404, "WALLET", "Кошелёк не найден")
	}
	acc := b.account(w.AccountID)
	amount := num(body["amount"])
	if str(body["iban"]) == "" || amount <= 0 {
		return 0, nil, fail(400, "FIELDS", "iban, amount")
	}
	avail := acc.Ledger - b.openHoldSum(acc.ID)
	if avail < amount {
		return 0, nil, fail(409, "INSUFFICIENT", "available", map[string]any{"available": avail})
	}
	id := b.nextID("out")
	hold := map[string]any{"id": b.nextID("hld"), "accountId": acc.ID, "paymentId": id, "amount": amount, "status": "OPEN", "createdAt": now()}
	b.state.Holds = append(b.state.Holds, hold)
	row := map[string]any{
		"id": id, "fromAccountId": acc.ID, "iban": body["iban"], "amount": amount, "currency": acc.Currency,
		"purpose": body["purpose"], "holdId": hold["id"], "status": "PENDING_BANK", "createdAt": now(),
	}
	b.state.Outgoing = append(b.state.Outgoing, row)
	b.save()
	out := clone(row)
	out["available"] = avail - amount
	return 201, out, nil
}

func (b *Bank) BankAck(id string, body map[string]any) (int, any, error) {
	row := find(b.state.Outgoing, id)
	if row == nil {
		return 0, nil, fail(404, "NOT_FOUND", "Исходящий не найден")
	}
	if str(row["status"]) != "PENDING_BANK" {
		return 0, nil, fail(409, "STATE", str(row["status"]))
	}
	hold := find(b.state.Holds, str(row["holdId"]))
	if str(body["result"]) == "rejected" {
		if hold != nil {
			hold["status"] = "RELEASED"
		}
		row["status"] = "REJECTED"
		b.save()
		return 200, row, nil
	}
	jid, err := b.postJournal("IBAN_OUT", str(row["id"]), str(row["iban"]), []map[string]any{
		{"accountId": row["fromAccountId"], "dc": "D", "amount": num(row["amount"])},
		{"accountId": "acc_nostro", "dc": "C", "amount": num(row["amount"])},
	})
	if err != nil {
		return 0, nil, err
	}
	if hold != nil {
		hold["status"] = "CAPTURED"
	}
	row["journalId"] = jid
	row["status"] = "POSTED"
	row["mt103"] = buildMt103(row, b.account(str(row["fromAccountId"])))
	b.save()
	return 200, row, nil
}

func buildMt103(row map[string]any, acc *Account) string {
	d := strings.ReplaceAll(now()[2:10], "-", "")
	iban, name := "", "MALO CUSTOMER"
	if acc != nil {
		iban = str(acc.IBAN)
		name = acc.Name
	}
	return strings.Join([]string{
		"{1:F01MALOUZ22AXXX0000000000}",
		"{2:I103ORNTUZ22XXXXN}",
		"{4:",
		":20:" + str(row["id"]),
		fmt.Sprintf(":32A:%s%s%g,", d, str(row["currency"]), num(row["amount"])),
		":50K:/" + iban,
		name,
		":59:/" + str(row["iban"]),
		":70:" + or(str(row["purpose"]), "TRANSFER"),
		"-}",
	}, "\n")
}

func (b *Bank) SalaryIngest(body map[string]any) (int, any, error) {
	rows, _ := body["rows"].([]any)
	if len(rows) == 0 {
		return 0, nil, fail(400, "ROWS", "rows[]")
	}
	fileName := or(str(body["fileName"]), b.nextID("salfile"))
	if b.state.Contract.SalaryFileOnce {
		for _, s := range b.state.Salary {
			if str(s["fileName"]) == fileName {
				return 0, nil, fail(409, "DUPLICATE_FILE", "Этот файл уже грузили", map[string]any{"existingId": s["id"]})
			}
		}
	}
	sum := 0.0
	for _, r := range rows {
		m, _ := r.(map[string]any)
		sum += num(m["amount"])
	}
	trailer := num(body["trailerSum"])
	if b.state.Contract.SalaryTrailerMustMatch && sum != trailer {
		return 0, nil, fail(400, "TRAILER", "сумма строк ≠ trailer", map[string]any{"sum": sum, "trailer": trailer})
	}
	id := b.nextID("sal")
	posted := []map[string]any{}
	for _, r := range rows {
		m, _ := r.(map[string]any)
		amount := num(m["amount"])
		var acc *Account
		for i := range b.state.Accounts {
			if str(b.state.Accounts[i].IBAN) == str(m["iban"]) {
				acc = &b.state.Accounts[i]
				break
			}
		}
		if acc == nil {
			posted = append(posted, map[string]any{"iban": m["iban"], "amount": amount, "status": "REJECTED", "reason": "IBAN"})
			continue
		}
		jid, err := b.postJournal("SALARY", id, fileName, []map[string]any{
			{"accountId": "acc_nostro", "dc": "D", "amount": amount},
			{"accountId": acc.ID, "dc": "C", "amount": amount},
		})
		if err != nil {
			return 0, nil, err
		}
		posted = append(posted, map[string]any{"iban": acc.IBAN, "amount": amount, "status": "POSTED", "journalId": jid})
	}
	file := map[string]any{"id": id, "fileName": fileName, "trailerSum": trailer, "rowSum": sum, "trailerOk": sum == trailer, "status": "POSTED", "rows": posted, "createdAt": now()}
	b.state.Salary = append(b.state.Salary, file)
	b.save()
	return 201, file, nil
}

func (b *Bank) FXConvert(body map[string]any) (int, any, error) {
	fromW, toW := b.wallet(str(body["fromWalletId"])), b.wallet(str(body["toWalletId"]))
	if fromW == nil || toW == nil {
		return 0, nil, fail(404, "WALLET", "from/to кошелёк")
	}
	from, to := b.account(fromW.AccountID), b.account(toW.AccountID)
	if from.Currency == to.Currency {
		return 0, nil, fail(409, "CCY", "Нужны разные валюты")
	}
	if !(from.Currency == "UZS" && to.Currency == "USD") {
		return 0, nil, fail(409, "PAIR", "Сейчас пет меняет только UZS → USD")
	}
	amountTo := num(body["amountTo"])
	if amountTo <= 0 {
		return 0, nil, fail(400, "AMOUNT", "amountTo > 0 (USD)")
	}
	rate := 12650.0
	for _, f := range b.state.FX {
		if str(f["pair"]) == "USD/UZS" {
			rate = num(f["rate"])
		}
	}
	spreadBps := b.state.Contract.FxSpreadBps
	mid := round(amountTo * rate)
	spread := round((mid * spreadBps) / 10000)
	uzsDebit := mid + spread
	avail := from.Ledger - b.openHoldSum(from.ID)
	if avail < uzsDebit {
		return 0, nil, fail(409, "INSUFFICIENT", "available UZS", map[string]any{"available": avail, "need": uzsDebit})
	}
	id := b.nextID("fx")
	var uzsLegs []map[string]any
	if spread > 0 {
		uzsLegs = []map[string]any{
			{"accountId": from.ID, "dc": "D", "amount": uzsDebit},
			{"accountId": "acc_fx_uzs", "dc": "C", "amount": mid},
			{"accountId": "acc_fee", "dc": "C", "amount": spread, "note": "fx spread"},
		}
	} else {
		uzsLegs = []map[string]any{
			{"accountId": from.ID, "dc": "D", "amount": mid},
			{"accountId": "acc_fx_uzs", "dc": "C", "amount": mid},
		}
	}
	uzsJ, err := b.postJournal("FX_UZS", id, fmt.Sprintf("USD %g @ %g", amountTo, rate), uzsLegs)
	if err != nil {
		return 0, nil, err
	}
	usdJ, err := b.postJournal("FX_USD", id, fmt.Sprintf("buy %g USD", amountTo), []map[string]any{
		{"accountId": "acc_fx_usd", "dc": "D", "amount": amountTo},
		{"accountId": to.ID, "dc": "C", "amount": amountTo},
	})
	if err != nil {
		return 0, nil, err
	}
	deal := map[string]any{
		"id": id, "pair": "USD/UZS", "rate": rate, "spreadBps": spreadBps, "amountUsd": amountTo,
		"amountUzs": uzsDebit, "spread": spread, "uzsJournalId": uzsJ, "usdJournalId": usdJ,
		"status": "POSTED", "createdAt": now(),
	}
	b.state.FXDeals = append(b.state.FXDeals, deal)
	b.save()
	return 201, map[string]any{"deal": deal, "trial": b.TrialBalance()}, nil
}

func (b *Bank) ClearingSettle() (int, any, error) {
	merch := b.account("acc_merchant")
	amount := merch.Ledger
	if amount <= 0 {
		return 0, nil, fail(409, "NOTHING", "Мерчанту нечего селить — сначала capture", map[string]any{"ledger": amount})
	}
	id := b.nextID("clr")
	jid, err := b.postJournal("CLEARING_SETTLE", id, "T+1 merchant payout", []map[string]any{
		{"accountId": "acc_merchant", "dc": "D", "amount": amount},
		{"accountId": "acc_nostro", "dc": "C", "amount": amount},
	})
	if err != nil {
		return 0, nil, err
	}
	row := map[string]any{"id": id, "amount": amount, "journalId": jid, "status": "SETTLED", "createdAt": now()}
	b.state.Clearing = append(b.state.Clearing, row)
	b.save()
	out := clone(row)
	out["trial"] = b.TrialBalance()
	return 201, out, nil
}

func (b *Bank) ReverseJournal(id string) (int, any, error) {
	j := find(b.state.Journals, id)
	if j == nil {
		return 0, nil, fail(404, "NOT_FOUND", "Журнал не найден")
	}
	if j["reversedBy"] != nil {
		return 0, nil, fail(409, "STATE", "Уже сторнирован")
	}
	legs := []map[string]any{}
	for _, l := range b.state.Ledger {
		if str(l["journalId"]) == id && liveLeg(l) {
			dc := "C"
			if str(l["dc"]) == "C" {
				dc = "D"
			}
			legs = append(legs, map[string]any{"accountId": l["accountId"], "dc": dc, "amount": num(l["amount"])})
		}
	}
	if len(legs) == 0 {
		return 0, nil, fail(409, "EMPTY", "Нет ног")
	}
	rid, err := b.postJournal(str(j["kind"])+"_REV", str(j["ref"]), "storno "+str(j["id"]), legs)
	if err != nil {
		return 0, nil, err
	}
	j["reversedBy"] = rid
	b.save()
	return 200, map[string]any{"original": j, "reverseId": rid, "trial": b.TrialBalance()}, nil
}

func (b *Bank) ExportRecon() map[string]any {
	rows := []map[string]any{}
	for _, p := range b.state.Payments {
		if str(p["kind"]) == "CARD" && in(str(p["status"]), "CAPTURED", "REFUNDED", "CHARGEBACK") && p["rrn"] != nil {
			amt := num(p["capturedAmount"])
			if amt == 0 {
				amt = num(p["amount"])
			}
			rows = append(rows, map[string]any{"rrn": p["rrn"], "amount": amt, "merchant": p["toAccountId"], "status": p["status"], "capturedAt": or(str(p["capturedAt"]), str(p["createdAt"]))})
		}
	}
	return map[string]any{"source": "malo", "generatedAt": now(), "rows": rows}
}

func (b *Bank) IngestRecon(body map[string]any) (int, any, error) {
	ours := b.ExportRecon()["rows"].([]map[string]any)
	incoming, _ := body["rows"].([]any)
	by := map[string]map[string]any{}
	for _, r := range ours {
		by[str(r["rrn"])] = r
	}
	items := []map[string]any{}
	seen := map[string]bool{}
	for _, row := range incoming {
		m, _ := row.(map[string]any)
		seen[str(m["rrn"])] = true
		mine := by[str(m["rrn"])]
		if mine == nil {
			items = append(items, merge(m, map[string]any{"result": "UNKNOWN_RRN"}))
			continue
		}
		if num(m["amount"]) != num(mine["amount"]) {
			items = append(items, merge(m, map[string]any{"result": "MISMATCH", "ours": mine["amount"]}))
			continue
		}
		items = append(items, merge(m, map[string]any{"result": "MATCH"}))
	}
	for _, o := range ours {
		if !seen[str(o["rrn"])] {
			items = append(items, map[string]any{"rrn": o["rrn"], "amount": o["amount"], "result": "MISSING_IN_FILE"})
		}
	}
	run := map[string]any{"id": b.nextID("rcn"), "source": or(str(body["source"]), "partner"), "at": now(), "items": items}
	b.state.Recon = append([]map[string]any{run}, b.state.Recon...)
	b.save()
	return 200, run, nil
}

func (b *Bank) StartKYC(key string, body map[string]any) (int, any, error) {
	if str(body["customer"]) == "" {
		return 0, nil, fail(400, "CUSTOMER", "customer обязателен")
	}
	if key != "" {
		for _, k := range b.state.KYC {
			if str(k["idempotencyKey"]) == key {
				return 200, k, nil
			}
		}
	}
	row := map[string]any{"id": b.nextID("kyc"), "customer": body["customer"], "status": "PENDING", "reason": "новая заявка", "idempotencyKey": nilIfEmpty(key), "createdAt": now()}
	b.state.KYC = append(b.state.KYC, row)
	b.save()
	return 201, row, nil
}

func (b *Bank) DecideKYC(id string, body map[string]any) (int, any, error) {
	row := find(b.state.KYC, id)
	if row == nil {
		return 0, nil, fail(404, "NOT_FOUND", "KYC не найден")
	}
	st := str(body["status"])
	if st != "APPROVED" && st != "REJECTED" {
		return 0, nil, fail(400, "STATUS", "APPROVED или REJECTED")
	}
	row["status"] = st
	row["reason"] = body["reason"]
	row["decidedAt"] = now()
	b.save()
	return 200, row, nil
}

func (b *Bank) CreateOrder(key string, body map[string]any) (int, any, error) {
	if b.state.Contract.OrderIdempotencyRequired && key == "" {
		return 0, nil, fail(400, "IDEMPOTENCY_KEY_REQUIRED", "Idempotency-Key обязателен")
	}
	if key != "" {
		for _, o := range b.state.Orders {
			if str(o["idempotencyKey"]) == key {
				if b.state.Contract.TransferDuplicateHTTP == 409 {
					return 409, map[string]any{"error": map[string]any{"code": "DUPLICATE"}, "existingId": o["id"]}, nil
				}
				return 200, o, nil
			}
		}
	}
	var sku map[string]any
	for _, s := range b.state.Stock {
		if str(s["sku"]) == str(body["sku"]) {
			sku = s
			break
		}
	}
	if sku == nil {
		return 0, nil, fail(404, "SKU", "Нет SKU")
	}
	qty := num(body["qty"])
	if qty == 0 {
		qty = 1
	}
	if num(sku["qty"]) < qty {
		return 0, nil, fail(409, "STOCK", "Нет остатка", map[string]any{"qty": sku["qty"]})
	}
	sku["qty"] = num(sku["qty"]) - qty
	order := map[string]any{"id": b.nextID("ord"), "sku": sku["sku"], "qty": qty, "sellerId": body["sellerId"], "status": "CREATED", "idempotencyKey": nilIfEmpty(key), "createdAt": now()}
	b.state.Orders = append(b.state.Orders, order)
	b.save()
	return 201, order, nil
}

func (b *Bank) GetPayment(id string) (map[string]any, error) {
	p := find(b.state.Payments, id)
	if p == nil {
		return nil, fail(404, "NOT_FOUND", "Платёж не найден")
	}
	return p, nil
}

func (b *Bank) GetWallet(id string) (map[string]any, error) {
	w := b.wallet(id)
	if w == nil {
		return nil, fail(404, "NOT_FOUND", "Кошелёк не найден")
	}
	return b.walletView(*w), nil
}

func (b *Bank) MT103(id string) (map[string]any, error) {
	row := find(b.state.Outgoing, id)
	if row == nil {
		return nil, fail(404, "NOT_FOUND", "Исходящий не найден")
	}
	if row["mt103"] == nil {
		return nil, fail(409, "STATE", "MT103 после posted")
	}
	return map[string]any{"id": row["id"], "mt103": row["mt103"], "status": row["status"]}, nil
}

func find(list []map[string]any, id string) map[string]any {
	for _, x := range list {
		if str(x["id"]) == id {
			return x
		}
	}
	return nil
}

func or(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func in(v string, opts ...string) bool {
	for _, o := range opts {
		if v == o {
			return true
		}
	}
	return false
}

func some(list []map[string]any, fn func(map[string]any) bool) bool {
	for _, x := range list {
		if fn(x) {
			return true
		}
	}
	return false
}

func somePay(list []map[string]any, fn func(map[string]any) bool) bool { return some(list, fn) }
func someLeg(list []map[string]any, fn func(map[string]any) bool) bool { return some(list, fn) }

func reconBad(list []map[string]any) bool {
	for _, r := range list {
		items, _ := r["items"].([]map[string]any)
		if items == nil {
			if raw, ok := r["items"].([]any); ok {
				for _, it := range raw {
					m, _ := it.(map[string]any)
					if in(str(m["result"]), "MISMATCH", "UNKNOWN_RRN") {
						return true
					}
				}
			}
			continue
		}
		for _, i := range items {
			if in(str(i["result"]), "MISMATCH", "UNKNOWN_RRN") {
				return true
			}
		}
	}
	return false
}

func merge(a, b map[string]any) map[string]any {
	o := clone(a)
	for k, v := range b {
		o[k] = v
	}
	return o
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
