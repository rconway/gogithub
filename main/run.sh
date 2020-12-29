#!/usr/bin/env bash

ORIG_DIR="$(pwd)"
cd "$(dirname "$0")"
BIN_DIR="$(pwd)"

trap "cd '${ORIG_DIR}'" EXIT

reflex -r '\.go$' -s sudo GO_LISTEN_ADDRESS=":80" $(which go) run .
