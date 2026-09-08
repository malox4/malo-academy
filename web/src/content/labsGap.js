const o = (text, good, why) => ({ text, good, why });

function lab(id, title, teaser, setting, blocks, gradeId = "intern") {
  return { id, gradeId, title, teaser, minutes: 9, xp: 70, setting, blocks };
}

/** One desk lab per gap topic. Intern-open unless middle architecture. */
export const GAP_LABS = [
  lab("lab-req-validate", "Проверка до спринта", "Тест, люди, how. Эмодзи — не торг.", "PO: берём WAL-4412 как есть. Kafka в заголовке.", [
    { kind: "spot", title: "SA · Можно ли брать", prompt: "Что ещё не проверено.", lines: [
      { id: "1", text: "Эмодзи в чате = утверждено", bad: true, why: "Не торг." },
      { id: "2", text: "Then: два callback, нет второй DR", bad: false, why: "Тест есть." },
      { id: "3", text: "Kafka приклеена к need", bad: true, why: "How не свободен." },
      { id: "4", text: "Журнал кивнул на одну сумму", bad: false, why: "Люди." },
    ]},
    { kind: "scene", title: "BA+SA · PO", setting: "Срок завтра.", steps: [{ from: "PO", line: "Не тормозите, все за.", options: [
      o("Ок, в спринт.", false, "Проверки нет."),
      o("Need — одна сумма. Kafka в парковку. Журнал подписывает Then.", true, "Торг."),
      o("Пусть IREB решит.", false, "Не смена."),
    ]}]},
  ]),
  lab("lab-req-lifecycle", "Состояние правила", "Черновик, согласовано, CR, ложный Done.", "Частичный capture кладут комментарием в Done-тикет.", [
    { kind: "sort", title: "SA · Куда класть", prompt: "Жизнь / CR / не жизнь.", buckets: [
      { id: "life", title: "Состояние правила" },
      { id: "cr", title: "CR" },
      { id: "no", title: "Не то" },
    ], items: [
      { id: "a", text: "Согласовали Then с журналом", bucket: "life", why: "Состояние." },
      { id: "b", text: "Частичный capture 8 000 после полной суммы", bucket: "cr", why: "Меняет Then." },
      { id: "c", text: "Закрыли Jira без теста", bucket: "no", why: "Ложный Done." },
      { id: "d", text: "Новый need NPS", bucket: "no", why: "Новая жизнь, не CR." },
      { id: "e", text: "Старый автотест надо сменить", bucket: "cr", why: "След CR." },
    ]},
  ]),
  lab("lab-req-test-analysis", "Пачка выстрелов", "Счастье уже есть. Найдите повтор, отказ, UNKNOWN.", "QA: один кейс по AC, pass.", [
    { kind: "spot", title: "SA · Какой выстрел забыли", prompt: "Брак пачки.", lines: [
      { id: "1", text: "Только счастливый capture 12 000", bad: true, why: "Нет ночи." },
      { id: "2", text: "Повтор ключа → 409, ног не две", bad: false, why: "Кейс." },
      { id: "3", text: "Timeout ACS → UNKNOWN, не SUCCESS", bad: false, why: "Ветка." },
      { id: "4", text: "Проверить что работает", bad: true, why: "Пустой кейс." },
    ]},
    { kind: "match", title: "SA · Чем смотреть", prompt: "Кейс → фонарик.", pairs: [
      { left: "Вторая DR", right: "/pet SQL окно rn=2" },
      { left: "409", right: "/pet?tab=api повтор" },
      { left: "OPEN без ног", right: "Вкладка холдов" },
    ]},
  ]),
  lab("lab-req-estimate", "Футболка без дыр", "S/M/L. Часы из буквы — брак.", "Dev: два стори. Then нет.", [
    { kind: "sort", title: "SA · Буква", prompt: "S / M / L.", buckets: [
      { id: "s", title: "S" },
      { id: "m", title: "M" },
      { id: "l", title: "L" },
    ], items: [
      { id: "a", text: "Then ясен, один шов, тест из двух выстрелов", bucket: "s", why: "Ясно." },
      { id: "b", text: "Шлюз + timeout в AC", bucket: "m", why: "Ветка." },
      { id: "c", text: "Need спорят, мастера id нет", bucket: "l", why: "Не брать." },
      { id: "d", text: "«Как у банка X» без нашей боли", bucket: "l", why: "Дыра." },
      { id: "e", text: "UI бейджа, журнал уже закрыт другим тикетом", bucket: "s", why: "Один шов UI." },
    ]},
  ]),
  lab("lab-req-gost", "Конверт или вода", "ГОСТ 34.602 — оглавление. Need внутри.", "Заказчик: ТЗ до пятницы. Команда пишет «обеспечивать».", [
    { kind: "spot", title: "SA · Что в конверт", prompt: "Вода / правило.", lines: [
      { id: "1", text: "Система должна обеспечивать удобство", bad: true, why: "Вода." },
      { id: "2", text: "Повтор ключа не пишет вторую DR", bad: false, why: "Then." },
      { id: "3", text: "Скопировать чужое ТЗ оптом", bad: true, why: "Чужая боль." },
      { id: "4", text: "Раздел надёжности: ключ обязателен", bad: false, why: "NFR." },
    ]},
  ]),
  lab("lab-db-layers", "Три слоя холда", "Концепт / логика / физика. Не сразу CREATE.", "PO: статус холда в JSON у кошелька.", [
    { kind: "sort", title: "SA · Слой", prompt: "Куда фраза.", buckets: [
      { id: "c", title: "Концепт" },
      { id: "l", title: "Логика" },
      { id: "p", title: "Физика" },
    ], items: [
      { id: "a", text: "Холд — не проводка", bucket: "c", why: "Смысл." },
      { id: "b", text: "Сущность Hold, ключ id, статус OPEN", bucket: "l", why: "ER." },
      { id: "c", text: "numeric(12,2) в Postgres", bucket: "p", why: "Диск." },
      { id: "d", text: "JSON в wallets «пока так»", bucket: "p", why: "Физика-мина." },
      { id: "e", text: "Журнал отдельно от холда", bucket: "l", why: "Сущности." },
    ]},
  ]),
  lab("lab-db-window", "Вторая нога в окне", "ROW_NUMBER на работе. В зале — GROUP BY на ledger_legs.", "Две DR 12 000. Карточка LIMIT 1.", [
    { kind: "spot", title: "SA · Запрос", prompt: "Что найдёт дубль.", lines: [
      { id: "1", text: "На работе: ROW_NUMBER по ключу, rn=2", bad: false, why: "Окно оставляет строки." },
      { id: "2", text: "UPDATE amount SET 0 WHERE rn=2", bad: true, why: "Штопка." },
      { id: "3", text: "SELECT * LIMIT 1", bad: true, why: "Та же ложь карточки." },
      { id: "4", text: "Смотреть на /pet?tab=sql", bad: false, why: "Пет." },
    ]},
    { kind: "sql", title: "SA · Какой SELECT", prompt: "Карточка LIMIT 1 врёт. Что запустите в /pet?tab=sql?", options: [
      o("SELECT payment_id FROM journal_entries; потом UPDATE.", false, "Таблицы journal_entries в зале нет."),
      o("SELECT wallet_id, amount, COUNT(*) AS cnt FROM ledger_legs WHERE dc = 'D' GROUP BY wallet_id, amount HAVING COUNT(*) > 1;", true, "Живые таблицы зала. GROUP BY, не выдуманное окно."),
      o("SELECT * FROM ledger_legs LIMIT 1;", false, "Та же ложь карточки."),
    ]},
  ]),
  lab("lab-db-cte", "Имя шага", "WITH dupes, не q1. Зерно не размножить.", "Матрёшка из четырёх SELECT. Валюта съехала.", [
    { kind: "order", title: "SA · Порядок WITH", prompt: "Сверху вниз.", items: [
      { id: "1", text: "dupes: ключи с count > 1", pos: 1 },
      { id: "2", text: "legs: ноги этих ключей", pos: 2 },
      { id: "3", text: "фильтр UZS именованным шагом", pos: 3 },
      { id: "4", text: "SELECT к кошельку, проверить зерно", pos: 4 },
    ]},
  ]),
  lab("lab-db-migrate", "Не сразу NOT NULL", "Расширить → заполнить → unique.", "Dev: UNIQUE payment_id сегодня. В журнале уже дубли.", [
    { kind: "order", title: "SA · Порядок миграции", prompt: "Живой журнал.", items: [
      { id: "1", text: "Найти дубли окном", pos: 1 },
      { id: "2", text: "Решить живую ногу по Then", pos: 2 },
      { id: "3", text: "Уникальность", pos: 3 },
      { id: "4", text: "Код больше не пишет вторую", pos: 4 },
    ]},
    { kind: "spot", title: "SA · Мина", prompt: "Что нельзя.", lines: [
      { id: "1", text: "NOT NULL сразу по грязным", bad: true, why: "Накат падает." },
      { id: "2", text: "Откат индекса продуман", bad: false, why: "CR." },
      { id: "3", text: "ALTER в том же PR, что цвет кнопки", bad: true, why: "Разный откат." },
    ]},
  ]),
  lab("lab-db-dwh", "Зерно витрины", "Холд не оборот. Строка = нога.", "Витрина «оборот» = 24 000 при одном capture 12 000.", [
    { kind: "sort", title: "SA · В витрину?", prompt: "Факт / не факт оборота.", buckets: [
      { id: "yes", title: "Строка склада" },
      { id: "no", title: "Не зерно оборота" },
    ], items: [
      { id: "a", text: "Нога журнала после capture", bucket: "yes", why: "Факт." },
      { id: "b", text: "OPEN hold 12 000", bucket: "no", why: "Резерв." },
      { id: "c", text: "Клик кассы SUCCESS на auth", bucket: "no", why: "Экран." },
      { id: "d", text: "DR/CR пары на 12 000", bucket: "yes", why: "Книга." },
    ]},
  ]),
  lab("lab-api-openapi", "YAML capture", "Ключ, 201, 409, 422. Не скрин Swagger.", "Принесли картинку UI. 200 на всё.", [
    { kind: "spot", title: "SA · Артефакт", prompt: "Что обязательно в YAML.", lines: [
      { id: "1", text: "Idempotency-Key required", bad: false, why: "Ключ." },
      { id: "2", text: "Только 200 default", bad: true, why: "Ретраи слепые." },
      { id: "3", text: "409 повтор ключа", bad: false, why: "Контракт." },
      { id: "4", text: "Скопировать чужой OpenAPI банка", bad: true, why: "Чужой холд." },
    ]},
    { kind: "yaml", title: "SA · YAML capture", prompt: "Отметьте брак в контракте. Живые строки не трогайте.", snippet: "paths:\n  /holds/{id}/capture:\n    post:\n      parameters:\n        - name: Idempotency-Key\n          required: true\n      responses:\n        '201': { description: ноги журнала }\n        '409': { description: повтор ключа }\n        '422': { description: сумма не число }", lines: [
      { id: "1", text: "Idempotency-Key required: true", bad: false, why: "Ключ." },
      { id: "2", text: "Только 200 default на всё", bad: true, why: "Ретраи слепые." },
      { id: "3", text: "409 повтор того же ключа", bad: false, why: "Контракт." },
      { id: "4", text: "Скрин Swagger вместо YAML", bad: true, why: "Картинка не артефакт." },
    ]},
    { kind: "scene", title: "BA+SA · Dev", setting: "Сгенерили YAML из кода без Then.", steps: [{ from: "Dev", line: "Документация сама из кода.", options: [
      o("Ок, так современнее.", false, "Задокументировали баг."),
      o("Сначала Then, потом файл. Сверить с /pet?tab=api.", true, "Артефакт."),
      o("Пусть QA в Postman сами.", false, "Нет контракта."),
    ]}]},
  ]),
  lab("lab-api-postman", "Пачка трёх выстрелов", "Auth, capture, повтор. /pet?tab=api.", "Онбординг: PDF со скрином 200.", [
    { kind: "order", title: "SA · Порядок пачки", prompt: "Живой контракт.", items: [
      { id: "1", text: "Auth → OPEN, ног нет", pos: 1 },
      { id: "2", text: "Capture с ключом → 201", pos: 2 },
      { id: "3", text: "Повтор того же ключа → 409", pos: 3 },
    ]},
  ]),
  lab("lab-api-asyncapi", "Повтор доставки", "Топик at-least-once. Дедуп по payment_id.", "Архитектор: события сами идемпотентны.", [
    { kind: "spot", title: "SA · Канал", prompt: "Что обязательно.", lines: [
      { id: "1", text: "Ключ payment_id в payload", bad: false, why: "Дедуп." },
      { id: "2", text: "Повтор доставки пишет вторую DR", bad: true, why: "Нет Then." },
      { id: "3", text: "PII в топике «для карточки»", bad: true, why: "Утечка." },
      { id: "4", text: "DLQ на яд схемы", bad: false, why: "Router." },
    ]},
  ]),
  lab("lab-api-graphql", "Деньги не в граф", "Query карточки ок. Capture — REST+409.", "Фронт: один GraphQL на всё, и оплату.", [
    { kind: "sort", title: "SA · Куда", prompt: "Граф / REST / нельзя.", buckets: [
      { id: "g", title: "Query" },
      { id: "r", title: "REST шов" },
      { id: "no", title: "Нельзя" },
    ], items: [
      { id: "a", text: "Поля карточки: available, OPEN, ноги", bucket: "g", why: "Чтение." },
      { id: "b", text: "Capture 12 000", bucket: "r", why: "409." },
      { id: "c", text: "Мутация без ключа", bucket: "no", why: "Дубль." },
      { id: "d", text: "Частичный data+errors на деньги", bucket: "no", why: "Враньё UI." },
    ]},
  ]),
  lab("lab-api-schema", "Забор JSON", "amount integer. Строка-сумма — 422.", "Callback: \"12000UZS\". Проглотили.", [
    { kind: "spot", title: "SA · Схема", prompt: "Что режет забор.", lines: [
      { id: "1", text: "amount не число → 422", bad: false, why: "Тип." },
      { id: "2", text: "Всё optional для гибкости", bad: true, why: "Нет забора." },
      { id: "3", text: "currency не из enum → 422", bad: false, why: "Валюта." },
      { id: "4", text: "Яд сразу в журнал", bad: true, why: "Нет DLQ." },
    ]},
  ]),
  lab("lab-api-pagination", "Курсор журнала", "Offset съезжает. SUCCESS ≠ холд.", "Выгрузка page=2 потеряла ногу на 12 000.", [
    { kind: "match", title: "SA · Нарезка", prompt: "Когда что.", pairs: [
      { left: "Живой журнал сверки", right: "Курсор after id" },
      { left: "Справочник валют", right: "Offset можно" },
      { left: "Фильтр success", right: "Не включать OPEN hold" },
    ]},
  ]),
  lab("lab-arch-styles", "Где две ноги", "Монолит книги vs MSA без саги.", "Вынесем wallet-api со своей БД.", [
    { kind: "spot", title: "SA · Стиль", prompt: "Что ломает trial.", lines: [
      { id: "1", text: "Ноги в двух БД без саги", bad: true, why: "Дыра." },
      { id: "2", text: "Учебный пет — монолит журнала", bad: false, why: "Видно книгу." },
      { id: "3", text: "MSA потому что LinkedIn", bad: true, why: "Мода." },
      { id: "4", text: "EDA + дедуп события", bad: false, why: "Then повтора." },
    ]},
  ], "middle"),
  lab("lab-arch-eip", "Выкройка на стрелке", "Pub-sub, translator, router. Не логотип Kafka.", "Три HTTP callback, три маппинга, три дубля.", [
    { kind: "match", title: "SA · EIP", prompt: "Имя выкройки.", pairs: [
      { left: "Один captured, карточка и сверка читают", right: "Pub-sub" },
      { left: "txnId → payment_id", right: "Translator" },
      { left: "Яд схемы → DLQ", right: "Router" },
    ]},
  ], "middle"),
  lab("lab-arch-saga", "Outbox не воздух", "Ноги и письмо вместе. Компенсация ≠ второй POST.", "Capture прошёл, топик лежал. Касса бьёт ещё.", [
    { kind: "spot", title: "SA · Сага", prompt: "Что правда.", lines: [
      { id: "1", text: "Outbox в той же транзакции, что ноги", bad: false, why: "Письмо." },
      { id: "2", text: "Компенсация = повторить capture", bad: true, why: "Вторая DR." },
      { id: "3", text: "Состояние саги видно ops", bad: false, why: "Где застряли." },
      { id: "4", text: "Событие руками после коммита", bad: true, why: "Процесс умер." },
    ]},
    { kind: "scene", title: "BA+SA · Касса", setting: "Топик лежал. Касса шлёт capture ещё.", steps: [{ from: "Касса", line: "Повторите списание, клиент ждёт.", options: [
      o("Второй POST capture, новый ключ.", false, "Вторая DR."),
      o("Идемпотентный повтор того же ключа. Компенсация — void или реверс ног, не новый capture.", true, "Сага."),
      o("Пусть ops руками вставит событие.", false, "Outbox умер."),
    ]}]},
  ], "middle"),
  lab("lab-arch-eventstorm", "Факт, не сервис", "Прошедшее время. Граница касса/книга.", "Клеят Wallet-api с первой минуты.", [
    { kind: "sort", title: "SA · Стикер", prompt: "Факт / экран / сервис.", buckets: [
      { id: "f", title: "Факт" },
      { id: "u", title: "Не факт" },
    ], items: [
      { id: "a", text: "Холд открыт", bucket: "f", why: "Случилось." },
      { id: "b", text: "wallet-api", bucket: "u", why: "Коробка." },
      { id: "c", text: "Capture записан", bucket: "f", why: "Книга." },
      { id: "d", text: "Система будет надёжной", bucket: "u", why: "Будущее." },
      { id: "e", text: "Callback повторён", bucket: "f", why: "Доставка." },
    ]},
  ]),
  lab("lab-arch-arc42", "ADR-12", "Решение с датой. Не 12 пустых разделов.", "Miro стёрли. Почему журнал не режем — никто не помнит.", [
    { kind: "spot", title: "SA · ADR", prompt: "Что класть.", lines: [
      { id: "1", text: "Контекст: trial в ту же секунду", bad: false, why: "Зачем." },
      { id: "2", text: "Все разделы arc42 пустые «на вырост»", bad: true, why: "Павлин." },
      { id: "3", text: "Отклонённый вариант: wallet-api своя БД", bad: false, why: "След." },
      { id: "4", text: "Цитата TOGAF вместо шва", bad: true, why: "Рамка." },
    ]},
  ], "middle"),
  lab("lab-proc-bpmn-collab", "Дорожка ACS", "Таймер, сообщение, void не refund.", "Одна дорожка «банк». 30% тишина ACS.", [
    { kind: "spot", title: "SA · BPMN 2.0", prompt: "Что добавить.", lines: [
      { id: "1", text: "ACS отдельной дорожкой", bad: false, why: "Collaboration." },
      { id: "2", text: "Таймер без Then", bad: true, why: "Касса сама." },
      { id: "3", text: "Компенсация void из OPEN", bad: false, why: "Откат." },
      { id: "4", text: "Refund из OPEN как компенсация", bad: true, why: "Дыра." },
    ]},
    { kind: "match", title: "SA · Сообщение", prompt: "Кто с кем.", pairs: [
      { left: "Таймер 30 с на ACS", right: "UNKNOWN, не SUCCESS" },
      { left: "Сообщение captured", right: "Дорожка журнала" },
      { left: "Refund из OPEN", right: "Запрещённый переход" },
    ]},
  ]),
  lab("lab-proc-uml-hold", "Запретный переход", "OPEN → void, не refund. SUCCESS не статус холда.", "Кнопка «Вернуть» на любом экране.", [
    { kind: "match", title: "SA · State", prompt: "Переход.", pairs: [
      { left: "OPEN + capture", right: "CLOSED + ноги" },
      { left: "OPEN + отмена", right: "VOID, журнал тих" },
      { left: "OPEN + refund", right: "Запрещено" },
    ]},
  ]),
  lab("lab-proc-figma", "Кадр UNKNOWN", "Handoff без состояния — угадайка.", "Три зелёных экрана. Спиннер = успех.", [
    { kind: "spot", title: "SA · Дыра макета", prompt: "Чего нет.", lines: [
      { id: "1", text: "Кадр UNKNOWN с текстом «не знаем»", bad: false, why: "Нужен." },
      { id: "2", text: "409 как «попробуйте снова»", bad: true, why: "Второй ключ." },
      { id: "3", text: "OPEN нарисован как оплачено", bad: true, why: "Ложь." },
      { id: "4", text: "Дыра в AC, не только коммент Figma", bad: false, why: "Тикет." },
    ]},
  ]),
  lab("lab-proc-mobile-web", "Дельта лифта", "Тот же Then. Нет авто-повтора без ключа.", "Курьер: приложение само повторило capture.", [
    { kind: "spot", title: "SA · Платформа", prompt: "Общее / дельта / брак.", lines: [
      { id: "1", text: "Журнал тот же, что веб", bad: false, why: "Need." },
      { id: "2", text: "Фоновый retry новым ключом", bad: true, why: "Две DR." },
      { id: "3", text: "Экран UNKNOWN + тот же ключ когда сеть", bad: false, why: "Дельта." },
      { id: "4", text: "Оптимистичный SUCCESS до 201", bad: true, why: "Ложь." },
    ]},
  ]),
  lab("lab-nfr-slo", "SLI на capture", "Не CPU. 409 — успех индикатора. Бюджет стопит бейдж.", "Дашборд главной зелёный. Касса 504.", [
    { kind: "sort", title: "SA · Куда", prompt: "SLI / не SLI / действие.", buckets: [
      { id: "s", title: "SLI шва" },
      { id: "n", title: "Не то" },
      { id: "b", title: "Бюджет" },
    ], items: [
      { id: "a", text: "Доля capture с 201/409 за 2с", bucket: "s", why: "Шов." },
      { id: "b", text: "CPU 40%", bucket: "n", why: "Железо." },
      { id: "c", text: "Съели бюджет — стоп покраски CRM", bucket: "b", why: "Стоп." },
      { id: "d", text: "Таймаут как «почти ок»", bucket: "n", why: "Враньё." },
    ]},
    { kind: "scene", title: "BA+SA · PO", setting: "Дашборд главной зелёный.", steps: [{ from: "PO", line: "NFR: система надёжная. Зачем вам SLO?", options: [
      o("Лозунг достаточно, архитектор допишет SLA.", false, "Неизмеримо."),
      o("SLO: доля capture 201/409 за 2с. 409 — успех индикатора. Бюджет стопит бейдж.", true, "Число."),
      o("CPU ниже 40% — вот надёжность.", false, "Не шов."),
    ]}]},
  ]),
  lab("lab-nfr-security", "STRIDE на топик", "PAN не в событии. Ключ ≠ личность.", "Payload captured «удобный»: ФИО и карта.", [
    { kind: "spot", title: "SA · Угроза I", prompt: "Что вычеркнуть.", lines: [
      { id: "1", text: "PAN в топике", bad: true, why: "Disclosure." },
      { id: "2", text: "payment_id и сумма", bad: false, why: "Можно." },
      { id: "3", text: "JWT есть — значит безопасно", bad: true, why: "Лог жив." },
      { id: "4", text: "Тест: grep PAN в логе capture", bad: false, why: "Then." },
    ]},
  ]),
  lab("lab-nfr-pci", "Карта не в Jira", "PAN/CVV не в лог, тикет, DWH.", "Скрин кассы с полной картой в тикете.", [
    { kind: "sort", title: "SA · Можно писать?", prompt: "Можно / нельзя.", buckets: [
      { id: "ok", title: "Можно" },
      { id: "no", title: "Нельзя" },
    ], items: [
      { id: "a", text: "payment_id", bucket: "ok", why: "Id." },
      { id: "b", text: "Полный PAN", bucket: "no", why: "PCI-взгляд." },
      { id: "c", text: "CVV «для повтора»", bucket: "no", why: "Повтор — ключ." },
      { id: "d", text: "Последние 4 на экране", bucket: "ok", why: "Маска." },
      { id: "e", text: "Сырой body в access-логе", bucket: "no", why: "Утечка." },
    ]},
  ]),
  lab("lab-nfr-wcag", "Не только цвет", "Tab до Capture. 409 текстом. Статус словом.", "Зелёный бейдж. Дальтоник прочитал auth как списание.", [
    { kind: "spot", title: "SA · NFR экрана", prompt: "Критерий живой?", lines: [
      { id: "1", text: "Статус OPEN словом, не только серый", bad: false, why: "Не цвет." },
      { id: "2", text: "Кнопка-иконка без имени", bad: true, why: "Скринридер." },
      { id: "3", text: "Путь клавиатурой auth→capture", bad: false, why: "Tab." },
      { id: "4", text: "Lighthouse 90 вместо AC", bad: true, why: "Балл." },
    ]},
  ]),
  lab("lab-bank-iso20022", "Не угадать IBAN", "Конверт ≠ книга. Suspense. Ключ файла.", "Парсер новой схемы взял не ту сумму. «Накинем похожему».", [
    { kind: "spot", title: "SA · Вход", prompt: "Правило.", lines: [
      { id: "1", text: "Непонятный счёт → suspense", bad: false, why: "Не угадать." },
      { id: "2", text: "SEPA-логика на UZS как есть", bad: true, why: "Чужой конверт." },
      { id: "3", text: "Повтор файла не двоит ноги", bad: false, why: "Ключ." },
      { id: "4", text: "Recall как новый кредит", bad: true, why: "Тип." },
    ]},
    { kind: "match", title: "SA · Поле конверта", prompt: "Куда маппить.", pairs: [
      { left: "InstdAmt / IntrBkSttlmAmt", right: "Сумма ног, не «похожее»" },
      { left: "CdtrAcct IBAN", right: "Клиент или suspense" },
      { left: "EndToEndId", right: "Ключ файла / дедуп" },
    ]},
  ]),
  lab("lab-bank-cp-cnp", "3DS не ноги", "CP vs CNP. ACS тишина → UNKNOWN.", "Сайт поставили поток терминала. Или ACS 504, касса зелёная.", [
    { kind: "match", title: "SA · Ветка", prompt: "Что значит.", pairs: [
      { left: "Успех 3DS", right: "Ещё не журнал" },
      { left: "Тишина ACS", right: "UNKNOWN, не capture" },
      { left: "CP обрыв терминала", right: "Тот же ключ auth" },
    ]},
  ]),
  lab("lab-bank-aml-kyc", "Мастер личности", "Не ACTIVE из кошелька. Порог числом. След не в Slack.", "Мобильный рисует ACTIVE на PENDING.", [
    { kind: "spot", title: "SA · Впуск", prompt: "Требование.", lines: [
      { id: "1", text: "Кошелёк не ACTIVE, пока KYC не APPROVED", bad: false, why: "Мастер." },
      { id: "2", text: "Подозрительно — без порога", bad: true, why: "Лозунг." },
      { id: "3", text: "Повтор /kyc/start с ключом", bad: false, why: "Идемпотентность." },
      { id: "4", text: "Пустить — в личке юриста", bad: true, why: "Нет следа." },
    ]},
  ]),
  lab("lab-bank-ledger-deep", "Три числа, не баланс", "UPDATE запрещён. Частичный — остаток холда. /pet.", "Накиньте 12 000. Trial потом. Или два частичных в минус.", [
    { kind: "sort", title: "SA · Какое число", prompt: "Кто врёт при UPDATE.", buckets: [
      { id: "bad", title: "Ломается" },
      { id: "rule", title: "Правило" },
    ], items: [
      { id: "a", text: "Ledger при штопке available", bucket: "bad", why: "Книга." },
      { id: "b", text: "available = ledger − OPEN", bucket: "rule", why: "Формула." },
      { id: "c", text: "UPDATE wallets.balance", bucket: "bad", why: "Запрет." },
      { id: "d", text: "Частичный: ноги 8, холд остаток 4", bucket: "rule", why: "Then." },
    ]},
    { kind: "scene", title: "BA+SA · Главбух", setting: "Diff 12 000 после «накинуть».", steps: [{ from: "Главбух", line: "Поле поправили, закроем смену.", options: [
      o("Ок, available совпал.", false, "Trial умрёт утром."),
      o("Сторно или реверс ног. Три числа с пета.", true, "Книга."),
      o("Это 1С.", false, "1С ест ваш журнал."),
    ]}]},
  ]),
  lab("lab-bank-recon", "Не клеить по 12 000", "Ключ, потом сумма. Хвост не накидывать.", "Два платежа по 12 000. Excel сошёлся LIMIT 1.", [
    { kind: "order", title: "SA · Матчинг", prompt: "Порядок.", items: [
      { id: "1", text: "Сопоставить ключ (e2e / payment_id)", pos: 1 },
      { id: "2", text: "Проверить сумму и валюту", pos: 2 },
      { id: "3", text: "Список unmatched с причиной", pos: 3 },
      { id: "4", text: "Не UPDATE, чтобы ноль", pos: 4 },
    ]},
  ]),
  lab("lab-car-portfolio", "Папка на собес", "Шов, не курсы. Без PAN и NDA.", "Кандидат: список вебинаров. Или скрин прод-Confluence.", [
    { kind: "spot", title: "SA · В папку?", prompt: "Кладём / нельзя.", lines: [
      { id: "1", text: "YAML capture учебного зала", bad: false, why: "Артефакт." },
      { id: "2", text: "Прод-лог с картой", bad: true, why: "Секрет." },
      { id: "3", text: "SQL rn=2 с /pet?tab=sql", bad: false, why: "Фонарик." },
      { id: "4", text: "Оглавление BABOK", bad: true, why: "Нет шва." },
    ]},
  ]),
  lab("lab-car-cert", "Буквы после шва", "IREB/IIBA если просят. Не зубрить вместо available.", "Вакансия CBAP. На capture — тишина.", [
    { kind: "scene", title: "BA+SA · Себе", setting: "Месяц до собеса.", steps: [{ from: "Вы", line: "Зубрить определения или пет?", options: [
      o("400 терминов IIBA, холд потом.", false, "На кейсе пусто."),
      o("Сначала три числа и 409. Сертификат — если вакансия must.", true, "Шов."),
      o("Купить дамп экзамена.", false, "Этика."),
    ]}]},
  ]),
  lab("lab-soft-raci", "Один A", "Два A — брак. Комплаенс не только I.", "Все за. Ноги нет. PO думал A у dev.", [
    { kind: "match", title: "SA · Буква", prompt: "WAL-4412 захват scope.", pairs: [
      { left: "PO на «берём в спринт»", right: "A" },
      { left: "SA на Then журнала", right: "R" },
      { left: "Два A на один шов", right: "Брак" },
    ]},
    { kind: "scene", title: "BA+SA · PO", setting: "Комплаенс только в копии письма.", steps: [{ from: "PO", line: "Dev Accountable за ключ, мы все согласны.", options: [
      o("Ок, два A спокойнее.", false, "Никто не решает."),
      o("Один A на захват scope — PO. Dev R на код. Комплаенс C, не только I.", true, "RACI."),
      o("Пусть юрист будет A «на всякий».", false, "Чужой A."),
    ]}]},
  ]),
  lab("lab-soft-workshop", "Парковка и скрипт", "90 минут. How в парковку. Не @channel.", "40 минут про бейдж. PO орёт срок.", [
    { kind: "spot", title: "SA · Фасилитация", prompt: "Ход.", lines: [
      { id: "1", text: "Need вслух к минуте 40", bad: false, why: "Таймер." },
      { id: "2", text: "Kafka с первой минуты", bad: true, why: "How." },
      { id: "3", text: "Эскалация: один A, одна фраза, срок ответа", bad: false, why: "Скрипт." },
      { id: "4", text: "@channel ночью", bad: true, why: "Шум." },
    ]},
  ]),
];

const JUNIOR_LABS = new Set([
  "lab-req-estimate",
  "lab-req-gost",
  "lab-db-window",
  "lab-db-cte",
  "lab-db-migrate",
  "lab-db-dwh",
  "lab-api-openapi",
  "lab-api-asyncapi",
  "lab-api-graphql",
  "lab-api-schema",
  "lab-arch-eventstorm",
  "lab-proc-bpmn-collab",
  "lab-nfr-slo",
  "lab-nfr-security",
  "lab-nfr-pci",
  "lab-bank-iso20022",
  "lab-bank-cp-cnp",
  "lab-bank-recon",
  "lab-soft-workshop",
]);

for (const l of GAP_LABS) {
  if (JUNIOR_LABS.has(l.id)) l.gradeId = "junior";
}
