#!/bin/bash
if [[ $# -ne 1 ]]; then
  echo "Wrong amount of arguments"
  exit 1
fi

sudo atlas migrate diff "$1" \
  --dir "file://ent/migrate/migrations" \
  --to "ent://ent/schema" \
  --dev-url "docker://postgres/15/test?search_path=public"