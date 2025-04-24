package messageTools

import (
	"github.com/disgoorg/disgo/discord"
)

func CreateEmbed(title string, description string, fields map[string]string) discord.Embed {

	embedFields := CreateFields(fields)

	// Construct the embed message
	embed := discord.Embed{
		Title:       title,
		Description: description,
		Color:       0x00ffcc, // Optional: Embed sidebar color
		Fields:      embedFields,
	}
	return embed

}

func CreateFields(fieldsMap map[string]string) []discord.EmbedField {
	var embedFields []discord.EmbedField
	for name, value := range fieldsMap {
		embedFields = append(embedFields, discord.EmbedField{
			Name:   name,
			Value:  value,
			Inline: BoolPtr(true), // or false, depending on your layout
		})
	}
	return embedFields
}

func CreateProgressFields(fieldsMap map[string]string, keys []string) []discord.EmbedField {
	var embedFields []discord.EmbedField
	// value contains statusInt and status message
	// e.g 2Failed
	// 0 = waiting / 1 = running / 2 = success / 3 = warning /4 = failed / 5 = skip bc of failed
	// for name, value := range fieldsMap {
	for _, key := range keys {
		var status string = string(fieldsMap[key][0])
		var statusSymbol string

		if status == "0" {
			statusSymbol = "🟡"
		} else if status == "1" {
			statusSymbol = "▶️"
		} else if status == "2" {
			statusSymbol = "🟢"
		} else if status == "3" {
			statusSymbol = "⚠️"
		} else if status == "4" {
			statusSymbol = "🔴"
		} else {
			statusSymbol = "⛔"
		}
		embedFields = append(embedFields, discord.EmbedField{
			Name:   statusSymbol + " " + key,
			Value:  fieldsMap[key][1:],
			Inline: BoolPtr(false), // or false, depending on your layout
		})
	}
	return embedFields
}

func BoolPtr(b bool) *bool {
	return &b
}
