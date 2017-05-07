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
	defaultFormat = `
		{{if eq .State 1}}🔌 {{.Percentage | printf "%.0f"}}% ({{.TimeToFull}}){{end}}
		{{if eq .State 2}}🔋 {{.Percentage | printf "%.0f"}}% ({{.TimeToEmpty}}){{end}}
		{{if eq .State 3}}🔋 EMPTY{{end}}
		{{if eq .State 4}}🔋 FULL{{end}}`
	defaultFormatOnAC = `
		{{if eq .State 1}}🔌 {{.Percentage | printf "%.0f"}}% ({{.TimeToFull}}){{end}}
		{{if eq .State 4}}🔌{{end}}`
)

var logger = log.New(os.Stderr, "["+name+"] ", log.LstdFlags)

type Config struct {
	model.BaseConfig
	FormatOnAC string
}

func (c *Config) ParseConfig(configTree *toml.TomlTree) {
	c.BaseConfig.Parse(name, configTree)
	c.Format = config.GetString(configTree, name+".format", defaultFormat)
	c.FormatOnAC = config.GetString(configTree, name+".format_on_ac", defaultFormatOnAC)
}

func (c *Config) Run(args *model.ModuleArgs) {
	sigCh := SignalChanged()
	outputBlock := model.NewBlock(moduleName, c.BaseConfig, args.Index)
	var outStr string

	// Cleanup template, newlines and tabs are not useful in i3bar.
	c.Format = strings.Replace(strings.Replace(c.Format, "\n", "", -1), "\t", "", -1)
	c.FormatOnAC = strings.Replace(strings.Replace(c.FormatOnAC, "\n", "", -1), "\t", "", -1)

	onBatTmpl := template.Must(template.New("onBattery").Parse(c.Format))
	onACTmpl := template.Must(template.New("onAC").Parse(c.FormatOnAC))

	for {
		buf := bytes.NewBufferString(outStr)
		data := getAllProperties("/org/freedesktop/UPower/devices/DisplayDevice")

		if isOnBattery() {
			if err := onBatTmpl.Execute(buf, data); err == nil {
				outputBlock.FullText = buf.String()
			} else {
				outputBlock.FullText = fmt.Sprint(err)
			}
		} else {
			if err := onACTmpl.Execute(buf, data); err == nil {
				outputBlock.FullText = buf.String()
			} else {
				outputBlock.FullText = fmt.Sprint(err)
			}
		}

		switch {
		case data.Percentage >= 66:
			outputBlock.Color = "#00ff00"
			outputBlock.Background = "000000"
		case data.Percentage < 66 && data.Percentage >= 33:
			outputBlock.Color = "#ffff00"
			outputBlock.Background = "000000"
		case data.Percentage < 33 && data.Percentage >= 10:
			outputBlock.Color = "#ff0000"
			outputBlock.Background = "000000"
		case data.Percentage < 10:
			outputBlock.Color = "#ffffff"
			outputBlock.Background = "ff0000"
		default:
			outputBlock.Color = "#ffffff"
			outputBlock.Background = "000000"
		}

		args.EventCh <- outputBlock

		// Block here until sth. happens. This is better than using in for { ... }
		// since this acts as a do { ... } while() loop and we do not have to
		// wait for the first event at startup. Also, we need to specify a
		// timeout, since the module could hang in some circumstances (don't
		// know why, but it happened).
		select {
		case <-sigCh:
		case <-time.After(10 * time.Second):
			continue
		}
	}
}
