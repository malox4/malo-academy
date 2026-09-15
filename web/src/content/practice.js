import { LABS } from "./labs";
import { topicForLab } from "./topics";

export { LABS };

export const PRACTICE_CLUSTERS = [
  { id: "requirements", title: "Требования", trackIds: ["requirements"] },
  { id: "database", title: "Данные", trackIds: ["database"] },
  { id: "api", title: "API", trackIds: ["restapi", "integrations"] },
  { id: "architecture", title: "Архитектура", trackIds: ["architecture"] },
  { id: "processes", title: "Процессы", trackIds: ["processes"] },
  { id: "quality", title: "Качество", trackIds: ["quality"] },
  { id: "bank", title: "Банк", trackIds: ["bank"] },
  { id: "softs", title: "Софты", trackIds: ["softs"] },
  { id: "career", title: "Карьера", trackIds: ["career"] },
];

const TRACK_CLUSTER = {};
for (const c of PRACTICE_CLUSTERS) {
  for (const t of c.trackIds) TRACK_CLUSTER[t] = c.id;
}

export function drillKind(lab) {
  const k = lab?.blocks?.[0]?.kind;
  const id = lab?.id || "";
  const teaser = lab?.teaser || "";
  if (id.includes("openapi") || k === "yaml") return "контракт YAML";
  if (k === "sql" || id.includes("sql") || id.includes("db-window") || id.includes("db-cte")) return "SQL";
  if (id.includes("raci") || k === "match") return "соответствие";
  if (teaser.includes("/pet") || id.includes("postman") || id.includes("ledger-deep")) return "учебный банк";
  if (k === "sort") return "стол";
  if (k === "spot") return "тикет";
  if (k === "scene" || k === "case") return "сценарий";
  if (k === "order") return "порядок шагов";
  return "стол";
}

export function clusterIdForLab(lab) {
  const topic = topicForLab(lab?.id);
  return TRACK_CLUSTER[topic?.trackId] || "";
}

export function labsByCluster() {
  const groups = PRACTICE_CLUSTERS.map((c) => ({ id: c.id, title: c.title, labs: [] }));
  const rest = [];
  for (const lab of LABS) {
    const cid = clusterIdForLab(lab);
    const g = groups.find((x) => x.id === cid);
    if (g) g.labs.push(lab);
    else rest.push(lab);
  }
  if (rest.length) groups.push({ id: "shift", title: "Смена", labs: rest });
  return groups.filter((g) => g.labs.length);
}

export function deskFromBlock(block) {
  if (!block) return null;
  if (block.kind === "sort") {
    return {
      kind: "desk",
      title: block.title,
      prompt: block.prompt,
      columns: block.buckets,
      cards: (block.items || []).map((i) => ({
        id: i.id,
        text: i.text,
        column: i.bucket,
        why: i.why,
        miss: i.why,
        hit: i.bucket,
      })),
    };
  }
  if (block.kind === "spot") {
    return { kind: "ticket", title: block.title, prompt: block.prompt, lines: block.lines };
  }
  if (block.kind === "scene") {
    return { kind: "chat", title: block.title, setting: block.setting, steps: block.steps };
  }
  if (block.kind === "case") {
    return {
      kind: "chat",
      title: block.title,
      setting: block.situation,
      prompt: block.debrief || block.question,
      steps: [{ from: "Кейс", line: block.question, options: block.options }],
    };
  }
  if (block.kind === "match") {
    return {
      kind: "map",
      title: block.title,
      prompt: block.prompt,
      slots: (block.pairs || []).map((p, i) => ({ id: "s" + i, title: p.right })),
      pins: (block.pairs || []).map((p, i) => ({
        id: "p" + i,
        text: p.left,
        slot: "s" + i,
        why: "Пара сошлась.",
        miss: "Не та сторона контракта.",
        hit: "s" + i,
      })),
    };
  }
  if (block.kind === "order") {
    return { kind: "order", title: block.title, prompt: block.prompt, items: block.items };
  }
  if (block.kind === "yaml") {
    return {
      kind: "ticket",
      title: block.title,
      prompt: block.prompt,
      snippet: block.snippet,
      lines: block.lines,
    };
  }
  if (block.kind === "sql") {
    return {
      kind: "chat",
      title: block.title,
      setting: block.prompt,
      steps: [{ from: "SQL · пет", line: block.prompt, options: block.options }],
    };
  }
  return block;
}

function flow(title, caption, nodes, notes) {
  const edges = [];
  for (let i = 0; i < nodes.length - 1; i++) {
    edges.push({ from: nodes[i].id, to: nodes[i + 1].id, note: notes[i] || "" });
  }
  return { kind: "flow", title, caption, nodes, edges };
}

export const LAB_BOARDS = {
  "lab-bank-ledger": {
    kind: "taccount",
    title: "Ночная книга",
    caption: "Refund без DR мерчанта разъезжает trial. Холд — не нога.",
    debit: [{ id: "dr", text: "DR клиента на capture", ok: true }],
    credit: [{ id: "cr", text: "CR мерчанта / CR клиента на дыре" }],
    hole: "Ищите ногу, которой нет — не поле balance.",
  },
  "lab-intern-1": flow(
    "KYC и кошелёк",
    "Личность и лимит — разные мастера.",
    [
      { id: "kyc", label: "KYC", sub: "личность", detail: "Статус личности. Не ACTIVE, пока не APPROVED." },
      { id: "wallet", label: "Кошелёк", sub: "лимит", detail: "Лимит после APPROVED. Не рисует KYC." },
      { id: "sms", label: "SMS", sub: "OTP", detail: "Доставка кода. Не решение впускать." },
    ],
    ["шов", "не мастер"],
  ),
  "lab-intern-2": {
    kind: "ticket",
    title: "Тикет ShopLine",
    caption: "Ключ заказа и 409. Не «как у WB».",
    lines: [
      { id: "1", text: "POST /orders с idempotency-key", why: "Контракт." },
      { id: "2", text: "Заказы должны синхронизироваться", why: "Нет ключа.", bad: true },
      { id: "3", text: "409 если order_id есть", why: "Fail." },
    ],
  },
};

export function boardForLab(lab) {
  return (
    LAB_BOARDS[lab.id] ||
    flow(
      lab.title,
      lab.teaser,
      [
        { id: "need", label: "Need", sub: "зачем", detail: lab.setting },
        { id: "seam", label: "Шов", sub: "контракт", detail: "Ключ, ошибка, мастер данных." },
        { id: "fail", label: "Fail", sub: "409 / дыра", detail: "Что будет, если повторить или потерять ack." },
      ],
      ["спека", "проверка"],
    )
  );
}

export const DRILLS = [
  {
    id: "hold",
    title: "Холд — ещё не проводка",
    teaser: "Касса авторизует 12 000. Журнал молчит, пока нет capture.",
    grade: "intern",
    write: true,
    explain:
      "Холд — банк временно откладывает сумму. Деньги ещё не списаны в журнал. Available уменьшается. Проводка появится на capture.",
    board: flow(
      "Холд",
      "OPEN холд. Available уменьшается. Ног нет.",
      [
        { id: "wallet", label: "Борис", sub: "available", detail: "До auth видит полный ledger." },
        { id: "hold", label: "Холд", sub: "OPEN", mood: "hold", detail: "12 000 отложены. Журнал пуст." },
        { id: "journal", label: "Журнал", sub: "позже", detail: "Ноги появятся, когда касса сделает capture." },
      ],
      ["auth", "capture"],
    ),
    read: [{ method: "GET", path: "/api/v1/holds", label: "Смотреть холды" }],
    writes: [
      {
        method: "PATCH",
        path: "/api/v1/contract",
        label: "Включить holdThenCapture",
        body: { holdThenCapture: true },
      },
      {
        method: "POST",
        path: "/api/v1/cards/authorize",
        label: "Auth 12 000 · карта Бориса",
        body: { walletId: "wal_boris", cardId: "crd_boris", amount: 12000, currency: "UZS" },
      },
    ],
  },
  {
    id: "p2p",
    title: "P2P: перевод и идемпотентность",
    teaser: "Девять кейсов в учебном банке: baseline, ключ, повтор, дубль без ключа, контракт, 409. Метод, путь и JSON на экране.",
    grade: "intern",
    write: true,
    explain:
      "P2P — внутренний перевод. В журнале: дебет отправителя, кредит получателя, один journalId. Без Idempotency-Key повтор создаёт вторую проводку — учебный дефект контракта.",
    board: {
      kind: "taccount",
      title: "P2P Анна → Борис",
      caption: "Пассив: отдала — дебет, получил — кредит.",
      debit: [{ id: "dr", text: "DR Анна 5 000", ok: true }],
      credit: [{ id: "cr", text: "CR Борис 5 000" }],
      hole: "Второй POST без ключа — вторая пара проводок. Это дефект контракта, не «обновите баланс».",
    },
    read: [{ method: "GET", path: "/api/v1/wallets", label: "Кошельки" }],
    writes: [
      {
        method: "POST",
        path: "/api/v1/transfers",
        label: "P2P 5 000",
        body: { fromWalletId: "wal_anna", toWalletId: "wal_boris", amount: 5000, currency: "UZS" },
      },
      {
        method: "POST",
        path: "/api/v1/transfers",
        label: "Повторить без ключа",
        body: { fromWalletId: "wal_anna", toWalletId: "wal_boris", amount: 5000, currency: "UZS" },
      },
    ],
  },
  {
    id: "iban",
    title: "Входящий IBAN",
    teaser: "25 000 на IBAN Анны. Nostro растёт, клиент кредитуется.",
    grade: "intern",
    write: true,
    explain:
      "IBAN — номер счёта, по которому приходит платёж из другого банка. Nostro — наш счёт у корреспондента (актив): входящий пишет DR nostro, CR клиента.",
    board: flow(
      "IBAN in",
      "DR nostro · CR Анна",
      [
        { id: "bank", label: "Банк", sub: "входящий", detail: "Платёж 25 000 UZS на IBAN Анны." },
        { id: "nostro", label: "Nostro", sub: "DR", detail: "Актив банка вырос." },
        { id: "anna", label: "Анна", sub: "CR", detail: "Пассив клиента вырос. Это журнал, не Excel." },
      ],
      ["зачисление", "книга"],
    ),
    read: [{ method: "GET", path: "/api/v1/bank/incoming", label: "Входящие" }],
    writes: [
      {
        method: "POST",
        path: "/api/v1/bank/incoming",
        label: "Зачислить 25 000 на Анну",
        body: { iban: "UZ12MALO000000000001", amount: 25000, currency: "UZS", paymentRef: "drill-iban" },
      },
    ],
  },
  {
    id: "hole",
    title: "Дыра в журнале",
    teaser: "Capture, затем refund без реверса. Trial UZS разъедется.",
    grade: "intern",
    write: true,
    explain:
      "Журнал — книга проводок: откуда ушли и куда пришли. Возврат без DR мерчанта — дыра: клиенту накинули, мерчант не уменьшили. Пробный баланс не 0.",
    board: {
      kind: "taccount",
      title: "Refund дырой",
      caption: "Контракт по умолчанию не пишет реверс. Это учебный баг.",
      debit: [{ id: "ghost", text: "DR мерчанта — этой ноги нет", ok: false }],
      credit: [{ id: "cr", text: "CR клиента (накинули)" }],
      hole: "Trial DR ≠ CR. Не UPDATE wallets — сторно или две ноги реверса.",
    },
    read: [{ method: "GET", path: "/api/v1/ledger/trial-balance", label: "Пробный баланс" }],
    writes: [
      {
        method: "POST",
        path: "/api/v1/cards/authorize",
        label: "Списание 8 000 (без холда в дефолте)",
        body: { walletId: "wal_boris", cardId: "crd_boris", amount: 8000, currency: "UZS" },
        keep: "payment",
      },
      {
        method: "POST",
        path: "/api/v1/payments/{id}/refund",
        label: "Refund без реверса",
        pathFrom: "payment",
        body: {},
      },
    ],
  },
];

export const QUESTS = [
  {
    id: "q-wallet",
    gradeId: "intern",
    title: "Двойное списание P2P",
    teaser: "Клиент видит два перевода по 5 000. PO предлагает поправить available вручную. Что сделает аналитик?",
    minutes: 8,
    board: {
      kind: "taccount",
      title: "Жалоба на P2P",
      caption: "Два SUCCESS. Один Need: не снимать дважды.",
      debit: [{ id: "dr", text: "DR Анна ×2 ?", ok: false }],
      credit: [{ id: "cr", text: "CR Борис ×2 ?" }],
      hole: "Ключ идемпотентности или дыра. Не поле balance.",
    },
    desk: {
      kind: "chat",
      title: "Ночной Slack",
      setting: "Анна видит два списания по 5 000.",
      steps: [
        {
          from: "PO",
          line: "Накиньте ей в кошелёк, завтра разберёмся с логом.",
          options: [
            { text: "UPDATE available и закрываем смену.", good: false, why: "Дыра останется. 1С ест журнал.", hit: "dr" },
            {
              text: "Сначала журнал: сколько ног с одним ключом. Потом сторно или дыра в контракте.",
              good: true,
              why: "Need — один перевод, одно движение.",
              hit: "dr",
            },
            { text: "Пусть саппорт извинится, деньги не трогаем.", good: false, why: "Клиент уже списан дважды.", hit: "cr" },
          ],
        },
        {
          from: "Dev",
          line: "Повтор без Idempotency-Key у нас 201. Так быстрее.",
          options: [
            { text: "Ок, ключ потом.", good: false, why: "Это и есть баг пета. На проде так нельзя оставлять «потом».", hit: "dr" },
            {
              text: "Контракт: ключ обязателен, повтор — 409, баланс тот же. Сейчас воспроизведём на /pet.",
              good: true,
              why: "Учебный баг оставляем в пете, в требованиях закрываем.",
              hit: "cr",
            },
          ],
        },
      ],
    },
  },
  {
    id: "q-hold",
    gradeId: "intern",
    title: "Касса SUCCESS — есть ли проводка?",
    teaser: "Терминал показал успех, ACS не ответил. Писать ли журнал или только холд?",
    minutes: 7,
    board: flow(
      "Касса vs журнал",
      "Auth может врать SUCCESS.",
      [
        { id: "till", label: "Касса", sub: "SUCCESS?", detail: "Терминал показал успех." },
        { id: "hold", label: "Холд", sub: "OPEN?", mood: "hold", detail: "Если контракт честный — только холд." },
        { id: "journal", label: "Журнал", sub: "ноги?", detail: "Нельзя писать SUCCESS в ledger на timeout." },
      ],
      ["резерв", "capture"],
    ),
    desk: {
      kind: "chat",
      title: "ACS timeout",
      setting: "3DS не ответил. Шлюз прислал неизвестность.",
      steps: [
        {
          from: "Касса",
          line: "Клиенту на экране успех. Проводите.",
          options: [
            { text: "Пишем SUCCESS в журнал, чтобы касса не краснела.", good: false, why: "Ложное списание. ACS timeout → UNKNOWN.", hit: "journal" },
            {
              text: "Холд или UNKNOWN. Журнал — на capture, когда деньги реально ушли.",
              good: true,
              why: "Касса ≠ книга.",
              hit: "hold",
            },
          ],
        },
      ],
    },
  },
  {
    id: "q-iban",
    gradeId: "intern",
    title: "Входящий на неизвестный IBAN",
    teaser: "Операции предлагают зачислить «похожему» клиенту. Куда класть деньги до идентификации?",
    minutes: 6,
    board: flow(
      "Неясные суммы",
      "Чужой IBAN — suspense.",
      [
        { id: "in", label: "Входящий", sub: "IBAN", detail: "Номер не наш." },
        { id: "suspense", label: "Suspense", sub: "неясные", detail: "DR nostro, CR suspense. Не сосед." },
        { id: "alloc", label: "Разнос", sub: "allocate", detail: "Когда узнали клиента — отдельная проводка." },
      ],
      ["неясные", "разнос"],
    ),
    desk: {
      kind: "chat",
      title: "Операции",
      setting: "IBAN на одну цифру мимо Анны.",
      steps: [
        {
          from: "Операции",
          line: "Похоже на Анну. Зачислите, утром поправим.",
          options: [
            { text: "На Анну, быстрее для NPS.", good: false, why: "Чужие деньги на чужом пассиве.", hit: "in" },
            { text: "Suspense. Разнос — когда ключ платежа и IBAN сходятся.", good: true, why: "Неясные суммы — отдельный счёт.", hit: "suspense" },
          ],
        },
      ],
    },
  },
  {
    id: "q-senior-recon",
    gradeId: "senior",
    title: "Сверка эквайера ночью",
    teaser: "Файл партнёра и наш реестр разошлись на 12 000. MISMATCH — не повод править ledger руками.",
    minutes: 10,
    board: {
      kind: "ticket",
      title: "Реестр vs файл",
      caption: "Ключ rrn+сумма. Не «примерно, копейки не смотрим».",
      lines: [
        { id: "1", text: "MATCH по rrn и сумме", why: "Шов." },
        { id: "2", text: "MISMATCH — холд выплаты, не UPDATE", why: "Деньги." },
        { id: "3", text: "Сверить глазами в Excel", why: "Ловушка.", bad: true },
      ],
    },
    desk: {
      kind: "chat",
      title: "Ночная сверка",
      setting: "Diff 12 000. Главбух в Slack.",
      steps: [
        {
          from: "Главбух",
          line: "Накиньте клиенту, файл партнёра master.",
          options: [
            { text: "Партнёр master, правим ledger руками.", good: false, why: "Файл не владеет статусом карты.", hit: "3" },
            {
              text: "Ключ rrn, допуск 0, MISMATCH → тикет и холд клиринга. Книгу не трогаем ночью.",
              good: true,
              why: "Сверка — контракт, не глаз.",
              hit: "2",
            },
          ],
        },
      ],
    },
  },
];
