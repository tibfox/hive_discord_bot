package commands

import (
	"fmt"

	"github.com/disgoorg/bot-template/era2bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var balance = discord.SlashCommandCreate{
	Name:        "balance",
	Description: "balance command",
}

func BalanceHandler(b *era2bot.Bot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		return e.CreateMessage(discord.MessageCreate{
			Content: fmt.Sprintf("just a test response"),
		})
	}
}
