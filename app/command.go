package main

import (
	"fmt"
	"net/url"

	"github.com/BruceJi7/fcc-bot-go/app/msg"
	"github.com/bwmarrin/discordgo"
)

type Commands struct {
	bot *Bot
}

func (c *Commands) Initialize() {
	c.create()
	c.bot.Session.AddHandler(c.AdminCommandGroup)
	c.bot.Session.AddHandler(c.RegularCommandGroup)
}

func (c *Commands) create() {

	allSuccessful := true

	if _, err := c.bot.Session.ApplicationCommandCreate(c.bot.Cfg.bot.id, c.bot.Cfg.server.guild, EraseCommand); err != nil {
		c.bot.SendLog(msg.LogError, "Whilst adding erase command:")
		c.bot.SendLog(msg.LogError, err.Error())
		allSuccessful = false
	}

	if _, err := c.bot.Session.ApplicationCommandCreate(c.bot.Cfg.bot.id, c.bot.Cfg.server.guild, ForceLogCommand); err != nil {
		c.bot.SendLog(msg.LogError, "Whilst adding forcelog command:")
		c.bot.SendLog(msg.LogError, err.Error())
		allSuccessful = false
	}

	if _, err := c.bot.Session.ApplicationCommandCreate(c.bot.Cfg.bot.id, c.bot.Cfg.server.guild, LearningResourceCommand); err != nil {
		c.bot.SendLog(msg.LogError, "Whilst adding learning resource command:")
		c.bot.SendLog(msg.LogError, err.Error())
		allSuccessful = false
	}

	if allSuccessful {
		c.bot.SendLog(msg.LogOnReady, "All commands successfully added")
	}
}

var EraseCommand = &discordgo.ApplicationCommand{
	Name:        "erase",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Erase messages in a channel",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "multiple",
			Type:        discordgo.ApplicationCommandOptionInteger,
			Description: "Specify amount to erase",
		},
	},
}

var ForceLogCommand = &discordgo.ApplicationCommand{
	Name:        "forcelog",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Force Bot to Log Something",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "message",
			Type:        discordgo.ApplicationCommandOptionString,
			Description: "Specify log message",
			Required:    true,
		},
	},
}

var LearningResourceCommand = &discordgo.ApplicationCommand{
	Name:        "learning-resource",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Submit a useful learning resource",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "resource-url",
			Type:        discordgo.ApplicationCommandOptionString,
			Description: "A valid url",
			Required:    true,
		},
		{
			Name:        "resource-description",
			Type:        discordgo.ApplicationCommandOptionString,
			Description: "A description of the resource. What language is it for? What can we learn from it?",
			Required:    true,
		},
	},
}

func (c *Commands) AdminCommandGroup(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	options := data.Options

	interactionID := i.Interaction.ID
	interactionChannel, _ := c.bot.Utils.GetChannelByID(i.ChannelID)
	interactionMember := i.Member

	interactionMemberIsAdmin, err := c.bot.Utils.IsUserAdmin(interactionMember.User.ID)
	if err != nil {
		c.bot.SendLog(msg.LogError, "Whilst evaluating admin privileges:")
		c.bot.SendLog(msg.LogError, err.Error())
		return
	} else {
		if !interactionMemberIsAdmin {
			c.bot.SendLog(msg.LogError, fmt.Sprintf("admin commands were exposed to %s", interactionMember.User.ID))
			return
		}
	}

	switch data.Name {
	case "erase":

		if len(options) == 0 {
			c.SingleErase(i, interactionChannel, interactionID, interactionMember)
		} else {
			c.MultiErase(i, options, interactionChannel, interactionID, interactionMember)
		}

	case "forcelog":

		err := s.InteractionRespond(i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{Content: "Log made in log channel", Flags: 1 << 6},
			})

		if err != nil {
			c.bot.SendLog(msg.LogError, "Whilst responding to command forcelog:")
			c.bot.SendLog(msg.LogError, err.Error())
		} else {
			logString := options[0].StringValue()
			c.bot.SendLog(msg.CommandForceLog, fmt.Sprintf("By User %s: %s", interactionMember.User.Username, logString))
		}
	}
}

func (c *Commands) RegularCommandGroup(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	options := data.Options
	interactionMember := i.Member

	switch data.Name {
	case "learning-resource":

		resourceUrl := options[0].StringValue()
		resourceDescription := options[1].StringValue()

		if _, err := url.ParseRequestURI(resourceUrl); err != nil {
			s.InteractionRespond(i.Interaction,
				&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{Content: "Whoops! It looks like your URL was invalid", Flags: 1 << 6},
				})

		} else {

			messageContents := fmt.Sprintf("Thanks, %s, who posted this resource:\n%s\nDescription: %s", interactionMember.User.Mention(), resourceUrl, resourceDescription)
			c.bot.Session.ChannelMessageSend(c.bot.Cfg.server.learningResources, messageContents)
			c.bot.SendLog(msg.LogLearning, fmt.Sprintf("%s submitted a Learning Resource via the bot", interactionMember.User.Username))

			s.InteractionRespond(i.Interaction,
				&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{Content: "Thanks for posting a learning resource!", Flags: 1 << 6},
				})
		}

	}
}

func (c *Commands) SingleErase(i *discordgo.InteractionCreate, interactionChannel *discordgo.Channel, interactionID string, interactionMember *discordgo.Member) {

	err := c.bot.Session.InteractionRespond(i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Messages Erased", Flags: 1 << 6},
		})

	if err != nil {
		c.bot.SendLog(msg.LogError, "Whilst responding to command erase (single):")
		c.bot.SendLog(msg.LogError, err.Error())
	} else {
		deleteErr := c.DeleteMessages(1, interactionChannel.ID, interactionID)
		if deleteErr != nil {
			c.bot.SendLog(msg.LogError, "Whilst attempting to delete:")
			logMessage := fmt.Sprintf("User %s | channel %s | %s", interactionMember.User.Username, interactionChannel.Name, deleteErr)
			c.bot.SendLog(msg.LogError, logMessage)
		} else {
			logMessage := fmt.Sprintf("User %s | channel %s", interactionMember.User.Username, interactionChannel.Name)
			c.bot.SendLog(msg.CommandErase, logMessage)
		}
	}

}

func (c *Commands) MultiErase(i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption, interactionChannel *discordgo.Channel, interactionID string, interactionMember *discordgo.Member) {

	eraseAmount := options[0].IntValue()
	err := c.bot.Session.InteractionRespond(i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Messages Erased", Flags: 1 << 6},
		})
	if err != nil {
		c.bot.SendLog(msg.LogError, "Whilst responding to command erase (multi):")
		c.bot.SendLog(msg.LogError, err.Error())
	} else {
		deleteErr := c.DeleteMessages(int(eraseAmount), interactionChannel.ID, interactionID)
		if deleteErr != nil {
			logMessage := fmt.Sprintf("User %s | channel %s | amount %d | %s", interactionMember.User.Username, interactionChannel.Name, eraseAmount, deleteErr)
			c.bot.SendLog(msg.LogError, "Whilst attempting to delete:")
			c.bot.SendLog(msg.LogError, logMessage)
		} else {
			logMessage := fmt.Sprintf("User %s | channel %s | amount %d | %s", interactionMember.User.Username, interactionChannel.Name, eraseAmount, deleteErr)
			c.bot.SendLog(msg.CommandErase, logMessage)
		}

	}
}

func (c *Commands) DeleteMessages(howMany int, channel string, messageID string) error {

	messages, err := c.bot.Session.ChannelMessages(channel, howMany, messageID, "", "")
	if err != nil {
		return err
	}
	var messageIDs []string

	for _, m := range messages {
		messageIDs = append(messageIDs, m.ID)
	}
	messageIDs = append(messageIDs, messageID)

	err = c.bot.Session.ChannelMessagesBulkDelete(channel, messageIDs)
	if err != nil {
		return err
	}

	return nil
}
