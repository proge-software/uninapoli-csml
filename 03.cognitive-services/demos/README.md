# Why So Serious Bot

It is a simple set of toy-bot organized in a monorepo for Telegram written in Go that makes use of Azure Cognitive Services.


The bot is not going to store any data anywhere.

The following bot are implemented:
- Simple Bot
- Vision Bot
- Lang Bot
- Moderate Bot

## Simple Bot

A simple example of an "Hello world!" bot for telegram.

Topics:
- Go
- Telegram ([tucnak/telebot](https://github.com/tucnak/telebot))
- dotenv ([joho/godotenv](https://github.com/joho/godotenv))

## Vision Bot

Based on Simple Bot, uses the Azure Cognitive Services for Vision to analyze the photos provided by the user.

Topics:
- Face (Azure Cognitive Service)
- Computer Vision (Azure Cognitive Service)
- Form Recognizer (preview) (Azure Cognitive Service)

## Language Bot

Based on Vision Bot, uses the Azure Cognitive Services for Language to analyze the text messages sent by the user.

Topics:
- Text Analytics (Azure Cognitive Service)
- Translator Text (Azure Cognitive Service)

## Moderate Bot

Based on Vision Bot, uses the Content Moderator Azure Cognitive Service to analyze the text messages and photos sent by the user.

Topics:
- Content Moderator (Azure Cognitive Service)


# Start one bot

The following procedures are the same for each bot implemented.
Let assume you want to execute the `simple-bot`.

## Locally

Follow the [telegram guide to create a bot](https://core.telegram.org/bots#3-how-do-i-create-a-bot) and copy the provided `token` at the end of the procedure.
Open the `.env ` file in `./cmd/04.simple-bot/.env` and set the value of the `WSS_TOKEN` environment variable to the token you copied before.

Install [go 1.14](https://golang.org/dl/) and run the following command:

```console
go run ./cmd/04.simple-bot/main.go
```

## Docker

Ensure you have docker installed and running.

Open a shell in the project root and run the following commands (substitute "{{ TG_TOKEN }}" with the token telegram provided you when registering the bot):

```console
docker build -f build/docker/04.simple-bot/Dockerfile -t simple-bot:latest .
docker run -it --rm -e WSS_TOKEN="{{ TG_TOKEN }}" simple-bot:latest
```
