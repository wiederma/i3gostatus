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

	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/config"
	"github.com/rumpelsepp/i3gostatus/lib/model"
)

const (
	name          = "battery"
	moduleName    = "i3gostatus.modules." + name
	defaultPeriod = 5000
	defaultFormat = `
		{{if eq .State 1}}BAT: 🔌 {{.Percentage | printf "%.0f"}}% ({{.TimeToFull}}){{end}}
		{{if eq .State 2}}BAT: {{.Percentage | printf "%.0f"}}% ({{.TimeToEmpty}}){{end}}
		{{if eq .State 3}}BAT: EMPTY{{end}}
		{{if eq .State 4}}BAT: FULL{{end}}`
)

var logger = log.New(os.Stderr, "["+name+"] ", log.LstdFlags)

type Config struct {
	model.BaseConfig
}

func (c *Config) ParseConfig(configTree *toml.TomlTree) {
	c.BaseConfig.Parse(name, configTree)
	c.Format = config.GetString(configTree, name+".format", defaultFormat)
	c.Period = config.GetDurationMs(configTree, name+".period", defaultPeriod)
}

func (c *Config) Run(args *model.ModuleArgs) {
	sigCh := SignalChanged()
	outputBlock := model.NewBlock(moduleName, c.BaseConfig, args.Index)
	var outStr string

	// Cleanup template, newlines and tabs are not useful in i3bar.
	c.Format = strings.Replace(c.Format, "\n", "", -1)
	c.Format = strings.Replace(c.Format, "\t", "", -1)

	var t = template.Must(template.New("upower").Parse(c.Format))

	for {
		buf := bytes.NewBufferString(outStr)
		data := getAllProperties("/org/freedesktop/UPower/devices/DisplayDevice")

		if err := t.Execute(buf, data); err == nil {
			outputBlock.FullText = buf.String()
		} else {
			outputBlock.FullText = fmt.Sprint(err)
		}

		switch {
		case data.Percentage >= 66:
			outputBlock.Color = "#00ff00"
		case data.Percentage < 66 && data.Percentage >= 33:
			outputBlock.Color = "#ffff00"
		case data.Percentage < 33 && data.Percentage >= 10:
			outputBlock.Color = "#ff0000"
		default:
			outputBlock.Color = "#ffffff"
			outputBlock.Background = "ff0000"
		}

		args.OutCh <- outputBlock

		// Block here until sth. happens.
		// This is better than using in for { ... } since
		// this acts as a do { ... } while() loop and we do
		// not have to wait for the first event at startup.
		<-sigCh
	}
}
