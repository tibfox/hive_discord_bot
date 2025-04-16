package messageTools

import (
	"github.com/disgoorg/disgo/discord"
)

func CreateEmbed(title string, description string, fields map[string]string) discord.Embed {

	var embedFields []discord.EmbedField
	for name, value := range fields {
		embedFields = append(embedFields, discord.EmbedField{
			Name:   name,
			Value:  value,
			Inline: BoolPtr(true), // or false, depending on your layout
		})
	}

	// Construct the embed message
	embed := discord.Embed{
		Title:       title,
		Description: description,
		Color:       0x00ffcc, // Optional: Embed sidebar color
		Fields:      embedFields,
	}
	return embed

}

func BoolPtr(b bool) *bool {
	return &b
}
