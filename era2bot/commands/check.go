package commands

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"

	"github.com/disgoorg/bot-template/era2bot"
	"github.com/disgoorg/bot-template/era2bot/hiveTools"
	"github.com/disgoorg/bot-template/era2bot/messageTools"
	"github.com/disgoorg/bot-template/era2bot/mySqlTools"
	"github.com/disgoorg/bot-template/era2bot/restTools"
)

// Command definition (registers it to Discord)
var checklink = discord.SlashCommandCreate{
	Name:        "checklink",
	Description: "Checks a curationlink issues.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "url",
			Description: "The URL to check",
			Required:    true,
		},
	},
}

// Handler for the command
func CheckLinkHandler(b *era2bot.Bot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		url := e.SlashCommandInteractionData().String("url")
		keys := []string{
			"link format",
			"post age",
			"already voted",
			"our Voting Power",
			"internal blacklist",
			"spaminator blacklist",
			"global queue",
			"your queue",
			"your curations",
			"word count",
			"plagiarism check"}
		fieldsMap := make(map[string]string)
		for i, key := range keys {
			if i == 0 {
				fieldsMap[key] = "1running..."
			} else {
				fieldsMap[key] = "0waiting..."
			}

		}

		embedFields := messageTools.CreateProgressFields(fieldsMap, keys)

		embed := discord.Embed{
			Title:       "🔍 Checking Link",
			Description: "Performing checks on: `" + url + "`",
			Color:       0x3498db,
			Fields:      embedFields,
		}

		// Respond with initial embed
		if err := e.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{embed},
		}); err != nil {
			return err
		}
		author := ""
		permlink := ""

		for i, k := range keys {

			if k == "link format" {
				a, p, error := restTools.ParseHiveLink(url)
				if error != nil {
					fieldsMap[keys[i]] = "4" + error.Error()
				} else {
					author = a
					permlink = p
					fieldsMap[keys[i]] = fmt.Sprintf("2we are checking %s/%s", author, permlink)

				}
			} else if k == "already voted" {
				alreadyVoted, _ := hiveTools.HasVoted(author, permlink, "diyhub")
				fieldsMap[keys[i]] = alreadyVoted

			} else if k == "our Voting Power" {
				enoughVP, _ := hiveTools.GetVotingPowerPercent("diyhub")
				fieldsMap[keys[i]] = enoughVP

			} else if k == "spaminator blacklist" {
				fieldsMap[keys[i]] = restTools.CheckSpaminator(&author)
			} else if k == "post age" {
				postAgeReturnValue, _ := hiveTools.PostAge(&author, &permlink)
				fieldsMap[keys[i]] = postAgeReturnValue
			} else if k == "internal blacklist" {
				// TODO: actually query the database here
				dbError := mySqlTools.PingDB()
				if dbError != nil {
					fieldsMap[keys[i]] = "4database is not online"
				} else {
					fieldsMap[keys[i]] = "2database is online"
				}

			} else if k == "global queue" {
				// TODO: actually query the database here
				dbError := mySqlTools.PingDB()
				if dbError != nil {
					fieldsMap[keys[i]] = "4database is not online"
				} else {
					fieldsMap[keys[i]] = "2database is online"
				}

			} else if k == "your queue" {
				// TODO: actually query the database here
				dbError := mySqlTools.PingDB()
				if dbError != nil {
					fieldsMap[keys[i]] = "4database is not online"
				} else {
					fieldsMap[keys[i]] = "2database is online"
				}

			} else if k == "your curations" {
				// TODO: actually query the database here
				dbError := mySqlTools.PingDB()
				if dbError != nil {
					fieldsMap[keys[i]] = "4database is not online"
				} else {
					fieldsMap[keys[i]] = "2database is online"
				}

			} else {
				print(permlink)
				// do some checks here
				time.Sleep(1 * time.Second)
				fieldsMap[keys[i]] = "2check passed"

				if i == 1 {
					fieldsMap[keys[i]] = "3warning: there is evidence that whatever reason"
				}

			}

			if i < len(keys)-1 {
				fieldsMap[keys[i+1]] = "1Running..."
			}

			// check if the current step failed
			if string(fieldsMap[keys[i]][0]) == "4" {
				for iRest := i; iRest < len(keys); iRest++ {
					if iRest > i {
						fieldsMap[keys[iRest]] = "5skipped because " + keys[i] + " failed"
					}
				}
			}

			embedFields = messageTools.CreateProgressFields(fieldsMap, keys)

			_, err := e.Client().Rest().UpdateInteractionResponse(
				e.ApplicationID(),
				e.Token(),
				discord.NewMessageUpdateBuilder().
					SetEmbeds(
						discord.Embed{
							Title:       embed.Title,
							Description: embed.Description,
							Color:       0xff0000,
							Fields:      embedFields,
						},
					).
					Build(),
			)
			if err != nil {
				return err
			}
			if string(fieldsMap[keys[i]][0]) == "4" {
				break
			}

		}

		return nil
	}
}
