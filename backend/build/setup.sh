#!/bin/bash
set -e

PACKAGES="ca-certificates"

# shellcheck disable=SC2086
apk add --no-cache ${PACKAGES} 

addgroup -g 1000 "${SERVICE_USER}"
adduser -D --no-create-home --home "${SERVICE_DIR}" --shell /bin/bash --uid 1000 --ingroup "${SERVICE_USER}" "${SERVICE_USER}"
mkdir -p "${SERVICE_DIR}"
chown -R "${SERVICE_USER}":"${SERVICE_USER}" "${SERVICE_DIR}"