#!/usr/bin/env bash

if [ -n $1 ]; then
  migrate create -ext sql -dir db/migrations -seq $1
else
  echo "you have to pass 'up' or 'down' as an argument"
  exit 1
fi
