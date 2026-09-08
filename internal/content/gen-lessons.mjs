import { writeFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { dirname, join } from "node:path";
import { intel } from "./intel.mjs";
import { TOPICS, LESSON_TOPIC } from "../../web/src/content/topics.js";

const root = dirname(fileURLToPath(import.meta.url));

const ids = [
  "intern-1-profession", "intern-1-sdlc", "intern-1-team", "intern-1-questions",
  "intern-2-requirement", "intern-2-levels", "intern-2-stakeholders", "intern-2-quality",
  "intern-3-elicitation", "intern-3-notes", "intern-3-story", "intern-3-meeting",
  "junior-1-stories", "junior-1-ac", "junior-1-usecase", "junior-1-clarify",
  "junior-2-brd", "junior-2-srs", "junior-2-trace", "junior-2-write",
  "junior-3-uml", "junior-3-bpmn", "junior-3-gap", "junior-3-explain",
  "middle-1-prio", "middle-1-usm", "middle-1-cjm", "middle-1-no",
  "middle-2-api", "middle-2-sql", "middle-2-integration", "middle-2-devtalk",
  "middle-3-nfr", "middle-3-arch", "middle-3-errors", "middle-3-workshop",
  "senior-1-strategy", "senior-1-options", "senior-1-eval", "senior-1-change",
  "senior-2-conflict", "senior-2-politics", "senior-2-exec", "senior-2-think",
  "senior-3-mentor", "senior-3-discovery", "senior-3-lead", "senior-3-studio",
];

const meta = {
  "intern-1-profession": ["Кто такой аналитик", "Бизнес-аналитик работает с потребностью. Системный аналитик — с правилами системы.", 10],
  "intern-1-sdlc": ["Как устроен IT-проект", "Путь от боли до прода. В платежах «разберёмся в проде» дорого.", 10],
  "intern-1-team": ["Кто в комнате", "PO, SM, dev, QA, ops — кто за что отвечает на жалобе.", 9],
  "intern-1-questions": ["Вопросы, которые открывают правду", "Не «какую кнопку». Кто страдает и чем измерим.", 10],
  "intern-2-requirement": ["Что такое требование", "Не цитата из Slack. Проверяемое правило для системы.", 10],
  "intern-2-levels": ["Уровни требований", "Бизнес → человек → правило → проверка. Не прыгайте сразу в технологию.", 10],
  "intern-2-stakeholders": ["Карта людей", "Кто пользуется, кто вето, кого забыли в ночную смену.", 9],
  "intern-2-quality": ["Качество формулировки", "INVEST-гигиена без лозунгов. Ловите брак в тикете.", 10],
  "intern-3-elicitation": ["Как достать факты", "Наблюдение смены бьёт опросник на 40 вопросов.", 10],
  "intern-3-notes": ["Заметки созвона", "Решение / открытый вопрос / обещание — три корзины.", 9],
  "intern-3-story": ["Первая история", "Роль, действие, ценность. Без «как обсуждали».", 10],
  "intern-3-meeting": ["Созвон без паники", "Тон, повтор слов, никакого «завтра починим».", 10],
  "junior-1-stories": ["User Story", "INVEST на живом платеже, не на шаблоне.", 11],
  "junior-1-ac": ["Acceptance Criteria", "Given/When/Then, которые QA не перепишет.", 11],
  "junior-1-usecase": ["Use Case", "Актёр, основной поток, альтернативы, ошибка шлюза.", 11],
  "junior-1-clarify": ["Уточнения", "Три вопроса, без которых разработку рано стартовать.", 10],
  "junior-2-brd": ["BRD", "Зачем меняемся. Не список экранов.", 11],
  "junior-2-srs": ["SRS", "Что система делает. Статусы, данные, ошибки.", 11],
  "junior-2-trace": ["Трассировка", "Need → история → правило → тест. Дырки видны.", 10],
  "junior-2-write": ["Ясный текст", "Глагол, условие, исход. Вычеркните «корректно обработать».", 10],
  "junior-3-uml": ["UML, который читают", "Состояния платежа, не павлин из стрелок.", 11],
  "junior-3-bpmn": ["BPMN смены", "As-is взыскания: кто ждёт, где ручной костыль.", 12],
  "junior-3-gap": ["As-Is / To-Be", "Разрыв — не «добавим кнопку».", 11],
  "junior-3-explain": ["Объяснить схеме", "PO и dev должны увидеть одно и то же.", 10],
  "middle-1-prio": ["Приоритеты", "MoSCoW на живом бэклоге двойных списаний.", 11],
  "middle-1-usm": ["User Story Mapping", "Хребет сценария коллектора, не гора стикеров.", 11],
  "middle-1-cjm": ["CJM жалобы", "Эмоция и канал на каждом шаге звонка.", 11],
  "middle-1-no": ["Как сказать нет", "Резать scope без войны в канале.", 10],
  "middle-2-api": ["Контракт API", "Поля, идемпотентность, 409. Дырка видна на схеме.", 12],
  "middle-2-sql": ["Данные, не отчёт", "Какая таблица врёт в карточке долга.", 11],
  "middle-2-integration": ["Интеграции", "Callback шлюза, повтор, timeout ACS.", 12],
  "middle-2-devtalk": ["Разговор с dev", "Правило, не «сделайте как в Тинькофф».", 10],
  "middle-3-nfr": ["NFR", "Идемпотентность, SLA ночи, объём сверки.", 11],
  "middle-3-arch": ["Контур, не слайд", "Где холд, где журнал, где nostro.", 12],
  "middle-3-errors": ["Ошибки", "Что видит клиент, что пишет журнал, кто дежурит.", 11],
  "middle-3-workshop": ["Воркшоп", "90 минут: need, правило, тест, go/no-go.", 10],
  "senior-1-strategy": ["Стратегия изменения", "Зачем портфель, не фича.", 11],
  "senior-1-options": ["Опции решения", "Процесс / система / ничего. Цена каждой.", 11],
  "senior-1-eval": ["Оценка", "Риск, срок, дыра в trial. Не «мне кажется».", 11],
  "senior-1-change": ["Impact", "Кого заденет ночной релиз.", 10],
  "senior-2-conflict": ["Конфликт", "PO жмёт срок, ops видит дубль. Ваш ход.", 11],
  "senior-2-politics": ["Политика", "Кто теряет KPI, если скрыть вторую строку.", 10],
  "senior-2-exec": ["C-level", "Одна минута: риск, деньги, решение.", 10],
  "senior-2-think": ["Мышление", "Отделить гипотезу от факта на доске.", 10],
  "senior-3-mentor": ["Менторство", "Стажёр в комнате. Как не сломать тон.", 10],
  "senior-3-discovery": ["Discovery", "Неделя фактов до решения, не после релиза.", 11],
  "senior-3-lead": ["Ownership", "Кто хозяин need, когда спринт уже начался.", 11],
  "senior-3-studio": ["Студия контура", "Соберите доску холда сами — это экзамен зала.", 12],
};

function flow(title, caption, nodes, edges) {
  return { kind: "flow", title, caption, nodes, edges };
}
function taccount(title, caption, left, right, hole) {
  return { kind: "taccount", title, caption, debit: left, credit: right, hole };
}
function ticketBoard(title, caption, lines) {
  return { kind: "ticket", title, caption, lines };
}
function compareBoard(title, caption, leftTitle, rightTitle, rows) {
  return { kind: "compare", title, caption, leftTitle, rightTitle, rows };
}
function desk(title, prompt, columns, cards, extra = {}) {
  return { kind: "desk", title, prompt, columns, cards, ...extra };
}
function chat(title, setting, steps, extra = {}) {
  return { kind: "chat", title, setting, steps, ...extra };
}
function ticketDesk(title, prompt, lines, extra = {}) {
  return { kind: "ticket", title, prompt, lines, ...extra };
}
function ledgerDesk(title, prompt, legs, missing, extra = {}) {
  return { kind: "ledger", title, prompt, legs, missing, ...extra };
}
function mapDesk(title, prompt, slots, pins, extra = {}) {
  return { kind: "map", title, prompt, slots, pins, ...extra };
}

const lessons = [];
const topicIndex = Object.fromEntries(TOPICS.map((t) => [t.id, t]));

function pack(id, board, desk, brief) {
  const i = ids.indexOf(id);
  const [title, teaser, minutes] = meta[id];
  const gradeId = id.split("-")[0];
  const levelId = id.split("-").slice(0, 2).join("-");
  const materialId = LESSON_TOPIC[id];
  const topic = materialId ? topicIndex[materialId] : null;
  lessons.push({
    id, gradeId, levelId, title, teaser, minutes,
    next: ids[i + 1] || "",
    prev: ids[i - 1] || "",
    board, desk,
    ...(intel[id] || {}),
    ...(brief ? { brief } : {}),
    ...(topic ? { materialId: topic.id, primer: topic.sections } : {}),
  });
}

pack("intern-1-profession",
  flow("Роль аналитика в команде", "Аналитик связывает потребность стейкхолдера с проверяемым требованием к системе.",
    [
      { id: "biz", label: "Стейкхолдер", sub: "заказчик, пользователь", detail: "Описывает проблему в работе. Часто предлагает готовое решение." },
      { id: "ba", label: "BA / SA", sub: "потребность → требование", detail: "Отделяет потребность от дизайна. Формулирует то, что можно проверить." },
      { id: "eng", label: "Разработка", sub: "dev, QA", detail: "Нужны границы, данные, ошибки — не фраза «чтобы больше так не было»." },
    ],
    [{ from: "biz", to: "ba", note: "потребность" }, { from: "ba", to: "eng", note: "требование" }]),
  desk("Что говорит аналитик на старте", "Разложите реплики: подходит к роли или нет.",
    [{ id: "do", title: "Роль аналитика" }, { id: "dont", title: "Не эта роль" }],
    [
      { id: "a", text: "Назвать роль и цель: понять, какая работа ломается — не чинить интеграцию с первого звонка.", column: "do", hit: "ba", why: "Сначала потребность, затем решение." },
      { id: "b", text: "«Я только стажёр, если что — поправьте, я ничего не знаю».", column: "dont", hit: "ba", miss: "Роль не требует извинений. Нужна рамка вопросов." },
      { id: "c", text: "Повторить слова стейкхолдера, чтобы зафиксировать факт, а не домысел.", column: "do", hit: "biz", why: "Потребность сформулирована словами источника." },
      { id: "d", text: "«Завтра починим» — срок разработки аналитик не назначает.", column: "dont", hit: "eng", miss: "Обещание даты — не требование и не оценка команды." },
      { id: "e", text: "Спросить, какой одной сумме должен верить оператор в разговоре с клиентом.", column: "do", hit: "biz", why: "Потребность измерена: одна сумма, не две." },
      { id: "f", text: "Пообещать конкретную технологию (очередь, шлюз), чтобы выглядеть компетентно.", column: "dont", hit: "eng", miss: "Решение раньше требования. Разработка получит технологию без правила." },
    ],
    { probes: [
      { from: "Владелец продукта", line: "Кнопку «Пересчитать», завтра в спринт. Пользователи жалуются.", options: [
        { text: "Сначала потребность: две суммы, какой одной верить в разговоре с клиентом?", good: true, why: "Потребность до интерфейса.", hit: "ba" },
        { text: "Открыть тикет на кнопку.", good: false, why: "Это оформление чужого решения, не аналитика.", hit: "ba" },
      ]},
      { from: "Оператор", line: "Я читаю верхнюю строку, иначе путаюсь.", options: [
        { text: "Зафиксировать как факт работы: верхняя строка = то, чему верят. Это потребность, не «ошибка человека».", good: true, why: "Стейкхолдер описал работу.", hit: "biz" },
        { text: "Посоветовать обновлять страницу, пока сумма не станет одной.", good: false, why: "Маскировка экрана. Требования к данным нет.", hit: "eng" },
      ]},
    ] }),
  { title: "Потребность одной фразой", prompt: "Кто стейкхолдер и какое одно число должно быть на карточке?", placeholder: "Оператор видит два списания…", minChars: 40 });

pack("intern-1-sdlc",
  flow("Жизнь платежа", "Понять → решить → построить → проверить → выпустить → смотреть ночь. Дырка на любом этапе — чужие деньги.",
    [
      { id: "need", label: "Need", sub: "кто страдает", detail: "Без изменения анализа нет." },
      { id: "build", label: "Сборка", sub: "правила и тесты", detail: "История закрывает кусок need, не «как в чате»." },
      { id: "rel", label: "Релиз", sub: "дежурство", detail: "Откат сложнее: часть платежей уже ушла в банк." },
    ],
    [{ from: "need", to: "build", note: "рамка" }, { from: "build", to: "rel", note: "не «в проде разберёмся»" }]),
  desk("Куда это класть", "Разложите фразы по этапам. Ошибка зажигает need или релиз красным.",
    [{ id: "need", title: "Need" }, { id: "ship", title: "Релиз" }],
    [
      { id: "a", text: "Одна сумма в карточке, которой можно верить в разговоре.", column: "need", hit: "need", why: "Это изменение работы, не релиз." },
      { id: "b", text: "Кто дежурит в ночь, если callback придёт дважды.", column: "ship", hit: "rel", why: "Дежурство — кусок релиза." },
      { id: "c", text: "Что считается аварией, а не «посмотрим завтра».", column: "ship", hit: "rel", why: "Авария ночи — не need-формулировка." },
      { id: "d", text: "Кто стейкхолдер — не «бизнес».", column: "need", hit: "need", why: "Имя человека до сборки." },
      { id: "e", text: "Чем QA провалит правило дубля.", column: "need", hit: "build", why: "Тест живёт рядом со сборкой, но вопрос задают в анализе." },
      { id: "f", text: "Критерий: повторный callback не создаёт вторую ногу.", column: "need", hit: "build", why: "Правило и тест формулируют до катки. Сборка на доске живая." },
    ],
    { probes: [
      { from: "Dev", line: "UI спрячет вторую строку, в проде посмотрим журнал.", options: [
        { text: "Релиз без правила дубля — 12 000 дыра и Excel до полудня. Сначала need и тест.", good: true, why: "Сборка горит. Релиз не врёт.", hit: "build" },
        { text: "Ок, катим, журнал потом.", good: false, why: "«В проде разберёмся». Узел релиза красный: деньги уже в банке.", hit: "rel" },
      ]},
    ] }));

pack("intern-1-team",
  flow("Комната на жалобе", "Один тред, четыре роли. Путаница ролей — пустой тикет.",
    [
      { id: "po", label: "PO", sub: "срок и ценность", detail: "Режет scope. Не пишет SRS." },
      { id: "ba", label: "Аналитик", sub: "need и правило", detail: "Не секретарь и не начальник разработки." },
      { id: "dev", label: "Dev / QA", sub: "как проверить", detail: "Без AC будут угадывать." },
      { id: "ops", label: "Ops", sub: "ночь и сверка", detail: "Видит дубль первым." },
    ],
    [{ from: "po", to: "ba", note: "ценность" }, { from: "ba", to: "dev", note: "правило" }, { from: "ops", to: "ba", note: "факт смены" }]),
  mapDesk("Кто хозяин чего", "Поставьте людей на зоны. Не туда — краснеет их узел.",
    [{ id: "need", title: "Need" }, { id: "code", title: "Код" }, { id: "night", title: "Ночь" }],
    [
      { id: "po", text: "PO", slot: "need", hit: "po", why: "PO режет ценность, не пишет код." },
      { id: "dev", text: "Разработчик", slot: "code", hit: "dev", why: "Код и тест — его зона." },
      { id: "ops", text: "Ночной ops", slot: "night", hit: "ops", why: "Факт дубля ночью." },
      { id: "ba", text: "Аналитик", slot: "need", hit: "ba", why: "Need и правило — шов ba." },
    ],
    { probes: [
      { from: "PO", line: "Пусть аналитик сам напишет SRS и задеплоит, мы в демо.", options: [
        { text: "SRS — правило для dev/QA. Деплой — не аналитик. Я держу need, ops несёт ночь.", good: true, why: "Роли на местах. Четыре узла живые.", hit: "ba" },
        { text: "Ок, я всё закрою, не тормозите демо.", good: false, why: "Пустой тикет. PO горит: он скинул свою работу и вашу.", hit: "po" },
      ]},
    ] }));

pack("intern-1-questions",
  flow("Вопрос до кнопки", "Сначала боль и факт. Потом решение.",
    [
      { id: "who", label: "Кто", sub: "конкретный человек", detail: "Не «бизнес»." },
      { id: "brk", label: "Что ломается", sub: "в работе сегодня", detail: "Не эмоция «жалуются»." },
      { id: "val", label: "Зачем", sub: "ценность", detail: "Риск, время, деньги, клиент." },
    ],
    [{ from: "who", to: "brk", note: "факт" }, { from: "brk", to: "val", note: "зачем менять" }]),
  chat("Первые 20 секунд", "Супервайзер взыскания, PO, два коллектора. Вас добавили стажёром. Каждый ход зажигает или гасит узел.",
    [
      { from: "Супервайзер", line: "Ну давайте уже кнопку Пересчитать, люди орут.", options: [
        { text: "Можно я сначала повторю, что ломается в звонке, и какая одна сумма должна быть на экране?", good: true, why: "Need до решения. Узел «что ломается» живой.", hit: "brk" },
        { text: "Ок, заведу тикет на кнопку, завтра в спринт.", good: false, why: "Вы почтовый ящик. Узел brk погас — боли в тикете нет.", hit: "brk" },
      ]},
      { from: "PO", line: "Бизнес хочет, чтобы больше так не было. Какую кнопку рисуем?", options: [
        { text: "Кто именно в звонке страдает — имя и роль, не «бизнес»?", good: true, why: "Узел «кто» загорелся. Есть человек, не абстракция.", hit: "who" },
        { text: "Рисуем Пересчитать, как просили.", good: false, why: "Кто = никто. Узел who красный.", hit: "who" },
      ]},
      { from: "Коллектор", line: "Мне всё равно как, лишь бы не орали.", options: [
        { text: "Зачем менять: чтобы в звонке была одна сумма 12 000, которой можно верить — иначе жалоба и риск.", good: true, why: "Ценность названа цифрой. Узел val живой.", hit: "val" },
        { text: "Ладно, как удобнее вам на экране.", good: false, why: "Нет зачем. Узел val красный, фича ниоткуда.", hit: "val" },
      ]},
    ]));

pack("intern-2-requirement",
  ticketBoard("Это ещё не требование", "Требование — правило, которое QA может провалить. Копипаст чата — логистика.",
    [
      { id: "1", text: "Починить двойные списания. Как обсуждали. Срочно.", bad: true, why: "Нет need, нет условия, нет исхода." },
      { id: "2", text: "По одному payment_id в карточке долга одна успешная сумма.", bad: false, why: "Можно проверить." },
      { id: "3", text: "Сделать как у конкурента.", bad: true, why: "Чужое решение без нашей боли." },
    ]),
  ticketDesk("Вычеркните брак", "Отметьте строки, которые нельзя отдавать в спринт.",
    [
      { id: "1", text: "Починить двойные списания. Как обсуждали. Срочно.", bad: true, why: "Нет need." },
      { id: "2", text: "Повторный callback с тем же payment_id не создаёт второе списание.", bad: false, why: "Правило." },
      { id: "3", text: "Ну чтобы больше так не было.", bad: true, why: "Непроверяемо." },
      { id: "4", text: "Если шлюз прислал дубль — в CRM остаётся одна сумма, статус DUPLICATE.", bad: false, why: "Есть условие и исход." },
      { id: "5", text: "Сделать как у конкурента, им же ок.", bad: true, why: "Чужое how без нашей боли." },
      { id: "6", text: "QA: два callback, COUNT ног = 1, карточка без второй строки SUCCESS.", bad: false, why: "Тест." },
    ]));

pack("intern-2-levels",
  flow("Четыре этажа", "Прыжок сразу на Kafka — типичная дыра. Сначала зачем, потом как.",
    [
      { id: "biz", label: "Бизнес", sub: "зачем портфелю", detail: "Меньше ложных звонков и жалоб." },
      { id: "stk", label: "Стейкхолдер", sub: "чья работа", detail: "Коллектор видит одну сумму." },
      { id: "sol", label: "Решение", sub: "правило системы", detail: "Идемпотентность callback." },
      { id: "tr", label: "Тест", sub: "как провалить", detail: "Повторный POST с тем же id." },
    ],
    [{ from: "biz", to: "stk", note: "кто" }, { from: "stk", to: "sol", note: "что" }, { from: "sol", to: "tr", note: "проверка" }]),
  desk("Разложите этажи", "Куда относится фраза?",
    [{ id: "biz", title: "Бизнес" }, { id: "sol", title: "Система" }],
    [
      { id: "a", text: "Меньше претензий «с меня уже сняли».", column: "biz" },
      { id: "b", text: "Ключ идемпотентности на callback.", column: "sol" },
      { id: "c", text: "Коллектор не звонит со второй суммой.", column: "biz" },
      { id: "d", text: "Статус DUPLICATE, журнал без второй ноги клиенту.", column: "sol" },
    ]));

pack("intern-2-stakeholders",
  flow("Кого забывают", "Ночной ops и комплаенс не пишут в Slack днём. Их всё равно заденет.",
    [
      { id: "use", label: "Пользуется", sub: "коллектор", detail: "Экран и скрипт звонка." },
      { id: "veto", label: "Вето", sub: "комплаенс", detail: "Нельзя врать сумму должнику." },
      { id: "night", label: "Ночь", sub: "ops", detail: "Callback в 03:12." },
    ],
    [{ from: "use", to: "veto", note: "риск" }, { from: "night", to: "use", note: "факт после смены" }]),
  mapDesk("Карта стейкхолдеров", "Куда поставить человека.",
    [{ id: "use", title: "Пользуется" }, { id: "veto", title: "Вето" }, { id: "side", title: "Рядом страдает" }],
    [
      { id: "col", text: "Коллектор", slot: "use" },
      { id: "cmp", text: "Комплаенс", slot: "veto" },
      { id: "sup", text: "Поддержка", slot: "side" },
      { id: "po", text: "PO", slot: "veto" },
    ]));

pack("intern-2-quality",
  ticketBoard("Тикет на ревью", "Ищите размытые глаголы и скрытые решения.",
    [
      { id: "1", text: "Система корректно обрабатывает повтор.", bad: true, why: "«Корректно» не тест." },
      { id: "2", text: "При повторном callback с тем же id вторая проводка не создаётся.", bad: false, why: "Можно провалить." },
    ]),
  ticketDesk("Флаг брака", "Отметьте строки, которые QA вернёт.",
    [
      { id: "1", text: "Система корректно обрабатывает повтор.", bad: true, why: "Пустой глагол." },
      { id: "2", text: "Available уменьшается на холд, ledger молчит до capture.", bad: false, why: "Два факта, можно проверить." },
      { id: "3", text: "Сделать удобно для коллектора.", bad: true, why: "Не измерить." },
      { id: "4", text: "В карточке по payment_id одна строка SUCCESS.", bad: false, why: "Счётный исход." },
    ]));

pack("intern-3-elicitation",
  flow("Откуда факты", "Смена > тред. Люди предлагают кнопку вместо боли.",
    [
      { id: "watch", label: "Смена", sub: "послушать звонок", detail: "Слышите «с меня уже сняли»." },
      { id: "tik", label: "Тикеты", sub: "поддержка", detail: "Повтор одной жалобы — паттерн." },
      { id: "reg", label: "Регламент", sub: "старый", detail: "Часто врёт относительно прода." },
    ],
    [{ from: "watch", to: "tik", note: "сверка" }, { from: "tik", to: "reg", note: "что устарело" }]),
  chat("Наблюдение", "Вы сидите рядом с коллектором. На экране два списания.",
    [{ from: "Коллектор", line: "Я просто верхнюю сумму читаю, иначе запутаюсь.", options: [
      { text: "Записать: в работе берут верхнюю строку, даже если она дубль. Это факт процесса.", good: true, why: "Наблюдение, не совет про кнопку. Узел смены живой.", hit: "watch" },
      { text: "Скажите ему жать F5, пока не станет одна.", good: false, why: "Вы чините экран за него. Узел смены красный.", hit: "watch" },
    ]}]));

pack("intern-3-notes",
  compareBoard("Три корзины заметок", "Если всё в одну кучу — решения потеряются в цитатах.",
    "Решение", "Ещё не решение",
    [
      { label: "Пример", left: "В карточке одна сумма SUCCESS.", right: "Кажется, виноват шлюз." },
      { label: "Хозяин", left: "Есть имя и срок.", right: "«Кто-то из ops»." },
    ]),
  desk("Разберите конспект", "Куда класть строку с созвона.",
    [{ id: "dec", title: "Решение" }, { id: "oq", title: "Открытый вопрос" }],
    [
      { id: "a", text: "Не стартуем разработку без need на полэкрана.", column: "dec" },
      { id: "b", text: "Повторяет ли шлюз callback с тем же id ночью?", column: "oq" },
      { id: "c", text: "В тикете запрещён копипаст Slack.", column: "dec" },
      { id: "d", text: "Кто видит карточку на ночной смене?", column: "oq" },
    ]));

pack("intern-3-story",
  flow("История = роль + действие + ценность", "Без ценности это просто задача в Jira.",
    [
      { id: "role", label: "Роль", sub: "коллектор", detail: "Конкретный человек." },
      { id: "want", label: "Действие", sub: "видеть одну сумму", detail: "Не «чтобы система»." },
      { id: "so", label: "Ценность", sub: "не врать должнику", detail: "Зачем в работе." },
    ],
    [{ from: "role", to: "want", note: "хочет" }, { from: "want", to: "so", note: "чтобы" }]),
  desk("Соберите историю", "Что в карточку истории, что в мусор.",
    [{ id: "in", title: "В историю" }, { id: "out", title: "Мусор" }],
    [
      { id: "a", text: "Как коллектор хочу одну сумму в карточке, чтобы не слышать претензию.", column: "in" },
      { id: "b", text: "Как обсуждали в Slack.", column: "out" },
      { id: "c", text: "Сделать Kafka.", column: "out" },
      { id: "d", text: "Чтобы звонок шёл с суммой, которой можно верить.", column: "in" },
    ]));

pack("intern-3-meeting",
  flow("Тон в комнате", "Тихий уверенный тон сильнее извинений. Молчать всю встречу — тоже ошибка.",
    [
      { id: "in", label: "Войти", sub: "роль одной фразой", detail: "Не секретарь, не сеньор." },
      { id: "rep", label: "Повтор", sub: "их словами", detail: "Проверка, что поняли." },
      { id: "out", label: "Выход", sub: "3 вопроса", detail: "Не обещание даты." },
    ],
    [{ from: "in", to: "rep", note: "слушать" }, { from: "rep", to: "out", note: "рамка" }]),
  chat("Первый созвон", "Все знают, что вы стажёр.",
    [{ from: "PO", line: "Короче заведи как обсудили, не тормози спринт.", options: [
      { text: "Тикет сегодня. В описании need и три вопроса. Без них разработку лучше не стартовать.", good: true, why: "Выход с рамкой. Узел out живой.", hit: "out" },
      { text: "По BABOK нужен воркшоп на два часа, тикет не заведу.", good: false, why: "Теория токсичной подачей. Вход сломан.", hit: "in" },
    ]}]));

pack("junior-1-stories",
  compareBoard("INVEST на платеже", "Независимая история всё ещё должна закрывать кусок need.",
    "Живая", "Мёртвая",
    [
      { label: "I", left: "Можно выпустить без «всего контура».", right: "Тянет Kafka, UI и сверку одним куском." },
      { label: "T", left: "QA знает, чем провалить.", right: "«Корректно»." },
    ]),
  desk("INVEST-стол", "Что оставить в истории.",
    [{ id: "keep", title: "Оставить" }, { id: "cut", title: "Вырезать" }],
    [
      { id: "a", text: "Повторный callback не создаёт второе списание.", column: "keep" },
      { id: "b", text: "И ещё перекрасить CRM.", column: "cut" },
      { id: "c", text: "И сразу MT103.", column: "cut" },
      { id: "d", text: "В карточке одна сумма SUCCESS.", column: "keep" },
    ]));

pack("junior-1-ac",
  ticketBoard("AC, которые живут", "Given — факт мира. When — действие. Then — наблюдаемый исход.",
    [
      { id: "1", text: "Given повторный callback с тем же payment_id When он приходит Then второй SUCCESS не создаётся.", bad: false, why: "Проваливаемо." },
      { id: "2", text: "Then система работает стабильно.", bad: true, why: "Не наблюдаемо." },
    ]),
  ticketDesk("Вычеркните пустой Then", "Отметьте бракованные критерии.",
    [
      { id: "1", text: "Then система работает стабильно.", bad: true, why: "Пусто." },
      { id: "2", text: "Then в журнале нет второй DR-ноги клиенту.", bad: false, why: "Журнал — факт." },
      { id: "3", text: "Then коллектору удобно.", bad: true, why: "Ощущение." },
      { id: "4", text: "Then available не меняется, если это дубль без нового холда.", bad: false, why: "Деньги на экране." },
    ]));

pack("junior-1-usecase",
  flow("Use case платежа", "Актёр не «система». Основной поток и то, что ломается.",
    [
      { id: "act", label: "Актёр", sub: "шлюз / касса", detail: "Кто инициирует." },
      { id: "main", label: "Основной", sub: "auth → hold", detail: "Счастливый путь." },
      { id: "alt", label: "Дубль", sub: "тот же id", detail: "Альтернатива, не сноска." },
    ],
    [{ from: "act", to: "main", note: "auth" }, { from: "main", to: "alt", note: "повтор" }]),
  desk("Куда класть ветку", "Основной поток или ошибка.",
    [{ id: "main", title: "Основной" }, { id: "err", title: "Ошибка / дубль" }],
    [
      { id: "a", text: "Касса шлёт auth, холд OPEN.", column: "main" },
      { id: "b", text: "ACS timeout.", column: "err" },
      { id: "c", text: "Capture ≤ auth.", column: "main" },
      { id: "d", text: "Повтор с тем же ключом.", column: "err" },
    ]));

pack("junior-1-clarify",
  flow("Три вопроса до кода", "Без них тикет — лотерея.",
    [
      { id: "idemp", label: "Ключ", sub: "что уникально", detail: "payment_id? RRN?" },
      { id: "night", label: "Ночь", sub: "кто видит дубль", detail: "Ops или коллектор утром." },
      { id: "fail", label: "Провал", sub: "как QA ломает", detail: "Два POST подряд." },
    ],
    [{ from: "idemp", to: "night", note: "когда" }, { from: "night", to: "fail", note: "тест" }]),
  chat("PO жмёт", "«Не надо усложнять, заведи как есть».",
    [{ from: "PO", line: "Три вопроса — это бюрократия.", options: [
      { text: "Без ключа дубля, хозяина ночи и теста мы снова попадём в прод. Это полэкрана, не воркшоп.", good: true, why: "Риск, не церемония. Узел fail живой.", hit: "fail" },
      { text: "Ладно, потом допишем.", good: false, why: "Потом = в проде. Узел fail красный.", hit: "fail" },
    ]}]));

pack("junior-2-brd",
  compareBoard("BRD ≠ список экранов", "BRD отвечает зачем. Экраны — позже.",
    "BRD", "Не BRD",
    [
      { label: "Вопрос", left: "Какое изменение в работе?", right: "Какой цвет кнопки?" },
      { label: "Успех", left: "Меньше ложных звонков за 4 недели.", right: "Тикет закрыт." },
    ]),
  desk("В BRD или нет", "Отсортируйте куски документа.",
    [{ id: "in", title: "В BRD" }, { id: "out", title: "Не сюда" }],
    [
      { id: "a", text: "Ценность: меньше претензий должникам.", column: "in" },
      { id: "b", text: "Swagger endpoint.", column: "out" },
      { id: "c", text: "Стейкхолдеры ночной смены.", column: "in" },
      { id: "d", text: "CSS карточки.", column: "out" },
    ]));

pack("junior-2-srs",
  ticketBoard("SRS живёт статусами", "Если нет таблицы состояний — это ещё BRD.",
    [
      { id: "1", text: "Платёж: AUTHORIZED → CAPTURED → REFUNDED.", bad: false, why: "Состояния." },
      { id: "2", text: "Система должна быть гибкой.", bad: true, why: "NFR-пустышка." },
    ]),
  ticketDesk("Брак в SRS", "Отметьте непригодные требования.",
    [
      { id: "1", text: "Система должна быть гибкой.", bad: true, why: "Пусто." },
      { id: "2", text: "Холд OPEN режет available, не ledger.", bad: false, why: "Правило денег." },
      { id: "3", text: "Как в старом файле, все знают.", bad: true, why: "Нет файла-истины." },
      { id: "4", text: "Refund без реверса мерчанта рвёт trial — запрещено в To-Be.", bad: false, why: "Явный запрет дыры." },
    ]));

pack("junior-2-trace",
  flow("Нить need → тест", "Дырка в нити = история в спринте ниоткуда.",
    [
      { id: "n", label: "Need", sub: "одна сумма", detail: "Боль коллектора." },
      { id: "s", label: "История", sub: "дубль callback", detail: "Кусок need." },
      { id: "r", label: "Правило", sub: "идемпотентность", detail: "SRS." },
      { id: "t", label: "Тест", sub: "два POST", detail: "QA." },
    ],
    [{ from: "n", to: "s", note: "режет" }, { from: "s", to: "r", note: "фиксирует" }, { from: "r", to: "t", note: "ломает" }]),
  mapDesk("Повесьте артефакт", "Куда он на нити.",
    [{ id: "n", title: "Need" }, { id: "r", title: "Правило" }, { id: "t", title: "Тест" }],
    [
      { id: "a", text: "Коллектор врёт сумму", slot: "n" },
      { id: "b", text: "Ключ payment_id", slot: "r" },
      { id: "c", text: "Повторный callback в Postman", slot: "t" },
    ]));

pack("junior-2-write",
  compareBoard("Глагол вместо тумана", "«Обеспечить корректность» не строится.",
    "Сильно", "Слабо",
    [
      { label: "Фраза", left: "Не создавать вторую проводку.", right: "Корректно обработать дубль." },
    ]),
  ticketDesk("Вычеркните туман", "Брак — размытый глагол.",
    [
      { id: "1", text: "Обеспечить корректную обработку.", bad: true, why: "Туман." },
      { id: "2", text: "Не создавать вторую DR-ногу клиенту.", bad: false, why: "Запрет действия." },
      { id: "3", text: "Оптимизировать процесс.", bad: true, why: "Нет метрики." },
      { id: "4", text: "Если ACS timeout — статус UNKNOWN, не SUCCESS.", bad: false, why: "Контрактный факт." },
    ]));

pack("junior-3-uml",
  flow("Состояния платежа", "Стрелка без события — украшение, не модель.",
    [
      { id: "a", label: "AUTHORIZED", sub: "hold OPEN", mood: "hold", detail: "Ledger молчит." },
      { id: "c", label: "CAPTURED", sub: "две ноги", detail: "Журнал живой." },
      { id: "r", label: "REFUNDED", sub: "реверс", detail: "Иначе дыра trial." },
    ],
    [{ from: "a", to: "c", note: "capture" }, { from: "c", to: "r", note: "refund" }]),
  desk("Событие на стрелку", "Какое событие двигает статус.",
    [{ id: "cap", title: "Capture" }, { id: "void", title: "Void" }],
    [
      { id: "a", text: "Холд → журнал, статус CAPTURED.", column: "cap" },
      { id: "b", text: "Холд RELEASED, журнала нет.", column: "void" },
      { id: "c", text: "Available возвращается без проводки.", column: "void" },
      { id: "d", text: "DR клиент CR мерчант.", column: "cap" },
    ]));

pack("junior-3-bpmn",
  flow("As-is звонка", "Ручной костыль — тоже шаг процесса, его надо нарисовать.",
    [
      { id: "call", label: "Звонок", sub: "читает верхнюю строку", detail: "Факт смены." },
      { id: "fight", label: "Спор", sub: "«уже сняли»", detail: "Клиент прав чаще, чем CRM." },
      { id: "esc", label: "Эскалация", sub: "поддержка + ops", detail: "Час сверки." },
    ],
    [{ from: "call", to: "fight", note: "две суммы" }, { from: "fight", to: "esc", note: "руками" }]),
  desk("Шаг или костыль", "Что рисуем как работу, что как обход.",
    [{ id: "work", title: "Работа" }, { id: "hack", title: "Костыль" }],
    [
      { id: "a", text: "Звонок по скрипту.", column: "work" },
      { id: "b", text: "Сверить шлюз в Excel ночью.", column: "hack" },
      { id: "c", text: "Закрыть контакт в CRM.", column: "work" },
      { id: "d", text: "Спрятать вторую строку руками.", column: "hack" },
    ]));

pack("junior-3-gap",
  compareBoard("Разрыв", "To-be не «добавим кнопку». Это другая работа смены.",
    "As-is", "To-be",
    [
      { label: "Экран", left: "Две строки SUCCESS.", right: "Одна сумма, дубль скрыт статусом." },
      { label: "Ночь", left: "Ops в Excel.", right: "Ключ не пускает вторую ногу." },
    ]),
  desk("Это разрыв?", "Кладём в gap или в «просто хотелку».",
    [{ id: "gap", title: "Gap" }, { id: "wish", title: "Хотелка" }],
    [
      { id: "a", text: "Две суммы vs одна достоверная.", column: "gap" },
      { id: "b", text: "Розовый бейдж в CRM.", column: "wish" },
      { id: "c", text: "Excel vs ключ идемпотентности.", column: "gap" },
      { id: "d", text: "Тёмная тема админки.", column: "wish" },
    ]));

pack("junior-3-explain",
  flow("Одна схема — две аудитории", "PO видит боль. Dev видит статус. Схема должна держать оба взгляда.",
    [
      { id: "po", label: "PO", sub: "ценность", detail: "Меньше жалоб." },
      { id: "sch", label: "Схема", sub: "статусы", detail: "AUTHORIZED не деньги мерчанта." },
      { id: "dev", label: "Dev", sub: "журнал", detail: "Где ноги." },
    ],
    [{ from: "po", to: "sch", note: "зачем" }, { from: "sch", to: "dev", note: "как" }]),
  chat("Покажите схему", "Dev: «это же просто баг UI».",
    [{ from: "Dev", line: "Спрячем вторую строку и всё.", options: [
      { text: "Схема: если ledger уже две ноги — прятать строку врёт коллектору ещё сильнее. Сначала журнал.", good: true, why: "Боль + контур. Узел схемы живой.", hit: "sch" },
      { text: "Ок, давайте CSS.", good: false, why: "UI закрыл симптом. Узел dev красный: журнал не тронут.", hit: "dev" },
    ]}]));

pack("middle-1-prio",
  compareBoard("MoSCoW на дублях", "Must — то, без чего снова позвоним с неверной суммой.",
    "Must", "Won't сейчас",
    [
      { label: "Кусок", left: "Идемпотентность callback.", right: "Новый дашборд NPS." },
    ]),
  desk("Must или later", "Разложите бэклог.",
    [{ id: "must", title: "Must" }, { id: "later", title: "Later" }],
    [
      { id: "a", text: "Ключ дубля callback.", column: "must" },
      { id: "b", text: "Анимация карточки.", column: "later" },
      { id: "c", text: "Одна сумма в карточке долга.", column: "must" },
      { id: "d", text: "Геймификация сбора.", column: "later" },
    ]));

pack("middle-1-usm",
  flow("Хребет смены", "Карта историй идёт вдоль работы человека, не вдоль микросервисов.",
    [
      { id: "open", label: "Открыл карточку", sub: "видит сумму", detail: "Первый шаг." },
      { id: "talk", label: "Говорит", sub: "скрипт", detail: "Доверяет экрану." },
      { id: "esc", label: "Эскалация", sub: "если спор", detail: "Редкий, дорогой." },
    ],
    [{ from: "open", to: "talk", note: "доверие" }, { from: "talk", to: "esc", note: "если враньё" }]),
  mapDesk("На хребет", "Куда клеить историю.",
    [{ id: "open", title: "Экран" }, { id: "talk", title: "Разговор" }, { id: "esc", title: "Эскалация" }],
    [
      { id: "a", text: "Одна сумма SUCCESS", slot: "open" },
      { id: "b", text: "Скрипт без «верхней строки»", slot: "talk" },
      { id: "c", text: "Тикет поддержки не открывается на дубль", slot: "esc" },
    ]));

pack("middle-1-cjm",
  flow("Путь жалобы", "Эмоция — данные. Канал — тоже.",
    [
      { id: "sms", label: "SMS банка", sub: "списание", detail: "Клиент уже «заплатил»." },
      { id: "call", label: "Звонок", sub: "вторая сумма", detail: "Гнев." },
      { id: "app", label: "Чат банка", sub: "тикет", detail: "Расход поддержки." },
    ],
    [{ from: "sms", to: "call", note: "не стыкуется" }, { from: "call", to: "app", note: "эскалация" }]),
  desk("Канал боли", "Где ломается опыт.",
    [{ id: "client", title: "Клиент" }, { id: "ops", title: "Операции" }],
    [
      { id: "a", text: "Слышит вторую сумму.", column: "client" },
      { id: "b", text: "Сверяет шлюз в Excel.", column: "ops" },
      { id: "c", text: "Пишет в чат «уже сняли».", column: "client" },
      { id: "d", text: "Ночной callback без хозяина.", column: "ops" },
    ]));

pack("middle-1-no",
  flow("Как резать", "Нет — это рамка, не хамство.",
    [
      { id: "val", label: "Ценность", sub: "что спасаем", detail: "Звонок с верной суммой." },
      { id: "cut", label: "Срез", sub: "что не в этом релизе", detail: "Дашборд NPS." },
      { id: "why", label: "Почему", sub: "риск", detail: "Иначе снова прод с дырой." },
    ],
    [{ from: "val", to: "cut", note: "жертва" }, { from: "cut", to: "why", note: "явно" }]),
  chat("PO: давайте ещё дашборд", "Спринт уже набит дублем.",
    [{ from: "PO", line: "Ну и график жалоб заодно, это же мелочь.", options: [
      { text: "Дашборд не закрывает need звонка. Если влезает — после ключа дубля. Иначе снова прячем симптом.", good: true, why: "Нет с рамкой. Узел cut живой.", hit: "cut" },
      { text: "Ок, накидаем, успеем.", good: false, why: "Мелочь съест ключ. Узел cut красный.", hit: "cut" },
    ]}]));

pack("middle-2-api",
  taccount("Идемпотентный POST", "Повтор с тем же ключом — тот же ответ, не вторая пара ног.",
    [{ id: "d1", text: "DR клиент 5 000", ok: true }, { id: "d2", text: "DR клиент 5 000 ещё раз", ok: false }],
    [{ id: "c1", text: "CR Борис 5 000", ok: true }],
    "Вторая DR — дыра контракта, пока ключ не обязателен."),
  ledgerDesk("Какая нога лишняя", "Нажмите на проводку, которой не должно быть при дубле ключа.",
    [
      { id: "1", text: "DR Анна 5 000 · ключ tap-1", hole: false },
      { id: "2", text: "CR Борис 5 000 · ключ tap-1", hole: false },
      { id: "3", text: "DR Анна 5 000 · тот же ключ повторно", hole: true },
    ], "3"));

pack("middle-2-sql",
  ticketBoard("Какая строка врёт", "Карточка долга читает не ту таблицу — классика.",
    [
      { id: "1", text: "SELECT * FROM payments WHERE status='SUCCESS' — без уникальности payment_id.", bad: true, why: "Две строки." },
      { id: "2", text: "Одна строка на payment_id, дубль в колонке duplicate_of.", bad: false, why: "Модель." },
    ]),
  ticketDesk("Брак в запросе", "Отметьте то, что даст две суммы.",
    [
      { id: "1", text: "Все SUCCESS без группировки по payment_id.", bad: true, why: "Дубли видны как живые." },
      { id: "2", text: "Один агрегат суммы по payment_id.", bad: false, why: "Одна цифра." },
      { id: "3", text: "ORDER BY created_at DESC LIMIT 1 без фильтра дубля.", bad: true, why: "Верхняя может быть ошибкой." },
    ]));

pack("middle-2-integration",
  flow("Callback шлюза", "Повтор и timeout — не «сеть пошалила», это ветки контракта.",
    [
      { id: "pos", label: "POS", sub: "auth", detail: "Ещё не деньги мерчанта." },
      { id: "acs", label: "ACS", sub: "timeout?", detail: "UNKNOWN ≠ SUCCESS." },
      { id: "cb", label: "Callback", sub: "тот же id", detail: "Идемпотентность." },
    ],
    [{ from: "pos", to: "acs", note: "3DS" }, { from: "acs", to: "cb", note: "повтор" }]),
  desk("Куда класть сбой", "Ветка контракта.",
    [{ id: "hold", title: "Холд жив" }, { id: "bug", title: "Баг завода" }],
    [
      { id: "a", text: "Timeout → UNKNOWN, холд не открыт.", column: "hold" },
      { id: "b", text: "Timeout → SUCCESS и сразу журнал.", column: "bug" },
      { id: "c", text: "Повтор ключа → 409.", column: "hold" },
      { id: "d", text: "Повтор ключа → вторая пара ног.", column: "bug" },
    ]));

pack("middle-2-devtalk",
  flow("Говорить правилом", "«Как у Тинькофф» — не спецификация. Журнал — да.",
    [
      { id: "ui", label: "Чужой UI", sub: "не наш ledger", detail: "Экран не спасает ноги." },
      { id: "rule", label: "Правило", sub: "ключ дубля", detail: "Можно проверить." },
    ],
    [{ from: "ui", to: "rule", note: "переводим" }]),
  chat("Стенд-ап", "Dev: «скажите как у Тинькофф».",
    [{ from: "Dev", line: "Ну сделаем как у них, зачем правило.", options: [
      { text: "Наше правило: повтор callback с тем же id не пишет вторую ногу. Их экран нас не спасает в журнале.", good: true, why: "Контур, не бренд. Узел rule живой.", hit: "rule" },
      { text: "Ок, скопируем UI.", good: false, why: "Чужой UI ≠ наш ledger. Узел ui красный.", hit: "ui" },
    ]}]));

pack("middle-3-nfr",
  flow("NFR денег", "Идемпотентность, ночной SLA, объём сверки — это требования, не «потом ops».",
    [
      { id: "id", label: "Идемпотентность", sub: "ключ", detail: "Обязателен на P2P и callback." },
      { id: "sla", label: "Ночь", sub: "кто дежурит", detail: "Callback 03:12." },
      { id: "vol", label: "Объём", sub: "сверка T+1", detail: "Файл Orient." },
    ],
    [{ from: "id", to: "sla", note: "когда" }, { from: "sla", to: "vol", note: "сколько" }]),
  desk("NFR или косметика", "Разложите.",
    [{ id: "nfr", title: "NFR" }, { id: "ui", title: "UI" }],
    [
      { id: "a", text: "Повтор ключа не пишет журнал.", column: "nfr" },
      { id: "b", text: "Скругление карточки 16px.", column: "ui" },
      { id: "c", text: "Сверка файла до 09:00.", column: "nfr" },
      { id: "d", text: "Иконка кафе.", column: "ui" },
    ]));

pack("middle-3-arch",
  flow("Где деньги", "Холд ≠ журнал ≠ nostro. Перепутаете — нарисуете баг как фичу.",
    [
      { id: "h", label: "Hold", sub: "available", detail: "OPEN режет, ledger молчит." },
      { id: "j", label: "Journal", sub: "две ноги", detail: "Capture." },
      { id: "n", label: "Nostro", sub: "банк", detail: "Settle / IBAN." },
    ],
    [{ from: "h", to: "j", note: "capture" }, { from: "j", to: "n", note: "T+1" }]),
  mapDesk("Поставьте слой", "Куда кладётся факт.",
    [{ id: "h", title: "Hold" }, { id: "j", title: "Journal" }, { id: "n", title: "Nostro" }],
    [
      { id: "a", text: "Auth 3DS", slot: "h" },
      { id: "b", text: "DR клиент CR кафе", slot: "j" },
      { id: "c", text: "Входящий IBAN", slot: "n" },
    ]));

pack("middle-3-errors",
  ticketBoard("Кто что видит", "Ошибка клиента, журнал и дежурный — три разных текста.",
    [
      { id: "1", text: "Клиенту: «не прошло». Журналу: ничего. Дежурному: ACS UNKNOWN.", bad: false, why: "Честно." },
      { id: "2", text: "Клиенту SUCCESS, журналу две ноги, ACS молчит.", bad: true, why: "Заводской баг." },
    ]),
  ticketDesk("Ложный успех", "Отметьте враньё статуса.",
    [
      { id: "1", text: "ACS timeout = SUCCESS + журнал.", bad: true, why: "Ложный capture." },
      { id: "2", text: "ACS timeout = UNKNOWN, холда нет.", bad: false, why: "Честный контракт." },
      { id: "3", text: "Клиенту «ок», ledger дыра.", bad: true, why: "Экран врёт." },
    ]));

pack("middle-3-workshop",
  flow("90 минут", "Need, правило, тест, go/no-go. Без романа.",
    [
      { id: "n", label: "15 мин need", sub: "кто страдает", detail: "Факт смены." },
      { id: "r", label: "правило", sub: "журнал", detail: "Что запрещено." },
      { id: "g", label: "go/no-go", sub: "релиз", detail: "Кто дежурит." },
    ],
    [{ from: "n", to: "r", note: "фиксируем" }, { from: "r", to: "g", note: "риск" }]),
  chat("Воркшоп срывается", "Все спорят про Kafka.",
    [{ from: "Архитектор", line: "Давайте сразу шину, иначе несерьёзно.", options: [
      { text: "Шина — опция после правила дубля. Сейчас go/no-go: ключ и ночное дежурство. Иначе снова Excel.", good: true, why: "Рамка. Узел g живой.", hit: "g" },
      { text: "Ок, давайте Kafka в этом воркшопе.", good: false, why: "Need потерян. Узел n красный.", hit: "n" },
    ]}]));

pack("senior-1-strategy",
  flow("Зачем портфелю", "Фича без стратегии — локальный костыль, который сломает сверку.",
    [
      { id: "port", label: "Портфель", sub: "риск жалоб", detail: "Регулятор рядом." },
      { id: "need", label: "Need", sub: "одна сумма", detail: "Работа взыскания." },
      { id: "opt", label: "Опции", sub: "процесс/система", detail: "Не сразу Kafka." },
    ],
    [{ from: "port", to: "need", note: "боль" }, { from: "need", to: "opt", note: "выбор" }]),
  desk("Стратегия или фича", "Куда относится.",
    [{ id: "st", title: "Стратегия" }, { id: "ft", title: "Фича" }],
    [
      { id: "a", text: "Снизить ложные контакты по долгу.", column: "st" },
      { id: "b", text: "Бейдж в CRM.", column: "ft" },
      { id: "c", text: "Честный trial после refund.", column: "st" },
      { id: "d", text: "Новый цвет холда.", column: "ft" },
    ]));

pack("senior-1-options",
  compareBoard("Три опции", "Ничего не делать — тоже опция с ценой.",
    "Система", "Процесс",
    [
      { label: "Цена", left: "Ключ, журнал, тесты.", right: "Скрипт: не читать верхнюю строку." },
      { label: "Риск", left: "Дольше, но держит ночь.", right: "Люди забудут в 03:12." },
    ]),
  desk("Опция на стол", "Куда класть ход.",
    [{ id: "sys", title: "Система" }, { id: "proc", title: "Процесс" }, { id: "no", title: "Ничего" }],
    [
      { id: "a", text: "Идемпотентность callback.", column: "sys" },
      { id: "b", text: "Скрипт супервайзера.", column: "proc" },
      { id: "c", text: "Ждать, пока банк починит.", column: "no" },
    ]));

pack("senior-1-eval",
  taccount("Цена дыры", "Refund без реверса: trial не бьётся. Это аргумент для C-level, не «мне кажется».",
    [{ id: "d", text: "DR нет (дыра)", ok: false }],
    [{ id: "c", text: "CR клиенту refund", ok: true }],
    "Пробный баланс по UZS красный — вот метрика опции."),
  ledgerDesk("Где дыра", "Нажмите ногу, которой не хватает при честном refund.",
    [
      { id: "1", text: "CR клиент (возврат)", hole: false },
      { id: "2", text: "DR мерчант (реверс capture)", hole: true },
    ], "2"));

pack("senior-1-change",
  flow("Кого заденет релиз", "Ночь, коллекторы, банк, комплаенс. Забыли — инцидент.",
    [
      { id: "col", label: "Смена", sub: "экран", detail: "Утро после релиза." },
      { id: "ops", label: "Ops", sub: "callback", detail: "03:12." },
      { id: "bnk", label: "Банк", sub: "файл", detail: "Сверка T+1." },
    ],
    [{ from: "ops", to: "col", note: "утром" }, { from: "bnk", to: "ops", note: "mismatch" }]),
  mapDesk("Impact", "Куда повесить риск.",
    [{ id: "people", title: "Люди" }, { id: "money", title: "Деньги" }, { id: "law", title: "Регулятор" }],
    [
      { id: "a", text: "Коллектор", slot: "people" },
      { id: "b", text: "Двойное списание", slot: "money" },
      { id: "c", text: "Жалоба должника", slot: "law" },
    ]));

pack("senior-2-conflict",
  flow("Срок против факта", "PO держит демо. Ops держит ночь. Рамка — need, не победитель спора.",
    [
      { id: "po", label: "PO", sub: "демо", detail: "Velocity." },
      { id: "you", label: "Вы", sub: "рамка", detail: "Need на слайд." },
      { id: "ops", label: "Ops", sub: "дубль", detail: "03:12." },
    ],
    [{ from: "po", to: "you", note: "срок" }, { from: "ops", to: "you", note: "факт" }]),
  chat("PO vs ops", "PO: «в срок». Ops: «ночью дубль». Вы в середине.",
    [{ from: "PO", line: "Не будем тормозить демо.", options: [
      { text: "Демо без ключа дубля снова покажет две суммы. Давайте need на слайд и no-go, если ключа нет.", good: true, why: "Конфликт в рамку. Узел you живой.", hit: "you" },
      { text: "Ops преувеличивает, катим.", good: false, why: "Сторона срока против факта. Узел ops красный.", hit: "ops" },
    ]}]));

pack("senior-2-politics",
  flow("Кто теряет KPI", "Спрятать строку спасает дашборд PO и топит комплаенс.",
    [
      { id: "po", label: "PO", sub: "velocity", detail: "Тикет закрыт." },
      { id: "cmp", label: "Комплаенс", sub: "правда суммы", detail: "Нельзя врать." },
      { id: "ops", label: "Ops", sub: "сверка", detail: "Excel растёт." },
    ],
    [{ from: "po", to: "cmp", note: "конфликт KPI" }, { from: "cmp", to: "ops", note: "факт" }]),
  mapDesk("Чей KPI", "Поставьте интерес.",
    [{ id: "speed", title: "Срок" }, { id: "truth", title: "Правда суммы" }],
    [
      { id: "po", text: "PO", slot: "speed" },
      { id: "cmp", text: "Комплаенс", slot: "truth" },
      { id: "col", text: "Коллектор", slot: "truth" },
    ]));

pack("senior-2-exec",
  flow("Минута наверху", "Риск, деньги, ход. Не методология.",
    [
      { id: "r", label: "Риск", sub: "жалоба", detail: "Врут сумму." },
      { id: "m", label: "Деньги", sub: "недели", detail: "Не квартал." },
      { id: "h", label: "Ход", sub: "ключ дубля", detail: "Потом дашборды." },
    ],
    [{ from: "r", to: "m", note: "цена" }, { from: "m", to: "h", note: "решение" }]),
  chat("Минута у C-level", "Вас спросили: «это баг или проект?»",
    [{ from: "CFO", line: "Коротко: сколько стоит и что будет, если нет.", options: [
      { text: "Сейчас коллектор врёт сумму. Риск жалоб. Ключ дубля — недели, не квартал. Без него refund и сверка останутся дырой.", good: true, why: "Риск, деньги, ход. Узел h живой.", hit: "h" },
      { text: "По BABOK мы в discovery, дайте слайды.", good: false, why: "Церемония вместо решения. Узел r красный.", hit: "r" },
    ]}]));

pack("senior-2-think",
  compareBoard("Факт и гипотеза", "На доске их надо развести, иначе спор про вкус.",
    "Факт", "Гипотеза",
    [
      { label: "Пример", left: "В журнале две DR при одном ключе.", right: "Наверное, виноват шлюз." },
    ]),
  desk("На две кучи", "Факт или догадка.",
    [{ id: "f", title: "Факт" }, { id: "h", title: "Гипотеза" }],
    [
      { id: "a", text: "Две ноги при одном Idempotency-Key.", column: "f" },
      { id: "b", text: "Банк точно шлёт дубль специально.", column: "h" },
      { id: "c", text: "Trial не бьётся после refund без реверса.", column: "f" },
      { id: "d", text: "Коллекторам просто лень смотреть.", column: "h" },
    ]));

pack("senior-3-mentor",
  flow("Не ломать тон", "Стажёр в комнате. Вы держите рамку, он учится повторять need.",
    [
      { id: "you", label: "Вы", sub: "рамка", detail: "Не отдаёте войну." },
      { id: "st", label: "Стажёр", sub: "повтор need", detail: "Голос в конце." },
      { id: "room", label: "Комната", sub: "PO", detail: "Срок живой." },
    ],
    [{ from: "you", to: "st", note: "роль" }, { from: "st", to: "room", note: "need словами" }]),
  chat("Стажёр в комнате", "Он молчит 15 минут «чтобы не мешать».",
    [{ from: "Вы", line: "Как втянуть его, не отдав ему войну с PO?", options: [
      { text: "Попросите его повторить need своими словами в конце. Вы держите рамку.", good: true, why: "Роль + обучение. Узел st живой.", hit: "st" },
      { text: "Пусть сам спорит с PO, так быстрее вырастет.", good: false, why: "Сломаете тон и его. Узел st красный.", hit: "st" },
    ]}]));

pack("senior-3-discovery",
  flow("Неделя фактов", "Решение после фактов. Не наоборот.",
    [
      { id: "w", label: "Смена", sub: "послушать", detail: "Звонки." },
      { id: "j", label: "Журнал", sub: "ноги", detail: "Пет / прод-выборка." },
      { id: "d", label: "Решение", sub: "опция", detail: "Только потом." },
    ],
    [{ from: "w", to: "j", note: "сверка" }, { from: "j", to: "d", note: "не раньше" }]),
  ticketDesk("Рано решать", "Отметьте решения без фактов.",
    [
      { id: "1", text: "Сразу Kafka, смену не слушали.", bad: true, why: "Решение первым." },
      { id: "2", text: "Две DR в журнале на одном ключе — факт.", bad: false, why: "Можно строить опции." },
      { id: "3", text: "«Как у конкурента», без нашей боли.", bad: true, why: "Чужая карта." },
    ]));

pack("senior-3-lead",
  flow("Хозяин need", "Когда спринт начался, need всё равно живой. Кто его не бросает?",
    [
      { id: "you", label: "Аналитик", sub: "рамка", detail: "Не отдаёт пустой тикет." },
      { id: "po", label: "PO", sub: "scope", detail: "Режет, не подменяет need." },
      { id: "team", label: "Команда", sub: "правило", detail: "Строит по нити." },
    ],
    [{ from: "you", to: "po", note: "need жив" }, { from: "po", to: "team", note: "кусок" }]),
  mapDesk("Ownership", "Кто за что не бросает.",
    [{ id: "need", title: "Need" }, { id: "ship", title: "Релиз" }],
    [
      { id: "ba", text: "Аналитик", slot: "need" },
      { id: "po", text: "PO", slot: "ship" },
      { id: "ops", text: "Ops", slot: "ship" },
    ]));

pack("senior-3-studio",
  flow("Соберите холд сами", "Auth режет available. Capture пишет журнал. Void не пишет.",
    [
      { id: "pos", label: "Касса", sub: "auth", detail: "Экран «успех» ещё не деньги мерчанта." },
      { id: "h", label: "Hold OPEN", sub: "available", mood: "hold", detail: "Ledger молчит." },
      { id: "j", label: "Journal", sub: "capture", detail: "Две ноги." },
    ],
    [{ from: "pos", to: "h", note: "не проводка" }, { from: "h", to: "j", note: "capture ≤ auth" }]),
  ledgerDesk("Лишняя нога на auth", "На auth журнала быть не должно. Найдите запрещённую проводку.",
    [
      { id: "1", text: "Hold OPEN 12 000", hole: false },
      { id: "2", text: "DR клиент 12 000 сразу на auth", hole: true },
      { id: "3", text: "Available − 12 000", hole: false },
    ], "2"));

writeFileSync(join(root, "lessons.json"), JSON.stringify(lessons, null, 2));
writeFileSync(join(root, "materials.json"), JSON.stringify({ topics: TOPICS }, null, 2));

const boardIntel = {
  "hold-capture": {
    why: "Касса может показать успех, пока деньги мерчанта ещё не проведены. Hold (авторизация) уменьшает доступный остаток и не создаёт проводку. Capture записывает дебет и кредит.",
    facts: [
      "Authorization (auth): статус hold OPEN, журнал пуст, available уменьшен.",
      "Capture: дебет клиента, кредит мерчанта, холд закрыт.",
      "Cancel / void: холд снят, проводок нет.",
      "В учебном контракте таймаут ACS может вернуть SUCCESS — это учебная ошибка, не норма платёжной схемы.",
    ],
    scene: {
      title: "Касса и журнал",
      body: "Клиент приложил карту. Экран кассы зелёный. Доступный остаток уже меньше. Журнал пуст. Если это принять за списание, товар отдают раньше факта в учёте.",
      chips: ["AUTHORIZED", "available −", "журнал пуст"],
    },
    trap: {
      title: "Типичная ошибка",
      body: "Назвать холд проводкой. Тогда void пытаются «откатить журнал», которого не было.",
    },
  },
  "t-account": {
    why: "Двойная запись: каждая операция — дебет и кредит на одну сумму в одной валюте. Пробный баланс (trial balance) сходится в ноль. Одна сторона без второй — ошибка журнала, не «баг экрана».",
    facts: [
      "P2P 5 000: дебет Анны, кредит Бориса.",
      "Два перевода без ключа идемпотентности — две дебетовые проводки Анны. Так устроен учебный контракт.",
      "Trial balance по валюте должен сходиться в 0.",
      "Правка поля баланса (UPDATE) уничтожает историю сторно.",
    ],
    scene: {
      title: "Перевод Анна → Борис",
      body: "На счёте Анны 150 000. Два запроса по 5 000 без ключа. Доступно 140 000, в журнале две дебетовые проводки. Trial balance может сходиться — и это всё равно дефект продукта.",
      chips: ["150 000", "2 × 5 000", "trial = 0"],
    },
    trap: {
      title: "Типичная ошибка",
      body: "Смотреть только available. Без T-счёта не видно расхождения в журнале.",
    },
  },
  "iban-suspense": {
    why: "Входящий платёж на наш IBAN кредитует клиента. Неизвестный IBAN — счёт невыясненных сумм (suspense), пока не выполнят allocate. HTTP 404 здесь теряет уже полученные деньги.",
    facts: [
      "Входящий: дебет nostro, кредит клиента.",
      "Неизвестный IBAN: кредит suspense до allocate.",
      "Ответ 404 на уже зачисленный платёж оставляет сумму без хозяина.",
      "Allocate — отдельное требование, не постоянная ручная правка.",
    ],
    scene: {
      title: "Входящий файл",
      body: "Строка на IBAN, которого нет в справочнике. Если ответить 404, nostro уже вырос, клиента нет. Правильный учёт — suspense, затем allocate.",
      chips: ["nostro DR", "suspense"],
    },
    trap: {
      title: "Типичная ошибка",
      body: "Считать неизвестный счёт ошибкой API. Это учётная ситуация, не код ответа.",
    },
  },
  "bpmn-as-is": {
    why: "Ночная сверка в Excel — шаг текущего процесса (as-is). Его рисуют на схеме. Иначе to-be автоматизирует не тот участок.",
    facts: [
      "Повтор callback с тем же id ночью — штатное поведение шлюза.",
      "CRM получает вторую строку SUCCESS.",
      "Сопровождение сверяет файл вручную.",
      "To-be: ключ идемпотентности не создаёт вторую проводку; ручная сверка этого дубля не нужна.",
    ],
    scene: {
      title: "Ночной повтор",
      body: "Шлюз прислал callback повторно. В карточке две успешные суммы. Сверка в Excel до полудня. Оператор уже говорит с клиентом.",
      chips: ["повтор", "Excel", "as-is"],
    },
    trap: {
      title: "Типичная ошибка",
      body: "Нарисовать to-be без ручного шага. Не видно, какой костыль должна убрать автоматизация.",
    },
  },
  "hall-intro": {
    why: "Учебный банк на /api/v1 — тренажёр контракта. Часть ошибок оставлена специально: повтор без ключа создаёт вторую проводку. Это материал для разбора, не поломка сайта.",
    facts: [
      "POST /transfers без Idempotency-Key можно выполнить дважды.",
      "Проверяйте trial balance и T-счёт отправителя.",
      "Правило: ключ обязателен, повтор — HTTP 409.",
      "Reset возвращает учебное состояние контракта.",
    ],
    scene: {
      title: "Два одинаковых перевода",
      body: "Два запроса по 5 000 без ключа. Две дебетовые проводки. Так показывают отсутствие идемпотентности.",
      chips: ["/transfers", "без ключа"],
    },
    trap: {
      title: "Типичная ошибка",
      body: "«Починить» учебный банк, чтобы запросы всегда проходили. Тогда не видно дефекта контракта.",
    },
  },
};

const boards = [
  {
    id: "hold-capture",
    title: "Hold и capture",
    term: "Hold / Capture",
    learn: "Чем бронь суммы отличается от проводки в журнале.",
    teaser: "Авторизация бронирует деньги. Проводка появляется только при capture.",
    grade: "Intern",
    plan: "free",
    minutes: 4,
    board: flow("Hold → Capture", "Экран кассы может показать успех, пока в журнале ещё нет проводки.",
      [
        { id: "pos", label: "POS", sub: "authorization", detail: "Клиент на кассе. Авторизация — ещё не списание." },
        { id: "h", label: "Hold", sub: "OPEN · available −", mood: "hold", detail: "Доступный остаток уменьшен. Журнал пуст, пока нет capture." },
        { id: "j", label: "Журнал", sub: "дебет и кредит", detail: "Capture: дебет клиента, кредит мерчанта. Холд закрыт." },
        { id: "v", label: "Void", sub: "отмена холда", detail: "Пока статус AUTHORIZED, холд снимают без проводок." },
      ],
      [{ from: "pos", to: "h", note: "auth" }, { from: "h", to: "j", note: "capture" }, { from: "h", to: "v", note: "cancel" }]),
  },
  {
    id: "t-account",
    title: "Двойная запись",
    term: "T-account",
    learn: "Как читать дебет и кредит одной операции.",
    teaser: "Каждая операция — дебет и кредит на одну сумму. Trial balance сходится в ноль.",
    grade: "Intern",
    plan: "free",
    minutes: 5,
    board: taccount("P2P Анна → Борис", "Пассив клиента: отдать — дебет, получить — кредит. Сторно — зеркальные проводки, не правка баланса.",
      [{ id: "d", text: "Дебет Анны 5 000", ok: true }],
      [{ id: "c", text: "Кредит Бориса 5 000", ok: true }],
      "Если есть только одна сторона, пробный баланс не сходится. Это ошибка журнала, не экрана."),
  },
  {
    id: "iban-suspense",
    title: "IBAN и suspense",
    term: "IBAN / suspense",
    learn: "Куда относить входящий платёж с неизвестным счётом.",
    teaser: "Неизвестный IBAN — счёт невыясненных сумм, не ответ HTTP 404.",
    grade: "Junior",
    plan: "pro",
    minutes: 6,
    board: flow("Входящий платёж", "Наш IBAN кредитует клиента. Чужой — suspense, пока не allocate.",
      [
        { id: "or", label: "Банк-корреспондент", sub: "файл / MT", detail: "Пришло кредитовое сообщение." },
        { id: "n", label: "Nostro", sub: "дебет", detail: "Актив банка вырос." },
        { id: "cl", label: "Клиент", sub: "наш IBAN", detail: "Кредит кошелька клиента." },
        { id: "s", label: "Suspense", sub: "неизвестный IBAN", detail: "Ждёт allocate — правило зачисления." },
      ],
      [{ from: "or", to: "n", note: "входящий" }, { from: "n", to: "cl", note: "наш IBAN" }, { from: "n", to: "s", note: "неизвестен" }]),
  },
  {
    id: "bpmn-as-is",
    title: "Сверка as-is",
    term: "As-is / reconciliation",
    learn: "Зачем в схему as-is включают ручной шаг сверки.",
    teaser: "Ручная сверка в Excel — шаг текущего процесса. Его включают в as-is.",
    grade: "Junior",
    plan: "pro",
    minutes: 5,
    board: flow("Ночная сверка", "Повтор callback → вторая строка в CRM → сверка вручную. To-be убирает дубль ключом, не лозунгом.",
      [
        { id: "cb", label: "Callback", sub: "тот же id", detail: "Шлюз прислал сообщение повторно." },
        { id: "crm", label: "CRM", sub: "две строки", detail: "Оператор видит два успеха." },
        { id: "xl", label: "Excel", sub: "сверка", detail: "Ручной шаг as-is, пока нет ключа." },
      ],
      [{ from: "cb", to: "crm", note: "повтор" }, { from: "crm", to: "xl", note: "вручную" }]),
  },
  {
    id: "hall-intro",
    title: "Учебный API",
    term: "Идемпотентность",
    learn: "Как увидеть дефект контракта на живом /api/v1.",
    teaser: "Повтор перевода без ключа создаёт вторую проводку. Это учебный пример, не поломка сайта.",
    grade: "Intern",
    plan: "free",
    minutes: 3,
    board: flow("Разбор контракта", "Два P2P без Idempotency-Key — две дебетовые проводки отправителя.",
      [
        { id: "p", label: "Запрос", sub: "POST /transfers", detail: "Перевод 5 000 без заголовка Idempotency-Key." },
        { id: "j", label: "Журнал", sub: "проводки", detail: "Проверьте дебет, кредит и trial balance." },
        { id: "c", label: "Правило", sub: "ключ + 409", detail: "Повтор с тем же ключом не создаёт вторую операцию." },
      ],
      [{ from: "p", to: "j", note: "эффект" }, { from: "j", to: "c", note: "требование" }]),
  },
];

writeFileSync(join(root, "boards.json"), JSON.stringify(boards.map((b) => ({ ...b, ...(boardIntel[b.id] || {}) })), null, 2));
writeFileSync(join(root, "materials.json"), JSON.stringify({ topics: TOPICS }, null, 2));
console.log("lessons", lessons.length, "boards", boards.length, "materials", TOPICS.length);
