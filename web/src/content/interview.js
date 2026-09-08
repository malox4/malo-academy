import { INTERVIEW as RAW } from "./interviewBank";

const SENIOR = new Set([
  "iv-fx",
  "iv-db",
  "iv-er",
  "iv-sql-nosql",
  "iv-joins",
  "iv-constraints",
  "iv-arch",
  "iv-integration",
  "iv-contract",
  "iv-brokers",
  "iv-esb",
  "iv-nfr",
  "iv-trace",
  "iv-val-ver",
  "iv-baccm",
  "iv-hard-req",
  "iv-saga",
  "iv-iso20022",
  "iv-slo",
  "iv-eventstorm",
  "iv-cr-impact",
]);

function flow(title, caption, nodes, notes) {
  const edges = [];
  for (let i = 0; i < nodes.length - 1; i++) {
    edges.push({ from: nodes[i].id, to: nodes[i + 1].id, note: notes[i] || "" });
  }
  return { kind: "flow", title, caption, nodes, edges };
}

const BOARD = {
  ledger: flow(
    "Книга, не поле balance",
    "Кошелёк показывает available. Журнал пишет ноги. Холд — не проводка.",
    [
      { id: "wallet", label: "Кошелёк", sub: "available", detail: "Клиент видит available = ledger − hold. Это витрина, не книга." },
      { id: "hold", label: "Холд", sub: "OPEN", mood: "hold", detail: "Auth откладывает сумму. В журнале ещё пусто." },
      { id: "journal", label: "Журнал", sub: "DR / CR", detail: "Каждое движение — две ноги, один journalId, одна валюта." },
      { id: "trial", label: "Trial", sub: "DR = CR", detail: "Пробный баланс по валюте. Дыра — это не «потом поправим»." },
    ],
    ["витрина", "capture", "сверка"],
  ),
  hold: flow(
    "Auth ещё не capture",
    "Касса может кричать SUCCESS, пока журнал молчит.",
    [
      { id: "till", label: "Касса", sub: "auth", detail: "Терминал показал успех. Это резерв, не списание." },
      { id: "hold", label: "Холд", sub: "OPEN", mood: "hold", detail: "Деньги отложены. Available уменьшается. Ног нет." },
      { id: "journal", label: "Журнал", sub: "пусто", detail: "Проводка появится на capture. Void снимает холд без ног." },
      { id: "settle", label: "Клиринг", sub: "T+1", detail: "Выплата мерчанту — другой журнал, не касса." },
    ],
    ["резерв", "capture", "settle"],
  ),
  iban: flow(
    "Входящий IBAN",
    "Nostro — наш счёт в чужом банке. Неизвестный IBAN не кладут соседу.",
    [
      { id: "bank", label: "Банк", sub: "MT103", detail: "Платёж пришёл по IBAN. Ключ и статус важнее «просто зачислить»." },
      { id: "nostro", label: "Nostro", sub: "актив", detail: "DR nostro: деньги сели к нам в корреспонденте." },
      { id: "client", label: "Клиент", sub: "пассив", detail: "CR клиента, если IBAN наш. Иначе suspense." },
      { id: "suspense", label: "Suspense", sub: "неясные", detail: "Неизвестный IBAN. Не соседний кошелёк." },
    ],
    ["ack", "зачисление", "если чужой"],
  ),
  ticket: {
    kind: "ticket",
    title: "Тикет с собеса",
    caption: "Строка, которую сейчас разбирают. Ход подсветит брак или шов.",
    lines: [
      { id: "1", text: "Need назван и проверяем", why: "Без Need это хотелка." },
      { id: "2", text: "Есть ключ, ошибка и кто мастер", why: "Контракт, не лозунг." },
      { id: "3", text: "«Сделайте как у банка»", why: "Чужой UI. Ловушка.", bad: true },
      { id: "4", text: "После релиза смотрим метрику, не только UAT", why: "Валидация ≠ верификация." },
    ],
  },
  taccount: {
    kind: "taccount",
    title: "T-счёт P2P Анна → Борис",
    caption: "Пассив клиента: отдал — дебет, получил — кредит. Один journalId.",
    debit: [{ id: "dr", text: "DR Анна 5 000 UZS", ok: true }],
    credit: [{ id: "cr", text: "CR Борис 5 000 UZS" }],
    hole: "Вторая такая же проводка без ключа — учебный баг пета: книга разъедется по смыслу продукта, trial может сойтись.",
  },
};

const TOPIC_BOARD = {
  Проводки: "ledger",
  Платежи: "hold",
  Банк: "iban",
};

function spoken(trap) {
  return trap
    .replace(/^Сказать\s+/i, "")
    .replace(/^«/, "")
    .replace(/»\.?$/, ".")
    .replace(/^Не уметь /, "Не знаю, ")
    .replace(/^Не /, "Не буду: ");
}

function repliesFor(item, board) {
  const nodes = board.nodes || board.debit || board.lines || [];
  const okHit = nodes[2]?.id || nodes[1]?.id || nodes[0]?.id || "journal";
  const badHit = nodes[0]?.id || okHit;
  const say = item.short.split(/(?<=\.)\s+/).slice(0, 2).join(" ");
  return [
    {
      id: "strong",
      text: say,
      good: true,
      why: item.strong,
      hit: okHit,
      lights: { [okHit]: "ok" },
    },
    {
      id: "trap0",
      text: spoken(item.traps[0] || "Сделаем как обычно."),
      good: false,
      why: "Ловушка. " + (item.traps[0] || ""),
      hit: badHit,
      broken: [badHit],
      tickets: board.kind === "ticket" ? { 3: "bad" } : undefined,
    },
    {
      id: "trap1",
      text: spoken(item.traps[1] || item.traps[0] || "Это не моя зона."),
      good: false,
      why: "Так звучит клише. Сильный ход: " + item.short.split(".")[0] + ".",
      hit: badHit,
      broken: [badHit],
    },
  ];
}

const CUSTOM = {
  "iv-ledger": (board) => [
    {
      id: "s",
      text: "Кошелёк — витрина: available = ledger − hold. Движение денег — журнал с DR и CR, один journalId.",
      good: true,
      why: "Книга и витрина разведены. Холд не нога.",
      hit: "journal",
      lights: { journal: "ok", wallet: "ok" },
    },
    {
      id: "t0",
      text: "Просто минусуем поле balance у отправителя.",
      good: false,
      why: "Это Excel. Без журнала 1С и сверка не сходятся.",
      hit: "wallet",
      broken: ["wallet", "trial"],
    },
    {
      id: "t1",
      text: "Hold, available и ledger — одно число, клиенту так проще.",
      good: false,
      why: "Три числа. Смешать их — не увидеть висящий холд.",
      hit: "hold",
      broken: ["hold"],
    },
  ],
  "iv-hold": () => [
    {
      id: "s",
      text: "Authorize открывает холд. Capture пишет ноги. Void снимает холд. Refund — реверс уже проведённого.",
      good: true,
      why: "Касса ≠ журнал. SUCCESS на auth — ложное списание, если ACS отвалился.",
      hit: "hold",
      lights: { hold: "hold", till: "ok" },
    },
    {
      id: "t0",
      text: "Для клиента auth и списание — одно и то же.",
      good: false,
      why: "Клиент видит кассу. Книга ещё пуста. Путать их — дыра.",
      hit: "till",
      broken: ["journal"],
    },
    {
      id: "t1",
      text: "Refund из холда, чтобы быстрее вернуть.",
      good: false,
      why: "Из AUTHORIZED делают void. Refund — после capture.",
      hit: "journal",
      broken: ["journal", "settle"],
    },
  ],
  "iv-iban": () => [
    {
      id: "s",
      text: "Входящий: DR nostro, CR клиента. Неизвестный IBAN — suspense. Исходящий до ack сидит в холде.",
      good: true,
      why: "Nostro — актив. Соседний клиент не получает чужие деньги.",
      hit: "nostro",
      lights: { nostro: "ok", client: "ok" },
    },
    {
      id: "t0",
      text: "IBAN — это просто строка в форме перевода.",
      good: false,
      why: "IBAN — маршрут и ключ зачисления, не поле UI.",
      hit: "bank",
      broken: ["bank"],
    },
    {
      id: "t1",
      text: "Неизвестный IBAN положим техническому клиенту.",
      good: false,
      why: "Это чужие деньги на чужом счёте. Suspense.",
      hit: "suspense",
      broken: ["client"],
      tickets: { 3: "bad" },
    },
  ],
  "iv-three-balances": () => [
    {
      id: "s",
      text: "ledger — ноги журнала, hold — OPEN, available = ledger − hold. UPDATE одного поля убьёт trial.",
      good: true,
      why: "Три числа. Ночной UPDATE — штопка экрана, не книга.",
      hit: "journal",
      lights: { journal: "ok", hold: "ok", wallet: "ok" },
    },
    {
      id: "t0",
      text: "Одно поле balance, так клиенту проще.",
      good: false,
      why: "Висящий холд станет невидимым.",
      hit: "wallet",
      broken: ["hold"],
    },
    {
      id: "t1",
      text: "Ночью UPDATE wallets.balance, смену закроем.",
      good: false,
      why: "1С ест журнал. Trial разъедется утром.",
      hit: "trial",
      broken: ["trial"],
    },
  ],
};

function boardFor(item) {
  if (item.id === "iv-ledger" || item.id === "iv-three-balances") return BOARD.ledger;
  if (item.id === "iv-hold" || item.id === "iv-3ds") return BOARD.hold;
  if (item.id === "iv-iban" || item.id === "iv-fx" || item.id === "iv-iso20022") return BOARD.iban;
  if (TOPIC_BOARD[item.topic] === "ledger") return BOARD.ledger;
  if (TOPIC_BOARD[item.topic] === "hold") return BOARD.hold;
  if (TOPIC_BOARD[item.topic] === "iban") return BOARD.iban;
  if (item.topic === "SRS" || item.topic === "BRD" || item.topic === "Документация") return BOARD.ticket;
  return BOARD.ticket;
}

function enrich(item) {
  const board = boardFor(item);
  const custom = CUSTOM[item.id];
  return {
    ...item,
    grade: SENIOR.has(item.id) ? "senior" : "intern",
    board,
    replies: custom ? custom(board) : repliesFor(item, board),
    follow: item.traps[2] || item.traps[0],
  };
}

export const INTERVIEW = RAW.map(enrich);

export const ORDERED_IDS = [
  "iv-role", "iv-clarify", "iv-elicitation", "iv-quality", "iv-req-types",
  "iv-sdlc", "iv-waterfall", "iv-agile", "iv-srs", "iv-brd", "iv-docs", "iv-cr-impact",
  "iv-db", "iv-er", "iv-sql-ops", "iv-sql-window", "iv-joins", "iv-constraints", "iv-sql-nosql",
  "iv-stakeholders", "iv-raci", "iv-raci-hold", "iv-prio", "iv-usecase", "iv-user-story", "iv-invest", "iv-usm", "iv-cjm",
  "iv-ledger", "iv-three-balances", "iv-hold", "iv-3ds", "iv-iban", "iv-iso20022", "iv-fx", "iv-api", "iv-openapi", "iv-soap", "iv-rest", "iv-soap-rest", "iv-contract",
  "iv-conflict", "iv-hard-req", "iv-nfr", "iv-slo", "iv-trace", "iv-val-ver", "iv-bpmn", "iv-eventstorm", "iv-uml",
  "iv-close", "iv-arch", "iv-saga", "iv-integration", "iv-brokers", "iv-esb", "iv-baccm",
];

const RANK = new Map(ORDERED_IDS.map((id, i) => [id, i]));

export function inArcOrder(items) {
  return [...items].sort((a, b) => (RANK.get(a.id) ?? 900) - (RANK.get(b.id) ?? 900));
}

export const TOPICS = ["Порядок", "Все", ...Array.from(new Set(INTERVIEW.map((i) => i.topic)))];

export function scoreAnswer(item, text) {
  const t = (text || "").toLowerCase();
  if (t.trim().length < 12) return { good: false, why: "Слишком коротко для стола. Назовите ноги, ключ или Need — не лозунг." };
  const trapHit = (item.traps || []).some((tr) => {
    const words = tr.toLowerCase().match(/[а-яa-z]{5,}/g) || [];
    return words.filter((w) => t.includes(w)).length >= 2;
  });
  const hints = (item.short + " " + item.strong)
    .toLowerCase()
    .match(/[а-яa-z]{5,}/g) || [];
  const uniq = [...new Set(hints)].slice(0, 14);
  const hits = uniq.filter((w) => t.includes(w)).length;
  if (trapHit && hits < 3) {
    return { good: false, why: "Это близко к ловушке. " + item.traps[0] };
  }
  if (hits >= 3) {
    return { good: true, why: item.strong };
  }
  return {
    good: false,
    why: "Фактов мало. Сильный ответ опирается на журнал, ключ или проверяемый Need — не на «как в BABOK».",
  };
}
