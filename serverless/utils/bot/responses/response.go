package bot_responses

import "github.com/bwmarrin/discordgo"

func RespondUnexpectedError(c *discordgo.Session, interaction *discordgo.Interaction) {
	c.FollowupMessageCreate(interaction, true, &discordgo.WebhookParams{
		Content: "Unexpected error occured. Please let Ant know.",
	})
}
