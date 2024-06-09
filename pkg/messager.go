package discogpt

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"slices"
	"strings"
	"time"

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
	Log             Logger
}

func NewDiscordMessager(ctx context.Context, tok string, trigger string, allowedChannels []string, g MessageGenerator, log Logger) (*discordMessager, error) {
	s, err := discordgo.New("Bot " + tok)
	if err != nil {
		return nil, err
	}
	return &discordMessager{
		S:               s,
		G:               g,
		AllowedChannels: allowedChannels,
		Trigger:         trigger + " ",
		Log:             log,
	}, nil
}

func (d *discordMessager) Run(ctx context.Context, errch chan error) {
	d.S.AddHandler(replyInjector(discordReply, d.Trigger, d.AllowedChannels, d.G, d.Log))
	d.S.Identify.Intents = discordgo.IntentsGuildMessages
	// Open a websocket connection to Discord and begin listening.
	err := d.S.Open()
	if err != nil {
		errch <- err
	}
	defer d.S.Close()
	errch <- nil
	d.Log.Debugf("Starting discord messager: Trigger: %v, AllowedChannels: %v", d.Trigger, d.AllowedChannels)
	<-ctx.Done()
}

func replyInjector(f func(*discordgo.Session, *discordgo.MessageCreate, string, []string, MessageGenerator, Logger),
	trigger string, allowedChannels []string, g MessageGenerator, l Logger) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		f(s, m, trigger, allowedChannels, g, l)
	}
}

func discordReply(s *discordgo.Session, m *discordgo.MessageCreate, trigger string, allowedChannels []string, g MessageGenerator, l Logger) {
	// Ignore all messages created by the bot or in channels not on the allow list
	if m.Author.ID == s.State.User.ID || !slices.Contains(allowedChannels, m.ChannelID) {
		return
	}
	index := strings.Index(strings.ToLower(m.Content), strings.ToLower(trigger))
	if index != -1 {
		_ = s.ChannelTyping(m.ChannelID)
		l.Debugf("Triggered by %s", m.Author.Username)
		reply, err := g.Generate(context.Background(), m.Content[index:], m.Author.Username)
		if err != nil {
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("error from generator %v", err))
			if err != nil {
				l.Errorf("Error sending error message: %v", err)
			}
			l.Errorf("error from generator %v", err)
			return
		}
		_, err = s.ChannelMessageSend(m.ChannelID, reply)
		if err != nil {
			l.Errorf("Error sending message: %v", err)
		}
		return
	}
}

// This messager response to message on io(stdio & etc)
// it responds to a line starting with trigger and ending with \n
// Caller is responsible for closing R and W
type ioMessager struct {
	R       io.Reader
	W       io.Writer
	G       MessageGenerator
	User    string // name to refer to the current user
	Trigger string
	Log     Logger
}

func NewIOMessager(r io.Reader, w io.Writer, trigger string, g MessageGenerator, user string, log Logger) *ioMessager {
	return &ioMessager{
		R:       r,
		W:       w,
		G:       g,
		User:    user,
		Trigger: trigger,
		Log:     log,
	}
}

func (i *ioMessager) Run(ctx context.Context, errch chan error) {
	_, _ = i.W.Write([]byte(fmt.Sprintf("[%v]: type a message...\n", time.Now().Local().String())))
	scanner := bufio.NewScanner(i.R)
	for scanner.Scan() {
		line := scanner.Text()
		index := strings.Index(strings.ToLower(line), strings.ToLower(i.Trigger))
		if index != -1 {
			_, _ = i.W.Write([]byte(fmt.Sprintf("[%v]: replying...\n", time.Now().Local().String())))
			prompt, err := i.G.Generate(ctx, line[index:], i.User)
			if err != nil {
				errch <- err
				return
			}
			_, _ = i.W.Write([]byte(fmt.Sprintf("[%v]: %v\n", time.Now().Local().String(), prompt)))
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
	errch <- scanner.Err()
}
