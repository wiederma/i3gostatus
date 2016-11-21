package load

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/config"
	"github.com/rumpelsepp/i3gostatus/lib/model"
)

const (
	name          = "load"
	moduleName    = "i3gostatus.modules." + name
	defaultFormat = "%s %s %s"
)

type Load struct {
	avg1  string
	avg5  string
	avg15 string
}

func getLoad() Load {
	loadfile := "/proc/loadavg"
	var loadStr string
	var loads []string

	if data, err := ioutil.ReadFile(loadfile); err == nil {
		loadStr = string(data)
		loads = strings.Split(loadStr, " ")
		return Load{loads[0], loads[1], loads[2]}
	}

	return Load{}
}

type Config struct {
	model.BaseConfig
}

func (c *Config) ParseConfig(configTree *toml.TomlTree) {
	c.BaseConfig.Parse(name, configTree)
	c.Format = config.GetString(configTree, name+".format", defaultFormat)
}

func (c *Config) Run(args *model.ModuleArgs) {
	outputBlock := model.NewBlock(moduleName, c.BaseConfig, args.Index)

	for range time.NewTicker(c.Period).C {
		load := getLoad()
		outputBlock.FullText = fmt.Sprintf(c.Format, load.avg1, load.avg5, load.avg15)
		args.OutCh <- outputBlock
	}
}
