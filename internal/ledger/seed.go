package ledger

import "time"

func seedBase() *State {
	return &State{
		Version: Version,
		Contract: Contract{
			TransferDuplicateHTTP: 200,
			ErrorInBody:           true,
			AcsTimeoutStatus:      "SUCCESS",
			DailyP2pLimit:         1000000,
		},
		Accounts: []Account{
			{ID: "acc_anna", WalletID: "wal_anna", IBAN: "UZ12MALO000000000001", Name: "Анна · текущий UZS", Type: "LIAB_CUSTOMER", Normal: "credit", Currency: "UZS", Ledger: 150000},
			{ID: "acc_boris", WalletID: "wal_boris", IBAN: "UZ12MALO000000000002", Name: "Борис · текущий UZS", Type: "LIAB_CUSTOMER", Normal: "credit", Currency: "UZS", Ledger: 80000},
			{ID: "acc_dina", WalletID: "wal_dina", IBAN: "UZ12MALO000000000003", Name: "Дина · текущий UZS", Type: "LIAB_CUSTOMER", Normal: "credit", Currency: "UZS", Ledger: 80000},
			{ID: "acc_elena", WalletID: "wal_elena", IBAN: "UZ12MALO000000000004", Name: "Елена · текущий UZS", Type: "LIAB_CUSTOMER", Normal: "credit", Currency: "UZS", Ledger: 42000},
			{ID: "acc_usd_anna", WalletID: "wal_anna_usd", IBAN: "UZ12MALO000000USD001", Name: "Анна · текущий USD", Type: "LIAB_CUSTOMER", Normal: "credit", Currency: "USD", Ledger: 200},
			{ID: "acc_merchant", WalletID: nil, IBAN: "UZ12CAFE000000000009", Name: "Cafe Orient · MCC 5812", Type: "LIAB_MERCHANT", Normal: "credit", Currency: "UZS", Ledger: 0},
			{ID: "acc_taxi", WalletID: nil, IBAN: "UZ12TAXI000000000011", Name: "Taxi Silk · MCC 4121", Type: "LIAB_MERCHANT", Normal: "credit", Currency: "UZS", Ledger: 0},
			{ID: "acc_pharma", WalletID: nil, IBAN: "UZ12PHAR000000000012", Name: "Pharma Plus · MCC 5912", Type: "LIAB_MERCHANT", Normal: "credit", Currency: "UZS", Ledger: 0},
			{ID: "acc_nostro", WalletID: nil, IBAN: "UZ12ORNT000000NOSTRO", Name: "Nostro Orient Bank UZS", Type: "ASSET_NOSTRO", Normal: "debit", Currency: "UZS", Ledger: 230000},
			{ID: "acc_nostro_usd", WalletID: nil, IBAN: "UZ12ORNT000000NOSTUSD", Name: "Nostro Orient Bank USD", Type: "ASSET_NOSTRO", Normal: "debit", Currency: "USD", Ledger: 500},
			{ID: "acc_fee", WalletID: nil, IBAN: nil, Name: "Доход комиссия Malo", Type: "INC_FEE", Normal: "credit", Currency: "UZS", Ledger: 0},
			{ID: "acc_cb", WalletID: nil, IBAN: nil, Name: "Расход chargeback", Type: "EXP_CHARGEBACK", Normal: "debit", Currency: "UZS", Ledger: 0},
			{ID: "acc_suspense", WalletID: nil, IBAN: nil, Name: "Невыясненные суммы", Type: "LIAB_SUSPENSE", Normal: "credit", Currency: "UZS", Ledger: 0},
			{ID: "acc_clearing", WalletID: nil, IBAN: nil, Name: "Клиринг карт T+1", Type: "LIAB_CLEARING", Normal: "credit", Currency: "UZS", Ledger: 0},
			{ID: "acc_fx_uzs", WalletID: nil, IBAN: nil, Name: "FX позиция UZS", Type: "POS_FX", Normal: "credit", Currency: "UZS", Ledger: 0},
			{ID: "acc_fx_usd", WalletID: nil, IBAN: nil, Name: "FX позиция USD", Type: "POS_FX", Normal: "debit", Currency: "USD", Ledger: 0},
		},
		Wallets: []Wallet{
			{ID: "wal_anna", Customer: "anna", Name: "Анна UZS", AccountID: "acc_anna", Currency: "UZS", Status: "ACTIVE"},
			{ID: "wal_anna_usd", Customer: "anna", Name: "Анна USD", AccountID: "acc_usd_anna", Currency: "USD", Status: "ACTIVE"},
			{ID: "wal_boris", Customer: "boris", Name: "Борис UZS", AccountID: "acc_boris", Currency: "UZS", Status: "ACTIVE"},
			{ID: "wal_dina", Customer: "dina", Name: "Дина UZS", AccountID: "acc_dina", Currency: "UZS", Status: "ACTIVE"},
			{ID: "wal_elena", Customer: "elena", Name: "Елена UZS", AccountID: "acc_elena", Currency: "UZS", Status: "ACTIVE"},
		},
		KYC: []map[string]any{
			{"id": "kyc_anna", "customer": "anna", "status": "PENDING", "reason": "ожидание документов"},
			{"id": "kyc_boris", "customer": "boris", "status": "APPROVED", "reason": nil},
			{"id": "kyc_dina", "customer": "dina", "status": "APPROVED", "reason": nil},
			{"id": "kyc_elena", "customer": "elena", "status": "APPROVED", "reason": nil},
		},
		Plastic: []map[string]any{
			{"id": "crd_anna", "walletId": "wal_anna", "last4": "4412", "scheme": "MIR", "status": "ACTIVE", "createdAt": "2026-01-01T00:00:00.000Z"},
			{"id": "crd_boris", "walletId": "wal_boris", "last4": "8821", "scheme": "VISA", "status": "ACTIVE", "createdAt": "2026-01-01T00:00:00.000Z"},
			{"id": "crd_dina", "walletId": "wal_dina", "last4": "1099", "scheme": "UZCARD", "status": "ACTIVE", "createdAt": "2026-03-12T00:00:00.000Z"},
			{"id": "crd_elena", "walletId": "wal_elena", "last4": "7740", "scheme": "HUMO", "status": "BLOCKED", "createdAt": "2026-04-02T00:00:00.000Z"},
		},
		FX:    []map[string]any{{"pair": "USD/UZS", "rate": 12650.0, "asOf": "2026-08-18"}},
		Stock: []map[string]any{{"sku": "SKU-100", "title": "Чайник", "qty": 12.0}, {"sku": "SKU-200", "title": "Кофемолка", "qty": 3.0}},
		Seq:   40,
	}
}

func (b *Bank) plantTeaching() {
	st := b.state
	st.Teaching = true
	st.Customers = []map[string]any{
		{"id": "anna", "name": "Анна Рахимова", "segment": "RETAIL", "city": "Ташкент", "openedAt": "2025-11-02T09:00:00.000Z"},
		{"id": "boris", "name": "Борис Ким", "segment": "RETAIL", "city": "Ташкент", "openedAt": "2025-12-18T11:20:00.000Z"},
		{"id": "dina", "name": "Дина Саидова", "segment": "RETAIL", "city": "Самарканд", "openedAt": "2026-02-09T08:10:00.000Z"},
		{"id": "elena", "name": "Елена Пак", "segment": "AFFLUENT", "city": "Ташкент", "openedAt": "2026-01-14T10:00:00.000Z"},
	}
	st.Merchants = []map[string]any{
		{"id": "mcc_cafe", "accountId": "acc_merchant", "name": "Cafe Orient", "mcc": "5812", "city": "Ташкент"},
		{"id": "mcc_taxi", "accountId": "acc_taxi", "name": "Taxi Silk", "mcc": "4121", "city": "Ташкент"},
		{"id": "mcc_pharma", "accountId": "acc_pharma", "name": "Pharma Plus", "mcc": "5912", "city": "Самарканд"},
	}

	day := time.Now().UTC().Add(-6 * time.Hour).Format("2006-01-02T15:04:05.000Z")
	night := time.Now().UTC().Add(-10 * time.Hour)
	if night.Hour() < 22 {
		night = time.Date(night.Year(), night.Month(), night.Day(), 22, 41, 0, 0, time.UTC)
		if night.After(time.Now().UTC()) {
			night = night.Add(-24 * time.Hour)
		}
	}
	nightAt := night.Format("2006-01-02T15:04:05.000Z")
	oldAt := time.Now().UTC().Add(-36 * time.Hour).Format("2006-01-02T15:04:05.000Z")

	j1, _ := b.postJournal("P2P", "tr_dup_1", "p2p without key", []map[string]any{
		{"accountId": "acc_anna", "dc": "D", "amount": 5000},
		{"accountId": "acc_boris", "dc": "C", "amount": 5000},
	})
	j2, _ := b.postJournal("P2P", "tr_dup_2", "p2p without key — повтор линии", []map[string]any{
		{"accountId": "acc_anna", "dc": "D", "amount": 5000},
		{"accountId": "acc_boris", "dc": "C", "amount": 5000},
	})
	b.stampJournal(j1, day, true, false)
	b.stampJournal(j2, day, true, false)
	st.Transfers = append(st.Transfers,
		map[string]any{"id": "tr_dup_1", "fromWalletId": "wal_anna", "toWalletId": "wal_boris", "fromAccountId": "acc_anna", "toAccountId": "acc_boris", "amount": 5000.0, "currency": "UZS", "status": "SUCCESS", "journalId": j1, "idempotencyKey": nil, "createdAt": day},
		map[string]any{"id": "tr_dup_2", "fromWalletId": "wal_anna", "toWalletId": "wal_boris", "fromAccountId": "acc_anna", "toAccountId": "acc_boris", "amount": 5000.0, "currency": "UZS", "status": "SUCCESS", "journalId": j2, "idempotencyKey": nil, "createdAt": day},
	)

	jNight, _ := b.postJournal("P2P", "tr_night", "перевод ночной смены", []map[string]any{
		{"accountId": "acc_elena", "dc": "D", "amount": 12000},
		{"accountId": "acc_boris", "dc": "C", "amount": 12000},
	})
	b.stampJournal(jNight, nightAt, false, true)
	st.Transfers = append(st.Transfers, map[string]any{
		"id": "tr_night", "fromWalletId": "wal_elena", "toWalletId": "wal_boris", "fromAccountId": "acc_elena", "toAccountId": "acc_boris",
		"amount": 12000.0, "currency": "UZS", "status": "SUCCESS", "journalId": jNight, "idempotencyKey": "night-manual", "createdAt": nightAt, "fourEyes": false, "nightShift": true,
	})

	jHole, _ := b.postJournal("REFUND_HOLE", "pay_hole", "возврат без CR мерчанта — дыра книги", []map[string]any{
		{"accountId": "acc_anna", "dc": "D", "amount": 3000},
	})
	b.stampJournal(jHole, day, true, false)

	jIn, _ := b.postJournal("SUSPENSE_IN", "in_old", "unknown iban salary leftover", []map[string]any{
		{"accountId": "acc_nostro", "dc": "D", "amount": 175000},
		{"accountId": "acc_suspense", "dc": "C", "amount": 175000},
	})
	b.stampJournal(jIn, oldAt, true, false)
	inc := map[string]any{"id": "in_old", "iban": "UZ00UNKNOWN00000000999", "amount": 175000.0, "currency": "UZS", "paymentRef": "SALARY-SEP", "journalId": jIn, "status": "SUSPENSE", "createdAt": oldAt}
	st.Incoming = append(st.Incoming, inc)
	sus := clone(inc)
	sus["allocated"] = false
	st.Suspense = append(st.Suspense, sus)

	st.Holds = append(st.Holds,
		map[string]any{"id": "hld_dina_over", "accountId": "acc_dina", "walletId": "wal_dina", "paymentId": "pay_dina_auth", "amount": 200000.0, "currency": "UZS", "status": "OPEN", "createdAt": day, "source": "CARD"},
		map[string]any{"id": "hld_boris_cafe", "accountId": "acc_boris", "walletId": "wal_boris", "paymentId": "pay_cafe_auth_1", "amount": 18000.0, "currency": "UZS", "status": "OPEN", "createdAt": day, "source": "CARD"},
	)

	plantPay := func(id, status, merchant, wallet, acc string, amount float64, at string, hold any) {
		st.Payments = append(st.Payments, map[string]any{
			"id": id, "kind": "CARD", "status": status, "amount": amount, "currency": "UZS",
			"fromAccountId": acc, "toAccountId": merchant, "merchantId": merchant, "walletId": wallet,
			"rrn": "rrn_" + id, "holdId": hold, "createdAt": at,
			"capturedAmount": map[bool]any{true: amount, false: 0.0}[status == "CAPTURED"],
		})
		authStatus := status
		if status == "AUTHORIZED" {
			authStatus = "AUTHORIZED"
		}
		st.Cards = append(st.Cards, map[string]any{"id": "auth_" + id, "paymentId": id, "walletId": wallet, "merchantId": merchant, "amount": amount, "rrn": "rrn_" + id, "acs": "OK", "status": authStatus, "createdAt": at})
	}
	plantPay("pay_cafe_auth_1", "AUTHORIZED", "acc_merchant", "wal_boris", "acc_boris", 18000, day, "hld_boris_cafe")
	plantPay("pay_cafe_auth_2", "AUTHORIZED", "acc_merchant", "wal_anna", "acc_anna", 9000, day, nil)
	plantPay("pay_cafe_cap_1", "CAPTURED", "acc_merchant", "wal_anna", "acc_anna", 14000, day, nil)
	jCafe, _ := b.postJournal("CARD_CAPTURE", "pay_cafe_cap_1", "RRN rrn_pay_cafe_cap_1", []map[string]any{
		{"accountId": "acc_anna", "dc": "D", "amount": 14000},
		{"accountId": "acc_merchant", "dc": "C", "amount": 14000},
	})
	b.stampJournal(jCafe, day, true, false)
	find(st.Payments, "pay_cafe_cap_1")["journalId"] = jCafe

	plantPay("pay_taxi_auth_1", "AUTHORIZED", "acc_taxi", "wal_elena", "acc_elena", 22000, day, nil)
	for i, amt := range []float64{8000, 11000, 6500, 4000} {
		id := "pay_taxi_cap_" + string(rune('1'+i))
		plantPay(id, "CAPTURED", "acc_taxi", "wal_boris", "acc_boris", amt, day, nil)
		jid, _ := b.postJournal("CARD_CAPTURE", id, "taxi capture", []map[string]any{
			{"accountId": "acc_boris", "dc": "D", "amount": amt},
			{"accountId": "acc_taxi", "dc": "C", "amount": amt},
		})
		b.stampJournal(jid, day, true, false)
		find(st.Payments, id)["journalId"] = jid
	}
	for i, amt := range []float64{15000, 9000, 7000} {
		id := "pay_pharma_" + string(rune('1'+i))
		stt := "AUTHORIZED"
		if i > 0 {
			stt = "CAPTURED"
		}
		plantPay(id, stt, "acc_pharma", "wal_anna", "acc_anna", amt, day, nil)
		if stt == "CAPTURED" {
			jid, _ := b.postJournal("CARD_CAPTURE", id, "pharma capture", []map[string]any{
				{"accountId": "acc_anna", "dc": "D", "amount": amt},
				{"accountId": "acc_pharma", "dc": "C", "amount": amt},
			})
			b.stampJournal(jid, day, true, false)
			find(st.Payments, id)["journalId"] = jid
		}
	}

	st.FXDeals = append(st.FXDeals,
		map[string]any{"id": "fx_naked", "pair": "USD/UZS", "rate": 12650.0, "amountUsd": 50.0, "amountUzs": 632500.0, "status": "BOOKED", "coverJournalId": nil, "uzsJournalId": nil, "usdJournalId": nil, "createdAt": day, "note": "сделка казначейства без покрытия"},
	)

	st.WalletUI = []map[string]any{
		{"walletId": "wal_anna", "shownAvailable": nil},
		{"walletId": "wal_boris", "shownAvailable": nil},
		{"walletId": "wal_dina", "shownAvailable": 80000.0},
		{"walletId": "wal_elena", "shownAvailable": 0.0},
	}

	st.Tickets = []map[string]any{
		{"id": "tkt_hold", "customer": "dina", "kind": "HOLD", "status": "OPEN", "priority": "P1", "title": "Клиент: «карта не проходит, а в приложении минус»", "body": "Дина на линии: available ушёл в минус после auth в аптеке. Операции просят сверку холда с книгой.", "openedAt": day, "slaHours": 4.0},
		{"id": "tkt_p2p", "customer": "anna", "kind": "P2P", "status": "OPEN", "priority": "P1", "title": "Два списания по одному переводу Анна → Борис", "body": "Саппорт: клиентка нажала «ещё раз» без ключа. В выписке два DR по 5000.", "openedAt": day, "slaHours": 8.0},
		{"id": "tkt_suspense", "customer": "", "kind": "IBAN", "status": "OPEN", "priority": "P2", "title": "Невыясненная сумма зарплатного файла старше суток", "body": "Операции ночи: IBAN не наш, сидит в suspense. HR звонит «где деньги».", "openedAt": oldAt, "slaHours": 24.0},
		{"id": "tkt_fx", "customer": "anna", "kind": "FX", "status": "OPEN", "priority": "P2", "title": "FX BOOKED без журнала покрытия", "body": "Казначейство: сделка 50 USD в книге сделок есть, ног нет.", "openedAt": day, "slaHours": 12.0},
		{"id": "tkt_night", "customer": "elena", "kind": "OPS", "status": "OPEN", "priority": "P2", "title": "Ночной P2P без четырёх глаз", "body": "Контроль: после 22:00 перевод 12 000 без второго подтверждения.", "openedAt": nightAt, "slaHours": 12.0},
	}

	b.syncWallets()
	st.Seq = 80
}

func (b *Bank) TeachingDump() map[string]any {
	accs := make([]map[string]any, 0, len(b.state.Accounts))
	for _, a := range b.state.Accounts {
		hold := b.openHoldSum(a.ID)
		accs = append(accs, map[string]any{
			"id": a.ID, "walletId": a.WalletID, "iban": a.IBAN, "name": a.Name,
			"type": a.Type, "normal": a.Normal, "currency": a.Currency,
			"ledger": a.Ledger, "hold": hold, "available": a.Ledger - hold,
		})
	}
	return map[string]any{
		"customers": b.state.Customers, "accounts": accs, "wallets": b.walletViews(),
		"holds": b.state.Holds, "journals": b.state.Journals, "ledger": b.state.Ledger,
		"transfers": b.state.Transfers, "payments": b.state.Payments, "cards": b.state.Plastic,
		"auths": b.state.Cards, "incoming": b.state.Incoming, "outgoing": b.state.Outgoing,
		"suspense": b.state.Suspense, "fxDeals": b.state.FXDeals, "merchants": b.state.Merchants,
		"tickets": b.state.Tickets, "audit": b.state.Audit, "walletUi": b.state.WalletUI,
		"kyc": b.state.KYC, "trial": b.TrialBalance(),
	}
}

func (b *Bank) stampJournal(jid, at string, fourEyes, night bool) {
	for _, j := range b.state.Journals {
		if str(j["id"]) == jid {
			j["at"] = at
			j["fourEyes"] = fourEyes
			j["nightShift"] = night
		}
	}
	for _, l := range b.state.Ledger {
		if str(l["journalId"]) == jid {
			l["at"] = at
			l["fourEyes"] = fourEyes
			l["nightShift"] = night
		}
	}
}
