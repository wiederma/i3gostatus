// Package static adds a static string to i3bar. Its main purpose is
// demonstrating the module API of `i3gostatus` and it acts as a template for
// new modules.
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
}

func (c *Config) ParseConfig(configTree *toml.TomlTree) {
	c.BaseConfig.Parse(moduleName, configTree)
}

func (c *Config) Run(args *model.ModuleArgs) {
	outputBlock := model.NewBlock(moduleName, c.BaseConfig, args.Index)

	for range time.NewTicker(c.Period).C {
		outputBlock.FullText = c.Format
		args.OutCh <- outputBlock
	}
}
