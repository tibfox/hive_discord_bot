package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"

	"github.com/disgoorg/bot-template/bottemplate"
)

func MessageHandler(b *bottemplate.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(e *events.MessageCreate) {
		if e.Message.Author.Bot {
			return
		}
		var message string
		if e.Message.Content == "ping" {
			message = "pong"
		} else if e.Message.Content == "pong" {
			message = "ping"
		}
		if message != "" {
			_, _ = e.Client().Rest().CreateMessage(e.ChannelID, discord.NewMessageCreateBuilder().SetContent(message).Build())
		}
	})
}
