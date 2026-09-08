package httpapi

import (
	"net/http"

	"github.com/malox4/malo-academy/internal/petdb"
)

func (s *Server) mountV1(mux *http.ServeMux) {
	wrap := func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !s.gateWrite(w, r) {
				return
			}
			s.Bank.Lock()
			defer s.Bank.Unlock()
			fn(w, r)
		}
	}

	mux.HandleFunc("GET /api/v1/catalog", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"catalog": petdb.ExplorerCatalog(), "contract": s.Bank.GetContract()})
	}))
	mux.HandleFunc("GET /api/openapi.json", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"openapi": "3.0.3", "info": map[string]any{"title": "Malo Core Banking / Wallet API", "version": "3.0.0"}})
	}))
	mux.HandleFunc("GET /api/v1/openapi.json", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"openapi": "3.0.3", "info": map[string]any{"title": "Malo Core Banking / Wallet API", "version": "3.0.0"}})
	}))
	mux.HandleFunc("GET /api/v1/state", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, s.Bank.Snapshot())
	}))
	mux.HandleFunc("POST /api/v1/reset", wrap(func(w http.ResponseWriter, r *http.Request) {
		out := s.Bank.Reset()
		s.Bank.Unlock()
		s.syncLab()
		s.Bank.Lock()
		writeJSON(w, 200, out)
	}))
	mux.HandleFunc("GET /api/v1/contract", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, s.Bank.GetContract())
	}))
	mux.HandleFunc("PUT /api/v1/contract", wrap(s.patchContract))
	mux.HandleFunc("PATCH /api/v1/contract", wrap(s.patchContract))
	mux.HandleFunc("GET /api/v1/accounts", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, s.Bank.AccountsTrial())
	}))
	mux.HandleFunc("GET /api/v1/accounts/{id}/statement", wrap(func(w http.ResponseWriter, r *http.Request) {
		row, err := s.Bank.Statement(r.PathValue("id"))
		s.bankJSON(w, r, 200, row, err)
	}))
	mux.HandleFunc("GET /api/v1/wallets", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"items": s.Bank.Wallets()})
	}))
	mux.HandleFunc("GET /api/v1/wallets/{id}", wrap(func(w http.ResponseWriter, r *http.Request) {
		row, err := s.Bank.GetWallet(r.PathValue("id"))
		s.bankJSON(w, r, 200, row, err)
	}))
	mux.HandleFunc("GET /api/v1/ledger", wrap(func(w http.ResponseWriter, r *http.Request) {
		acc := r.URL.Query().Get("accountId")
		writeJSON(w, 200, map[string]any{"items": s.Bank.LegsFilter(acc), "trial": s.Bank.TrialBalance()})
	}))
	mux.HandleFunc("GET /api/v1/ledger/journals", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"items": s.Bank.Journals(), "trial": s.Bank.TrialBalance()})
	}))
	mux.HandleFunc("POST /api/v1/ledger/journals/{id}/reverse", wrap(func(w http.ResponseWriter, r *http.Request) {
		st, body, err := s.Bank.ReverseJournal(r.PathValue("id"))
		s.bankJSON(w, r, st, body, err)
	}))
	mux.HandleFunc("GET /api/v1/ledger/t-accounts/{id}", wrap(func(w http.ResponseWriter, r *http.Request) {
		row, err := s.Bank.TAccount(r.PathValue("id"))
		s.bankJSON(w, r, 200, row, err)
	}))
	mux.HandleFunc("GET /api/v1/ledger/trial-balance", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, s.Bank.TrialBalance())
	}))
	mux.HandleFunc("GET /api/v1/holds", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"items": s.Bank.Holds()})
	}))
	mux.HandleFunc("GET /api/v1/payments", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"items": s.Bank.Payments()})
	}))
	mux.HandleFunc("GET /api/v1/payments/{id}", wrap(func(w http.ResponseWriter, r *http.Request) {
		row, err := s.Bank.GetPayment(r.PathValue("id"))
		s.bankJSON(w, r, 200, row, err)
	}))
	mux.HandleFunc("POST /api/v1/transfers", wrap(func(w http.ResponseWriter, r *http.Request) {
		body, _ := readJSON(r)
		st, out, err := s.Bank.CreateTransfer(header(r, "Idempotency-Key"), body)
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("POST /api/v1/payments/{id}/capture", wrap(func(w http.ResponseWriter, r *http.Request) {
		body, _ := readJSON(r)
		st, out, err := s.Bank.CapturePayment(r.PathValue("id"), body)
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("POST /api/v1/payments/{id}/cancel", wrap(func(w http.ResponseWriter, r *http.Request) {
		st, out, err := s.Bank.CancelPayment(r.PathValue("id"))
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("POST /api/v1/payments/{id}/refund", wrap(func(w http.ResponseWriter, r *http.Request) {
		body, _ := readJSON(r)
		st, out, err := s.Bank.RefundPayment(r.PathValue("id"), body)
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("POST /api/v1/payments/{id}/chargeback", wrap(func(w http.ResponseWriter, r *http.Request) {
		body, _ := readJSON(r)
		st, out, err := s.Bank.ChargebackPayment(r.PathValue("id"), body)
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("GET /api/v1/kyc", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"items": s.Bank.KYCList()})
	}))
	mux.HandleFunc("POST /api/v1/kyc/start", wrap(func(w http.ResponseWriter, r *http.Request) {
		body, _ := readJSON(r)
		st, out, err := s.Bank.StartKYC(header(r, "Idempotency-Key"), body)
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("POST /api/v1/kyc/{id}/decide", wrap(func(w http.ResponseWriter, r *http.Request) {
		body, _ := readJSON(r)
		st, out, err := s.Bank.DecideKYC(r.PathValue("id"), body)
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("GET /api/v1/cards", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"items": s.Bank.Plastic()})
	}))
	mux.HandleFunc("POST /api/v1/cards/issue", wrap(func(w http.ResponseWriter, r *http.Request) {
		body, _ := readJSON(r)
		st, out, err := s.Bank.IssueCard(body)
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("POST /api/v1/cards/{id}/block", wrap(func(w http.ResponseWriter, r *http.Request) {
		st, out, err := s.Bank.BlockCard(r.PathValue("id"))
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("GET /api/v1/auths", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"items": s.Bank.CardsAuth()})
	}))
	mux.HandleFunc("POST /api/v1/cards/authorize", wrap(func(w http.ResponseWriter, r *http.Request) {
		body, _ := readJSON(r)
		timeout := header(r, "X-ACS-Timeout") == "1" || body["acsTimeout"] == true
		st, out, err := s.Bank.AuthorizeCard(timeout, body)
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("GET /api/v1/bank/incoming", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"items": s.Bank.Incoming()})
	}))
	mux.HandleFunc("POST /api/v1/bank/incoming", wrap(func(w http.ResponseWriter, r *http.Request) {
		body, _ := readJSON(r)
		st, out, err := s.Bank.BankIncoming(body)
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("GET /api/v1/bank/outgoing", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"items": s.Bank.Outgoing()})
	}))
	mux.HandleFunc("POST /api/v1/bank/outgoing", wrap(func(w http.ResponseWriter, r *http.Request) {
		body, _ := readJSON(r)
		st, out, err := s.Bank.BankOutgoing(body)
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("POST /api/v1/bank/outgoing/{id}/ack", wrap(func(w http.ResponseWriter, r *http.Request) {
		body, _ := readJSON(r)
		st, out, err := s.Bank.BankAck(r.PathValue("id"), body)
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("GET /api/v1/bank/outgoing/{id}/mt103", wrap(func(w http.ResponseWriter, r *http.Request) {
		row, err := s.Bank.MT103(r.PathValue("id"))
		s.bankJSON(w, r, 200, row, err)
	}))
	mux.HandleFunc("GET /api/v1/bank/salary", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"items": s.Bank.Salary()})
	}))
	mux.HandleFunc("POST /api/v1/bank/salary", wrap(func(w http.ResponseWriter, r *http.Request) {
		body, _ := readJSON(r)
		st, out, err := s.Bank.SalaryIngest(body)
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("GET /api/v1/bank/suspense", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"items": s.Bank.Suspense()})
	}))
	mux.HandleFunc("POST /api/v1/bank/suspense/{id}/allocate", wrap(func(w http.ResponseWriter, r *http.Request) {
		body, _ := readJSON(r)
		st, out, err := s.Bank.AllocateSuspense(r.PathValue("id"), body)
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("GET /api/v1/fx", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, s.Bank.FXRates())
	}))
	mux.HandleFunc("POST /api/v1/fx/convert", wrap(func(w http.ResponseWriter, r *http.Request) {
		body, _ := readJSON(r)
		st, out, err := s.Bank.FXConvert(body)
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("GET /api/v1/fx/deals", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"items": s.Bank.FXDeals()})
	}))
	mux.HandleFunc("GET /api/v1/clearing", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"items": s.Bank.Clearing()})
	}))
	mux.HandleFunc("POST /api/v1/clearing/settle", wrap(func(w http.ResponseWriter, r *http.Request) {
		st, out, err := s.Bank.ClearingSettle()
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("GET /api/v1/recon/export", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, s.Bank.ExportRecon())
	}))
	mux.HandleFunc("POST /api/v1/recon/ingest", wrap(func(w http.ResponseWriter, r *http.Request) {
		body, _ := readJSON(r)
		st, out, err := s.Bank.IngestRecon(body)
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("GET /api/v1/shopline/stock", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"items": s.Bank.Stock()})
	}))
	mux.HandleFunc("GET /api/v1/customers", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"items": s.Bank.Customers()})
	}))
	mux.HandleFunc("GET /api/v1/merchants", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"items": s.Bank.Merchants()})
	}))
	mux.HandleFunc("GET /api/v1/tickets", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"items": s.Bank.Tickets()})
	}))
	mux.HandleFunc("GET /api/v1/tickets/{id}", wrap(func(w http.ResponseWriter, r *http.Request) {
		row, err := s.Bank.GetTicket(r.PathValue("id"))
		s.bankJSON(w, r, 200, row, err)
	}))
	mux.HandleFunc("POST /api/v1/tickets", wrap(func(w http.ResponseWriter, r *http.Request) {
		body, _ := readJSON(r)
		st, out, err := s.Bank.CreateTicket(body)
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("POST /api/v1/tickets/{id}/comment", wrap(func(w http.ResponseWriter, r *http.Request) {
		body, _ := readJSON(r)
		st, out, err := s.Bank.CommentTicket(r.PathValue("id"), body)
		s.bankJSON(w, r, st, out, err)
	}))
	mux.HandleFunc("GET /api/v1/holds/{id}", wrap(func(w http.ResponseWriter, r *http.Request) {
		row, err := s.Bank.GetHold(r.PathValue("id"))
		s.bankJSON(w, r, 200, row, err)
	}))
	mux.HandleFunc("GET /api/v1/audit", wrap(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"items": s.Bank.AuditLog()})
	}))
	mux.HandleFunc("POST /api/v1/shopline/orders", wrap(func(w http.ResponseWriter, r *http.Request) {
		body, _ := readJSON(r)
		st, out, err := s.Bank.CreateOrder(header(r, "Idempotency-Key"), body)
		s.bankJSON(w, r, st, out, err)
	}))
}

func (s *Server) patchContract(w http.ResponseWriter, r *http.Request) {
	body, _ := readJSON(r)
	c := s.Bank.PatchContract(body)
	writeJSON(w, 200, map[string]any{"ok": true, "contract": c})
}

func v1Catalog() []map[string]any {
	return petdb.ExplorerCatalog()
}

