package services

import (
	"fmt"
	"log"
	"strings"

	"askeladden/internal/bot"
	"askeladden/internal/database"
	"github.com/bwmarrin/discordgo"
)

// ApprovalService handles logic for question approval.
type ApprovalService struct {
	Bot *bot.Bot
}

// PostNewQuestionToApprovalQueue posts a newly created question to the approval queue.
func (s *ApprovalService) PostNewQuestionToApprovalQueue(questionID int64) {
	// Get the question from database
	question, err := s.Bot.Database.GetPendingQuestionByID(int(questionID))
	if err != nil {
		log.Printf("Failed to get question for approval queue posting: %v", err)
		return
	}

	s.postToApprovalQueue(s.Bot.Session, question)
}

// postToApprovalQueue posts a question to the approval queue channel.
func (s *ApprovalService) postToApprovalQueue(session *discordgo.Session, question *database.Question) {
	if s.Bot.Config.Approval.QueueChannelID == "" {
		log.Println("Approval queue channel not configured")
		return
	}

	// Get the author's user info
	author, err := session.User(question.AuthorID)

	approvalEmbed := CreateApprovalEmbed(question.Question, "⏳ Opplysar-godkjenning: ventar", author)

	approvalMessage, err := session.ChannelMessageSendEmbed(s.Bot.Config.Approval.QueueChannelID, approvalEmbed)
	if err != nil {
		log.Printf("Failed to post to approval queue: %v", err)
		return
	}

	// Add thumbs up reaction
	session.MessageReactionAdd(s.Bot.Config.Approval.QueueChannelID, approvalMessage.ID, "👍")

	// Update the database with the approval message ID
	err = s.Bot.Database.UpdateApprovalMessageID(question.ID, approvalMessage.ID)
	if err != nil {
		log.Printf("Failed to update approval message ID: %v", err)
	}
}

// UserHasOpplysarRole checks if a user has the opplysar role.
func (s *ApprovalService) UserHasOpplysarRole(session *discordgo.Session, guildID, userID string) bool {
	if s.Bot.Config.Approval.OpplysarRoleID == "" {
		return false
	}

	member, err := session.GuildMember(guildID, userID)
	if err != nil {
		log.Printf("Failed to get guild member: %v", err)
		return false
	}

	for _, roleID := range member.Roles {
		if roleID == s.Bot.Config.Approval.OpplysarRoleID {
			return true
		}
	}

	return false
}

// UserHasRettskrivarRole checks if a user has the rettskrivar role.
func (s *ApprovalService) UserHasRettskrivarRole(session *discordgo.Session, guildID, userID string) bool {
	if s.Bot.Config.BannedWords.RettskrivarRoleID == "" {
		return false
	}

	member, err := session.GuildMember(guildID, userID)
	if err != nil {
		log.Printf("Failed to get guild member: %v", err)
		return false
	}

	for _, roleID := range member.Roles {
		if roleID == s.Bot.Config.BannedWords.RettskrivarRoleID {
			return true
		}
	}

	return false
}

// NotifyUserApproval notifies the user that their question was approved.
func (s *ApprovalService) NotifyUserApproval(session *discordgo.Session, question *database.Question, approverID string) {
	privateChannel, err := session.UserChannelCreate(question.AuthorID)
	if err != nil {
		log.Printf("Failed to create private channel for approval notification: %v", err)
		return
	}

	embed := CreateBotEmbed(session, "🎉 Gratulerer! 🎉", fmt.Sprintf("Spørsmålet ditt er vorte godkjent av <@%s>!\n\n**\"%s\"**\n\nDet er no tilgjengeleg for daglege spørsmål! ✨", approverID, question.Question), EmbedTypeSuccess)
	_, err = session.ChannelMessageSendEmbed(privateChannel.ID, embed)
	if err != nil {
		log.Printf("Failed to send approval notification to user: %v", err)
	}
}

// PostPendingBannedWordToRettingChannel posts a newly created banned word to the retting channel for approval.
func (s *ApprovalService) PostPendingBannedWordToRettingChannel(bannedWordID int64) {
	// Get the specific banned word by ID
	bannedWord, err := s.Bot.Database.GetBannedWordByID(int(bannedWordID))
	if err != nil {
		log.Printf("Failed to get banned word for retting channel posting: %v", err)
		return
	}

	if bannedWord == nil {
		log.Printf("No banned word found with ID %d", bannedWordID)
		return
	}

	channelID := s.Bot.Config.BannedWords.ApprovalChannelID
	if channelID == "" {
		log.Println("Retting channel is not configured")
		return
	}

	// Get the hammer user info
	hammerUser, err := s.Bot.Session.User(bannedWord.AuthorID)

	approvalEmbed := CreateApprovalEmbed(bannedWord.Word, "⏳ Opplysar-godkjenning: ventar\n⏳ Rettskrivar-godkjenning: ventar", hammerUser)

	message, err := s.Bot.Session.ChannelMessageSendEmbed(channelID, approvalEmbed)
	if err != nil {
		log.Printf("Failed to post to retting channel: %v", err)
		return
	}

	// Add reaction emoji
	s.Bot.Session.MessageReactionAdd(channelID, message.ID, "👍")

	// Update the database with the approval message ID
	err = s.Bot.Database.UpdateBannedWordApprovalMessageID(int(bannedWord.ID), message.ID)
	if err != nil {
		log.Printf("Failed to update approval message ID: %v", err)
	}
}

// PostBannedWordReport creates a forum post in the grammar channel for banned word discussion
// Returns the forum thread if a new one was created, or nil if referencing existing threads
func (s *ApprovalService) PostBannedWordReport(session *discordgo.Session, words []string, reporterID string, guildID string, originalChannelID string, originalMessageID string) *discordgo.Channel {
	if s.Bot.Config.Grammar.ChannelID == "" {
		log.Println("Grammar channel not configured")
		return nil
	}

	// Get reporter info
	reporter, err := session.User(reporterID)
	reporterName := "Unknown User"
	if err == nil {
		reporterName = reporter.Username
	}

	// For newly approved banned words, always create a forum thread
	// Check if any words already have forum threads (for logging purposes)
	var existingThreads []string

	for _, word := range words {
		isBanned, bannedWord, err := s.Bot.Database.IsBannedWord(word)
		if err != nil {
			log.Printf("Error checking if word '%s' is banned: %v", word, err)
			continue
		}

		if isBanned && bannedWord.ForumThreadID != nil && *bannedWord.ForumThreadID != "" {
			// Word already exists with a forum thread
			existingThreads = append(existingThreads, *bannedWord.ForumThreadID)
			log.Printf("Word '%s' already exists with forum thread %s", word, *bannedWord.ForumThreadID)
		}
	}

	// Always create forum threads for newly approved words
	newWords := words

	// Create forum post for new words - use just the word as title
	postTitle := strings.Join(newWords, ", ")
	if len(postTitle) > 100 {
		postTitle = postTitle[:97] + "..."
	}

	// Note: We no longer need guild ID since we simplified the forum message

	// Create forum post (thread in forum channel) with minimal initial message
	initialMessage := "🔨 Grammatikkdiskusjon"
	thread, err := session.ForumThreadStart(s.Bot.Config.Grammar.ChannelID, postTitle, 60, initialMessage)
	if err != nil {
		log.Printf("Failed to create forum post (may require approval): %v", err)
		// Forum post creation failed - likely requires manual approval
		// Return nil so words are still added to database without thread reference
		return nil
	}

	// Now send a proper embed as the main discussion starter
	var originalInfo string
	if originalChannelID != "" && originalMessageID != "" {
		// Create a Discord link for the original message
		originalInfo = fmt.Sprintf("[Hopp til opphavleg melding](https://discord.com/channels/%s/%s/%s)", guildID, originalChannelID, originalMessageID)
	} else if originalMessageID != "" {
		// We have message ID but not channel ID - show what we can
		originalInfo = fmt.Sprintf("Meldings-ID: `%s`", originalMessageID)
	} else {
		originalInfo = "Informasjon om opphavleg melding ikkje tilgjengeleg"
	}

	// Get additional reporter info for embed author (reuse existing reporter variable)
	var reporterAvatarURL string
	if err == nil {
		reporterAvatarURL = reporter.AvatarURL("")
	} else {
		reporterName = "Ukjend brukar"
		reporterAvatarURL = ""
	}

	// Create discussion embed
	discussionEmbed := NewEmbedBuilder().
		SetTitle("📝 Grammatikkdiskusjon: "+strings.Join(words, ", ")).
		SetDescription("Dette ordet/desse orda har vorte rapporterte som grammatisk feil.").
		SetColor(0xff6b35). // Orange color
		SetAuthor("Rapportert av "+reporterName, reporterAvatarURL).
		AddField("📍 Opphavleg melding", originalInfo, false).
		AddField("💡 Diskusjonsrettleiing", "• Forklar kvifor ordet er feil\n• Gje korrekte alternativ\n• Del relevante reglar eller kjelder", false).
		SetFooter("Ver snill og diskuter på ein konstruktiv måte", "").
		Build()

	// Send the embed to the thread
	_, err = session.ChannelMessageSendEmbed(thread.ID, discussionEmbed)
	if err != nil {
		log.Printf("Failed to send discussion embed to forum thread: %v", err)
		// Continue anyway, the thread was created successfully
	}

	log.Printf("Created forum post %s (%s) for banned word discussion", thread.Name, thread.ID)
	return thread
}

// NotifyUserRejection notifies the user that their question was rejected.
func (s *ApprovalService) NotifyUserRejection(session *discordgo.Session, question *database.Question, rejectorID string) {
	privateChannel, err := session.UserChannelCreate(question.AuthorID)
	if err != nil {
		log.Printf("Failed to create private channel for rejection notification: %v", err)
		return
	}

	rejector, err := session.User(rejectorID)
	var rejectorName string
	if err != nil {
		rejectorName = "ein opplysar"
	} else {
		rejectorName = rejector.Username
	}

	embed := CreateBotEmbed(session, "❌ Spørsmål avvist", fmt.Sprintf("Spørsmålet ditt har blitt avvist av %s.\n\n**\"%s\"**\n\nDu kan prøve å sende inn eit anna spørsmål som passar betre.", rejectorName, question.Question), EmbedTypeError)
	_, err = session.ChannelMessageSendEmbed(privateChannel.ID, embed)
	if err != nil {
		log.Printf("Failed to send rejection notification to user: %v", err)
	}
}
