#!/bin/bash
set -e

# shellcheck disable=SC2046
SCRIPT_PATH="$(dirname $(realpath $0))"

golangci-lint run pkg/... internal/... cmd/...