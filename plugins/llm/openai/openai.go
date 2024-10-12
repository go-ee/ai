package openai

import (
	"context"
	"errors"
	"fabricng/api"
	"fabricng/core"
	"fmt"
	"github.com/samber/lo"
	"github.com/sashabaranov/go-openai"
	"io"
	"log/slog"
)

var pluginBase = core.PluginBase{Name: "OpenAI", Type: api.PluginTypeLLM}

func NewLLM() (ret *LLM) {
	return NewLLMClientCompatible("OpenAI", "https://core.openai.com/v1", nil)
}

func NewLLMClientCompatible(vendorName string, defaultBaseUrl string, configureCustom func() error) (ret *LLM) {
	ret = &LLM{PluginBase: pluginBase}

	if configureCustom == nil {
		configureCustom = ret.configure
	}
	return
}

type LLM struct {
	core.PluginBase
	ApiKey     string
	ApiBaseURL string
	ApiClient  *openai.Client
}

func (o *LLM) configure() (ret error) {
	config := openai.DefaultConfig(o.ApiKey)
	if o.ApiBaseURL != "" {
		config.BaseURL = o.ApiBaseURL
	}
	o.ApiClient = openai.NewClientWithConfig(config)
	return
}

func (o *LLM) ListModels() (ret []string, err error) {
	var models openai.ModelsList
	if models, err = o.ApiClient.ListModels(context.Background()); err != nil {
		return
	}

	model := models.Models
	for _, mod := range model {
		ret = append(ret, mod.ID)
	}
	return
}

func (o *LLM) ChatStream(
	ctx context.Context, msgs []*api.Message, opts *api.ChatOptions, channel chan string,
) (err error) {
	req := o.buildChatCompletionRequest(msgs, opts)
	req.Stream = true

	var stream *openai.ChatCompletionStream
	if stream, err = o.ApiClient.CreateChatCompletionStream(ctx, req); err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}

	defer stream.Close()

	for {
		var response openai.ChatCompletionStreamResponse
		if response, err = stream.Recv(); err == nil {
			if len(response.Choices) > 0 {
				channel <- response.Choices[0].Delta.Content
			} else {
				channel <- "\n"
				close(channel)
				break
			}
		} else if errors.Is(err, io.EOF) {
			channel <- "\n"
			close(channel)
			err = nil
			break
		} else if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			break
		}
	}
	return
}

func (o *LLM) Chat(ctx context.Context, msgs []*api.Message, opts *api.ChatOptions) (ret []*api.Message, err error) {
	req := o.buildChatCompletionRequest(msgs, opts)

	var resp openai.ChatCompletionResponse
	if resp, err = o.ApiClient.CreateChatCompletion(ctx, req); err != nil {
		return
	}
	if len(resp.Choices) > 0 {
		ret = append(ret, &api.Message{Content: resp.Choices[0].Message.Content})
		slog.Debug("SystemFingerprint: " + resp.SystemFingerprint)
	}
	return
}

func (o *LLM) buildChatCompletionRequest(
	msgs []*api.Message, opts *api.ChatOptions,
) (ret openai.ChatCompletionRequest) {
	messages := lo.Map(msgs, func(message *api.Message, _ int) openai.ChatCompletionMessage {
		return openai.ChatCompletionMessage{Role: message.Role, Content: message.Content}
	})

	if opts.Raw {
		ret = openai.ChatCompletionRequest{
			Model:    opts.Model,
			Messages: messages,
		}
	} else {
		if opts.Seed == 0 {
			ret = openai.ChatCompletionRequest{
				Model:            opts.Model,
				Temperature:      float32(opts.Temperature),
				TopP:             float32(opts.TopP),
				PresencePenalty:  float32(opts.PresencePenalty),
				FrequencyPenalty: float32(opts.FrequencyPenalty),
				Messages:         messages,
			}
		} else {
			ret = openai.ChatCompletionRequest{
				Model:            opts.Model,
				Temperature:      float32(opts.Temperature),
				TopP:             float32(opts.TopP),
				PresencePenalty:  float32(opts.PresencePenalty),
				FrequencyPenalty: float32(opts.FrequencyPenalty),
				Messages:         messages,
				Seed:             &opts.Seed,
			}
		}
	}
	return
}

func NewFactory() (ret *Factory) {
	ret = &Factory{PluginFactoryBase: core.PluginFactoryBase[LLM]{PluginBase: pluginBase}}
	return
}

type Factory struct {
	core.PluginFactoryBase[LLM]
}

func (o *Factory) Setup(instanceName string) (settings map[string]string, err error) {
	return
}

func (o *Factory) Create(instanceName string, settings map[string]string) (ret api.Plugin, err error) {
	item := NewLLM()
	item.ApiKey = settings["ApiKey"]
	item.ApiBaseURL = settings["ApiBaseUrl"]
	ret = item
	return
}
