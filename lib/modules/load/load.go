package load

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/config"
	"github.com/rumpelsepp/i3gostatus/lib/model"
)

const (
	name               = "load"
	moduleName         = "i3gostatus.modules." + name
	defaultFormat      = "{{.Avg1}} {{.Avg5}} {{.Avg15}}"
	defaultUrgentColor = "#ff0000"
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
	UrgentLoad  float64
	UrgentColor string
}

func (c *Config) ParseConfig(configTree *toml.TomlTree) {
	c.BaseConfig.Parse(name, configTree)
	c.Format = config.GetString(configTree, name+".format", defaultFormat)
	c.UrgentLoad = config.GetFloat64(configTree, name+".urgent_load", float64(runtime.NumCPU()))
	c.UrgentColor = config.GetString(configTree, name+".urgent_color", defaultUrgentColor)
}

func (c *Config) Run(args *model.ModuleArgs) {
	var outStr string
	t := template.Must(template.New("load").Parse(c.Format))

	for range time.NewTicker(c.Period).C {
		outputBlock := model.NewBlock(moduleName, c.BaseConfig, args.Index)
		buf := bytes.NewBufferString(outStr)
		load := getLoad()

		if err := t.Execute(buf, load); err == nil {
			if loadVal, err := strconv.ParseFloat(load.Avg1, 64); err == nil {
				if loadVal >= float64(c.UrgentLoad) {
					outputBlock.Color = c.UrgentColor
					outputBlock.Urgent = true
				}
			}
			outputBlock.FullText = buf.String()
		} else {
			outputBlock.FullText = fmt.Sprint(err)
		}

		args.OutCh <- outputBlock
	}
}
