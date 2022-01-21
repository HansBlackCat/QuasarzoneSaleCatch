package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	zlog "github.com/rs/zerolog/log"
)

func InitBot(tok string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(tok)
	OsPanic(err)
	zlog.Info().Str("Username", bot.Self.UserName).Msg("Authorized on account")

	return bot
}
