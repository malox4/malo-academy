/** Каталог /live — карта блоков сессии. Банк вопросов живёт на /interview. */

export const LIVE_TRACKS = [
  {
    id: "warmup",
    title: "Разминка",
    kicker: "роль и уточнения",
    blurb: "Кем вы за столом и какие три вопроса зададите, пока PO орёт «как у банка».",
    cards: [
      { id: "iv-role", title: "Роль BA / SA", teaser: "Need и границы. Не секретарь чата и не копия чужого UI.", minutes: 6, free: true },
      { id: "iv-clarify", title: "Уточнения", teaser: "Лозунг «надёжные платежи» разложить на ключ, SLA и хозяина.", minutes: 6, free: true },
    ],
  },
  {
    id: "process",
    title: "Процесс",
    kicker: "SDLC, не Scrum-плакат",
    blurb: "Где аналитик живёт в цикле платежа. CR после Done — не комментарий в тикете.",
    cards: [
      { id: "iv-sdlc", title: "Жизненный цикл", teaser: "Discovery → приёмка → оценка после релиза. Ночной инцидент кормит следующий анализ.", minutes: 8, free: true },
      { id: "iv-cr-impact", title: "CR и частичный capture", teaser: "8 000 после полной суммы. Together: тест, YAML, экран.", minutes: 7, free: false },
    ],
  },
  {
    id: "data",
    title: "Данные и SQL",
    kicker: "книга, не Excel",
    blurb: "Поймать двойное списание запросом к журналу. В зале — GROUP BY. Окно ROW_NUMBER — формулировка для работы. Три числа, не одно поле.",
    cards: [
      { id: "iv-sql-ops", title: "Дубли в журнале", teaser: "GROUP BY ключу, HAVING COUNT(*) > 1. UPDATE в проде — не ваш инструмент.", minutes: 8, free: true },
      { id: "iv-sql-window", title: "Окно SQL", teaser: "В зале — GROUP BY по ключу. На работе окно оставляет сами строки. Карточка LIMIT 1 врёт.", minutes: 7, free: true },
      { id: "iv-three-balances", title: "Три числа", teaser: "ledger, hold, available. UPDATE balance убьёт trial.", minutes: 7, free: true },
    ],
  },
  {
    id: "product",
    title: "Стейкхолдеры",
    kicker: "день и ночь",
    blurb: "Коллекторы днём, сверка ночью. Один A на захват scope.",
    cards: [
      { id: "iv-stakeholders", title: "Карта людей", teaser: "Кто орёт в чате ≠ кто Accountable. Комплаенс и ночь на карте.", minutes: 7, free: true },
      { id: "iv-raci-hold", title: "RACI на холде", teaser: "PO — один A. Dev не владеет политикой денег. Два A — брак.", minutes: 6, free: true },
    ],
  },
  {
    id: "case",
    title: "Кейс на доске",
    kicker: "холд и журнал",
    blurb: "Касса SUCCESS, available меньше, журнал пуст. 3DS ≠ ноги. OpenAPI ≠ скрин.",
    cards: [
      { id: "iv-hold", title: "Холд на кассе", teaser: "Холд OPEN: available меньше, журнала ещё нет.", minutes: 8, free: true },
      { id: "iv-3ds", title: "Таймаут 3DS", teaser: "ACS тишина → UNKNOWN. Успех ACS ещё не capture.", minutes: 7, free: true },
      { id: "iv-ledger", title: "Двойной P2P", teaser: "Два журнала. Trial может сойтись — продукт нет. Ключ, не disabled кнопки.", minutes: 8, free: true },
      { id: "iv-openapi", title: "OpenAPI vs скрин", teaser: "YAML: ключ, 201, 409. Картинка Swagger не контракт.", minutes: 7, free: true },
      { id: "iv-iso20022", title: "Поле ISO 20022", teaser: "Чужой конверт. Suspense, не «похожий» IBAN.", minutes: 7, free: false },
    ],
  },
  {
    id: "softs",
    title: "Софты",
    kicker: "конфликт вслух",
    blurb: "Риск, коллекторы и продукт тянут в разные стороны. Оранжевый стикер — факт, не сервис.",
    cards: [
      { id: "iv-conflict", title: "Спор о списаниях", teaser: "Факты, три опции, один Accountable. Не «всем будет хорошо».", minutes: 7, free: true },
      { id: "iv-eventstorm", title: "Оранжевый стикер", teaser: "HoldOpened, не wallet-api. Граница касса / книга.", minutes: 7, free: false },
      { id: "iv-slo", title: "SLO vs лозунг", teaser: "SLI на capture. 409 — успех индикатора. Не CPU.", minutes: 7, free: false },
      { id: "iv-saga", title: "Компенсация саги", teaser: "Outbox с ногами. Не второй POST capture.", minutes: 8, free: false },
    ],
  },
  {
    id: "close",
    title: "Закрытие",
    kicker: "первая неделя",
    blurb: "Какие вопросы задать нам и что делать первую неделю на контуре — не переписывать ядро.",
    cards: [
      { id: "iv-close", title: "Вопросы нам", teaser: "Последний инцидент capture, кто закрывает suspense. As-is, не рефакторинг.", minutes: 6, free: true },
    ],
  },
];

export function liveTrackById(id) {
  return LIVE_TRACKS.find((t) => t.id === id) || LIVE_TRACKS[0];
}
