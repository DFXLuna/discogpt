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
	config, err := discogpt.GenerateConfig(configFile)
	if err != nil {
		panic(err)
	}
	c, err := conf.String(&config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Using config %s\n", c)
	mods := []discogpt.RequestModifier{}
	if config.OAIToken != "" {
		mods = append(mods, func(req *http.Request) error {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.OAIToken))
			return nil
		})
	}

	mg, err := discogpt.NewOpenAIGenerator(config.OAIHost, config.OAIModel, config.OAISystemPrompt, mods...)
	if err != nil {
		panic(err)
	}

	dm, err := discogpt.NewDiscordMessager(context.Background(), config.BotToken, config.Trigger, discogpt.GetAllowedChannels(config), mg)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errch := make(chan error)
	go dm.Run(ctx, errch)
	err = <-errch
	if err != nil {
		panic(err)
	}
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

// token := os.Getenv(envStr)
// if token == "" {
// 	fmt.Println("ERROR: No token")
// 	return
// }

// apiHost := os.Getenv(hostStr)
// if apiHost == "" {
// 	fmt.Println("ERROR: No token")
// 	return
// }

// oaiToken := os.Getenv(oaiStr)
// if apiHost == "" {
// 	fmt.Println("ERROR: No token")
// 	return
// }
