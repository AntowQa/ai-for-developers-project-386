#!/bin/sh
set -eu

root=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
cd "$root"

export DATABASE_URL="${DATABASE_URL:-postgres://booking:booking@127.0.0.1:54329/booking?sslmode=disable}"

(cd web && pnpm install --frozen-lockfile)
(cd tools/sdkgen && pnpm install --frozen-lockfile)

generated=$(mktemp -d)
trap 'rm -rf "$generated"' EXIT
pnpm --dir tools/sdkgen exec tsp format --check main.tsp
GEN_DIR="$generated" pnpm --dir tools/sdkgen run generate
for file in api/gen/openapi/openapi.yaml web/src/gen/schema.ts; do
  if ! cmp "$file" "$generated/$file"; then
    echo "Regenerate $file with 'mise run generate'." >&2
    exit 1
  fi
done
(cd api && golangci-lint fmt --diff && golangci-lint run)
(cd web && pnpm exec oxlint . && pnpm exec oxfmt --check . && pnpm exec tsc -b && pnpm exec vite build)

docker compose -f compose.yaml -f compose.check.yaml up -d --wait postgres
(cd api && go test -count=1 ./...)
