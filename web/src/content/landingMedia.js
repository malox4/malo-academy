/** Free Unsplash photos, stored locally. License: https://unsplash.com/license */

export const LANDING_CREDITS =
  "Фото: Unsplash · William Iven, Austin Distel, Green Chameleon, Fotis Fotopoulos, Campaign Creators, Carlos Muza, Amy Hirschi, Nastuh Abootalebi, Scott Graham";

export const LANDING_ROOMS = [
  {
    id: "bank",
    n: "01",
    title: "Учебный банк",
    sell: "Живой HTTP ядра. Intern читает кошельки и журнал. PRO вызывает POST.",
    proof: "GET /api/v1/wallets · POST /api/v1/cards/authorize · POST /api/v1/transfers",
    img: "/media/landing/bank-pay.jpg",
    alt: "Касса принимает бесконтактную оплату: тот контур, что учебный холд",
  },
  {
    id: "lessons",
    n: "02",
    title: "48 уроков intern → senior",
    sell: "Intern открыт сразу. Junior, Middle, Senior — в PRO и в триале.",
    proof: "Термин первым, затем пример с цифрой и тикетом. Дальше не запирается.",
    img: "/media/landing/lesson-notes.jpg",
    alt: "Рука чертит схему на бумаге — так начинается урок аналитика",
  },
  {
    id: "labs",
    n: "03",
    title: "78 лаб и квесты смены",
    sell: "Intern-лабы бесплатно. Остальные — когда пишете в банк.",
    proof: "Лаборатория к статье, не тест из десяти галочек.",
    img: "/media/landing/lab-code.jpg",
    alt: "Экран с кодом: SQL и контракт API рядом с постановкой",
  },
  {
    id: "interview",
    n: "04",
    title: "Собес и Live",
    sell: "Intern-вопросы открыты. Middle+ — 57 вопросов в PRO. /live без оплаты.",
    proof: "Три числа кошелька, идемпотентность, дыра журнала.",
    img: "/media/landing/interview-table.jpg",
    alt: "Стол собеседования: ноутбуки, экран, вопрос со смены",
  },
];

export const LANDING_STRIP = [
  {
    img: "/media/landing/bank-pay.jpg",
    stamp: "касса",
    alt: "Оплата картой на кассе: холд до проводки",
    cap: "Холд до проводки. Intern видит GET, PRO жмёт POST.",
  },
  {
    img: "/media/landing/analytics.jpg",
    stamp: "данные",
    alt: "Дашборд на ноутбуке: три числа читают как контракт, не как картинку",
    cap: "Ledger, холд, available — те же три числа, что в банке.",
  },
  {
    img: "/media/landing/workshop.jpg",
    stamp: "смена",
    alt: "Стена историй: команда разбирает постановку, не вебинар",
    cap: "Стена историй. Дальше — лаба и собес, не Zoom.",
  },
];

export const LANDING_PATH_SHOTS = [
  "/media/landing/lesson-notes.jpg",
  "/media/landing/documents.jpg",
  "/media/landing/analytics.jpg",
  "/media/landing/lab-code.jpg",
  "/media/landing/bank-pay.jpg",
  "/media/landing/workshop.jpg",
];

export const LANDING_PATH_ALTS = [
  "Конспект урока: карандаш и схема на столе",
  "Постановка на бумаге: правило, которое можно проверить",
  "Экран с метриками: ledger, холд, available",
  "Лаба у экрана: запрос к ядру, не слайд",
  "Касса: POST холда в учебный банк",
  "Стена историй: решить, остаётесь ли в PRO",
];
