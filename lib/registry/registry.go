package registry

// Every single module package must be imported and added to the registry!
import (
	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/model"
	"github.com/rumpelsepp/i3gostatus/lib/modules/backlight"
	"github.com/rumpelsepp/i3gostatus/lib/modules/datetime"
	"github.com/rumpelsepp/i3gostatus/lib/modules/disk"
	"github.com/rumpelsepp/i3gostatus/lib/modules/load"
	"github.com/rumpelsepp/i3gostatus/lib/modules/static"
	"github.com/rumpelsepp/i3gostatus/lib/modules/temperature"
)

var availableModules map[string]model.Module

func init() {
	availableModules = make(map[string]model.Module)

	// Add all available modules here!
	availableModules["backlight"] = &backlight.Config{}
	availableModules["datetime"] = &datetime.Config{}
	availableModules["load"] = &load.Config{}
	availableModules["disk"] = &disk.Config{}
	availableModules["static"] = &static.Config{}
	availableModules["temperature"] = &temperature.Config{}
}

func Initialize(configTree *toml.TomlTree) []model.Module {
	configuredModules := configTree.Get("modules").([]interface{})
	enabledModules := make([]model.Module, 0, len(configuredModules))

	for _, module := range configuredModules {
		moduleStr := module.(string)

		if val, ok := availableModules[moduleStr]; ok {
			enabledModules = append(enabledModules, val)
		}
	}

	return enabledModules
}
