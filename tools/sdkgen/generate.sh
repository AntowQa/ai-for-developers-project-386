#!/bin/sh
set -eu

root=$(cd ../.. && pwd)
out=${GEN_DIR:-$root}
openapi="$out/api/gen/openapi"
schema="$out/web/src/gen/schema.ts"

tsp compile main.tsp --warn-as-error --emit @typespec/openapi3 \
  --option "@typespec/openapi3.emitter-output-dir=$openapi"
mkdir -p "$(dirname "$schema")"
openapi-typescript "$openapi/openapi.yaml" -o "$schema"
pnpm --dir "$root/web" exec oxfmt "$schema"
