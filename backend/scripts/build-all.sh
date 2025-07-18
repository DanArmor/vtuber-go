#!/bin/bash
set -e

# shellcheck disable=SC2046
SCRIPT_PATH="$(dirname $(realpath $0))"

echo "Start build" &&
sudo docker build -t local/vtuber-go-worker:0.0.2 -f build/worker/Dockerfile . &&
sudo docker build -t local/vtuber-go-backend:0.0.2 -f build/app/Dockerfile . &&
echo "Build done"