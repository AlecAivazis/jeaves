package bot

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// CommandContext holds the contextual information for a message that we receive
type CommandContext struct {
	context.Context
	GuildID   string
	ChannelID string
	Message   *discordgo.Message
}

// CommandHandler handles the parsing and dispatching of commands for Jeeves
func (b *JeevesBot) CommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	// only look at commands
	if message.Content[0] != '!' {
		return
	}

	// since the message is presumably text, we care about words, not letters
	words := strings.SplitN(message.Content[1:], " ", 2)
	command := words[0]

	// construct the context object
	ctx := &CommandContext{
		GuildID:   message.GuildID,
		ChannelID: message.ChannelID,
		Context:   context.Background(),
		Message:   message.Message,
	}

	var err error
	// check the command against our known strings
	switch command {
	case CommandAssignBankChannel:
		err = b.InitializeBankChannel(ctx)
	case CommandDeposit:
		// there could be trailing white space around the word
		trimmed := strings.Trim(words[1], ", ")

		err = b.DepositItems(ctx, strings.Split(trimmed, ","))
	case CommandWithdraw:
		// there could be trailing white space around the word
		trimmed := strings.Trim(words[1], ", ")

		err = b.WithdrawItems(ctx, strings.Split(trimmed, ","))
	}
	// if the command failed
	if err != nil {
		// send the error to the channel we received the message on
		b.ReportError(message.ChannelID, err)
		return
	}
}
