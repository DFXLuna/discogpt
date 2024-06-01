package main

import (
	"context"
	discogpt "egrant/disco-gpt/pkg"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ardanlabs/conf/v3"
)

const (
	configFile = "config.yaml"
)

func main() {
	log, err := discogpt.NewLogger()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = log.Sync()
	}()

	config, mods, err := SetupConfig(configFile, log)
	if err != nil {
		panic(err)
	}

	mg, err := discogpt.NewOpenAIGenerator(config.OAIHost, config.OAIModel, config.OAISystemPrompt, log, mods...)
	if err != nil {
		panic(err)
	}

	dm, err := discogpt.NewDiscordMessager(context.Background(), config.BotToken, config.Trigger, discogpt.GetAllowedChannels(config), mg, log)
	if err != nil {
		panic(err)
	}

	err = Run(mg, dm, log)
	if err != nil {
		panic(err)
	}
}

func SetupConfig(configFilePath string, log discogpt.Logger) (discogpt.Config, []discogpt.RequestModifier, error) {
	config, err := discogpt.GenerateConfig(configFilePath, log)
	if err != nil {
		return discogpt.Config{}, []discogpt.RequestModifier{}, err
	}
	c, err := conf.String(&config)
	if err != nil {
		return discogpt.Config{}, []discogpt.RequestModifier{}, err
	}
	log.Infof("Using config %s\n", c)
	mods := []discogpt.RequestModifier{}
	if config.OAIToken != "" {
		mods = append(mods, func(req *http.Request) error {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.OAIToken))
			return nil
		})
	}
	return config, mods, nil
}

func Run(g discogpt.MessageGenerator, m discogpt.GeneratorMessager, log discogpt.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errch := make(chan error)
	go m.Run(ctx, errch)
	err := <-errch
	if err != nil {
		return err
	}
	// Wait here until CTRL-C or other term signal is received.
	log.Infof("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	return nil
}
