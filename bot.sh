#!/bin/bash

COMMAND=$1

BOT_BIN=bot.exe

case "$COMMAND" in 
"run")
# run go project
go run bot.go db.go tg.go forismatic.go
;;
"build")
# build go project only
go build -o $BOT_BIN bot.go db.go tg.go forismatic.go
;;
"build-run")
# build go project and run binary
./bot.sh build
./$BOT_BIN
;;
"-ct")
# create tables in db
./bot.sh build
./$BOT_BIN -ct
;;
esac
