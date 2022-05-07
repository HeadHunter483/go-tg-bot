#!/bin/bash

COMMAND=$1

BOT_BIN=bin/bot
GOOSE_BIN=bin/goose

case "$COMMAND" in 
"run")
# run go project
go run cmd/server/main.go
;;
"build")
# build go project only
go build -o $BOT_BIN server/main.go
;;
"build-run")
# build go project and run binary
./bot.sh build
./$BOT_BIN
;;
"migrate")
# build goose binary and run migrations
go build -o $GOOSE_BIN goose/main.go
./$GOOSE_BIN up
;;
esac
