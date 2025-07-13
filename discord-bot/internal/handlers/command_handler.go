package handlers

import (
	"fmt"
	"log"

	"sukimise-discord-bot/internal/services"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler struct {
	discordService *services.DiscordService
}

func NewCommandHandler(discordService *services.DiscordService) *CommandHandler {
	return &CommandHandler{
		discordService: discordService,
	}
}

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "connect",
		Description: "Connect your Discord account to Sukimise",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "username",
				Description: "Your Sukimise username",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "password",
				Description: "Your Sukimise password",
				Required:    true,
			},
		},
	},
	{
		Name:        "disconnect",
		Description: "Disconnect your Discord account from Sukimise",
	},
	{
		Name:        "add",
		Description: "Add a store from Google Maps URL",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "url",
				Description: "Google Maps URL of the store",
				Required:    true,
			},
		},
	},
	{
		Name:        "help",
		Description: "Show help information about Sukimise Discord Bot",
	},
}

func (h *CommandHandler) HandleReady(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Bot is ready! Logged in as %s", s.State.User.Username)

	// Register slash commands
	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			log.Printf("Cannot create '%s' command: %v", cmd.Name, err)
		}
	}
}

func (h *CommandHandler) HandleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == "" {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "connect":
		h.handleConnectCommand(s, i)
	case "disconnect":
		h.handleDisconnectCommand(s, i)
	case "add":
		h.handleAddCommand(s, i)
	case "help":
		h.handleHelpCommand(s, i)
	}
}

func (h *CommandHandler) handleConnectCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Defer response to avoid timeout
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	options := i.ApplicationCommandData().Options
	username := options[0].StringValue()
	password := options[1].StringValue()

	discordID := i.Member.User.ID

	// Connect Discord user to Sukimise
	link, err := h.discordService.ConnectDiscordUser(discordID, username, password)
	if err != nil {
		content := fmt.Sprintf("❌ **Connection Failed**\n%s", err.Error())
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		return
	}

	content := fmt.Sprintf("✅ **Successfully Connected!**\n"+
		"Discord account linked to Sukimise user: **%s**\n"+
		"You can now use `/add <google_maps_url>` to register stores!", link.Username)

	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: content,
		Flags:   discordgo.MessageFlagsEphemeral,
	})
}

func (h *CommandHandler) handleDisconnectCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	discordID := i.Member.User.ID

	err := h.discordService.DisconnectDiscordUser(discordID)
	if err != nil {
		content := fmt.Sprintf("❌ **Disconnection Failed**\n%s", err.Error())
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	content := "✅ **Successfully Disconnected!**\n" +
		"Your Discord account has been disconnected from Sukimise.\n" +
		"Use `/connect` to link again."

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func (h *CommandHandler) handleAddCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Defer response to avoid timeout
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	options := i.ApplicationCommandData().Options
	googleMapURL := options[0].StringValue()

	discordID := i.Member.User.ID

	// Add store from Google Maps URL
	storeResp, err := h.discordService.AddStoreFromGoogleMaps(discordID, googleMapURL)
	if err != nil {
		content := fmt.Sprintf("❌ **Store Registration Failed**\n%s", err.Error())
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: content,
		})
		return
	}

	// Success message with detailed log
	content := fmt.Sprintf("✅ **Store Successfully Registered!**\n\n"+
		"**Store Details:**\n"+
		"• **Name:** %s\n"+
		"• **Address:** %s\n"+
		"• **Store ID:** %s\n\n"+
		"**Registration Info:**\n"+
		"• **Registered by:** <@%s>\n"+
		"• **Google Maps URL:** %s\n\n"+
		"The store has been added to the Sukimise database and is now available for reviews and management.",
		storeResp.Name, storeResp.Address, storeResp.ID, discordID, googleMapURL)

	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: content,
	})

	// Log the successful registration
	log.Printf("Store registered successfully - ID: %s, Name: %s, Discord User: %s", 
		storeResp.ID, storeResp.Name, discordID)
}

func (h *CommandHandler) handleHelpCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	helpContent := `**🏪 Sukimise Discord Bot - Help**

**Available Commands:**

**🔗 /connect <username> <password>**
Connect your Discord account to your Sukimise account.
• Required: Your Sukimise username and password
• Note: Each Discord account can only be linked to one Sukimise account

**🔌 /disconnect**
Disconnect your Discord account from Sukimise.
• Removes the link between your Discord and Sukimise accounts

**🏪 /add <google_maps_url>**
Add a store from Google Maps URL to Sukimise database.
• Required: Google Maps URL (must start with https://www.google.com/maps/place/)
• Note: You must be connected to Sukimise first using /connect

**❓ /help**
Show this help information.

**📋 How to Use:**
1. First, use '/connect' to link your Discord account to Sukimise
2. Use '/add' with Google Maps URLs to register stores
3. Use '/disconnect' if you want to unlink your accounts

**🔗 Supported Google Maps URLs:**
• https://www.google.com/maps/place/Store+Name/...
• https://maps.google.com/maps/place/Store+Name/...
• https://goo.gl/maps/...
• https://maps.app.goo.gl/...

**💡 Tips:**
• Store information is extracted automatically from Google Maps
• All registered stores are tagged with "discord"
• Registration logs are posted to the channel for transparency

**🔒 Privacy:**
• Connection information is stored securely
• Passwords are only used for authentication and not stored
• Each Sukimise account can only be linked to one Discord account`

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: helpContent,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}