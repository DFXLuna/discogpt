package discogpt

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ardanlabs/conf/v3"
	"github.com/ardanlabs/conf/v3/yaml"
)

type Config struct {
	OAIHost         string `yaml:"OAIHost" conf:"env:DISCOGPT_OAI_HOST"`
	OAIToken        string `yaml:"OAIToken" conf:"env:DISCOGPT_OAI_TOKEN,mask"`
	OAIModel        string `yaml:"OAIModel" conf:"env:DISCOGPT_OAI_MODEL"`
	OAISystemPrompt string `yaml:"OAISystemPrompt" conf:"env:DISCOGPT_OAI_SYSTEM_PROMPT"`

	BotToken        string `yaml:"BotToken" conf:"env:DISCOGPT_BOT_TOKEN,mask"`
	AllowedChannels string `yaml:"AllowedChannels" conf:"env:DISCOGPT_ALLOWED_CHANNELS"` //comma separated list of channel IDs for bot to operate in
	Trigger         string `yaml:"Trigger" conf:"env:DISCOGPT_TRIGGER"`

	Mode string `yaml:"Mode" conf:"DISCOGPT_MODE"` // This is used to select between the discord messager and the ioMessager for local testing

	ChromaURL            string `yaml:"ChromaURL" conf:"env:DISCOGPT_CHROMA_URL"`        //URL for your chroma server
	ChromaTEIURL         string `yaml:"ChromaTEIURL" conf:"env:DISCOGPT_CHROMA_TEI_URL"` //URL for your hugging face TEI server
	ChromaCollectionName string `yaml:"ChromaCollectionName" conf:"env:DISCOGPT_CHROMA_COLLECTION_NAME"`

	Debug bool `yaml:"Debug" conf:"env:DISCOGPT_DEBUG"` //turns on the debug logger
}

func GetAllowedChannels(c Config) []string {
	return strings.Split(c.AllowedChannels, ",")
}

func GenerateConfig(configFile string, log Logger) (Config, error) {
	path, err := filepath.Abs(configFile)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if file, err := os.ReadFile(path); err == nil {
		log.Infof("Using config from %s\n", configFile)
		help, err := conf.Parse("", &cfg, yaml.WithData(file))
		if err != nil {
			if errors.Is(err, conf.ErrHelpWanted) {
				fmt.Println(help)
				return Config{}, err
			}
			return Config{}, err
		}
	} else {
		log.Infof("Using config env")
		help, err := conf.Parse("", &cfg)
		if err != nil {
			if errors.Is(err, conf.ErrHelpWanted) {
				fmt.Println(help)
			}
			return Config{}, err
		}
	}
	return cfg, nil
}
