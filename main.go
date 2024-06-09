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
	"go.uber.org/zap"
)

const (
	configFile = "config.yaml"
)

func main() {
	config, reqMods, genMods, log, err := SetupConfig(configFile)
	if err != nil {
		panic(err)
	}

	mg, err := discogpt.NewOpenAIGenerator(config.OAIHost, config.OAIModel, config.OAISystemPrompt, log, reqMods, genMods)
	if err != nil {
		panic(err)
	}
	if config.Mode == string(discogpt.DiscordMode) {
		dm, err := discogpt.NewDiscordMessager(context.Background(), config.BotToken, config.Trigger, discogpt.GetAllowedChannels(config), mg, log)
		if err != nil {
			panic(err)
		}

		err = Run(mg, dm, log)
		if err != nil {
			panic(err)
		}

	} else if config.Mode == string(discogpt.IOMode) {
		im := discogpt.NewIOMessager(os.Stdin, os.Stdout, config.Trigger, mg, "user", log)

		err = Run(mg, im, log)
		if err != nil {
			panic(err)
		}
	}
	panic(fmt.Errorf("bad mode %v", config.Mode))

}

func SetupConfig(configFilePath string) (
	discogpt.Config, []discogpt.HTTPRequestModifier, []discogpt.GenerationRequestModifier,
	*zap.SugaredLogger, error) {
	log, err := discogpt.NewProdLogger()
	if err != nil {
		return discogpt.Config{}, []discogpt.HTTPRequestModifier{}, []discogpt.GenerationRequestModifier{}, nil, err
	}

	config, err := discogpt.GenerateConfig(configFilePath, log)
	if err != nil {
		return discogpt.Config{}, []discogpt.HTTPRequestModifier{}, []discogpt.GenerationRequestModifier{}, nil, err
	}
	if config.Debug {
		_ = log.Sync() // Clear original logger before replacing it
		log, err = discogpt.NewDebugLogger()
		if err != nil {
			return discogpt.Config{}, []discogpt.HTTPRequestModifier{}, []discogpt.GenerationRequestModifier{}, nil, err
		}
	}

	c, err := conf.String(&config)
	if err != nil {
		return discogpt.Config{}, []discogpt.HTTPRequestModifier{}, []discogpt.GenerationRequestModifier{}, nil, err
	}
	log.Infof("Using config %s\n", c)
	reqMods := []discogpt.HTTPRequestModifier{}
	if config.OAIToken != "" {
		reqMods = append(reqMods, func(req *http.Request) error {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.OAIToken))
			return nil
		})
	}

	genMods := []discogpt.GenerationRequestModifier{}
	if config.ChromaURL != "" {
		log.Infof("Adding chroma mod\n")
		chromaMod, err := discogpt.NewChromaMod(config.ChromaURL, config.ChromaTEIURL, config.ChromaCollectionName, log)
		if err != nil {
			return discogpt.Config{}, []discogpt.HTTPRequestModifier{}, []discogpt.GenerationRequestModifier{}, nil, err
		}
		genMods = append(genMods, chromaMod)
	}
	log.Debugf("len genMods: %d", len(genMods))
	return config, reqMods, genMods, log, nil
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
