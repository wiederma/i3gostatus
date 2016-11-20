package i3gostatus

import (
	"flag"
	"github.com/rumpelsepp/i3gostatus/lib/config"
)

type runtimeOptions struct {
	configPath string
}

func ParseOptions() *runtimeOptions {
	options := &runtimeOptions{}
	flag.StringVar(&options.configPath, "config", config.Path(), "Set config path")
	flag.Parse()
	return options
}
