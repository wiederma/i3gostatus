package static

import (
	"time"

	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/model"
)

const (
	name       = "static"
	moduleName = "i3gostatus.modules." + name
)

type Config struct {
	model.BaseConfig
	Output string
}

func (config *Config) Run(out chan *model.I3BarBlockWrapper, index int) {
	ticker := time.NewTicker(config.Period)
	outputBlock := model.NewBlock(moduleName, config.BaseConfig, index)

	for range ticker.C {
		outputBlock.FullText = config.Output
		out <- outputBlock
	}
}

func (config *Config) ReadConfig(configTree *toml.TomlTree) {
	config.BaseConfig.ReadConfig(name, configTree)

	config.Output = configTree.GetDefault(name+".output", "").(string)
}
