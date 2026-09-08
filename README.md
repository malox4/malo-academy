# Analyst Hall

Зал самообучения **BA / SA**. Intern бесплатно. PRO — когда нужно **писать** в учебный банк (`POST /api/v1`), открыть junior+ и полный SQL. Новым — **3 дня PRO**. Не Zoom-школа, сертификата нет.

Гость видит лендинг на [http://localhost:8080](http://localhost:8080). После входа — [зал](http://localhost:8080/hall). Сравнение: [http://localhost:8080/pricing](http://localhost:8080/pricing).

| | |
|---|---|
| Главная | http://localhost:8080 |
| Зал | http://localhost:8080/hall |
| Материалы | http://localhost:8080/materials |
| Практика | http://localhost:8080/practice |
| Собес | http://localhost:8080/interview · http://localhost:8080/live |
| Пет | http://localhost:8080/pet · `/api/v1` |
| Health | http://localhost:8080/api/health |

Сиды админа и PRO задаются только переменными `ADMIN_EMAIL` / `ADMIN_PASSWORD` и `PRO_EMAIL` / `PRO_PASSWORD`. На сайте и на экране входа пароли не показываем.

## Docker

Нужны Docker Desktop и свободный порт **8080**. Образ собирает Vue (включая фото лендинга) и Go. Postgres в той же compose-сети, порт 5432 наружу не открыт.

```bash
git clone https://github.com/malox4/analyst-hall.git
cd analyst-hall
cp .env.example .env
```

В `.env` поставьте свои `ADMIN_PASSWORD`, `PRO_PASSWORD`, `POSTGRES_PASSWORD`. На публичном HTTPS: `COOKIE_SECURE=1`.

```bash
docker compose up --build -d
curl http://127.0.0.1:8080/api/health
```

Откройте [http://localhost:8080](http://localhost:8080). Логи: `docker compose logs -f academy`. Стоп: `docker compose down`. `down -v` сотрёт учебный банк и Postgres.

Образ: Node 22 → `web/dist`, Go 1.23 → бинарь, Debian slim слушает `:8080`. `.env` в слой образа не копируется.

## Без Docker

Нужны **Go 1.23+** и **Node 22**.

```bash
cd web && npm install && npm run build && cd ..
go run ./cmd/academy
```

Без `DATABASE_URL` учётки пишутся в `data/academy-users.json`.

## Пет

База: `http://localhost:8080/api/v1`

Дважды P2P без ключа идемпотентности — две DR-ноги Анны. Это учебный баг, не поломка.
