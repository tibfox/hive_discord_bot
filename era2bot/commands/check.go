package commands

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"

	"github.com/disgoorg/bot-template/era2bot"
	"github.com/disgoorg/bot-template/era2bot/hiveTools"
	"github.com/disgoorg/bot-template/era2bot/internal/enums"
	"github.com/disgoorg/bot-template/era2bot/internal/models"
	"github.com/disgoorg/bot-template/era2bot/internal/services"
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

var AllCheckSteps = []enums.PossibleCheckStep{
	enums.CheckLinkFormat,
	enums.CheckPostAge,
	enums.CheckAlreadyVoted,
	enums.CheckOurVP,
	enums.CheckInternalBlacklist,
	enums.CheckExternalBlacklist,
	enums.CheckGlobalQueue,
	enums.CheckPersonalQueue,
	enums.CheckPersonalCurations,
	enums.CheckWordCount,
	enums.CheckPlagiarismCheck,
}

// Handler for the command
func CheckLinkHandler(b *era2bot.Bot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {
		url := e.SlashCommandInteractionData().String("url")
		curatorDiscordUserId := e.Member().User.ID.String()
		curatorService := services.NewCuratorService()
		curator, _ := curatorService.GetCurator(nil, &curatorDiscordUserId)

		fieldsMap := make(map[string]string)
		keys := make([]string, 0, len(fieldsMap))

		for i, key := range AllCheckSteps {
			keys = append(keys, key.String())
			if i == 0 {
				fieldsMap[key.String()] = "1running..."
			} else {
				fieldsMap[key.String()] = "0waiting..."
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
		// TODO: extract logic into a service(?)
		for i, k := range AllCheckSteps {
			switch k {
			case enums.CheckLinkFormat:
				a, p, error := restTools.ParseHiveLink(url)
				if error != nil {
					fieldsMap[keys[i]] = "4" + error.Error()
				} else {
					author = a
					permlink = p
					fieldsMap[keys[i]] = fmt.Sprintf("2we are checking %s/%s", author, permlink)

				}
			case enums.CheckAlreadyVoted:
				alreadyVoted, _ := hiveTools.HasVoted(author, permlink, "diyhub")
				fieldsMap[keys[i]] = alreadyVoted

			case enums.CheckOurVP:
				enoughVP, _ := hiveTools.GetVotingPowerPercent("diyhub")
				fieldsMap[keys[i]] = enoughVP

			case enums.CheckExternalBlacklist:
				fieldsMap[keys[i]] = restTools.CheckSpaminator(&author)
			case enums.CheckPostAge:
				postAgeReturnValue, _ := hiveTools.PostAge(&author, &permlink)
				fieldsMap[keys[i]] = postAgeReturnValue

			case enums.CheckInternalBlacklist:
				authorService := services.NewAuthorService()
				authorObject, err := authorService.GetAuthor(&author)
				if authorObject == nil {
					fieldsMap[keys[i]] = fmt.Sprintf("5skip: author not found: %w", err)
				} else {
					if authorObject.BlacklistedBy != "" {
						fieldsMap[keys[i]] = fmt.Sprintf("4author blacklisted by %s: %s", authorObject.BlacklistedBy, authorObject.BlacklistReason)
					} else {
						fieldsMap[keys[i]] = "2author is okay"
					}
				}

			case enums.CheckGlobalQueue:
				// Create service
				queueService := services.NewQueueService()

				// Call service
				queue, err := queueService.GetQueue(models.GlobalQueue, nil)
				if err != nil {
					fieldsMap[keys[i]] = "5queue not available"
				} else {
					if queue.Count > 5 {
						fieldsMap[keys[i]] = "4queue is full - please try again soon"
					} else {
						fieldsMap[keys[i]] = "2queue is okay"
					}
				}
			case enums.CheckPersonalQueue:
				// get curator name by discord id

				if curator == nil {
					fieldsMap[keys[i]] = "5curator not found"
				} else {
					// Create service
					queueService := services.NewQueueService()

					// Call service
					queue, err := queueService.GetQueue(models.GlobalQueue, curator)
					if err != nil {
						fieldsMap[keys[i]] = "5queue not available"
					} else {
						if queue.Count > 5 {
							fieldsMap[keys[i]] = "4your queue is full - please try again soon"
						} else {
							fieldsMap[keys[i]] = "2your queue is okay"
						}
					}
				}
			case enums.CheckPersonalCurations:
				if curator == nil {
					fieldsMap[keys[i]] = "5curator not found"
				} else {
					if curator.CountCurationsInInterval > 5 {
						fieldsMap[keys[i]] = "4you made too many curations today"
					} else {
						fieldsMap[keys[i]] = "2you can still curate"
					}
				}
			case enums.CheckWordCount:
				// TODO: actually query the database here
				dbError := mySqlTools.PingDB()
				if dbError != nil {
					fieldsMap[keys[i]] = "4database is not online"
				} else {
					fieldsMap[keys[i]] = "2database is online"
				}
			case enums.CheckPlagiarismCheck:
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
