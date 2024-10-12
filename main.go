package main

import (
	"fabricng/api"
	"fabricng/cli"
	"fabricng/core"
	"fabricng/plugins"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	registry := plugins.GetDefaultPluginRegistry()

	configurator := core.ConfiguratorEnv{
		FileName: filepath.Join(homedir, ".config/fabric/.env"),
		Plugins:  registry.GetPluginsAll(),
	}

	llmPlugin, err := registry.GetByName("OpenAI")
	if err != nil {
		fmt.Println(err)
		return
	}
	config, err := configurator.Load(api.DefaultPluginInstance, llmPlugin)
	if err != nil {
		fmt.Println(err)
		return
	}

	llm, err := llmPlugin.Create(config.GetInstanceName(), config.GetSettings())
	if err != nil {
		fmt.Println(err)
		return
	}
	println(llm.GetName())

	registry.PrintPlugins()
	return
}

func processWorkflowCall() {
	builder := cli.WorkflowBuilderCLI{Args: os.Args}
	workflow, err := builder.Build()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = workflow.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
