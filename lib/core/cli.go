package i3gostatus

import (
	"github.com/rumpelsepp/i3gostatus/lib/config"
	flag "github.com/spf13/pflag"
)

type runtimeOptions struct {
	configPath string
}

func ParseOptions() *runtimeOptions {
	options := &runtimeOptions{}
	flag.StringVarP(&options.configPath, "config", "c", config.Path(), "Set config path")
	flag.Parse()
	return options
}
