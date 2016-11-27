package syncthing

import (
	"os/exec"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/config"
	"github.com/rumpelsepp/i3gostatus/lib/model"
	"github.com/rumpelsepp/i3gostatus/lib/utils"
)

const (
	name       = "static"
	moduleName = "i3gostatus.modules." + name
)

type Config struct {
	model.BaseConfig
	STUrl      string
	UpString   string
	UpColor    string
	DownString string
	DownColor  string
}

func (c *Config) ParseConfig(configTree *toml.TomlTree) {
	c.BaseConfig.Parse(name, configTree)
	c.BaseConfig.Period = config.GetDurationMs(configTree, c.Name+".period", 10000)
	c.STUrl = config.GetString(configTree, name+".st_url", "http://localhost:8384")
	c.UpString = config.GetString(configTree, name+".up_string", "ST UP")
	c.DownString = config.GetString(configTree, name+".down_string", "ST DOWN")
	c.UpColor = config.GetString(configTree, name+".down_string", "#00FF00")
	c.DownColor = config.GetString(configTree, name+".down_string", "#FF0000")
}

func (c *Config) Run(args *model.ModuleArgs) {
	outputBlock := model.NewBlock(moduleName, c.BaseConfig, args.Index)
	stUp := false
	initHTTPSession(c.STUrl)

	go func() {
		// TODO: Make a wrapper method for this. It is duplicated in several modules.
		xdgOpen, err := utils.Which("xdg-open")
		if err != nil {
			switch err.(type) {
			case utils.CommandNotAvailError:
				// TODO: Log a warning here (once the logging system is there...)
				return
			default:
				panic(err)
			}
		}

		// TODO: Make this configurable; DO NOT depend on systemd by design!!
		systemctl, err := utils.Which("systemctl")
		for event := range args.InCh {
			switch event.Button {
			case model.MouseButtonLeft:
				exec.Command(xdgOpen, c.STUrl).CombinedOutput()
			case model.MouseButtonRight:
				if stUp {
					exec.Command(systemctl, "--user", "stop", "syncthing.service").CombinedOutput()
				} else {
					exec.Command(systemctl, "--user", "start", "syncthing.service").CombinedOutput()
				}
			default:
				continue
			}
		}
	}()

	for range time.NewTicker(c.Period).C {
		if resp, err := stGet(c.STUrl, "/rest/system/ping"); err == nil {
			// I do not feel motivated to parse JSON now...
			// This should suffice in most cases.
			if resp == `{"ping":"pong"}` {
				stUp = true
			} else {
				stUp = false
			}
		} else if _, ok := err.(noActiveSessionError); ok {
			initHTTPSession(c.STUrl)
			continue
		} else {
			stUp = false
		}

		if stUp {
			outputBlock.Color = c.UpColor
			outputBlock.FullText = c.UpString
		} else {
			outputBlock.Color = c.DownColor
			outputBlock.FullText = c.DownString
		}

		args.OutCh <- outputBlock
	}
}
