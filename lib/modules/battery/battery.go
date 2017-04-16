// Package static adds a static string to i3bar. Its main purpose is
// demonstrating the module API of `i3gostatus` and it acts as a template for
// new modules.
package battery

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/config"
	"github.com/rumpelsepp/i3gostatus/lib/model"
)

const (
	name          = "battery"
	moduleName    = "i3gostatus.modules." + name
	defaultPeriod = 5000
	defaultFormat = `
	{{ range $_, $val := .}}
		{{if eq .State 1}}{{.NativePath}}: 🔌 {{.Percentage}}% {{end}}
		{{if eq .State 2}}{{.NativePath}}: {{.Percentage}}% {{end}}
		{{if eq .State 3}}{{.NativePath}}: EMPTY {{end}}
		{{if eq .State 4}}{{.NativePath}}: FULL {{end}}
	{{end}}`
)

var logger = log.New(os.Stderr, "["+name+"] ", log.LstdFlags)

type Config struct {
	model.BaseConfig
}

// TODO: Make queried devices configurable.
func (c *Config) ParseConfig(configTree *toml.TomlTree) {
	c.BaseConfig.Parse(name, configTree)
	c.Format = config.GetString(configTree, name+".format", defaultFormat)
	c.Period = config.GetDurationMs(configTree, name+".period", defaultPeriod)
}

func (c *Config) Run(args *model.ModuleArgs) {
	outputBlock := model.NewBlock(moduleName, c.BaseConfig, args.Index)
	var outStr string

	// Cleanup template, newlines and tabs are not useful in i3bar.
	c.Format = strings.Replace(c.Format, "\n", "", -1)
	c.Format = strings.Replace(c.Format, "\t", "", -1)
	var t = template.Must(template.New("upower").Parse(c.Format))

	// FIXME: Do not spam dbus, instead subscribe to signals.
	for range time.NewTicker(c.Period).C {
		buf := bytes.NewBufferString(outStr)
		devs := enumerateDevices()
		var dev_data []Properties

		for _, dev := range devs {
			// TODO: First query type, then query everything.
			d := getAllProperties(dev)

			if d.Type == Battery && d.IsPresent {
				dev_data = append(dev_data, d)
			}
		}

		if err := t.Execute(buf, dev_data); err == nil {
			outputBlock.FullText = buf.String()
		} else {
			outputBlock.FullText = fmt.Sprint(err)
		}
		args.OutCh <- outputBlock
	}
}
