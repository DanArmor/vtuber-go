#!/bin/bash

set -e

prj_name="backend"
context_name=$(basename "$PWD")

# shellcheck disable=SC2046
SCRIPT_PATH="$(dirname $(realpath $0))"

go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert ./ent/schema
