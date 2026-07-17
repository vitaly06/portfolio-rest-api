# Developer Portfolio API

Бэкенд-сервис для лендинг-презентации разработчика с полноценным REST API, файловой аналитикой, защитой от спама (Rate Limiting) и интеграцией с искусственным интеллектом (Yandex GPT).

Проект разработан в рамках тестового задания с фокусом на чистоту архитектуры, слоистую структуру кода и надёжность (Graceful Fallback).

> **🚀 Живой деплой:** API развёрнут и доступен по адресу:
> - Swagger UI: **http://109.69.22.44:3003/swagger**
> - Base API URL: **http://109.69.22.44:3003**

---

## 1. Как запустить проект

### Требования

- Go 1.26.5 или выше
- Аккаунт Yandex Cloud (и IAM / API ключ для Yandex GPT)
- Настройки SMTP (для отправки email-уведомлений)

### Инструкция по запуску локально

**1. Клонируйте репозиторий:**

```bash
git clone https://github.com/vitaly06/portfolio-rest-api
cd portfolio-rest-api
```

**2. Настройте переменные окружения:**

Создайте файл `.env` в корневой директории на основе примера:

```bash
cp .env.example .env
```

Заполните `.env`:

```env
PORT=3000
YANDEX_API_KEY=your_yandex_iam_or_api_key
SMTP_HOST=smtp.beget.com
SMTP_PORT=587
SMTP_USER=your_email@example.com
SMTP_PASSWORD=your_smtp_password
OWNER_EMAIL=your_email@example.com
```

**3. Установите зависимости:**

```bash
go mod tidy
```

**4. Запустите сервер:**

```bash
go run cmd/main.go
```

Сервер запустится по адресу: `http://localhost:3000`

Swagger UI: `http://localhost:3000/swagger`

### Запуск через Docker

```bash
# Собрать и запустить
docker-compose up --build -d

# Посмотреть логи
docker-compose logs -f
```

Контейнер доступен на `http://localhost:3000`

---

## 2. Стек технологий

| Категория | Технология |
|-----------|------------|
| **Язык** | Go (Golang) — высокая производительность, горутины, минимальные аллокации |
| **Веб-фреймворк** | [Fiber v3](https://gofiber.io/) — построен на `fasthttp`, один из самых быстрых фреймворков для Go |
| **AI-провайдер** | [Yandex GPT (YandexGPT Lite)](https://cloud.yandex.ru/docs/foundation-models/) — нативный HTTP-клиент без сторонних SDK |
| **Валидация** | `github.com/go-playground/validator/v10` — индустриальный стандарт |
| **SMTP** | Стандартная библиотека Go `net/smtp` |
| **Конфигурация** | `github.com/joho/godotenv` |
| **Документация** | Swagger / OpenAPI через `swaggo/swag` + `gofiber/contrib/swaggo` |
| **Rate Limiting** | Встроенный `fiber/v3/middleware/limiter` (скользящее окно) |
| **Хранение данных** | Файловая система (`data/app.log`, `data/stats.json`) |
| **Контейнеризация** | Docker + Docker Compose |

---

## 3. Архитектура проекта

Проект построен по принципу **слоистой архитектуры (Layered Architecture)**:

```
portfolio-rest-api/
├── cmd/
│   └── main.go                  # Точка входа: DI, роутер, запуск сервера
├── internal/
│   ├── config/
│   │   └── config.go            # Загрузка .env переменных окружения
│   ├── deliviry/http/
│   │   ├── handler.go           # HTTP-хендлеры (Delivery Layer)
│   │   └── middleware.go        # Middleware: логирование каждого запроса
│   ├── domain/
│   │   └── models.go            # Доменные модели и структуры данных
│   ├── repository/
│   │   └── file_repo.go         # Работа с файловой системой (Repository Layer)
│   └── usecase/
│       ├── ai_service.go        # AI-интеграция с Yandex GPT (Service Layer)
│       └── contact_usecase.go   # Бизнес-логика обработки контакта (Use Case)
├── pkg/
│   └── mailer/
│       └── mailer.go            # Отправка email через net/smtp
├── docs/                        # Сгенерированная Swagger-документация
├── data/
│   ├── app.log                  # JSON-логи всех HTTP-запросов
│   └── stats.json               # Статистика тональности обращений
├── Dockerfile
├── docker-compose.yml
└── .env.example
```

### Паттерны проектирования

- **Dependency Injection** — все зависимости явно передаются в конструкторы (`NewHandler`, `NewContactUsecase`, `NewAIService`), что упрощает тестирование и замену компонентов.
- **Repository Pattern** — слой работы с данными (`FileRepository`) абстрагирован: при необходимости можно заменить файловое хранилище на PostgreSQL без изменения бизнес-логики.
- **Async Workers (Горутины)** — отправка email вынесена в отдельную горутину (`go func()`), клиент получает HTTP-ответ мгновенно, не ожидая работы SMTP-сервера.

### Поток запроса

```
HTTP Request
    ↓
Rate Limiter Middleware (≤ 3 req/min per IP)
    ↓
Logger Middleware (async write to data/app.log)
    ↓
Handler (валидация JSON, go-playground/validator)
    ↓
ContactUsecase.ProcessContactForm()
    ├─→ AIService.AnalyzeAndReply() → Yandex GPT API (или Fallback)
    ├─→ go mailer.SendContactEmails() [async]
    └─→ repo.UpdateMetrics()
    ↓
JSON Response (200 OK)
```

---

## 4. API: эндпоинты и примеры

Полная интерактивная документация: **http://109.69.22.44:3003/swagger**

---

### `POST /api/contact` — Форма обратной связи

**Запрос:**

```bash
curl -X POST http://109.69.22.44:3003/api/contact \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Виталий",
    "email": "user@example.com",
    "phone": "+79991112233",
    "comment": "Отличный бэкенд на Go! Хотим пригласить на собеседование."
  }'
```

**Ответ `200 OK`:**

```json
{
  "success": true,
  "sentiment": "positive",
  "ai_reply": "Здравствуйте, Виталий! Спасибо за высокую оценку. С удовольствием пройду собеседование — напишите удобное время."
}
```

**Ошибка валидации `400 Bad Request`:**

```json
{
  "error": "Key: 'ContactRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
}
```

**Rate Limit `429 Too Many Requests`** (более 3 запросов в минуту с одного IP):

```json
{
  "error": "Слишком много запросов. Пожалуйста, подождите минуту."
}
```

**Правила валидации входных данных:**

| Поле | Тип | Правило |
|------|-----|---------|
| `name` | string | обязательное, мин. 2 символа |
| `email` | string | обязательное, валидный email |
| `phone` | string | обязательное |
| `comment` | string | обязательное, мин. 5 символов |

---

### `GET /api/health` — Health Check

```bash
curl http://109.69.22.44:3003/api/health
```

```json
{
  "status": "OK"
}
```

---

### `GET /api/metrics` — Статистика обращений

```bash
curl http://109.69.22.44:3003/api/metrics
```

```json
{
  "total_requests": 12,
  "sentiment_stats": {
    "positive": 8,
    "neutral": 3,
    "negative": 1
  },
  "last_update": "2026-07-17T13:40:00Z"
}
```

---

### `GET /swagger/*` — Swagger UI

Интерактивная документация OpenAPI: `http://109.69.22.44:3003/swagger`

---

## 5. AI-интеграция и Fallback

### Используемый провайдер: Yandex GPT

Интеграция реализована через **нативный `net/http` клиент** без сторонних SDK — это минимизирует зависимости и даёт полный контроль над запросами.

- **Модель:** `yandexgpt-lite/latest`
- **Эндпоинт:** `https://llm.api.cloud.yandex.net/foundationModels/v1/completion`
- **Таймаут:** 7 секунд
- **Температура:** 0.3 (предсказуемые, деловые ответы)

### Что делает AI

AI выполняет **две задачи одновременно**:

1. **Анализ тональности (Sentiment Analysis)** — определяет, является ли комментарий `positive`, `neutral` или `negative`.
2. **Генерация персонализированного автоответа** — пишет вежливый ответ от имени разработчика.

### Промпт (System + User)

```
System:
"Ты — AI-ассистент бэкенд-разработчика Виталия. Твоя задача — проанализировать
комментарий пользователя и сгенерировать автоответ. Определи тональность
комментария (доступно только 3 варианта: 'positive', 'neutral', 'negative').
Напиши вежливый, профессиональный краткий ответ от лица Виталия.
Верни ответ СТРОГО в формате JSON без markdown-разметки, содержащий поля
'sentiment' и 'reply'.
Пример структуры: {"sentiment": "positive", "reply": "текст"}."

User:
"Комментарий для анализа: "[текст_комментария]""
```

### Механизм Graceful Fallback

Если Yandex GPT недоступен (нет ключа, сетевая ошибка, плохой статус-код, ошибка парсинга JSON) — **сервис не падает**:

1. Ошибка выводится в stdout для мониторинга (`[YANDEX ERROR] ...`).
2. Система мгновенно возвращает статичный шаблон:

```json
{
  "sentiment": "neutral",
  "reply": "Здравствуйте! Спасибо за ваше обращение. Я получил ваше сообщение и свяжусь с вами в ближайшее время для обсуждения деталей."
}
```

3. Пользователь получает корректный ответ, email всё равно отправляется.

---

## 6. Что сделано с помощью AI

При разработке ИИ использовался как инструмент повышения продуктивности:

- **Генерация бойлерплейта:** Базовые скелеты структур данных (`domain/models.go`) набрасывались нейросетью, после чего дорабатывались вручную.
- **Генерация Swagger-аннотаций:** `godoc`-комментарии к хендлерам генерировались с помощью AI и затем выверялись.

**Что пришлось исправлять вручную:**

- Нейросети регулярно возвращали JSON внутри Markdown-блоков (` ```json ... ``` `), что ломало `json.Unmarshal`. Добавлена ручная зачистка строки: `strings.TrimPrefix(rawJSON, "` ` ` `json")`.
- Интеграция изначально предлагалась через сторонние SDK с избыточными зависимостями. Реализация переписана на чистый `net/http` клиент для получения компактного Go-бинарника.
- Структура ответа Yandex GPT (`result.alternatives[0].message.text`) отличается от OpenAI, что потребовало ручного написания структуры для десериализации.

---

## 7. Хранение данных, логирование и Rate Limiting

### Логирование запросов

Реализовано через кастомную Middleware в `internal/deliviry/http/middleware.go`.

- Каждый HTTP-запрос логируется в `data/app.log` в формате JSON (одна запись — одна строка).
- Запись происходит **асинхронно** (`go func()`), не блокируя ответ.

Пример записи в `data/app.log`:

```json
{"timestamp":"2026-07-17T13:40:00+05:00","method":"POST","path":"/api/contact","ip":"91.220.xx.xx","status":200}
```

### Rate Limiting

Реализован встроенным `fiber/v3/middleware/limiter` (скользящее окно), настроен в `cmd/main.go`:

- **Максимум:** 3 запроса с одного IP в минуту на эндпоинт `POST /api/contact`
- При превышении: `429 Too Many Requests` без вызова AI и SMTP
- Ключ ограничения: IP-адрес клиента (`c.IP()`)

### Статистика (Metrics)

- Хранится в `data/stats.json` в человекочитаемом JSON-формате.
- Обновляется при каждом успешном `POST /api/contact`.
- Операции защищены `sync.Mutex` (потокобезопасный инкремент).

```json
{
  "total_requests": 12,
  "sentiment_stats": {
    "positive": 8,
    "neutral": 3,
    "negative": 1
  },
  "last_update": "2026-07-17T13:40:00.000Z"
}
```
