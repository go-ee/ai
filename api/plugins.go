package api

import "context"

const DefaultPluginInstance = "default"

type PluginType string

// Enum values for PluginType
const (
	PluginTypeChatter     PluginType = "Chatter"
	PluginTypeLLM         PluginType = "LLM"
	PluginTypeInput       PluginType = "Inputs"
	PluginTypeOutput      PluginType = "Output"
	PluginTypeTransformer PluginType = "Transformer"
	PluginTypeMeta        PluginType = "Meta"
)

// Plugin defines the interface for all plugins
type Plugin interface {
	GetName() string
	GetType() PluginType // Type Models, Input, Transformer, Output
}

type Models interface {
	Plugin
	ListModels() ([]string, error)
	ChatStream(context.Context, []*Message, *ChatOptions, chan string) error
	Chat(context.Context, []*Message, *ChatOptions) ([]*Message, error)
}

type Chatter interface {
	Plugin
	Chat(context.Context, []*Message, *ChatOptions) ([]*Message, error)
}

type Input interface {
	Plugin
	GetMessages() ([]*Message, error)
}

type Transformer interface {
	Plugin
	Transform([]*Message) ([]*Message, error)
}

type Output interface {
	Plugin
	Output([]*Message) error
}

// PluginFactory defines the interface for all plugin factories
type PluginFactory interface {
	Plugin
	Setup(instanceName string) (settings map[string]string, err error)
	Create(instanceName string, settings map[string]string) (Plugin, error)
}

type WorkflowBuild interface {
	Build() (*Workflow, error)
}

type Configurator interface {
	Load(instanceName string, plugin Plugin) (PluginConfiguration, error)
	Store(PluginConfiguration) error

	LoadAll() ([]PluginConfiguration, error)
	StoreAll([]PluginConfiguration) error
}

type PluginConfiguration interface {
	Plugin
	GetInstanceName() string
	GetSettings() map[string]string
}
