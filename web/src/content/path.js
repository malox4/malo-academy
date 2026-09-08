/** Path nodes on intern→senior levels. Does not replace the 48 lesson ids. */

export const PATH_EXTRAS = {
  "intern-1": ["soft-raci", "car-portfolio", "nfr-wcag", "car-cert"],
  "intern-2": ["req-validate", "req-test-analysis", "req-lifecycle"],
  "intern-3": ["db-layers", "proc-figma", "api-postman", "bank-ledger-deep", "proc-uml-hold", "proc-mobile-web"],
  "junior-1": ["api-schema", "api-pagination"],
  "junior-2": ["api-openapi", "db-window", "db-cte", "db-migrate", "api-asyncapi", "api-graphql"],
  "junior-3": ["proc-bpmn-collab", "nfr-security", "nfr-pci-log", "bank-cp-cnp", "db-dwh"],
  "middle-1": ["req-gost-34602"],
  "middle-2": ["bank-iso20022", "bank-recon", "bank-aml-kyc"],
  "middle-3": ["arch-saga", "nfr-slo", "arch-eventstorm", "arch-styles", "arch-eip"],
  "senior-1": ["req-estimate"],
  "senior-2": ["soft-workshop"],
  "senior-3": ["arch-arc42"],
};

export const LANDING_PATH = [
  { n: "01", t: "Открываете Intern", d: "Регистрация. 12 уроков: кто SA, Need, RACI. Потока и куратора нет." },
  { n: "02", t: "Пишете постановку", d: "Given / When / Then для QA. Эмодзи в чате — не торг." },
  { n: "03", t: "Смотрите три числа", d: "ledger, холд, available на учебном кошельке. Чтение открыто." },
  { n: "04", t: "Разбираете intern-лабу", d: "Проводки журнала, вызов API, ошибка холда. SQL — витрины intern_*." },
  { n: "05", t: "Берёте 3 дня PRO", d: "POST в /api/v1, полный SQL, junior+ уроки, все 78 лаб." },
  { n: "06", t: "Оставляете PRO или Intern", d: "Триал кончился — Intern остаётся. Платите, если нужна запись в банк." },
];

export function extrasForLevel(levelId) {
  return PATH_EXTRAS[levelId] || [];
}
