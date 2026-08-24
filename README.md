# Malo Academy

Играбельная академия **Business Analyst / System Analyst**: путь Intern → Senior, ночные смены, квесты и **живой core-banking пет**, в который можно стрелять из Postman.

Не учебник. Контур меняется от ваших запросов: P2P пишет журнал, холд режет available, refund без реверса рвёт пробный баланс.

| | |
|---|---|
| Открыть | [http://localhost:8080](http://localhost:8080) |
| Пет / Postman | [http://localhost:8080/pet](http://localhost:8080/pet) · `/api/v1` |
| Справочник проводок | [http://localhost:8080/book](http://localhost:8080/book) |
| Спека | [http://localhost:8080/api/v1/openapi.json](http://localhost:8080/api/v1/openapi.json) |

Прогресс ученика — в браузере (LocalStorage + IndexedDB). Авторизации нет. Ledger пета — в `data/` (Docker volume `academy-data`).

## За 30 секунд

```bash
git clone https://github.com/malox4/malo-academy.git
cd malo-academy
docker compose up --build
```

```bash
curl -sS http://localhost:8080/api/health
# {"ok":true,"product":"malo-core-banking","version":3}

curl -sS -X POST http://localhost:8080/api/v1/transfers \
  -H 'Content-Type: application/json' \
  -d '{"fromWalletId":"wal_anna","toWalletId":"wal_boris","amount":5000}'
```

В консоли [/pet](http://localhost:8080/pet) дважды стрельните P2P без ключа — две DR-ноги Анны. Это заводской баг контракта, не «сломалось». Выкатите контракт из каталога, когда поймёте шов.

## Что внутри

**Путь** — 48 модулей Intern → Senior. Вкладка «Играть»: сортировки, споты, сцены, AC. Разбор — рядом, не вместо.

**BPMN** — живой редактор процесса банкомата: [http://localhost:8080/bpmn](http://localhost:8080/bpmn). Схема в `public/bpmn/atm-process.bpmn` открывается в Camunda Modeler и bpmn.io.

**API · пет** — двойная запись, T-счета, холд → capture → refund, IBAN / nostro / MT103, FX двумя журналами, зарплатный файл, клиринг T+1, сверка Orient. Контракт с завода дырявый: нет идемпотентности, ACS timeout = SUCCESS, refund без реверса.

**Квесты** — пять контуров: ночь на ledger, Malo Wallet, ShopLine, MedQueue, CityPark. Пишете ноги журнала, 409, SMS, NFR. Исход считается из метрик, не из «правильного теста».

**Зал** — смены как в 03:12: KYC, журнал, файл банка, go/no-go.

**Справочник** — дебет, кредит, холд ≠ ledger, nostro, trailer. С примером и кнопкой «вставить» в квесте.

**Собес / симуляция / печати** — голос на интервью, ситуационный экзамен 70%+, XP и бейджи.

## Docker

Образ — Node: статика `dist` + тот же `server/malo-api.mjs`, что в `npm run dev`. Порт **8080**.

```bash
docker compose up --build
```

Остановка: `docker compose down`. Ledger в volume; `down -v` сотрёт проводки.

## Без Docker

```bash
npm install
npm run dev
```

Тот же [http://localhost:8080](http://localhost:8080).

Прод-сборка:

```bash
npm run build
npm start
# или: npm run preview
```

## Пет: куда стрелять

База: `http://localhost:8080/api/v1`

| Контур | Примеры |
|---|---|
| Счета | `GET /accounts` · `GET /ledger/trial-balance` · `GET /ledger/t-accounts/acc_anna` |
| P2P | `POST /transfers` `{ fromWalletId, toWalletId, amount }` |
| Карты | `POST /cards/authorize` → `POST /payments/{id}/capture` · refund · chargeback |
| Банк | `POST /bank/incoming` · outgoing + `ack` · `GET .../mt103` |
| FX | `POST /fx/convert` `{ fromWalletId: wal_anna, toWalletId: wal_anna_usd, amountTo: 10 }` |
| Зарплата | `POST /bank/salary` — trailer = сумма строк |
| Клиринг | `POST /clearing/settle` после capture |
| Контракт | `PUT /contract` · `POST /reset` вернёт заводские баги |

Кошельки с завода: `wal_anna` (UZS), `wal_anna_usd`, `wal_boris`. IBAN Анны: `UZ12MALO000000000001`.

## Стек

Vite · React · TypeScript · Tailwind · Framer Motion · Zustand · Lucide  
Пет: Node HTTP (`server/malo-api.mjs` + `server/serve.mjs`)  
Шрифты: Inter + Fraunces

## Источники идей

IIBA BABOK v3, Karl Wiegers *Software Requirements*, IREB CPRE, SkillMap BA/SA, банк собеседований. Не дословные копии.
