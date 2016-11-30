package i3gostatus

import (
	"flag"
	"github.com/rumpelsepp/i3gostatus/lib/config"
)

type runtimeOptions struct {
	configPath  string
	clickEvents bool
}

func ParseOptions() *runtimeOptions {
	options := &runtimeOptions{}
	flag.StringVar(&options.configPath, "config", config.Path(), "Set config path")
	flag.BoolVar(&options.clickEvents, "no-click-events", true, "Disable click events")
	flag.Parse()
	return options
}
