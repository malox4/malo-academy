package interview

import (
	"strings"
	"time"
	"unicode"
)

type Stage struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	IDs   []string `json:"ids"`
}

func Arc() []Stage {
	return []Stage{
		{ID: "warmup", Title: "Разминка", IDs: []string{"iv-role", "iv-clarify", "iv-elicitation", "iv-quality", "iv-req-types"}},
		{ID: "process", Title: "Процесс / SDLC", IDs: []string{"iv-sdlc", "iv-waterfall", "iv-agile", "iv-srs", "iv-brd", "iv-docs", "iv-cr-impact"}},
		{ID: "data", Title: "Данные и SQL", IDs: []string{"iv-db", "iv-er", "iv-sql-ops", "iv-sql-window", "iv-joins", "iv-constraints", "iv-sql-nosql"}},
		{ID: "product", Title: "Продукт и стейкхолдеры", IDs: []string{"iv-stakeholders", "iv-raci", "iv-raci-hold", "iv-prio", "iv-usecase", "iv-user-story", "iv-invest", "iv-usm", "iv-cjm"}},
		{ID: "case", Title: "Кейс на доске", IDs: []string{"iv-ledger", "iv-three-balances", "iv-hold", "iv-3ds", "iv-iban", "iv-iso20022", "iv-fx", "iv-api", "iv-openapi", "iv-soap", "iv-rest", "iv-soap-rest", "iv-contract"}},
		{ID: "softs", Title: "Софты", IDs: []string{"iv-conflict", "iv-hard-req", "iv-nfr", "iv-slo", "iv-trace", "iv-val-ver", "iv-bpmn", "iv-eventstorm", "iv-uml"}},
		{ID: "close", Title: "Закрытие", IDs: []string{"iv-close", "iv-arch", "iv-saga", "iv-integration", "iv-brokers", "iv-esb", "iv-baccm"}},
	}
}

func OrderedIDs() []string {
	var out []string
	for _, s := range Arc() {
		out = append(out, s.IDs...)
	}
	return out
}

func LiveIDs() []string {
	return []string{"iv-role", "iv-clarify", "iv-sdlc", "iv-sql-ops", "iv-sql-window", "iv-stakeholders", "iv-hold", "iv-3ds", "iv-ledger", "iv-conflict", "iv-close"}
}

type Question struct {
	ID         string   `json:"id"`
	Stage      string   `json:"stage"`
	StageTitle string   `json:"stageTitle"`
	Topic      string   `json:"topic"`
	Question   string   `json:"question"`
	Follow     string   `json:"follow"`
	Need       []string `json:"need"`
	Traps      []string `json:"traps"`
	Strong     string   `json:"strong"`
	Skills     []string `json:"skills"`
	Kind       string   `json:"kind"` // hard | soft | mix
	Replies    []Reply  `json:"replies"`
}

type Reply struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Good bool   `json:"good"`
	Why  string `json:"why"`
	Hit  string `json:"hit,omitempty"`
}

func LiveQuestions() []Question {
	return []Question{
		{
			ID: "iv-role", Stage: "warmup", StageTitle: "Разминка", Topic: "Роль", Kind: "soft",
			Skills:   []string{"clarifying", "clarity"},
			Question: "Кем вы хотите быть в банке — BA, SA или смесь — и как вы начнёте первый разговор, если задача звучит «сделайте как у Тинькофф»?",
			Follow:   "А если продукт говорит «нам некогда уточнять, рисуйте сразу макет»?",
			Need:     []string{"уточн", "need", "контекст", "зачем", "границ", "стейкхолдер", "как у", "чужой"},
			Traps:    []string{"сразу нарисую экраны", "скопируем банк", "это не моя зона"},
			Strong:   "Сначала роль и Need: зачем меняемся и для кого. «Как у Тинькофф» — чужой UI, не требование. Уточняю контекст, ограничения, кто Accountable.",
			Replies: []Reply{
				{ID: "s", Text: "Спрошу Need, кто страдает и какие границы. «Как у банка X» — референс, не спецификация. Роль: я снимаю неоднозначность, не копирую чужой UI.", Good: true, Why: "Разминка про Clarifying. Вы не бросились рисовать.", Hit: "need"},
				{ID: "t0", Text: "Сразу нарисую экраны как у них — так быстрее согласуем.", Good: false, Why: "Чужой UI. Без Need команда идеально закроет чужую проблему.", Hit: "ui"},
				{ID: "t1", Text: "Пусть продукт сам скажет, BA тут ни при чём.", Good: false, Why: "Аналитик как раз начинает разговор. Это не уход от роли.", Hit: "role"},
			},
		},
		{
			ID: "iv-clarify", Stage: "warmup", StageTitle: "Разминка", Topic: "Уточнения", Kind: "soft",
			Skills:   []string{"clarifying"},
			Question: "Вам сказали: «надо, чтобы платежи были надёжными». Какие три уточнения вы зададите, прежде чем писать историю?",
			Follow:   "Стейкхолдер отвечает «ну чтобы не глючило». Что дальше?",
			Need:     []string{"идемпотент", "повтор", "ключ", "двой", "списан", "окно", "метри", "кто", "need", "sla", "409"},
			Traps:    []string{"сделаем kafka", "добавим мониторинг", "надежность это nfr и всё"},
			Strong:   "Что значит надёжно: нет второго capture, повтор с ключом, какой SLA, кто видит отказ. Перевожу лозунг в проверяемый инвариант.",
			Replies: []Reply{
				{ID: "s", Text: "Что именно ломается: второй capture, таймаут шлюза, ночная сверка? Какой ключ идемпотентности? Какая метрика после релиза?", Good: true, Why: "Лозунг разобран на инвариант и измерение.", Hit: "key"},
				{ID: "t0", Text: "Поставим Kafka — так надёжнее доставка.", Good: false, Why: "Технология вместо Need. Надёжность платежа — про книгу и ключ, не про брокер.", Hit: "kafka"},
				{ID: "t1", Text: "Это NFR, пусть архитектор напишет «высокая доступность».", Good: false, Why: "Неизмеримое NFR на приёмке превращается в вкус.", Hit: "nfr"},
			},
		},
		{
			ID: "iv-sdlc", Stage: "process", StageTitle: "Процесс / SDLC", Topic: "SDLC", Kind: "hard",
			Skills:   []string{"process"},
			Question: "Где аналитик живёт в жизненном цикле платёжного контура — и какая дыра самая частая?",
			Follow:   "Вас зовут только «написать stories на спринт». Что ответите?",
			Need:     []string{"need", "релизе", "эксплуатац", "оценк", "sdlc", "после", "kpi", "сопровожд", "discovery", "приёмк"},
			Traps:    []string{"только анализ во втором квадрате", "после релиза не моё", "scrum это sdls"},
			Strong:   "От discovery до оценки после релиза. Частая дыра — исчезнуть после UAT и не проверить, закрылся ли Need. Ночные сбои — новые требования.",
			Replies: []Reply{
				{ID: "s", Text: "Discovery, спецификация среза, переход, приёмка и оценка после релиза. Исчезнуть после UAT — обрезанный цикл: не узнаем, упали ли двойные списания.", Good: true, Why: "SDLC как цикл ценности, не квадрат «анализ».", Hit: "eval"},
				{ID: "t0", Text: "Аналитик работает до разработки, дальше QA и поддержка.", Good: false, Why: "Эксплуатация кормит следующий анализ. Платёжный контур так и живёт.", Hit: "handoff"},
				{ID: "t1", Text: "У нас Scrum, SDLC не нужен.", Good: false, Why: "События Scrum не заменяют цикл Need → решение → оценка.", Hit: "scrum"},
			},
		},
		{
			ID: "iv-sql-ops", Stage: "data", StageTitle: "Данные и SQL", Topic: "SQL", Kind: "hard",
			Skills:   []string{"sql", "data"},
			Question: "Как запросом поймать подозрение на двойное списание P2P и чем WHERE отличается от HAVING?",
			Follow:   "В выборке 2 миллиона строк. Что сделаете, прежде чем UPDATE в проде?",
			Need:     []string{"group", "having", "count", "идемпотент", "where", "ключ", "дубл", "select"},
			Traps:    []string{"поправлю одним скриптом", "sql аналитику не нужен", "where после группировки"},
			Strong:   "GROUP BY ключу/сумме HAVING COUNT(*) > 1. WHERE режет до группировки, HAVING — после. UPDATE в проде — не инструмент аналитика без окна и бэкапа.",
			Replies: []Reply{
				{ID: "s", Text: "Группирую переводы без ключа по from/to/amount и смотрю HAVING COUNT(*) > 1. WHERE — фильтр строк, HAVING — фильтр групп. В прод руками не лезу.", Good: true, Why: "Это тот запрос, который есть в SQL-лабе пета.", Hit: "having"},
				{ID: "t0", Text: "Напишу UPDATE и уберу лишние строки — так быстрее клиенту вернём.", Good: false, Why: "Правка книги втихую. Сначала гипотеза SELECT, потом согласованное сторно.", Hit: "update"},
				{ID: "t1", Text: "Аналитику SQL не нужен, это DBA.", Good: false, Why: "Без выборки вы спорите мнениями. На столе BA SQL — инструмент гипотезы.", Hit: "dba"},
			},
		},
		{
			ID: "iv-stakeholders", Stage: "product", StageTitle: "Продукт и стейкхолдеры", Topic: "Стейкхолдеры", Kind: "soft",
			Skills:   []string{"stakeholders", "escalation"},
			Question: "В двойных списаниях днём кричат коллекторы. Кого ещё вы обязаны увидеть на карте — и что будет, если забудете ночную смену?",
			Follow:   "Комплаенс не пришёл на встречу. Как вовлекаете без театра?",
			Need:     []string{"ночн", "операц", "комплаенс", "бухгал", "шлюз", "клиент", "риск", "забыл", "карту"},
			Traps:    []string{"только заказчик и po", "кто громче в чате", "карту один раз"},
			Strong:   "День — коллекторы и продукт. Ночь — сверка и шлюз. Комплаенс, бухгалтерия, вендор. Забыл ночную смену — кнопка сломает их регламент.",
			Replies: []Reply{
				{ID: "s", Text: "Коллекторы, риск, ночные операции, шлюз, комплаенс, бухгалтерия, клиент. Забуду ночь — получу решение, которое днём красиво, а в 22:00 плодит дубли.", Good: true, Why: "Карта не оргструктура. Забытый стейкхолдер = сюрприз на релизе.", Hit: "night"},
				{ID: "t0", Text: "Стейкхолдеры — заказчик и PO, остальные исполнители.", Good: false, Why: "Так теряют комплаенс и ночь. Это управление риском, не вежливость.", Hit: "po"},
				{ID: "t1", Text: "Кто пишет в чат громче — того и слушаем.", Good: false, Why: "Громкость ≠ влияние. RACI и карта как раз про это.", Hit: "chat"},
			},
		},
		{
			ID: "iv-hold", Stage: "case", StageTitle: "Кейс на доске", Topic: "Холд", Kind: "hard",
			Skills:   []string{"bank"},
			Question: "Клиент на кассе видит SUCCESS, в приложении available меньше, в журнале пусто. Что произошло и чем void отличается от refund?",
			Follow:   "ACS отвалился, а контракт пишет SUCCESS в ledger. Что в требованиях?",
			Need:     []string{"холд", "hold", "auth", "capture", "void", "refund", "available", "журнал", "провод"},
			Traps:    []string{"auth и списание одно", "refund из холда", "settle на auth"},
			Strong:   "Authorize открывает холд. Capture пишет ноги. Void снимает холд. Refund — реверс capture. ACS timeout → UNKNOWN, не SUCCESS.",
			Replies: []Reply{
				{ID: "s", Text: "Это auth: холд OPEN, available уменьшается, журнала ещё нет. Void из AUTHORIZED. Refund только после capture. SUCCESS на таймауте ACS — ложное списание.", Good: true, Why: "Касса ≠ книга. Доска холда должна зажечься, журнал — нет.", Hit: "hold"},
				{ID: "t0", Text: "Для клиента auth и списание — одно и то же, можно сразу минусовать balance.", Good: false, Why: "Тогда ACS-таймаут превращается в дыру. Это Excel, не ядро.", Hit: "balance"},
				{ID: "t1", Text: "Сделаем refund из холда, чтобы быстрее вернуть.", Good: false, Why: "Из холда — void. Refund без capture плодит ноги из воздуха.", Hit: "refund"},
			},
		},
		{
			ID: "iv-ledger", Stage: "case", StageTitle: "Кейс на доске", Topic: "Журнал", Kind: "hard",
			Skills:   []string{"bank", "data"},
			Question: "P2P Анна → Борис 5000 без Idempotency-Key прошёл дважды. Что в книге, что в trial и что вы напишете в AC?",
			Follow:   "Разработчик предлагает «просто не давать кнопку жать дважды». Почему этого мало?",
			Need:     []string{"журнал", "debit", "кредит", "dr", "cr", "journal", "идемпотент", "409", "ключ", "trial"},
			Traps:    []string{"минусуем поле balance", "кнопка disabled", "trial всегда сойдётся значит ок"},
			Strong:   "Два журнала DR Анна CR Борис. Trial может сойтись — книга по суммам ровная, продукт нет. AC: повтор с ключом → 409, без ключа — 400 если контракт требует.",
			Replies: []Reply{
				{ID: "s", Text: "Две проводки: дважды DR Анна, дважды CR Борис. Trial может быть 0, а клиент потерял 10 000. AC на Idempotency-Key и 409, не на disabled кнопки.", Good: true, Why: "Учебный баг пета. Сходимость trial ≠ корректность продукта.", Hit: "journal"},
				{ID: "t0", Text: "Просто дважды минусанули balance у Анны — сверка с 1С не нужна.", Good: false, Why: "Без journalId это Excel. 1С и ночь сверки разъедутся по смыслу.", Hit: "wallet"},
				{ID: "t1", Text: "Заблокируем кнопку на фронте — дублей не будет.", Good: false, Why: "Повтор придёт с ретрая сети и с Postman. Ключ на границе API.", Hit: "ui"},
			},
		},
		{
			ID: "iv-conflict", Stage: "softs", StageTitle: "Софты", Topic: "Конфликт", Kind: "soft",
			Skills:   []string{"escalation", "stakeholders", "clarity"},
			Question: "Риск хочет меньше попыток списания, коллекторы — больше денег, продукт — «чтобы не ныли». Ваш ход?",
			Follow:   "Вас просят «реши сам, ты аналитик». Берёте политику на себя?",
			Need:     []string{"accountable", "эскал", "опци", "факт", "риск", "не сглаж", "политик", "raci"},
			Traps:    []string{"всем угодим", "решу сам как логичнее", "по-тихому в кулуарах"},
			Strong:   "Конфликт явным: факты, три опции, один Accountable. Не сглаживаю и не становлюсь владельцем политики списаний.",
			Replies: []Reply{
				{ID: "s", Text: "Кладу факты: доля дублей, recovery, жалобы. Три опции с ценой. Эскалирую на того, кто Accountable за политику. Своё — сделать спор видимым.", Good: true, Why: "Софт эскалации без цинизма. Аналитик не владелец правила «когда списывать».", Hit: "raci"},
				{ID: "t0", Text: "Сделаем так, чтобы всем было хорошо — найдём компромисс без потерь.", Good: false, Why: "Без потерь не бывает. Это отказ решать.", Hit: "peace"},
				{ID: "t1", Text: "Я аналитик, решу как логичнее и запишу в ТЗ.", Good: false, Why: "Подмена владельца политики. На релизе конфликт вернётся.", Hit: "hero"},
			},
		},
		{
			ID: "iv-close", Stage: "close", StageTitle: "Закрытие", Topic: "Закрытие", Kind: "mix",
			Skills:   []string{"clarifying", "clarity", "process"},
			Question: "Интервью заканчивается. Какие вопросы вы зададите нам и что сделаете в первую неделю на контуре холд / журнал / P2P?",
			Follow:   "Вам говорят «оффера нет, спасибо». Как закроете разговор, не сгорев?",
			Need:     []string{"need", "сверк", "журнал", "идемпотент", "стейкхолдер", "метрик", "перв", "недел", "вопрос", "инцидент"},
			Traps:    []string{"сколько платите", "я всё умею без вопросов", "перепишу ядро с нуля"},
			Strong:   "Спрашиваю, где болит книга сейчас, кто Accountable за ретраи, как выглядит ночная сверка. Первая неделя: карта стейкхолдеров, as-is холда, выборка дублей, не «перепишу ядро».",
			Replies: []Reply{
				{ID: "s", Text: "Спрошу, какой инцидент последний по двойному capture и кто закрывает ночной suspense. В первую неделю — as-is, выборка дублей, карта людей. Не рефакторинг ядра.", Good: true, Why: "Закрытие как взрослый аналитик: вопросы про контур, план без супергероя.", Hit: "week"},
				{ID: "t0", Text: "Вопросов нет, я всё внедрю с первого дня, перепишу API.", Good: false, Why: "Красный флаг. Без as-is перепись ядра — новый инцидент.", Hit: "hero"},
				{ID: "t1", Text: "Спрошу только про вилку и удалёнку.", Good: false, Why: "Нормальные бытовые вопросы — после интереса к контуру, не вместо.", Hit: "money"},
			},
		},
		{
			ID: "iv-sql-window", Stage: "data", StageTitle: "Данные и SQL", Topic: "SQL", Kind: "hard",
			Skills:   []string{"sql", "data"},
			Question: "В журнале две DR по 12 000 на один payment_id. Карточка показывает LIMIT 1. Какой запрос вы напишете — и почему не UPDATE?",
			Follow:   "PO просит «занулить вторую строку, чтобы карточка сошлась».",
			Need:     []string{"group", "ledger_legs", "having", "окно", "row_number", "update", "пет", "дубл"},
			Traps:    []string{"update второй строки", "limit 1 дубля нет", "это задача dba"},
			Strong:   "В зале: GROUP BY wallet_id, amount на ledger_legs, HAVING COUNT(*) > 1. На работе — ROW_NUMBER. UPDATE не пишу.",
			Replies: []Reply{
				{ID: "s", Text: "SELECT wallet_id, amount, COUNT(*) AS cnt FROM ledger_legs WHERE dc = 'D' GROUP BY wallet_id, amount HAVING COUNT(*) > 1. UPDATE не пишу. На работе — ROW_NUMBER по ключу.", Good: true, Why: "Живая таблица зала. Книгу не штопают.", Hit: "window"},
				{ID: "t0", Text: "UPDATE amount = 0 WHERE вторая строка — карточка сойдётся.", Good: false, Why: "Штопка журнала. Сторно, не затирание.", Hit: "update"},
				{ID: "t1", Text: "SELECT * LIMIT 1 — дубля нет, карточка права.", Good: false, Why: "Та же ложь карточки.", Hit: "limit"},
			},
		},
		{
			ID: "iv-three-balances", Stage: "case", StageTitle: "Кейс на доске", Topic: "Проводки", Kind: "hard",
			Skills:   []string{"bank", "data"},
			Question: "Назовите три числа на кошельке. Что сломается, если ночью сделать UPDATE wallets.balance?",
			Follow:   "Главбух: поле поправили, закроем смену.",
			Need:     []string{"ledger", "hold", "available", "update", "trial", "ног"},
			Traps:    []string{"одно поле balance", "update закроем смену", "available и есть книга"},
			Strong:   "ledger — ноги, hold — OPEN, available = ledger − hold. UPDATE убьёт trial.",
			Replies: []Reply{
				{ID: "s", Text: "ledger, hold, available = ledger − hold. UPDATE одного поля спрячет холд и разъедет trial.", Good: true, Why: "Три числа. Книга не поле.", Hit: "journal"},
				{ID: "t0", Text: "Одно поле balance, клиенту так проще.", Good: false, Why: "Висящий холд невидим.", Hit: "wallet"},
				{ID: "t1", Text: "Ночью UPDATE, смену закроем.", Good: false, Why: "1С ест журнал.", Hit: "trial"},
			},
		},
		{
			ID: "iv-openapi", Stage: "case", StageTitle: "Кейс на доске", Topic: "OpenAPI", Kind: "hard",
			Skills:   []string{"integrations"},
			Question: "Вам принесли скрин Swagger: «это контракт capture». Чем OpenAPI 3 отличается от картинки?",
			Follow:   "Dev: YAML сам из кода.",
			Need:     []string{"yaml", "ключ", "409", "201", "идемпотент", "скрин", "openapi"},
			Traps:    []string{"скрин это контракт", "только 200", "скопировать чужой openapi"},
			Strong:   "YAML: ключ required, 201, 409, 422. Скрин не фиксирует повтор. Сверить с /pet API.",
			Replies: []Reply{
				{ID: "s", Text: "В YAML ключ обязателен, 201 ноги, 409 повтор, 422 строка-сумма. Скрин Swagger — не артефакт.", Good: true, Why: "Контракт, не картинка.", Hit: "yaml"},
				{ID: "t0", Text: "Скрин Swagger достаточно, так быстрее согласуем.", Good: false, Why: "Нет 409 и ключа.", Hit: "screen"},
				{ID: "t1", Text: "Документация сама из кода — современнее.", Good: false, Why: "Задокументировали баг.", Hit: "codefirst"},
			},
		},
		{
			ID: "iv-iso20022", Stage: "case", StageTitle: "Кейс на доске", Topic: "Банк", Kind: "hard",
			Skills:   []string{"bank", "integrations"},
			Question: "Парсер ISO 20022 взял не ту сумму. Что в контракте разбора — и куда непонятный IBAN?",
			Follow:   "Операции: накинем похожему счёту.",
			Need:     []string{"iso", "suspense", "iban", "сумм", "ключ", "конверт", "endtoend"},
			Traps:    []string{"накинуть похожему", "iso это наша книга", "без ключа файла"},
			Strong:   "Конверт ≠ книга. Поле суммы и IBAN названы. Непонятный счёт — suspense. Повтор файла не двоит ноги.",
			Replies: []Reply{
				{ID: "s", Text: "Маппинг InstdAmt и IBAN. Непонятный счёт — suspense. Ключ файла, не вторая пара ног.", Good: true, Why: "Чужой конверт, своя книга.", Hit: "suspense"},
				{ID: "t0", Text: "Похожий IBAN — зачислим, утром поправим.", Good: false, Why: "Чужие деньги.", Hit: "guess"},
				{ID: "t1", Text: "ISO 20022 и есть наш журнал.", Good: false, Why: "Конверт, не ноги.", Hit: "book"},
			},
		},
		{
			ID: "iv-saga", Stage: "close", StageTitle: "Закрытие", Topic: "Сага", Kind: "hard",
			Skills:   []string{"bank", "integrations"},
			Question: "Capture прошёл, топик лежал. Касса бьёт ещё. Что такое компенсация саги?",
			Follow:   "Ops хочет вставить событие руками.",
			Need:     []string{"outbox", "компенс", "сага", "void", "409", "ключ", "ног"},
			Traps:    []string{"второй capture новым ключом", "событие руками после коммита", "msa сама чинит"},
			Strong:   "Outbox с ногами. Компенсация — void или реверс, не второй POST capture.",
			Replies: []Reply{
				{ID: "s", Text: "Письмо в той же транзакции, что ноги. Повтор того же ключа. Компенсация void/реверс, не новый capture.", Good: true, Why: "Сага на шве.", Hit: "outbox"},
				{ID: "t0", Text: "Второй POST capture, новый ключ — клиент ждёт.", Good: false, Why: "Вторая DR.", Hit: "repost"},
				{ID: "t1", Text: "Ops руками вставит событие, если процесс жив.", Good: false, Why: "Outbox для этого.", Hit: "hands"},
			},
		},
		{
			ID: "iv-slo", Stage: "softs", StageTitle: "Софты", Topic: "SLO", Kind: "hard",
			Skills:   []string{"process"},
			Question: "PO пишет NFR «система надёжная». Чем SLO отличается — какой SLI на capture?",
			Follow:   "Дашборд главной зелёный, касса 504.",
			Need:     []string{"slo", "sli", "409", "capture", "бюджет", "cpu"},
			Traps:    []string{"высокая доступность без цифры", "sli на cpu", "409 это ошибка индикатора"},
			Strong:   "SLI — доля capture 201/409 за 2с. 409 — успех. Бюджет стопит бейдж. Не CPU.",
			Replies: []Reply{
				{ID: "s", Text: "SLO: доля capture 201 и 409 за 2 секунды. 409 — хороший код. CPU главной не SLI шва.", Good: true, Why: "Число на шве, не лозунг.", Hit: "sli"},
				{ID: "t0", Text: "NFR «надёжно» достаточно, архитектор допишет.", Good: false, Why: "Неизмеримо.", Hit: "slogan"},
				{ID: "t1", Text: "CPU 40% — вот надёжность. 409 красим красным.", Good: false, Why: "Не шов. 409 — успех ключа.", Hit: "cpu"},
			},
		},
		{
			ID: "iv-raci-hold", Stage: "product", StageTitle: "Продукт и стейкхолдеры", Topic: "RACI", Kind: "soft",
			Skills:   []string{"stakeholders", "escalation"},
			Question: "На захвате холда все за, ног нет. PO думал A у разработки. Как разложить RACI?",
			Follow:   "Комплаенс только в копии письма.",
			Need:     []string{"accountable", "один", "raci", "комплаенс", "po", "два a"},
			Traps:    []string{"dev accountable за деньги", "два a спокойнее", "комплаенс только cc"},
			Strong:   "Один A на scope — PO. SA R на Then. Комплаенс C. Два A — брак.",
			Replies: []Reply{
				{ID: "s", Text: "Один A на захват scope — PO. Я R на Then журнала. Dev R на код. Комплаенс C, не только I.", Good: true, Why: "Буквы на шов, не на отделы.", Hit: "a"},
				{ID: "t0", Text: "Пусть dev будет A, они пишут ключ.", Good: false, Why: "Политика денег не у разработки.", Hit: "dev"},
				{ID: "t1", Text: "Два A — PO и юрист, так спокойнее.", Good: false, Why: "Никто не решает.", Hit: "two"},
			},
		},
		{
			ID: "iv-cr-impact", Stage: "process", StageTitle: "Процесс / SDLC", Topic: "CR", Kind: "hard",
			Skills:   []string{"process"},
			Question: "После Done кладут частичный capture 8 000. Комментарий или CR?",
			Follow:   "QA: старый тест на полную сумму ещё зелёный.",
			Need:     []string{"cr", "then", "частичн", "тест", "openapi", "done"},
			Traps:    []string{"комментарий в done", "только ui", "любой need это cr"},
			Strong:   "CR: меняется Then. Вместе — автотест, YAML, экран. Ложный Done без теста.",
			Replies: []Reply{
				{ID: "s", Text: "Это CR: ноги 8, холд остаток 4. Сменяю Then, тест и OpenAPI. Комментарий в Done не книга.", Good: true, Why: "Impact правила.", Hit: "cr"},
				{ID: "t0", Text: "Допишем в тот же Done-тикет комментарием.", Good: false, Why: "Правило уже «закрыто» враньём.", Hit: "done"},
				{ID: "t1", Text: "Поправим только бейдж кассы.", Good: false, Why: "Тест и контракт останутся врать.", Hit: "ui"},
			},
		},
		{
			ID: "iv-3ds", Stage: "case", StageTitle: "Кейс на доске", Topic: "3DS", Kind: "hard",
			Skills:   []string{"bank"},
			Question: "ACS не ответил 30 секунд, касса зелёная. Что в требованиях на 3DS timeout?",
			Follow:   "Касса: клиенту уже успех, проводите.",
			Need:     []string{"acs", "3ds", "unknown", "timeout", "hold", "void", "ноги"},
			Traps:    []string{"success в журнал на таймауте", "3ds это проводка", "refund из open"},
			Strong:   "Тишина ACS → UNKNOWN, не SUCCESS. 3DS не пишет ноги. Void из OPEN.",
			Replies: []Reply{
				{ID: "s", Text: "UNKNOWN, не SUCCESS. Журнал молчит. Успех ACS ещё не capture. Из OPEN — void, не refund.", Good: true, Why: "Касса ≠ книга. ACS на своей дорожке.", Hit: "hold"},
				{ID: "t0", Text: "Пишем SUCCESS, клиент уже видел зелёное.", Good: false, Why: "Ложное списание.", Hit: "success"},
				{ID: "t1", Text: "3DS прошёл — значит ноги можно писать.", Good: false, Why: "ACS не журнал.", Hit: "legs"},
			},
		},
		{
			ID: "iv-eventstorm", Stage: "softs", StageTitle: "Софты", Topic: "Event Storming", Kind: "hard",
			Skills:   []string{"process", "data"},
			Question: "Клеят оранжевый стикер wallet-api. Что должно быть на оранжевом?",
			Follow:   "Фасилитатор орёт «сначала сервисы, так быстрее».",
			Need:     []string{"событ", "оранж", "прошлом", "холд", "касса", "книга", "факт"},
			Traps:    []string{"сервис на оранжевый", "будущее как событие", "смешать success кассы и posting"},
			Strong:   "Оранжевый — факт в прошедшем: HoldOpened, CapturePosted. Сервис — не факт. Касса и книга — разный SUCCESS.",
			Replies: []Reply{
				{ID: "s", Text: "Оранжевый: холд открыт, capture записан. wallet-api — не событие. Граница касса/книга сначала.", Good: true, Why: "Прошедшее время, не коробка.", Hit: "event"},
				{ID: "t0", Text: "Оранжевый — имена сервисов, так карта быстрее.", Good: false, Why: "Потеряли границу денег.", Hit: "service"},
				{ID: "t1", Text: "Стикер «система будет надёжной» — тоже событие.", Good: false, Why: "Будущее не факт.", Hit: "future"},
			},
		},
	}
}

func QuestionByID(id string) *Question {
	for _, q := range LiveQuestions() {
		if q.ID == id {
			qq := q
			return &qq
		}
	}
	return nil
}

func ScoreText(q Question, text string) (int, string) {
	t := strings.ToLower(strings.TrimSpace(text))
	if len([]rune(t)) < 20 {
		return 0, "Коротко для стола. Назовите ноги, ключ или Need."
	}
	trapHits := 0
	for _, tr := range q.Traps {
		words := words5(strings.ToLower(tr))
		n := 0
		for _, w := range words {
			if strings.Contains(t, w) {
				n++
			}
		}
		if n >= 2 {
			trapHits++
		}
	}
	needHits := 0
	for _, n := range q.Need {
		if strings.Contains(t, strings.ToLower(n)) {
			needHits++
		}
	}
	if trapHits > 0 && needHits < 2 {
		return 0, "Близко к ловушке. " + q.Traps[0]
	}
	if needHits >= 3 {
		return 2, q.Strong
	}
	if needHits >= 1 {
		return 1, "Часть фактов есть. Сильный ответ: " + q.Strong
	}
	return 0, "Фактов мало. Опирайтесь на журнал, ключ, стейкхолдера или проверяемый Need."
}

func words5(s string) []string {
	var b strings.Builder
	var out []string
	flush := func() {
		w := b.String()
		b.Reset()
		if len([]rune(w)) >= 5 {
			out = append(out, w)
		}
	}
	for _, r := range s {
		if unicode.IsLetter(r) {
			b.WriteRune(r)
		} else {
			flush()
		}
	}
	flush()
	return out
}

func LevelOf(score, max int) (string, float64) {
	if max <= 0 {
		return "intern", 0.35
	}
	pct := float64(score) / float64(max)
	switch {
	case pct >= 0.80:
		return "senior", clamp(0.55 + pct*0.4)
	case pct >= 0.62:
		return "middle", clamp(0.5 + (pct - 0.62))
	case pct >= 0.40:
		return "junior", clamp(0.45 + (pct - 0.40))
	default:
		return "intern", clamp(0.4 + pct)
	}
}

func clamp(v float64) float64 {
	if v > 0.92 {
		return 0.92
	}
	if v < 0.35 {
		return 0.35
	}
	return v
}

const LiveMinutes = 30

func DefaultDeadline() time.Time {
	return time.Now().UTC().Add(LiveMinutes * time.Minute)
}
