/** Учебный банк: операции и практика P2P. Термин → объяснение → зачем → кейс с видимым API. */

export const PET_TABS = [
  { id: "console", label: "Операции", about: "Состояние банка и живые вызовы API" },
  { id: "sql", label: "SQL", about: "Запросы к таблицам ядра" },
  { id: "api", label: "API", about: "Справочник методов: путь, тело, ответ" },
  { id: "p2p", label: "P2P", about: "Внутренний перевод и идемпотентность" },
  { id: "hold", label: "Холд", about: "Авторизация карты: резерв без проводки" },
  { id: "iban", label: "IBAN", about: "Входящий платёж из другого банка" },
  { id: "hole", label: "Журнал", about: "Пробный баланс и дыра возврата" },
];

export const LESSON_TABS = ["p2p", "hold", "iban", "hole"];

export const BANK_TRACKS = [
  {
    id: "hold",
    tab: "hold",
    term: "Hold / Authorize",
    method: "POST",
    path: "/api/v1/cards/authorize",
    title: "Холд: резерв, ещё не проводка",
    about: "Касса резервирует 12 000 на карте Бориса. Available падает. Журнал молчит, пока нет capture.",
    purpose: "Отличить SUCCESS кассы от списания в книге. Тикет: «карта не проходит, в приложении минус».",
    steps: 7,
  },
  {
    id: "p2p",
    tab: "p2p",
    term: "P2P",
    method: "POST",
    path: "/api/v1/transfers",
    title: "P2P: перевод и идемпотентность",
    about: "Анна → Борис 5 000. С ключом повтор безопасен. Без ключа учебный банк проводит второй перевод.",
    purpose: "Описать AC: Idempotency-Key обязателен, повтор даёт 409, книга одна.",
    steps: 9,
  },
  {
    id: "iban",
    tab: "iban",
    term: "Incoming IBAN",
    method: "POST",
    path: "/api/v1/bank/incoming",
    title: "Входящий платёж по IBAN",
    about: "25 000 на IBAN Анны: дебет nostro, кредит клиента. Чужой IBAN уходит в невыясненные, не «похожему» клиенту.",
    purpose: "Не путать внутренний P2P и платёж из другого банка. Другой путь, другие счета.",
    steps: 7,
  },
  {
    id: "hole",
    tab: "hole",
    term: "Trial balance / refund",
    method: "POST",
    path: "/api/v1/payments/{id}/refund",
    title: "Возврат и дыра в журнале",
    about: "Возврат без обратной проводки по мерчанту: клиенту начислили, мерчанта не уменьшили. Пробный баланс расходится.",
    purpose: "Доказать дыру по trial.diff и journal kind, не правкой поля balance.",
    steps: 7,
  },
];

export const PET_INTRO = {
  title: "Учебный банк",
  lead: "Это не набор кнопок. Это живое ядро: кошельки, холды, журнал проводок. Вы вызываете те же HTTP-методы, что вызвал бы продукт — и сразу видите путь, заголовки, тело и ответ.",
  what: "Учебный банк хранит клиентов, витрину кошельков, открытые холды и журнал double-entry. Intern читает. Запись (POST/PATCH) открывается в PRO или в 3-дневном триале — отказ приходит кодом PLAN, кнопка не «серая».",
  why: "BA и SA на работе не «нажимают перевод». Они читают контракт: какой метод, какое тело, какой статус, что изменилось в книге. Здесь это можно сделать руками, а не по слайду.",
  how: [
    { n: "1", title: "Снять baseline", text: "GET кошельки и журнал — зафиксировать available и последние проводки." },
    { n: "2", title: "Вызвать операцию", text: "POST с явным JSON. Смотрите HTTP-статус, не только экран." },
    { n: "3", title: "Сверить книгу", text: "Снова GET: available, холд, journalId. Расхождение — предмет требования, не «обновите баланс»." },
  ],
};

export const P2P_LESSON = {
  name: "P2P",
  en: "peer-to-peer transfer",
  term: "P2P — внутренний перевод между кошельками одного банка. Не путать с IBAN: деньги не уходят в другой банк.",
  about: [
    "Клиентский счёт — пассив банка. Списание у отправителя — дебет (DR). Зачисление получателю — кредит (CR). Одна операция — один journalId и две проводки одной суммы.",
    "Идемпотентность: заголовок Idempotency-Key говорит ядру «этот запрос уже обрабатывали». Повтор с тем же ключом не должен создать вторую проводку. Без ключа учебный банк по умолчанию примет повтор как новый перевод — это учебный дефект контракта, его нужно увидеть и описать в AC.",
  ],
  example: [
    "Анна переводит Борису 5 000 UZS: POST /api/v1/transfers, тело fromWalletId / toWalletId / amount / currency.",
    "Тот же Idempotency-Key второй раз: тот же transfer id, без второй проводки (или 409 DUPLICATE — если так задан контракт).",
    "Два POST без ключа: две проводки, Анна теряет 10 000, пробный баланс при этом может сходиться.",
  ],
  purpose: "Научиться читать перевод как контракт API + книга, а не как кнопку. Это вопрос собеса и тикета «списали дважды».",
  watch: "После каждого вызова сравните available Анны и Бориса и последние журналы kind = P2P. Сходимость пробного баланса не доказывает, что клиента не списали дважды.",
};

function hdr(key) {
  return key ? { "Idempotency-Key": key } : {};
}

export function consoleCases(key) {
  return [
    {
      id: "wallets",
      title: "Прочитать витрину кошельков",
      term: "Витрина кошелька",
      about: "Available — то, что клиент может потратить: ledger минус сумма открытых холдов. Это экран приложения, не книга проводок.",
      purpose: "Снять baseline до любой операции. Без этого нельзя доказать, что перевод или холд что-то изменили.",
      expect: "Список кошельков: id, имя, валюта, ledger, hold, available.",
      method: "GET",
      path: "/api/v1/wallets",
    },
    {
      id: "trial",
      title: "Проверить пробный баланс",
      term: "Пробный баланс (trial balance)",
      about: "Контроль книги: сумма дебетов равна сумме кредитов по валюте. Сходимость ≠ отсутствие дубля перевода.",
      purpose: "Отличить дыру в журнале (дебет ≠ кредит) от корректной, но вредной для клиента двойной проводки.",
      expect: "Поля debit, credit, diff, balanced. В сиде может быть учебная дыра REFUND_HOLE.",
      method: "GET",
      path: "/api/v1/ledger/trial-balance",
    },
    {
      id: "journals",
      title: "Прочитать журнал проводок",
      term: "Журнал (ledger journal)",
      about: "Каждая успешная денежная операция оставляет журнал: kind, сумма, journalId. Холд (auth) журнала не создаёт.",
      purpose: "Связать HTTP-ответ с книгой. BA на разборе тикета ищет journalId, не скрин кассы.",
      expect: "Список журналов. Ищите kind P2P, CARD_CAPTURE, входящий IBAN.",
      method: "GET",
      path: "/api/v1/ledger/journals",
    },
    {
      id: "p2p-key",
      title: "P2P 5 000 с ключом идемпотентности",
      term: "POST /api/v1/transfers",
      about: "Внутренний перевод. Заголовок Idempotency-Key фиксирует попытку. Тело задаёт отправителя, получателя и сумму.",
      purpose: "Увидеть, куда уходит запрос: метод, путь, JSON, заголовок. После 201 — available Анны −5000, Бориса +5000, один journalId.",
      expect: "201 и объект перевода с journalId. Повтор этого кейса с тем же ключом не должен списать ещё 5 000.",
      method: "POST",
      path: "/api/v1/transfers",
      headers: hdr(key),
      body: { fromWalletId: "wal_anna", toWalletId: "wal_boris", amount: 5000, currency: "UZS" },
      write: true,
    },
    {
      id: "hold",
      title: "Авторизация карты (холд)",
      term: "Hold / authorize",
      about: "Касса резервирует сумму. Холд OPEN уменьшает available. Журнала ещё нет — проводка появится на capture.",
      purpose: "Отличить «касса SUCCESS» от «деньги списаны в книге». Классический тикет: клиент не может платить, capture не было.",
      expect: "201: hold OPEN, payment AUTHORIZED. GET /holds покажет резерв. Журнал не вырастет.",
      method: "POST",
      path: "/api/v1/cards/authorize",
      body: { walletId: "wal_boris", cardId: "crd_boris", amount: 12000, currency: "UZS" },
      write: true,
    },
    {
      id: "iban",
      title: "Входящий платёж по IBAN",
      term: "Incoming credit / IBAN",
      about: "Деньги пришли из другого банка на IBAN клиента. Это не P2P: в журнале дебет nostro и кредит клиента.",
      purpose: "Увидеть другой контур: другой путь, другое тело, другие счета. IBAN Анны в сиде: UZ12MALO000000000001.",
      expect: "201 и журнал зачисления. Available Анны +25 000, если IBAN найден.",
      method: "POST",
      path: "/api/v1/bank/incoming",
      body: { iban: "UZ12MALO000000000001", amount: 25000, currency: "UZS", paymentRef: "MT103-LAB" },
      write: true,
    },
    {
      id: "holds-read",
      title: "Список открытых холдов",
      term: "Hold status OPEN",
      about: "Резерв, который ещё не capture и не void. Он есть в таблице холдов и в поле hold кошелька.",
      purpose: "После авторизации проверить, что резерв живой. Пустой список — либо capture прошёл, либо auth не создал холд.",
      expect: "Массив холдов: id, amount, status, wallet/account.",
      method: "GET",
      path: "/api/v1/holds",
    },
  ];
}

export function p2pCases(key, key2) {
  return [
    {
      id: "p2p-base",
      n: 1,
      title: "Baseline: кошельки до перевода",
      term: "Available",
      about: "Зафиксируйте available Анны и Бориса. Все следующие кейсы сравниваются с этой точкой.",
      purpose: "Без baseline нельзя доказать эффект перевода.",
      expect: "GET 200, в items есть wal_anna и wal_boris.",
      method: "GET",
      path: "/api/v1/wallets",
    },
    {
      id: "p2p-send",
      n: 2,
      title: "Перевод 5 000 с Idempotency-Key",
      term: "P2P + Idempotency-Key",
      about: "Один POST, один ключ. Ядро должно создать один перевод и одну пару проводок: DR Анна, CR Борис.",
      purpose: "Связать JSON запроса с journalId в ответе.",
      expect: "201 Created. В теле: id, journalId, amount 5000. Available Анны −5000, Бориса +5000.",
      method: "POST",
      path: "/api/v1/transfers",
      headers: hdr(key),
      body: { fromWalletId: "wal_anna", toWalletId: "wal_boris", amount: 5000, currency: "UZS" },
      write: true,
    },
    {
      id: "p2p-replay",
      n: 3,
      title: "Повтор с тем же ключом",
      term: "Идемпотентный повтор",
      about: "Клиент нажал «ещё раз», шлюз повторил запрос. Тот же Idempotency-Key. Новой проводки быть не должно.",
      purpose: "Увидеть безопасный повтор. На работе это AC: «повтор с тем же ключом не создаёт второй перевод».",
      expect: "200 с тем же id (контракт по умолчанию) либо 409 DUPLICATE, если включён transferDuplicateHttp: 409. Available не должен уйти ещё на 5 000.",
      method: "POST",
      path: "/api/v1/transfers",
      headers: hdr(key),
      body: { fromWalletId: "wal_anna", toWalletId: "wal_boris", amount: 5000, currency: "UZS" },
      write: true,
    },
    {
      id: "p2p-naked",
      n: 4,
      title: "Перевод без ключа",
      term: "Запрос без Idempotency-Key",
      about: "По умолчанию ключ в учебном банке не обязателен. Первый POST без ключа создаст перевод как обычную операцию.",
      purpose: "Отделить «ключ не передали» от «ключ обязателен в контракте».",
      expect: "201 и новый journalId. Это ещё не баг — баг проявится на следующем шаге.",
      method: "POST",
      path: "/api/v1/transfers",
      body: { fromWalletId: "wal_anna", toWalletId: "wal_boris", amount: 5000, currency: "UZS" },
      write: true,
    },
    {
      id: "p2p-dup",
      n: 5,
      title: "Второй POST без ключа — дубль",
      term: "Учебный дефект контракта",
      about: "Повтор без ключа учебный банк принимает как новый перевод. Две проводки, клиент потерял 10 000. Пробный баланс при этом может остаться сходящимся.",
      purpose: "Сформулировать требование: ключ обязателен, повтор даёт 409, а не «кнопку disabled». Это пишут в AC, не в чат разработки.",
      expect: "201 ещё раз. Два разных id. Два журнала P2P. Available Анны −10 000 относительно шага 4.",
      method: "POST",
      path: "/api/v1/transfers",
      body: { fromWalletId: "wal_anna", toWalletId: "wal_boris", amount: 5000, currency: "UZS" },
      write: true,
    },
    {
      id: "p2p-books",
      n: 6,
      title: "Сверить журнал",
      term: "journalId",
      about: "Книга подтверждает, сколько переводов реально прошло. Считайте записи kind = P2P после ваших вызовов.",
      purpose: "Доказать дубль выборкой, не скрином приложения.",
      expect: "В ленте журналов несколько P2P. У дубля без ключа — разные journalId.",
      method: "GET",
      path: "/api/v1/ledger/journals",
    },
    {
      id: "p2p-contract",
      n: 7,
      title: "Включить обязательный ключ",
      term: "Контракт ядра",
      about: "PATCH /api/v1/contract меняет правила учебного банка. transferIdempotencyRequired: true — POST без ключа должен получить 400.",
      purpose: "Увидеть, что идемпотентность — требование к API, а не «поведение кнопки».",
      expect: "200 и обновлённый контракт. Следующий шаг проверит отказ.",
      method: "PATCH",
      path: "/api/v1/contract",
      body: { transferIdempotencyRequired: true, transferDuplicateHttp: 409 },
      write: true,
    },
    {
      id: "p2p-required",
      n: 8,
      title: "POST без ключа после контракта",
      term: "IDEMPOTENCY_KEY_REQUIRED",
      about: "Тот же JSON, что в шаге 4, но контракт уже требует заголовок. Ядро должно отказать, а не провести третий перевод.",
      purpose: "Проверить AC: нет ключа → 400, книга не меняется.",
      expect: "400 и код IDEMPOTENCY_KEY_REQUIRED. Available не меняется.",
      method: "POST",
      path: "/api/v1/transfers",
      body: { fromWalletId: "wal_anna", toWalletId: "wal_boris", amount: 5000, currency: "UZS" },
      write: true,
    },
    {
      id: "p2p-safe",
      n: 9,
      title: "Новый ключ, затем повтор → 409",
      term: "409 DUPLICATE",
      about: "С новым Idempotency-Key перевод создаётся один раз. Второй вызов с тем же ключом при transferDuplicateHttp: 409 возвращает конфликт, не новую проводку.",
      purpose: "Закрепить полный контракт: ключ обязателен, дубль — 409, книга одна.",
      expect: "Первый вызов 201, повтор с тем же ключом 409 и existingId. Нажмите «Выполнить» дважды.",
      method: "POST",
      path: "/api/v1/transfers",
      headers: hdr(key2),
      body: { fromWalletId: "wal_anna", toWalletId: "wal_boris", amount: 5000, currency: "UZS" },
      write: true,
    },
  ];
}

export const HOLD_LESSON = {
  name: "Холд",
  en: "authorization hold",
  term: "Холд — резерв суммы на счёте. Деньги ещё не списаны в журнал. Available уменьшается на OPEN hold. Проводка появится на capture.",
  about: [
    "Authorize (auth) — касса просит банк отложить сумму. Если контракт holdThenCapture включён, ядро пишет холд OPEN и платёж AUTHORIZED. Журнала нет.",
    "Capture — фактическое списание: холд закрывается, в журнале дебет клиента и кредит мерчанта. Void до capture снимает холд без проводки. Путать SUCCESS кассы с проводкой — типичный тикет.",
  ],
  example: [
    "POST /api/v1/cards/authorize { walletId: wal_boris, cardId: crd_boris, amount: 12000, currency: UZS }.",
    "GET /api/v1/holds — статус OPEN, amount 12000. GET /api/v1/wallets — available Бориса меньше, ledger тот же.",
    "GET /api/v1/ledger/journals — новой проводки нет, пока не будет capture.",
  ],
  purpose: "На собесе и в тикете отделить три числа: ledger, hold, available. Экран клиента показывает available, 1С ест журнал.",
  watch: "После auth сравните hold кошелька и список OPEN. Если journalId появился сразу — контракт пишет capture без холда.",
};

export const IBAN_LESSON = {
  name: "Входящий IBAN",
  en: "incoming credit transfer",
  term: "IBAN — международный номер счёта. Входящий платёж приходит из другого банка. Это не P2P: в журнале дебет nostro (наш счёт у корреспондента) и кредит клиента.",
  about: [
    "Известный IBAN клиента: POST /api/v1/bank/incoming создаёт журнал IBAN_IN, available клиента растёт.",
    "Неизвестный IBAN: деньги нельзя класть «соседу, номер похож». Контракт кладёт их на suspense (невыясненные). Разнос — отдельный POST, когда IBAN подтверждён.",
  ],
  example: [
    "IBAN Анны в сиде: UZ12MALO000000000001. POST amount 25000 → CR Анна, DR nostro.",
    "Чужой IBAN: GET /api/v1/bank/suspense покажет очередь. SLA старения — отдельное требование.",
  ],
  purpose: "Развести контуры: внутренний перевод, входящий, исходящий. BA пишет разные пути и разные счета в AC.",
  watch: "После зачисления смотрите journal kind IBAN_IN и available Анны. P2P сюда не при чём.",
};

export const HOLE_LESSON = {
  name: "Дыра журнала",
  en: "unbalanced refund",
  term: "Пробный баланс — контроль книги: сумма дебетов равна сумме кредитов по валюте. Дыра — проводка только с одной стороны. Сходимость ≠ «клиенту ничего не должны».",
  about: [
    "Правильный возврат после capture: дебет мерчанта и кредит клиента на ту же сумму (реверс). Контракт refundPostsReversal это включает.",
    "В учебном банке по умолчанию возврат может начислить клиенту без дебета мерчанта (kind REFUND_HOLE). Пробный баланс расходится. Чинить поле balance нельзя: 1С ест журнал.",
  ],
  example: [
    "GET /api/v1/ledger/trial-balance — в сиде уже может быть diff из REFUND_HOLE.",
    "POST /api/v1/payments/pay_cafe_cap_1/refund при refundPostsReversal: false увеличит дыру.",
    "Требование: refundPostsReversal true, две проводки, trial сходится по этой операции.",
  ],
  purpose: "Увидеть, что «вернули клиенту» без реверса мерчанта ломает книгу. Это AC на возврат, не Excel.",
  watch: "Смотрите trial.diff до и после refund. Ищите kind REFUND_HOLE в журнале и в ledger.",
};

export function holdCases() {
  return [
    {
      id: "hold-wallets",
      n: 1,
      title: "Baseline: кошелёк Бориса",
      term: "Available vs ledger",
      about: "Зафиксируйте available и hold Бориса. После auth available должен упасть, ledger — нет.",
      purpose: "Без baseline нельзя доказать, что холд что-то изменил.",
      expect: "GET 200, wal_boris в списке, поле hold.",
      method: "GET",
      path: "/api/v1/wallets",
    },
    {
      id: "hold-list",
      n: 2,
      title: "Список холдов до операции",
      term: "GET /api/v1/holds",
      about: "Открытый холд уменьшает available. В сиде уже может быть учебный холд Дины.",
      purpose: "Отличить старый OPEN от того, который создадите на следующем шаге.",
      expect: "Массив холдов: id, amount, status.",
      method: "GET",
      path: "/api/v1/holds",
    },
    {
      id: "hold-contract",
      n: 3,
      title: "Включить holdThenCapture",
      term: "Контракт holdThenCapture",
      about: "Если флаг выключен, authorize сразу пишет capture в журнал. Для этого разбора нужен именно холд.",
      purpose: "Правило «сначала резерв, потом списание» — требование к API, не настройка кассы.",
      expect: "200 и contract.holdThenCapture: true.",
      method: "PATCH",
      path: "/api/v1/contract",
      body: { holdThenCapture: true },
      write: true,
    },
    {
      id: "hold-auth",
      n: 4,
      title: "Авторизация 12 000 · карта Бориса",
      term: "POST /api/v1/cards/authorize",
      about: "Терминал запрашивает резерв. Тело: кошелёк, карта, сумма, валюта. Ответ содержит payment AUTHORIZED и hold OPEN.",
      purpose: "Увидеть полный путь и JSON. Это не «кнопка холд».",
      expect: "201: hold.status OPEN, payment.status AUTHORIZED, journalId нет.",
      method: "POST",
      path: "/api/v1/cards/authorize",
      body: { walletId: "wal_boris", cardId: "crd_boris", amount: 12000, currency: "UZS" },
      write: true,
    },
    {
      id: "hold-after",
      n: 5,
      title: "Проверить OPEN hold",
      term: "Hold OPEN",
      about: "Резерв должен появиться в списке. Available Бориса −12 000 относительно шага 1.",
      purpose: "Связать HTTP 201 с таблицей холдов.",
      expect: "Новый холд OPEN на 12000. Ledger Бориса не изменился.",
      method: "GET",
      path: "/api/v1/holds",
    },
    {
      id: "hold-journals",
      n: 6,
      title: "Журнал после auth",
      term: "Журнал молчит",
      about: "Пока нет capture, новой проводки CARD_CAPTURE быть не должно. Касса может показывать SUCCESS — книга ещё пуста.",
      purpose: "Доказать тикет «касса зелёная, денег нет» выборкой журнала.",
      expect: "Последние журналы без новой CARD_CAPTURE на эти 12 000.",
      method: "GET",
      path: "/api/v1/ledger/journals",
    },
    {
      id: "hold-wallets-2",
      n: 7,
      title: "Снова витрина",
      term: "Три числа",
      about: "ledger — книга, hold — резерв, available = ledger − OPEN holds. Клиенту показывают available.",
      purpose: "Закрепить формулу. На собесе называют три числа, не «баланс».",
      expect: "Борис: hold больше, available меньше, ledger как в шаге 1.",
      method: "GET",
      path: "/api/v1/wallets",
    },
  ];
}

export function ibanCases() {
  return [
    {
      id: "iban-wallets",
      n: 1,
      title: "Baseline: кошелёк Анны",
      term: "Available",
      about: "Зафиксируйте available Анны в UZS до входящего.",
      purpose: "После зачисления сравнить +25 000.",
      expect: "GET 200, wal_anna.",
      method: "GET",
      path: "/api/v1/wallets",
    },
    {
      id: "iban-list",
      n: 2,
      title: "Реестр входящих",
      term: "GET /api/v1/bank/incoming",
      about: "Очередь зачислений, которые уже прошли через ядро.",
      purpose: "Увидеть контур до своего POST.",
      expect: "Список incoming. В сиде может быть невыясненная сумма.",
      method: "GET",
      path: "/api/v1/bank/incoming",
    },
    {
      id: "iban-in",
      n: 3,
      title: "Зачислить 25 000 на IBAN Анны",
      term: "POST /api/v1/bank/incoming",
      about: "Тело: iban, amount, currency, paymentRef. IBAN Анны: UZ12MALO000000000001. Журнал: DR nostro, CR клиент.",
      purpose: "Связать JSON с книгой. Это не POST /transfers.",
      expect: "201, journal IBAN_IN. Available Анны +25 000.",
      method: "POST",
      path: "/api/v1/bank/incoming",
      body: { iban: "UZ12MALO000000000001", amount: 25000, currency: "UZS", paymentRef: "MT103-LAB" },
      write: true,
    },
    {
      id: "iban-after",
      n: 4,
      title: "Витрина после зачисления",
      term: "Кредит клиента",
      about: "Пассив клиента вырос. Nostro (актив банка) тоже изменился, но витрина кошелька показывает клиента.",
      purpose: "Подтвердить эффект на available, не на «ощущении».",
      expect: "wal_anna available больше на 25 000.",
      method: "GET",
      path: "/api/v1/wallets",
    },
    {
      id: "iban-journals",
      n: 5,
      title: "Журнал IBAN_IN",
      term: "journalId",
      about: "Найдите свежий kind IBAN_IN. Одна операция — две проводки одной суммы.",
      purpose: "Доказать зачисление книгой.",
      expect: "Журнал IBAN_IN, debit = credit.",
      method: "GET",
      path: "/api/v1/ledger/journals",
    },
    {
      id: "iban-unknown",
      n: 6,
      title: "Входящий на чужой IBAN",
      term: "Suspense",
      about: "Номер не из справочника клиентов. Нельзя зачислить «похожему». Ядро кладёт на невыясненные, если так задан контракт.",
      purpose: "Увидеть другой исход того же метода: тот же POST, другой IBAN.",
      expect: "201 в suspense или 404 IBAN — смотрите тело ответа и контракт unknownIbanToSuspense.",
      method: "POST",
      path: "/api/v1/bank/incoming",
      body: { iban: "UZ99XXXX000000000099", amount: 10000, currency: "UZS", paymentRef: "UNKNOWN-IBAN" },
      write: true,
    },
    {
      id: "iban-suspense",
      n: 7,
      title: "Очередь невыясненных",
      term: "GET /api/v1/bank/suspense",
      about: "Чужие деньги ждут идентификации. Разнос — отдельный метод, не UPDATE кошелька.",
      purpose: "Закрепить: неизвестный IBAN ≠ соседний клиент.",
      expect: "Записи со статусом SUSPENSE, в том числе учебная in_old.",
      method: "GET",
      path: "/api/v1/bank/suspense",
    },
  ];
}

export function holeCases() {
  return [
    {
      id: "hole-trial",
      n: 1,
      title: "Пробный баланс до операции",
      term: "GET /api/v1/ledger/trial-balance",
      about: "Сумма дебетов и кредитов по UZS. В сиде уже может быть учебная дыра REFUND_HOLE — diff ≠ 0.",
      purpose: "Зафиксировать diff. Дальше вы либо увеличите дыру, либо увидите, что она уже в книге.",
      expect: "Поля debit, credit, diff, balanced.",
      method: "GET",
      path: "/api/v1/ledger/trial-balance",
    },
    {
      id: "hole-journals",
      n: 2,
      title: "Найти REFUND_HOLE в журнале",
      term: "kind REFUND_HOLE",
      about: "Односторонняя запись: клиенту начислили, мерчанта не дебетовали. Ищите kind в ленте журналов.",
      purpose: "Доказать дыру выборкой, не скрином Excel.",
      expect: "В сиде есть журнал или проводка с kind REFUND_HOLE.",
      method: "GET",
      path: "/api/v1/ledger/journals",
    },
    {
      id: "hole-pay",
      n: 3,
      title: "Карточка уже захваченного платежа",
      term: "GET /api/v1/payments/{id}",
      about: "pay_cafe_cap_1 в сиде — capture Cafe Orient. Refund разрешён только из CAPTURED.",
      purpose: "Увидеть статус до возврата. Если уже REFUNDED — откройте API → «Сбросить учебный банк» и повторите шаги.",
      expect: "status CAPTURED, есть amount / capturedAmount.",
      method: "GET",
      path: "/api/v1/payments/pay_cafe_cap_1",
    },
    {
      id: "hole-contract",
      n: 4,
      title: "Выключить реверс на возврате",
      term: "refundPostsReversal: false",
      about: "Учебный дефект: возврат начисляет клиенту без дебета мерчанта. Так показывают дыру, не «поломку случайно».",
      purpose: "Понять, что поведение возврата — флаг контракта, его пишут в AC.",
      expect: "200, refundPostsReversal false.",
      method: "PATCH",
      path: "/api/v1/contract",
      body: { refundPostsReversal: false },
      write: true,
    },
    {
      id: "hole-refund",
      n: 5,
      title: "Refund без реверса мерчанта",
      term: "POST /api/v1/payments/pay_cafe_cap_1/refund",
      about: "Тот же платёж. Тело может задать amount. При выключенном реверсе книга расходится сильнее.",
      purpose: "Увидеть живой путь возврата и эффект на trial.",
      expect: "200, payment REFUNDED. trial.diff изменился. Не 409 STATE, если платёж ещё CAPTURED.",
      method: "POST",
      path: "/api/v1/payments/pay_cafe_cap_1/refund",
      body: { amount: 1000 },
      write: true,
    },
    {
      id: "hole-trial-2",
      n: 6,
      title: "Пробный баланс после дыры",
      term: "diff ≠ 0",
      about: "Сравните с шагом 1. Правильный вывод: возврат без двух проводок ломает книгу. Лечение — сторно/реверс, не UPDATE wallets.",
      purpose: "Сформулировать требование: refundPostsReversal true.",
      expect: "balanced false либо diff больше, чем в шаге 1.",
      method: "GET",
      path: "/api/v1/ledger/trial-balance",
    },
    {
      id: "hole-fix",
      n: 7,
      title: "Включить правильный реверс",
      term: "refundPostsReversal: true",
      about: "Следующий возврат должен писать две проводки. Уже созданную дыру этот флаг не затирает — её закрывают сторно.",
      purpose: "Отделить «включить контракт» от «починить прошлую проводку».",
      expect: "200, флаг true. Старый REFUND_HOLE в книге остаётся — это честно.",
      method: "PATCH",
      path: "/api/v1/contract",
      body: { refundPostsReversal: true },
      write: true,
    },
  ];
}

export function lessonFor(tab) {
  if (tab === "hold") return HOLD_LESSON;
  if (tab === "iban") return IBAN_LESSON;
  if (tab === "hole") return HOLE_LESSON;
  return P2P_LESSON;
}

export function lessonCases(tab, key, key2) {
  if (tab === "hold") return holdCases();
  if (tab === "iban") return ibanCases();
  if (tab === "hole") return holeCases();
  if (tab === "p2p") return p2pCases(key, key2);
  return consoleCases(key);
}

export function formatRequest(c) {
  const lines = [`${c.method} ${c.path}`];
  if (c.headers) {
    for (const [k, v] of Object.entries(c.headers)) lines.push(`${k}: ${v}`);
  }
  if (c.body) lines.push("", JSON.stringify(c.body, null, 2));
  else if (c.method === "GET") lines.push("", "(тела нет)");
  return lines.join("\n");
}
