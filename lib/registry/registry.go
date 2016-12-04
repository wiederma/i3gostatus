package registry

// Every single module package must be imported and added to the registry!
import (
	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/config"
	"github.com/rumpelsepp/i3gostatus/lib/model"
	"github.com/rumpelsepp/i3gostatus/lib/modules/backlight"
	"github.com/rumpelsepp/i3gostatus/lib/modules/cpu"
	"github.com/rumpelsepp/i3gostatus/lib/modules/datetime"
	"github.com/rumpelsepp/i3gostatus/lib/modules/disk"
	"github.com/rumpelsepp/i3gostatus/lib/modules/load"
	"github.com/rumpelsepp/i3gostatus/lib/modules/pulseaudio"
	"github.com/rumpelsepp/i3gostatus/lib/modules/static"
	"github.com/rumpelsepp/i3gostatus/lib/modules/syncthing"
	"github.com/rumpelsepp/i3gostatus/lib/modules/temperature"
	"github.com/rumpelsepp/i3gostatus/lib/modules/xkblayout"
)

var availableModules map[string]model.Module

func init() {
	availableModules = make(map[string]model.Module)

	// Add all available modules here!
	availableModules["backlight"] = &backlight.Config{}
	availableModules["cpu"] = &cpu.Config{}
	availableModules["datetime"] = &datetime.Config{}
	availableModules["load"] = &load.Config{}
	availableModules["disk"] = &disk.Config{}
	availableModules["pulseaudio"] = &pulseaudio.Config{}
	availableModules["static"] = &static.Config{}
	availableModules["syncthing"] = &syncthing.Config{}
	availableModules["temperature"] = &temperature.Config{}
	availableModules["xkblayout"] = &xkblayout.Config{}
}

func Initialize(configTree *toml.TomlTree) []model.Module {
	configuredModules := config.GetStringList(configTree, "modules", []string{})
	enabledModules := make([]model.Module, len(configuredModules))

	for i, module := range configuredModules {
		if val, ok := availableModules[module]; ok {
			enabledModules[i] = val
		}
	}

	return enabledModules
}
