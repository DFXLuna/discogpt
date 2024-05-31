# Discogpt
A little discord bot to connect your OpenAI compliant LLM API to a Discord server.

![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/dfxluna/discogpt/total)
![GitHub License](https://img.shields.io/github/license/dfxluna/discogpt)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/49ec0462c97949edbd8719d813f415d8)](https://app.codacy.com/gh/DFXLuna/discogpt/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)

## Requirements
- [A discord bot token](https://discord.com/developers/applications)
    - Create an application
    - Click bot and enable `Message Content Intent`
    - Copy the bot token into your config.yaml
- An OpenAI compatible LLM Inference API
    - [OpenAI API](https://openai.com/api/)
    - [Cloudflare AI Workers](https://developers.cloudflare.com/workers-ai/)
    - [Text-Generation-Webui](https://github.com/oobabooga/text-generation-webui)
- An API key for your OAI API if applicable

## Configuration
Discogpt uses a [config.yaml](./example-config.yaml) file to provide application configuration. See the example [here](./example-config.yaml). Table of configurations is here

The easiest way to use Discogpt is as a container. Write a config.yaml and mount it into a container

### Compose
```yaml
version: "3"

services:
  discogpt:
    image: dfxluna/discogpt:latest
    volumes:
      - ./config.yaml:/discogpt/config.yaml
    restart: "unless-stopped"
```

### Command line 
```sh
docker run -d -v ./config.yaml:/discogpt/config.yaml dfxluna/discogpt:latest
```

### Config.yaml fields
See [example-config.yaml](./example-config.yaml).

| Field | Comment | Example |
| ----- | ------- | ------- |
|OAIHost| The **base url** of your OpenAI API host.| `https://api.cloudflare.com/client/v4/accounts/{your_account_id}/ai`|
|OAIToken| If applicable, a bearer token to be provided with requests. Leave empty if not used.| Service dependant.|
|OAISystemPrompt|A prompt to include from user "System" as a message before the user's prompt.| `[You are ChadBot. The life of the party. ]` |
|OAIModel| Used to specify which model to use if multiple are available. Service dependant. If empty, either the service will choose (text-generation-webui) or can error. |`@hf/mistral/mistral-7b-instruct-v0.2`|
|BotToken| Your discord bot token | 
|AllowedChannels| A comma delimited list of [Discord Channel IDs](https://support.discord.com/hc/en-us/articles/206346498-Where-can-I-find-my-User-Server-Message-ID#h_01HRSTXPS5FMK2A5SMVSX4JW4E) for the bot to operate in.| `1137824512383429025,976152812312351829`
|Trigger| The case insensitive phrase that will trigger your bot. A space is automatically inserted after your trigger phrase. | `Hey ChadBot,`

## Development Reqs
- Go 1.22+
- golangci-lint

## Contributions
Feel free to file PRs.

## Support
Feel free to file an issue if something is broken or missing a feature.
A small donation is a large motivator for feature requests: [![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/A0A8GTT67)