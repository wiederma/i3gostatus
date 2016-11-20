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
	c.BaseConfig.Parse(name, configTree)
}

func (c *Config) Run(out chan *model.I3BarBlockWrapper, in chan *model.I3ClickEvent, index int) {
	outputBlock := model.NewBlock(moduleName, c.BaseConfig, index)

	for range time.NewTicker(c.Period).C {
		outputBlock.FullText = c.Format
		out <- outputBlock
	}
}
