# Simple Go Telegram Bot

Example of Telegram Bot written in Go programming language. There is base implementation of [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api) with usage of postgresql database and API calls to external resource.

To run the bot there must be specified `.env` file as shown in `.env-example`, there must exist a postgresql database and applied migrations for it. Migrations are applied using custom `goose` binary. It has to be built from [goose/main.go](goose/main.go) entrypoint:
```sh
$ go build -o bin/goose goose/main.go
$ bin/goose up # running db migrations to the latest version
```

The bot can be manually launched from [server/main.go](server/main.go) file:

```sh
$ go run server/main.go
```

or from manually built binary:

```sh
$ go build -o bin/server server/main.go
$ bin/server
```
