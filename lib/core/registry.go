package i3gostatus

// Every single module package must be imported and added to the registry!
import (
	"github.com/rumpelsepp/i3gostatus/lib/model"
	"github.com/rumpelsepp/i3gostatus/lib/modules/datetime"
	"github.com/rumpelsepp/i3gostatus/lib/modules/static"
	"github.com/rumpelsepp/i3gostatus/lib/modules/temperature"
)

var AvailableModules map[string]model.Module
var EnabledModules []model.Module

func init() {
	AvailableModules = make(map[string]model.Module)

	// Add all available modules here!
	AvailableModules["datetime"] = &datetime.Config{}
	AvailableModules["static"] = &static.Config{}
	AvailableModules["temperature"] = &temperature.Config{}
}
