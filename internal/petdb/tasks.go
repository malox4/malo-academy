package petdb

import (
	"encoding/json"
	"sort"
	"strings"
)

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Story       string `json:"story"`
	Question    string `json:"question"`
	Hint        string `json:"hint"`
	Why         string `json:"why"`
	Gold        string `json:"-"`
	GoldPublic  string `json:"gold,omitempty"`
	InternSafe  bool   `json:"internSafe"`
	Grade       string `json:"grade"`
}

func Tasks() []Task {
	return []Task{
		{
			ID: "sql-intern-wallets", Title: "Витрина кошельков", Grade: "intern", InternSafe: true,
			Story: "Клиент видит в приложении доступный остаток. На разборе тикета BA сначала читает витрину intern_wallets — без холдов, поле available.",
			Question: "Выберите id, имя, валюту и available из intern_wallets.",
			Hint: "SELECT id, name, currency, available FROM intern_wallets",
			Why: "Available — то, что клиент может потратить. Это витрина. Книга (ledger) может отличаться из‑за холдов.",
			Gold: "SELECT id, name, currency, available FROM intern_wallets",
		},
		{
			ID: "sql-intern-holds", Title: "Открытые холды", Grade: "intern", InternSafe: true,
			Story: "Касса вернула SUCCESS, клиент не может потратить деньги. Проверьте, есть ли открытый холд — резерв без проводки в журнале.",
			Question: "Найдите холды со статусом OPEN в intern_holds.",
			Hint: "SELECT id, wallet_id, amount, status FROM intern_holds WHERE status = 'OPEN'",
			Why: "Холд — не проводка. Нужно отличить резерв (OPEN) от списания в журнале.",
			Gold: "SELECT id, wallet_id, amount, status FROM intern_holds WHERE status = 'OPEN'",
		},
		{
			ID: "sql-hold-negative", Title: "Available ушёл в минус", Grade: "junior",
			Story: "Тикет tkt_hold: Дина не может заплатить, в приложении отрицательный available. Нужна сверка холда с книгой.",
			Question: "Найдите кошельки, у которых есть открытый холд и available < 0 (ledger − hold).",
			Hint: "wallets.hold > 0 AND ledger - hold < 0. Не путайте с wallet_ui.",
			Why: "Классика линии: authorize больше остатка из‑за гонки или забытого холда. Клиенту показывают available, не ledger.",
			Gold: `SELECT id, customer_id, ledger, hold, (ledger - hold) AS available
FROM wallets
WHERE hold > 0 AND (ledger - hold) < 0`,
		},
		{
			ID: "sql-ui-mismatch", Title: "Экран не совпадает с книгой", Grade: "junior",
			Story: "Елена: в приложении ноль, в отделении говорят, что деньги есть. Таблица wallet_ui — то, что рисует фронт.",
			Question: "Найдите кошельки, где shown_available не равен ledger − hold (и shown_available задан).",
			Hint: "JOIN wallets и wallet_ui. NULL в витрине означает «считай сам», это не ошибка.",
			Why: "Частый дефект витрины: кэш available. Аналитик сверяет экран с книгой.",
			Gold: `SELECT w.id, (w.ledger - w.hold) AS book_available, u.shown_available
FROM wallets w
JOIN wallet_ui u ON u.wallet_id = w.id
WHERE u.shown_available IS NOT NULL
  AND u.shown_available != (w.ledger - w.hold)`,
		},
		{
			ID: "sql-p2p-dup", Title: "P2P без ключа: два перевода", Grade: "junior",
			Story: "Тикет tkt_p2p: Анна отправила перевод дважды. Idempotency-Key пустой. Учебный дефект ядра: оба запроса прошли.",
			Question: "Найдите пары from/to/amount/currency, где переводов больше одного и ключ пустой.",
			Hint: "GROUP BY ... HAVING COUNT(*) > 1. idempotency_key IS NULL.",
			Why: "Инцидент «списали дважды». Аналитик доказывает дубль запросом, не скрином чата.",
			Gold: `SELECT from_wallet_id, to_wallet_id, amount, currency, COUNT(*) AS cnt
FROM transfers
WHERE idempotency_key IS NULL
GROUP BY from_wallet_id, to_wallet_id, amount, currency
HAVING COUNT(*) > 1`,
		},
		{
			ID: "sql-suspense-age", Title: "IBAN в невыясненных старше 12 часов", Grade: "junior",
			Story: "Ночная смена: зарплатный платёж с неизвестным IBAN лежит в suspense. HR спрашивает, где деньги сотрудников.",
			Question: "Покажите suspense со статусом SUSPENSE, открытые больше 12 часов.",
			Hint: "created_at <= datetime('now', '-12 hours')",
			Why: "Невыясненные — чужие деньги. SLA старения — операционное требование.",
			Gold: `SELECT id, iban, amount, status, created_at
FROM suspense
WHERE status = 'SUSPENSE'
  AND created_at <= datetime('now', '-12 hours')`,
		},
		{
			ID: "sql-trial", Title: "Пробный баланс не сходится", Grade: "middle",
			Story: "Контроль: по UZS сумма дебетов не равна сумме кредитов. Возврат провели одной стороной.",
			Question: "Посчитайте сумму DR, CR и diff по живым проводкам ledger_legs в UZS.",
			Hint: "status != 'REVERSED'. diff = debit − credit. Ищите kind REFUND_HOLE.",
			Why: "Пробный баланс — первый контроль книги. Односторонняя проводка не чинится правкой Excel.",
			Gold: `SELECT
  SUM(CASE WHEN dc = 'D' THEN amount ELSE 0 END) AS debit,
  SUM(CASE WHEN dc = 'C' THEN amount ELSE 0 END) AS credit,
  SUM(CASE WHEN dc = 'D' THEN amount ELSE -amount END) AS diff
FROM ledger_legs
WHERE currency = 'UZS' AND status != 'REVERSED'`,
		},
		{
			ID: "sql-merchant-auth-cap", Title: "Мерчанты: AUTHORIZED vs CAPTURE", Grade: "middle",
			Story: "Эквайринг: Cafe Orient сообщает, что касса зелёная, а выплат нет. Сравните количество AUTHORIZED и CAPTURED по мерчантам.",
			Question: "По каждому merchant_account_id: сколько AUTHORIZED и сколько CAPTURED (и рефанды/чарджбек как captured-контур).",
			Hint: "GROUP BY merchant_account_id. CASE по status.",
			Why: "Auth без capture — деньги в холде, мерчант ещё не заработал. Путать их на сверке — ложный payout.",
			Gold: `SELECT merchant_account_id,
  SUM(CASE WHEN status = 'AUTHORIZED' THEN 1 ELSE 0 END) AS authorized,
  SUM(CASE WHEN status IN ('CAPTURED', 'REFUNDED', 'CHARGEBACK') THEN 1 ELSE 0 END) AS captured
FROM payments
WHERE kind = 'CARD'
GROUP BY merchant_account_id
ORDER BY merchant_account_id`,
		},
		{
			ID: "sql-fx-naked", Title: "FX без покрытия", Grade: "middle",
			Story: "Тикет tkt_fx: казначейство видит сделку 50 USD, в журнале нет покрытия. Позиция открыта без проводок.",
			Question: "Найдите FX-сделки без cover/uzs журнала или со статусом BOOKED.",
			Hint: "cover_journal_id IS NULL AND uzs_journal_id IS NULL, либо status = 'BOOKED'.",
			Why: "FX без журнала по валютам разъезжает позицию. На собесе это «сделка в чате казначейства».",
			Gold: `SELECT id, pair, amount_usd, status, cover_journal_id, uzs_journal_id
FROM fx_deals
WHERE status = 'BOOKED'
   OR ((cover_journal_id IS NULL OR cover_journal_id = '')
       AND (uzs_journal_id IS NULL OR uzs_journal_id = ''))`,
		},
		{
			ID: "sql-night-4eyes", Title: "Ночная смена без четырёх глаз", Grade: "middle",
			Story: "Контроль: после 22:00 перевод без второго подтверждения. Политика банка — four_eyes на ночь.",
			Question: "Журналы, где час ≥ 22 и four_eyes = 0.",
			Hint: "substr(at, 12, 2) >= '22'. four_eyes хранится как 0/1.",
			Why: "Ночные операции — отдельный контур риска. Аналитик описывает контроль, не «доверьтесь смене».",
			Gold: `SELECT id, kind, at, four_eyes, night_shift, note
FROM journals
WHERE substr(at, 12, 2) >= '22' AND four_eyes = 0`,
		},
		{
			ID: "sql-tickets-open", Title: "Открытые инциденты линии", Grade: "intern", InternSafe: true,
			Story: "Утро операционного зала: какие тикеты ещё OPEN, чтобы не потерять SLA.",
			Question: "Список открытых тикетов из intern_tickets.",
			Hint: "SELECT id, kind, status, title FROM intern_tickets WHERE status = 'OPEN'",
			Why: "Аналитик на линии начинает с очереди, не с ER-диаграммы.",
			Gold: "SELECT id, kind, status, title FROM intern_tickets WHERE status = 'OPEN'",
		},
	}
}

func TaskByID(id string) *Task {
	for _, t := range Tasks() {
		if t.ID == id {
			tt := t
			return &tt
		}
	}
	return nil
}

func PublicTasks(showGold map[string]bool) []map[string]any {
	out := []map[string]any{}
	for _, t := range Tasks() {
		row := map[string]any{
			"id": t.ID, "title": t.Title, "story": t.Story, "question": t.Question,
			"hint": t.Hint, "why": t.Why, "internSafe": t.InternSafe, "grade": t.Grade,
		}
		if showGold[t.ID] {
			row["gold"] = t.Gold
		}
		out = append(out, row)
	}
	return out
}

func MatchGold(got, gold []map[string]any) (bool, string) {
	if len(gold) == 0 && len(got) == 0 {
		return true, "Пусто как в эталоне — тоже ответ, если так в книге."
	}
	gkeys := keysOf(gold)
	skeys := keysOf(got)
	if len(gkeys) == 0 {
		if len(got) == len(gold) {
			return true, "Число строк совпало."
		}
		return false, "Эталон без ключа id: число строк не совпало."
	}
	missing := []string{}
	for k := range gkeys {
		if !skeys[k] {
			missing = append(missing, k)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		return false, "Не хватает строк эталона: " + strings.Join(missing, ", ")
	}
	if len(got) > len(gold)+3 {
		return false, "Слишком много лишних строк — фильтр, скорее всего, широкий."
	}
	return true, "Набор ключевых строк совпал с эталоном. На линии так и показывают выборку, не «я уверен»."
}

func keysOf(rows []map[string]any) map[string]bool {
	m := map[string]bool{}
	for _, r := range rows {
		if v, ok := r["id"]; ok && v != nil && str(v) != "" {
			m[str(v)] = true
			continue
		}
		if v, ok := r["from_wallet_id"]; ok {
			m[str(v) + "|" + str(r["to_wallet_id"]) + "|" + str(r["amount"])] = true
			continue
		}
		if v, ok := r["merchant_account_id"]; ok {
			m[str(v)] = true
			continue
		}
		if v, ok := r["debit"]; ok {
			m["trial:" + str(v) + "/" + str(r["credit"])] = true
			continue
		}
		if v, ok := r["book_available"]; ok {
			m[str(r["id"]) + "|" + str(v)] = true
			continue
		}
		b, _ := json.Marshal(r)
		m[string(b)] = true
	}
	return m
}
