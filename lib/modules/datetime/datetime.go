package datetime

import (
	"strings"
	"time"

	"github.com/cactus/gostrftime"
	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/model"
)

const (
	name       = "datetime"
	moduleName = "i3gostatus.modules." + name
)

type Config struct {
	model.BaseConfig
	Format string
}

func (config *Config) Run(out chan *model.I3BarBlockWrapper, index int) {
	ticker := time.NewTicker(config.Period)
	outputBlock := model.NewBlock(moduleName, config.BaseConfig, index)
	var strftime bool

	if strings.HasPrefix(config.Format, "strftime::") {
		strftime = true
		config.Format = strings.TrimPrefix(config.Format, "strftime::")
	} else {
		strftime = false
	}

	for range ticker.C {
		now := time.Now()
		var str string

		if strftime {
			str = gostrftime.Format(config.Format, now)
		} else {
			str = now.Format(config.Format)
		}

		outputBlock.FullText = str
		out <- outputBlock
	}
}

func (config *Config) ReadConfig(configTree *toml.TomlTree) {
	config.BaseConfig.ReadConfig(name, configTree)

	// http://fuckinggodateformat.com/
	// The golang dateformat string is a mess... So, let's support the classic
	// strftime syntax as well. This prefix must be present in the config.Format
	// string: `strftime::`
	config.Format = configTree.GetDefault(name+".format", "2006-01-02 15:04:05").(string)
}
