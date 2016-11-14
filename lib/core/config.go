package i3gostatus

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/model"
	"os"
	"time"
)

func getConfigPath() string {
	// TODO: Support this XDG Basedirectory thing.
	return os.Getenv("HOME") + "/.config/i3gostatus/config.toml"
}

func findFastestPeriod(configTree *toml.TomlTree) time.Duration {
	res := configTree.Get("modules").([]interface{})
	var smallest int64
	var current int64

	smallest = model.DefaultPeriod

	for _, module := range res {
		moduleStr := module.(string)
		current = configTree.GetDefault(moduleStr+".period", int64(model.DefaultPeriod)).(int64)
		if current < smallest {
			smallest = current
		}
	}

	return time.Duration(smallest) * time.Millisecond
}

func registerEnabledModules(configTree *toml.TomlTree) {
	res := configTree.Get("modules").([]interface{})
	EnabledModules = make([]model.Module, 0, len(res))

	for _, module := range res {
		moduleStr := module.(string)

		if val, ok := AvailableModules[moduleStr]; ok {
			EnabledModules = append(EnabledModules, val)
		}
	}
}

// TODO: Default config
func loadConfig(path string) *toml.TomlTree {
	configTree, err := toml.LoadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	registerEnabledModules(configTree)

	return configTree
}
