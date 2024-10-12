package core

import (
	"bufio"
	"fabricng/api"
	"os"
	"strings"
)

type ConfiguratorEnv struct {
	ConfiguratorBase
	FileName string
	Plugins  []api.Plugin

	envFileLines []*EnvFileLine
}

func (o *ConfiguratorEnv) Load(instanceName string, plugin api.Plugin) (ret api.PluginConfiguration, err error) {
	if o.envFileLines == nil {
		_, err = o.LoadAll()
	}
	ret, err = o.ConfiguratorBase.Load(instanceName, plugin)
	return
}

func (o *ConfiguratorEnv) LoadAll() (ret []api.PluginConfiguration, err error) {
	if o.envFileLines == nil {
		if err = o.loadEnvFile(); err != nil {
			return
		}

		o.derivePluginConfigurations()
	}
	return o.ConfiguratorBase.LoadAll()
}

func (o *ConfiguratorEnv) derivePluginConfigurations() {
	for _, plugin := range o.Plugins {
		pluginEnvVariablePrefix := BuildEnvVariablePrefix(plugin.GetName())
		var pluginConfiguration *PluginConfiguration

		for _, line := range o.envFileLines {
			if strings.HasPrefix(line.Key, pluginEnvVariablePrefix) {
				//TODO for now only default instances, we need an extra field for instances
				if pluginConfiguration == nil {
					pluginConfiguration = &PluginConfiguration{
						PluginBase: PluginBase{
							Name: plugin.GetName(),
							Type: plugin.GetType(),
						},
						InstanceName: api.DefaultPluginInstance,
						Settings:     map[string]string{},
					}
					o.pluginConfigs = append(o.pluginConfigs, pluginConfiguration)
				}
				key := strings.TrimPrefix(line.Key, pluginEnvVariablePrefix)
				pluginConfiguration.Settings[SnakeCaseToCamelcase(key)] = line.Value
			}
		}
	}
}

func (o *ConfiguratorEnv) loadEnvFile() (err error) {
	var file *os.File
	if file, err = os.Open(o.FileName); err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText := scanner.Text()
		line := &EnvFileLine{Line: lineText}
		o.envFileLines = append(o.envFileLines, line)

		if strings.TrimSpace(lineText) == "" || strings.HasPrefix(lineText, "#") {
			continue
		}

		parts := strings.SplitN(lineText, "=", 2)
		if len(parts) != 2 {
			continue
		}

		line.Key = strings.TrimSpace(parts[0])
		valuePart := parts[1]
		valueParts := strings.SplitN(valuePart, "#", 2)
		line.Value = strings.TrimSpace(valueParts[0])

		if len(valueParts) > 1 {
			line.Comment = strings.TrimSpace(valueParts[1])
		}

		if _, exists := os.LookupEnv(line.Key); exists {
			line.ValueExternal = os.Getenv(line.Key)
		}
	}

	if err = scanner.Err(); err != nil {
		return
	}
	return
}

type EnvFileLine struct {
	Line          string
	Key           string
	Value         string
	ValueExternal string
	Comment       string
}

func BuildEnvVariablePrefix(name string) (ret string) {
	ret = BuildEnvVariable(name)
	if ret != "" {
		ret += "_"
	}
	return
}

func BuildEnvVariable(name string) string {
	name = strings.TrimSpace(name)
	return strings.ReplaceAll(strings.ToUpper(name), " ", "_")
}

func SnakeCaseToCamelcase(key string) string {
	parts := strings.Split(strings.ToLower(key), "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}
