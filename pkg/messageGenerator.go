package discogpt

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type MessagerMode string

const (
	DiscordMode MessagerMode = "Discord"
	IOMode      MessagerMode = "IO"
)

// This file encapsulates an inference providing backend
// Backends must provide an object that fulfils the interface

//go:generate mockgen -source ./messageGenerator.go -destination ./mock/messageGenerator.go
type MessageGenerator interface {
	Generate(ctx context.Context, prompt string, user string) (string, error)
}

var (
	ErrOAIHTTP = errors.New("error code from oai server")

	oaiCompletionsEndpoint = "/v1/chat/completions"
	oaiInstruct            = "instruct"
	oaiUser                = "user"
	oaiSystem              = "system"
)

// fulfills MessageGenerator for OpenAI compatible APIs like textgen
// in instruct mode
type oaiGenerator struct {
	CompletionsURL      string                      // will be constructed by parsing with url & appending /v1/chat/completions
	SystemPrompt        string                      // gets inserted before the provided messages as a prompt with role System
	RequestModifiers    []HTTPRequestModifier       // Will be called on the http request made for generation, mostly used for auth
	GenerationModifiers []GenerationRequestModifier // Will be called on the oai generation request prior to it being sent to the API
	Model               string                      //The model to include in the OAI chat completions call
	Log                 Logger
}

type oaiCompletionsReq struct {
	Model    string       `json:"model"`
	Messages []oaiMessage `json:"messages"`
	Mode     string       `json:"mode"`
}

type oaiCompletionsResp struct {
	ID      string      `json:"id"`
	Choices []oaiChoice `json:"choices"`
	Usage   oaiUsage    `json:"usage"`
}

type oaiChoice struct {
	Index        int        `json:"index"`
	FinishReason string     `json:"finish_reason"`
	Message      oaiMessage `json:"message"`
}

type oaiUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
}

type oaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type HTTPRequestModifier func(*http.Request) error
type GenerationRequestModifier func(*oaiCompletionsReq) error

func NewOpenAIGenerator(baseURL string, model string, promptPrefix string, log Logger,
	reqMods []HTTPRequestModifier, genMods []GenerationRequestModifier) (*oaiGenerator, error) {
	completions, err := url.JoinPath(baseURL, oaiCompletionsEndpoint)
	if err != nil {
		return nil, err
	}
	return &oaiGenerator{
		CompletionsURL:      completions,
		SystemPrompt:        promptPrefix + "\n",
		RequestModifiers:    reqMods,
		GenerationModifiers: genMods,
		Model:               model,
		Log:                 log,
	}, nil
}

func (o *oaiGenerator) Generate(ctx context.Context, prompt string, user string) (string, error) {
	o.Log.Debugf("Generating for %v", user)
	completionReq := oaiCompletionsReq{
		Model: o.Model,
		Mode:  oaiInstruct,
		Messages: []oaiMessage{
			{
				Role:    oaiSystem,
				Content: o.SystemPrompt + "\n[The user that sent the following message is " + user + "]",
			},
			{
				Role:    oaiUser,
				Content: prompt,
			},
		},
	}
	for _, mod := range o.GenerationModifiers {
		err := mod(&completionReq)
		if err != nil {
			return "", err
		}
	}
	o.Log.Debugf("modified json: %+v\n", completionReq)
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(completionReq)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.CompletionsURL, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	for _, mod := range o.RequestModifiers {
		err = mod(req)
		if err != nil {
			return "", err
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 100 || resp.StatusCode > 299 {
		return "", fmt.Errorf("%w: %s", ErrOAIHTTP, resp.Status)
	}
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}
	var respDecode oaiCompletionsResp
	err = json.Unmarshal(bs, &respDecode)
	if err != nil {
		return "", nil
	}
	return respDecode.Choices[0].Message.Content, nil
}
