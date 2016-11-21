package load

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/config"
	"github.com/rumpelsepp/i3gostatus/lib/model"
)

const (
	name          = "load"
	moduleName    = "i3gostatus.modules." + name
	defaultFormat = "{{.Avg1}} {{.Avg5}} {{.Avg15}}"
)

type Load struct {
	Avg1  string
	Avg5  string
	Avg15 string
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
	t := template.Must(template.New("load").Parse(c.Format))
	var outStr string

	for range time.NewTicker(c.Period).C {
		buf := bytes.NewBufferString(outStr)

		if err := t.Execute(buf, getLoad()); err == nil {
			outputBlock.FullText = buf.String()
		} else {
			outputBlock.FullText = fmt.Sprint(err)
		}

		args.OutCh <- outputBlock
	}
}
