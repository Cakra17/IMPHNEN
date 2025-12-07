#!/usr/bin/env bash

set -a
source .env
set +a

case "$1" in
  up|down)
    migrate -verbose -path=./db/migrations -database $DSN $1
    ;;
  *)
    echo "you have to pass 'up' or 'down' as an argument"
    exit 1
    ;;
esac
