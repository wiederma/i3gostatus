// Package datetime can show the current date and time depending on a
// configurable format string. The formatstring has to be in the Golang specific
// datetime format (see the official docs!). Since it is quite annoying and
// absolutely not recognizeable, it is possible to use the well known strftime
// format. For that purpose a "strftime::" prefix can be added to the format
// string.
package datetime

import (
	"strings"
	"time"

	"github.com/cactus/gostrftime"
	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/config"
	"github.com/rumpelsepp/i3gostatus/lib/model"
)

const (
	name          = "datetime"
	moduleName    = "i3gostatus.modules." + name
	defaultFormat = "2006-01-02 15:04:05"
)

type Config struct {
	model.BaseConfig
}

func (c *Config) ParseConfig(configTree *toml.TomlTree) {
	c.BaseConfig.Parse(moduleName, configTree)
	// http://fuckinggodateformat.com/
	// The golang dateformat string is a mess... So, let's support the classic
	// strftime syntax as well. This prefix must be present in the c.Format
	// string: `strftime::`
	c.Format = config.GetString(configTree, name+".format", defaultFormat)
}

func (c *Config) Run(args *model.ModuleArgs) {
	outputBlock := model.NewBlock(moduleName, c.BaseConfig, args.Index)
	var strftime bool

	if strings.HasPrefix(c.Format, "strftime::") {
		strftime = true
		c.Format = strings.TrimPrefix(c.Format, "strftime::")
	} else {
		strftime = false
	}

	for range time.NewTicker(c.Period).C {
		now := time.Now()
		var str string

		if strftime {
			str = gostrftime.Format(c.Format, now)
		} else {
			str = now.Format(c.Format)
		}

		outputBlock.FullText = str
		args.OutCh <- outputBlock
	}
}
