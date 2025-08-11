package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"askeladden/internal/bot"
	"askeladden/internal/bot/services"
	"askeladden/internal/commands"
	"askeladden/internal/reactions"
)

// Handler struct holds the bot instance and services.
type Handler struct {
	Bot            *bot.Bot
	Services       *services.BotServices
	warnedChannels map[string]bool
}

// New creates a new Handler instance.
func New(b *bot.Bot) *Handler {
	// Create services instance
	botServices := &services.BotServices{
		Approval: &services.ApprovalService{Bot: b},
	}
	
	return &Handler{
		Bot:      b,
		Services: botServices,
	}
}

// Ready handles the ready event.
func (h *Handler) Ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("[BOT] Askeladden is connected and ready.")
	if h.Bot.Config.Discord.LogChannelID != "" {
		embed := services.CreateBotEmbed(s, "🟢 Online", "Askeladden is online and ready! ✨", services.EmbedTypeSuccess)
		s.ChannelMessageSendEmbed(h.Bot.Config.Discord.LogChannelID, embed)
	}
}

// MessageCreate handles new messages.
func (h *Handler) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if s.State.User.ID == m.Author.ID {
		return
	}

	log.Printf("[DEBUG] Received message: '%s', prefix: '%s'", m.Content, h.Bot.Config.Discord.Prefix)

	// Handle commands (messages with prefix)
	if strings.HasPrefix(m.Content, h.Bot.Config.Discord.Prefix) {
		// Extract command and arguments
		commandWithPrefix := strings.Split(m.Content, " ")[0]
		log.Printf("[DEBUG] Command with prefix: '%s'", commandWithPrefix)

		// Check if the command is admin-only
		if commands.IsAdminCommand(commandWithPrefix) {
			log.Printf("[DEBUG] Command is admin-only, checking permissions")
			if !h.Services.Approval.UserHasOpplysarRole(s, m.GuildID, m.Author.ID) {
				log.Printf("[DEBUG] User doesn't have admin role, ignoring")
				return // Silently ignore admin commands from non-admins
			}
		}

		// Run the command
		log.Printf("[DEBUG] Running command: '%s'", commandWithPrefix)
		commands.MatchAndRunCommand(commandWithPrefix, s, m, h.Bot)
		return
	}

	// Handle non-command messages (replies, etc.)
	h.handleNonCommandMessage(s, m)

	// Check for banned words in the message
	h.checkForBannedWords(s, m)
}


// ReactionAdd handles when a user reacts to a message.
func (h *Handler) ReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}

	// Handle hammer emoji for reporting incorrect words
	if r.Emoji.Name == "🔨" {
		h.promptForIncorrectWord(s, r)
		return
	}

	// Check if the reaction is admin-only
	if reactions.IsAdminReaction(r.Emoji.Name) {
		if !h.Services.Approval.UserHasOpplysarRole(s, r.GuildID, r.UserID) {
			return // Silently ignore admin reactions from non-admins
		}
	}

	// Run the reaction handler
	reactions.MatchAndRunReaction(r.Emoji.Name, s, r, h.Bot)
}

// promptForIncorrectWord prompts the user to provide the incorrect word(s)
func (h *Handler) promptForIncorrectWord(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	log.Printf("User %s reported an incorrect word in message %s", r.UserID, r.MessageID)
	_, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Printf("Error fetching message: %v", err)
		return
	}

	// Store the original hammered message info in the prompt for later reference
	promptEmbed := services.NewEmbedBuilder().
		SetTitle("🚨 Rapporter feil ord").
		SetDescription(fmt.Sprintf("Ver snill og svar med ord som er feil, skilde med komma viss det er fleire.\n\n[Hopp til opphavleg melding](https://discord.com/channels/%s/%s/%s)", r.GuildID, r.ChannelID, r.MessageID)).
		SetColorByType(services.EmbedTypeError).
		Build()

	s.ChannelMessageSendEmbed(r.ChannelID, promptEmbed)
}

// handleNonCommandMessage processes non-command messages like replies
func (h *Handler) handleNonCommandMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("[DEBUG] Processing non-command message: %s", m.Content)
	
	// Check if this is a reply to a bot message (indicating user is responding to hammer emoji prompt)
	if m.ReferencedMessage != nil && m.ReferencedMessage.Author.ID == s.State.User.ID {
		log.Printf("[DEBUG] User %s replied to bot message with: %s", m.Author.ID, m.Content)
		
		// Check if the referenced message was a "Rapporter feil ord" prompt
		if len(m.ReferencedMessage.Embeds) > 0 && 
		   (strings.Contains(m.ReferencedMessage.Embeds[0].Title, "Report Incorrect Word") || strings.Contains(m.ReferencedMessage.Embeds[0].Title, "Rapporter feil ord")) {
			h.processIncorrectWordReport(s, m)
		}
	}
}

// processIncorrectWordReport processes the user's response to the incorrect word prompt
func (h *Handler) processIncorrectWordReport(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("[DEBUG] Processing incorrect word report from user %s: %s", m.Author.ID, m.Content)
	
	// Parse the words from the message (comma-separated)
	words := strings.Split(m.Content, ",")
	
	for i, word := range words {
		words[i] = strings.TrimSpace(word)
	}
	
	// Filter out empty words
	var validWords []string
	for _, word := range words {
		if word != "" {
			validWords = append(validWords, word)
		}
	}
	
	if len(validWords) == 0 {
		log.Printf("No valid words found in report")
		return
	}
	
	// Extract original hammered message info from the embed description jump link
	originalChannelID := ""
	originalMessageID := ""
	if m.ReferencedMessage != nil && len(m.ReferencedMessage.Embeds) > 0 {
		// Parse jump link from description: https://discord.com/channels/GUILD_ID/CHANNEL_ID/MESSAGE_ID
		description := m.ReferencedMessage.Embeds[0].Description
		
		// Look for Discord jump link pattern
		if strings.Contains(description, "discord.com/channels/") {
			// Find the URL in the description
			startIdx := strings.Index(description, "https://discord.com/channels/")
			if startIdx != -1 {
				// Extract everything after the base URL
				urlPart := description[startIdx+len("https://discord.com/channels/"):]
				// Find the end of the URL (usually a closing parenthesis for markdown link)
				endIdx := strings.Index(urlPart, ")")
				if endIdx != -1 {
					urlPart = urlPart[:endIdx]
				}
				
				// Split by / to get guild/channel/message IDs
				parts := strings.Split(urlPart, "/")
				if len(parts) == 3 {
					// parts[0] = guild_id, parts[1] = channel_id, parts[2] = message_id
					originalChannelID = parts[1]
					originalMessageID = parts[2]
				}
			}
		}
		
		// Fallback: try to parse from footer if jump link parsing failed
		if originalChannelID == "" && originalMessageID == "" && m.ReferencedMessage.Embeds[0].Footer != nil {
			// Parse footer text: "CHANNEL_ID|MESSAGE_ID" or old format "Channel: CHANNEL_ID | Message: MESSAGE_ID"
			footerText := m.ReferencedMessage.Embeds[0].Footer.Text
			if strings.Contains(footerText, "|") && !strings.Contains(footerText, " | ") {
				// New format: "CHANNEL_ID|MESSAGE_ID"
				parts := strings.Split(footerText, "|")
				if len(parts) == 2 {
					originalChannelID = parts[0]
					originalMessageID = parts[1]
				}
			} else if strings.Contains(footerText, " | ") {
				// Old format: "Channel: CHANNEL_ID | Message: MESSAGE_ID"
				footerParts := strings.Split(footerText, " | ")
				if len(footerParts) == 2 {
					if strings.HasPrefix(footerParts[0], "Channel: ") {
						originalChannelID = strings.TrimPrefix(footerParts[0], "Channel: ")
					}
					if strings.HasPrefix(footerParts[1], "Message: ") {
						originalMessageID = strings.TrimPrefix(footerParts[1], "Message: ")
					}
				}
			}
		}
	}
	log.Printf("Extracted original message info - Channel: %s, Message: %s", originalChannelID, originalMessageID)
	
	// Create a variable to hold the thread, but don't create it yet
	// thread := h.Services.Approval.PostBannedWordReport(s, validWords, m.Author.ID, originalMessageID)
	
	// Add words to database and update with forum thread ID if one was created
	var newWords []string
	var existingWords []string
	
	for _, word := range validWords {
		// Check if word already exists
		isBanned, _, err := h.Bot.Database.IsBannedWord(word)
		if err != nil {
			log.Printf("Error checking if word '%s' exists: %v", word, err)
			continue
		}
		
		if isBanned {
			existingWords = append(existingWords, word)
			// Update existing word with forum thread ID if we created one
			// if thread != nil {
			// 	err = h.Bot.Database.UpdateBannedWordThread(word, thread.ID)
			// 	if err != nil {
			// 		log.Printf("Error updating forum thread for word '%s': %v", word, err)
			// 	}
			// }
		} else {
			newWords = append(newWords, word)
			// Add new word as pending approval with forum thread ID if we created one
			// if thread != nil {
			// 	forumThreadID = thread.ID
			// }
			
			wordID, err := h.Bot.Database.AddBannedWordPending(word, "Reported via hammer emoji", m.Author.ID, m.Author.Username, "", fmt.Sprintf("%s|%s", originalChannelID, originalMessageID))
			if err != nil {
				log.Printf("Error adding pending banned word '%s': %v", word, err)
			} else {
				log.Printf("Added pending banned word: %s with ID %d", word, wordID)
				// Post to retting channel for approval
				h.Services.Approval.PostPendingBannedWordToRettingChannel(wordID)
			}
		}
	}
	
	// Send confirmation with appropriate messaging
	var confirmText string
	if len(newWords) > 0 && len(existingWords) > 0 {
		confirmText = fmt.Sprintf("Takk! Nye ord lagt til: %s. Finst allereie: %s", strings.Join(newWords, ", "), strings.Join(existingWords, ", "))
	} else if len(newWords) > 0 {
		confirmText = fmt.Sprintf("Takk! Desse orda har blitt lagt til som forbodne: %s", strings.Join(newWords, ", "))
	} else {
		confirmText = fmt.Sprintf("Alle orda finst allereie i lista over forbodne ord: %s", strings.Join(existingWords, ", "))
	}
	
	if len(newWords) > 0 {
		confirmText += "\n\nEi diskusjonstråd vil bli oppretta etter godkjenning. Sjekk grammatikkforumet seinare."
	} else {
		confirmText += "\n\nSjå eksisterande diskusjonar i grammatikkforumet for desse orda."
	}
	
	confirmEmbed := services.CreateSuccessEmbed("Ord rapporterte", confirmText)
	
	s.ChannelMessageSendEmbed(m.ChannelID, confirmEmbed)
}

// checkForBannedWords checks if a message contains banned words and shows warnings
func (h *Handler) checkForBannedWords(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Skip checking replies to bot messages (to avoid warning on reporting flow)
	if m.ReferencedMessage != nil && m.ReferencedMessage.Author.ID == s.State.User.ID {
		return
	}

	// Convert message to lowercase and split into words
	messageWords := strings.Fields(strings.ToLower(m.Content))
	var foundBannedWords []string
	var forumThreads []string

	for _, word := range messageWords {
		// Clean word of punctuation
		cleanWord := strings.Trim(word, ".,!?;:()[]{}\"'")
		
		isBanned, bannedWord, err := h.Bot.Database.IsBannedWord(cleanWord)
		if err != nil {
			log.Printf("Error checking banned word '%s': %v", cleanWord, err)
			continue
		}
		
		if isBanned {
			foundBannedWords = append(foundBannedWords, cleanWord)
			if bannedWord.ForumThreadID != nil {
				forumThreads = append(forumThreads, *bannedWord.ForumThreadID)
			}
			log.Printf("Detected banned word '%s' in message from user %s", cleanWord, m.Author.ID)
		}
	}

	if len(foundBannedWords) > 0 {
		h.sendBannedWordWarning(s, m, foundBannedWords, forumThreads)
	}
}

// sendBannedWordWarning sends a warning about detected banned words
func (h *Handler) sendBannedWordWarning(s *discordgo.Session, m *discordgo.MessageCreate, bannedWords []string, forumThreads []string) {
	warningEmbed := services.CreateBannedWordWarningEmbed(bannedWords, forumThreads)

	// Send as a reply to the original message
	reply := &discordgo.MessageSend{
		Embed: warningEmbed,
		Reference: &discordgo.MessageReference{
			MessageID: m.ID,
			ChannelID: m.ChannelID,
			GuildID:   m.GuildID,
		},
	}

	_, err := s.ChannelMessageSendComplex(m.ChannelID, reply)
	if err != nil {
		log.Printf("Failed to send banned word warning: %v", err)
	}
}

// InteractionCreate handles button clicks and other interactions
func (h *Handler) InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent {
		customID := i.MessageComponentData().CustomID

		if customID == "confirm_clear_database" {
			// Check if the user is an admin
			if !h.Services.Approval.UserHasOpplysarRole(s, i.GuildID, i.Member.User.ID) {
				// Respond to the interaction with an error message
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Du har ikkje tilgang til å tømme databasen.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					log.Printf("Failed to send interaction response: %v", err)
				}
				return
			}

			// Clear the database
			if err := h.Bot.Database.ClearDatabase(); err != nil {
				log.Printf("Failed to clear database: %v", err)
				// Let the user know something went wrong
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Ein feil oppstod under tømming av databasen.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				return
			}

			// Respond to the interaction
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "✅ Databasen har blitt tømt.",
				},
			})
			if err != nil {
				log.Printf("Failed to send interaction response: %v", err)
			}

			// Delete the original confirmation message
			s.ChannelMessageDelete(i.ChannelID, i.Message.ID)
		}
	}
}


