package i3gostatus

import flag "github.com/spf13/pflag"

type runtimeOptions struct {
	configPath string
}

func ParseOptions() *runtimeOptions {
	options := &runtimeOptions{}
	flag.StringVarP(&options.configPath, "config", "c", getConfigPath(), "Set config path")
	flag.Parse()
	return options
}
