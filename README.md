# Планировщик задач

Веб-приложение для управления задачами с поддержкой повторений. Написано на Go, база данных — SQLite, фронтенд уже готов.

## Описание

REST API сервер с готовым фронтендом. Позволяет добавлять задачи с датой и правилом повторения, отмечать как выполненные, редактировать и удалять. Повторяющиеся задачи автоматически переносятся на следующую дату.

## Выполненные задания со звёздочкой

- `TODO_PORT` — порт сервера через переменную окружения
- `TODO_DBFILE` — путь к файлу БД через переменную окружения
- Правила повторения `w` (дни недели) и `m` (числа месяца)
- Поиск задач по тексту и по дате (`/api/tasks?search=`)
- Аутентификация через JWT (`TODO_PASSWORD`)
- Docker-образ и docker-compose

## Запуск локально

```bash
git clone https://github.com/JakimPit/go_final_project.git
cd go_final_project

go run main.go
```

Открыть в браузере: http://localhost:7540

### Переменные окружения

| Переменная | По умолчанию | Описание |
|------------|-------------|----------|
| `TODO_PORT` | `7540` | Порт сервера |
| `TODO_DBFILE` | `scheduler.db` | Путь к файлу базы данных |
| `TODO_PASSWORD` | _(не задан)_ | Пароль для входа (если не задан — авторизация отключена) |

Пример запуска с параметрами:
```bash
TODO_PORT=8080 TODO_PASSWORD=secret go run main.go
```

## Запуск тестов

Перед запуском тестов сервер должен быть запущен.

```bash
# В одном терминале
go run main.go

# В другом — запустить все тесты
go test ./tests/

# Или по отдельности
go test -run ^TestApp$ ./tests/
go test -run ^TestDB$ ./tests/
go test -run ^TestNextDate$ ./tests/
go test -run ^TestAddTask$ ./tests/
go test -run ^TestTasks$ ./tests/
go test -run ^TestTask$ ./tests/
go test -run ^TestEditTask$ ./tests/
go test -run ^TestDone$ ./tests/
go test -run ^TestDelTask$ ./tests/
```

### Параметры в `tests/settings.go`

| Переменная | Значение | Описание |
|------------|---------|----------|
| `Port` | `7540` | Порт сервера |
| `DBFile` | `../scheduler.db` | Путь к файлу БД |
| `FullNextDate` | `true` | Включает тест правил `w` и `m` |
| `Search` | `true` | Включает тест поиска задач |
| `Token` | _(пусто)_ | JWT-токен для тестов с авторизацией |

Если задан `TODO_PASSWORD`, перед запуском тестов нужно получить токен:
```bash
curl -X POST http://localhost:7540/api/signin \
  -H "Content-Type: application/json" \
  -d '{"password":"yourpassword"}'
```
Полученный токен вставить в `tests/settings.go` → `Token`.

## Сборка и запуск через Docker

```bash
# Собрать образ и запустить
docker-compose up --build

# С паролем
TODO_PASSWORD=secret docker-compose up --build
```

База данных сохраняется в папке `./data` на хосте — данные не теряются при перезапуске контейнера.

Запуск без docker-compose:
```bash
docker build -t todo-app .
docker run -p 7540:7540 -v $(pwd)/data:/data todo-app
```

## Структура проекта

```
├── main.go              — точка входа
├── pkg/
│   ├── api/
│   │   ├── api.go       — регистрация маршрутов
│   │   ├── nextdate.go  — функция NextDate и /api/nextdate
│   │   ├── addtask.go   — POST /api/task
│   │   ├── tasks.go     — GET /api/tasks
│   │   ├── task.go      — GET, PUT, DELETE /api/task
│   │   ├── done.go      — POST /api/task/done
│   │   └── auth.go      — JWT аутентификация
│   ├── db/
│   │   ├── db.go        — подключение к SQLite
│   │   └── task.go      — CRUD операции
│   └── server/
│       └── server.go    — запуск HTTP-сервера
├── web/                 — фронтенд (готовый)
├── tests/               — тесты (готовые)
├── Dockerfile
└── docker-compose.yml
```
