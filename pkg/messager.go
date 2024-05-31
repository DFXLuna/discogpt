package discogpt

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// This package encapsulates some kind of output method
// to forward messages to

//go:generate mockgen -source ./messager.go -destination ./mock/messager.go
type GeneratorMessager interface {
	// Calls g.Generate inresponse to the trigger phrase
	Run(ctx context.Context, errch chan error)
}

type discordMessager struct {
	// Token string
	G               MessageGenerator
	S               *discordgo.Session
	Trigger         string
	AllowedChannels []string
}

func NewDiscordMessager(ctx context.Context, tok string, trigger string, allowedChannels []string, g MessageGenerator) (*discordMessager, error) {
	s, err := discordgo.New("Bot " + tok)
	if err != nil {
		return nil, err
	}
	return &discordMessager{
		S:               s,
		G:               g,
		AllowedChannels: allowedChannels,
		Trigger:         trigger + " ",
	}, nil
}

func (d *discordMessager) Run(ctx context.Context, errch chan error) {
	d.S.AddHandler(replyInjector(discordReply, d.Trigger, d.AllowedChannels, d.G))
	d.S.Identify.Intents = discordgo.IntentsGuildMessages
	// Open a websocket connection to Discord and begin listening.
	err := d.S.Open()
	if err != nil {
		errch <- err
	}
	defer d.S.Close()
	errch <- nil
	<-ctx.Done()
}

func replyInjector(f func(*discordgo.Session, *discordgo.MessageCreate, string, []string, MessageGenerator),
	trigger string, allowedChannels []string, g MessageGenerator) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		f(s, m, trigger, allowedChannels, g)
	}
}

func discordReply(s *discordgo.Session, m *discordgo.MessageCreate, trigger string, allowedChannels []string, g MessageGenerator) {
	// Ignore all messages created by the bot or in channels not on the allow list
	if m.Author.ID == s.State.User.ID || !slices.Contains(allowedChannels, m.ChannelID) {
		return
	}
	index := strings.Index(strings.ToLower(m.Content), strings.ToLower(trigger))
	if index != -1 {
		_ = s.ChannelTyping(m.ChannelID)
		reply, err := g.Generate(context.Background(), m.Content[index:], m.Author.Username)
		if err != nil {
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("error from generator %v", err))
			if err != nil {
				fmt.Printf("Error sending error message: %v", err)
			}
			fmt.Printf("error from generator %v", err)
			return
		}
		_, err = s.ChannelMessageSend(m.ChannelID, reply)
		if err != nil {
			fmt.Printf("Error sending message: %v", err)
		}
		return
	}
}
