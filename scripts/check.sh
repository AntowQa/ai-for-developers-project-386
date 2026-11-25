#!/bin/sh
set -eu

root=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
cd "$root"

export PATH="$(go env GOPATH)/bin:${PATH}"
export DATABASE_URL="${DATABASE_URL:-postgres://booking:booking@127.0.0.1:54329/booking?sslmode=disable}"

(cd web && pnpm install --frozen-lockfile)
(cd tools/sdkgen && pnpm install --frozen-lockfile)

(cd api && go install tool)
buf generate --path proto/booking/v1/health.proto
openapi=$(mktemp)
trap 'rm -f "$openapi"' EXIT
(
  cd tools/sdkgen
  pnpm exec swagger2openapi "$root/api/gen/openapiv2/booking/v1/health.swagger.json" -o "$openapi"
  pnpm exec openapi-typescript "$openapi" -o "$root/web/src/gen/schema.ts"
)
(cd web && pnpm exec oxfmt src/gen/schema.ts)

if ! git diff --exit-code -- api/gen web/src/gen; then
  echo "generated files differ from git" >&2
  exit 1
fi
untracked=$(git ls-files --others --exclude-standard -- api/gen web/src/gen)
if [ -n "$untracked" ]; then
  echo "generated files are not committed:" >&2
  echo "$untracked" >&2
  exit 1
fi

buf lint
(cd api && golangci-lint fmt --diff && golangci-lint run)
(cd web && pnpm exec oxlint . && pnpm exec oxfmt --check . && pnpm exec tsc -b && pnpm exec vite build)

docker compose -f compose.yaml -f compose.check.yaml up -d --wait postgres
(cd api && go test ./...)
