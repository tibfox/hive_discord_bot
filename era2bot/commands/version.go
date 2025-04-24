package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"

	"github.com/disgoorg/bot-template/era2bot"
	"github.com/disgoorg/bot-template/era2bot/messageTools"
)

var version = discord.SlashCommandCreate{
	Name:        "version",
	Description: "version command",
}

func VersionHandler(b *era2bot.Bot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {

		owner := "tibfox"
		repo := "hive_discord_bot"
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", owner, repo)

		resp, err := http.Get(url)
		if err != nil {
			return e.CreateMessage(discord.MessageCreate{
				Content: noCommitFound(),
			})
		}
		defer resp.Body.Close()

		var commits []Commit
		if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
			return e.CreateMessage(discord.MessageCreate{
				Content: noCommitFound(),
			})
		}

		if len(commits) == 0 {
			fmt.Println("No commits found.")
			return e.CreateMessage(discord.MessageCreate{
				Content: noCommitFound(),
			})
		}

		latest := commits[0]

		fieldsMap := map[string]string{
			"SHA":     latest.SHA,
			"Message": latest.Commit.Message,
			"Date":    latest.Commit.Author.Date.Format(time.RFC1123),
		}

		// Construct the embed message
		embed := messageTools.CreateEmbed("version information", fmt.Sprintf("Details about github.com/%s/%s", owner, repo), fieldsMap)

		return e.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{embed},
		})
	}
}

func noCommitFound() string {
	return "Commit: no commit found"
}

type Commit struct {
	Commit struct {
		Author struct {
			Date time.Time `json:"date"`
		} `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
	SHA string `json:"sha"`
}

// func main() {
// 	owner := "vercel"
// 	repo := "next.js"
// 	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", owner, repo)

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		fmt.Println("Error fetching data:", err)
// 		os.Exit(1)
// 	}
// 	defer resp.Body.Close()

// 	var commits []Commit
// 	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
// 		fmt.Println("Error decoding JSON:", err)
// 		os.Exit(1)
// 	}

// 	if len(commits) == 0 {
// 		fmt.Println("No commits found.")
// 		return
// 	}

// 	latest := commits[0]
// 	fmt.Printf("Latest Commit:\n")
// 	fmt.Printf("SHA: %s\n", latest.SHA)
// 	fmt.Printf("Message: %s\n", latest.Commit.Message)
// 	fmt.Printf("Date: %s\n", latest.Commit.Author.Date.Format(time.RFC1123))
// }
