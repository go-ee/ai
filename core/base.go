package core

import (
	"fabricng/api"
	"fmt"
)

type PluginBase struct {
	Name string
	Type api.PluginType
}

func (o PluginBase) GetName() string {
	return o.Name
}

func (o PluginBase) GetType() api.PluginType {
	return o.Type
}

type PluginFactoryBase[T api.Plugin] struct {
	PluginBase
}

type PluginConfiguration struct {
	PluginBase
	InstanceName string
	Settings     map[string]string
}

func (o PluginConfiguration) GetInstanceName() string {
	return o.InstanceName
}

func (o PluginConfiguration) GetSettings() map[string]string {
	return o.Settings
}

type ConfiguratorBase struct {
	pluginConfigs []api.PluginConfiguration
}

func (o *ConfiguratorBase) Load(instanceName string, plugin api.Plugin) (ret api.PluginConfiguration, err error) {
	for _, config := range o.pluginConfigs {
		if config.GetName() == plugin.GetName() &&
			config.GetType() == plugin.GetType() &&
			config.GetInstanceName() == instanceName {

			ret = config
			break
		}
	}

	if ret == nil {
		err = fmt.Errorf("plugin configuration not found")
	}
	return
}

func (o *ConfiguratorBase) Store(pluginConfig api.PluginConfiguration) (err error) {
	replaced := false
	for i, config := range o.pluginConfigs {
		if config.GetName() == pluginConfig.GetName() &&
			config.GetType() == pluginConfig.GetType() &&
			config.GetInstanceName() == pluginConfig.GetInstanceName() {

			o.pluginConfigs[i] = pluginConfig
			replaced = true
			break
		}
	}

	if !replaced {
		o.pluginConfigs = append(o.pluginConfigs, pluginConfig)
	}
	return
}

func (o *ConfiguratorBase) LoadAll() (ret []api.PluginConfiguration, err error) {
	ret = o.pluginConfigs
	return
}

func (o *ConfiguratorBase) StoreAll(pluginConfigs []api.PluginConfiguration) (err error) {
	o.pluginConfigs = pluginConfigs
	return
}
