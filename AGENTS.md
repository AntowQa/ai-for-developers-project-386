# Агенты

Команды запуска, проверок и первого входа описаны в [README](README.md#команды).

## Устройство

- `tools/sdkgen/main.tsp` — источник HTTP-контракта. Сейчас это `GET /api/health`.
- `api/` — HTTP-сервер на Go. Перед прослушиванием порта процесс проверяет соединение с PostgreSQL и завершается с ошибкой, если база недоступна.
- `web/` — Vite, React, TypeScript, shadcn/ui. nginx отдаёт сборку и проксирует `/api/` на API.
- `compose.yaml` поднимает PostgreSQL, API и nginx. Порт базы наружу не публикуется. `compose.check.yaml` открывает её только на `127.0.0.1:54329` для `go test`.
- Сгенерированные OpenAPI и типы TypeScript для клиента коммитятся вместе с изменениями контракта.

## Коммиты

Сообщения коммитов — [Conventional Commits](https://www.conventionalcommits.org/). Первая строка: `type: subject` или `type(scope): subject`. Типы: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`, `ci`, `build`. Хук `commit-msg` отклоняет другие сообщения.
