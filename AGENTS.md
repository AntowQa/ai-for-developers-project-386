# Агенты

Команды запуска, проверок и первого входа описаны в [README](README.md#команды).

## Устройство

- `proto/` — источник контракта. Сейчас это `GET /api/health`.
- `api/` — Go-процесс: gRPC внутри контейнера, HTTP отдаёт grpc-gateway. Перед прослушиванием порта процесс пингует PostgreSQL и выходит с ошибкой, если базы нет.
- `web/` — Vite, React, TypeScript, shadcn/ui. nginx отдаёт сборку и проксирует `/api/` на API.
- `tools/sdkgen/` — генератор TypeScript SDK из OpenAPI. В проекте TypeScript 7, а `openapi-typescript` 7.13 работает на TypeScript 5, поэтому генератор лежит отдельно.
- `compose.yaml` поднимает PostgreSQL, API и nginx. Порт базы наружу не публикуется. `compose.check.yaml` открывает её только на `127.0.0.1:54329` для `go test`.
- Сгенерированные Go, OpenAPI и SDK коммитятся.

## Коммиты

Сообщения коммитов — [Conventional Commits](https://www.conventionalcommits.org/). Первая строка: `type: subject` или `type(scope): subject`. Типы: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`, `ci`, `build`. Хук `commit-msg` отклоняет другие сообщения.
