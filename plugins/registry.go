package plugins

import (
	"fabricng/api"
	"fabricng/plugins/input/youtube"
	"fabricng/plugins/llm/openai"
	"fmt"
	"github.com/elliotchance/orderedmap"
)

var defaultRegistry *PluginRegistry

func GetDefaultPluginRegistry() *PluginRegistry {
	if defaultRegistry == nil {
		defaultRegistry = NewPluginRegistry()

		//registering LLM plugins
		defaultRegistry.AddPluginFactory(openai.NewFactory())

		//registering Tools plugins
		defaultRegistry.AddPluginFactory(youtube.NewFactory())
	}
	return defaultRegistry
}

func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{factories: orderedmap.NewOrderedMap()}
}

// PluginRegistry to store all enabled factories grouped by type
type PluginRegistry struct {
	factories *orderedmap.OrderedMap
}

// AddPluginFactory Add a plugin to the registry
func (pr *PluginRegistry) AddPluginFactory(plugin api.PluginFactory) {
	pr.factories.Set(plugin.GetName(), plugin)
}

func (pr *PluginRegistry) PrintPlugins() (ret int) {
	counter := 1
	lastType := api.PluginTypeMeta
	for _, name := range pr.factories.Keys() {
		plugin := pr.getFactoryByKey(name)
		if plugin.GetType() != lastType {
			fmt.Printf("\n\n%v Plugins:\n\n", plugin.GetType())
			lastType = plugin.GetType()
		}
		fmt.Printf("%d. %v\n", counter, name)
		counter++
	}
	return
}

func (pr *PluginRegistry) GetPluginByIndex(index int) (ret api.PluginFactory, err error) {
	names := pr.factories.Keys()
	if len(names) < index {
		err = fmt.Errorf("there is no plugin with the index %v", index)
		return
	}
	ret = pr.getFactoryByKey(names[index-1])
	return
}

func (pr *PluginRegistry) GetFactoryByName(name interface{}) (ret api.PluginFactory, err error) {
	if plugin, ok := pr.factories.Get(name); ok {
		ret = plugin.(api.PluginFactory)
	} else {
		err = fmt.Errorf("plugin %v not found", name)
	}
	return
}

func (pr *PluginRegistry) GetPluginsAll() (ret []api.Plugin) {
	for _, name := range pr.factories.Keys() {
		plugin := pr.getFactoryByKey(name)
		ret = append(ret, plugin)
	}
	return
}

func (pr *PluginRegistry) getFactoryByKey(name interface{}) api.PluginFactory {
	plugin, _ := pr.factories.Get(name)
	return plugin.(api.PluginFactory)
}
