package disk

import (
	"bytes"
	"fmt"
	"syscall"
	"text/template"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/rumpelsepp/i3gostatus/lib/config"
	"github.com/rumpelsepp/i3gostatus/lib/model"
	"github.com/rumpelsepp/i3gostatus/lib/utils"
)

const (
	name          = "disk"
	moduleName    = "i3gostatus.modules." + name
	defaultFormat = "{{.Used}}/{{.Total}} [{{.Avail}}]"
	defaultPath   = "/"
	defaultPeriod = 5000
)

type Config struct {
	model.BaseConfig
	Path string
}

func (c *Config) ParseConfig(configTree *toml.TomlTree) {
	c.BaseConfig.Parse(name, configTree)
	c.BaseConfig.Period = config.GetDurationMs(configTree, c.Name+".period", defaultPeriod)
	c.Format = config.GetString(configTree, name+".format", defaultFormat)
	c.Path = config.GetString(configTree, name+".path", defaultPath)
}

type space struct {
	Avail string
	Used  string
	Total string
}

func getSpace(path string) *space {
	var stat syscall.Statfs_t
	syscall.Statfs(path, &stat)
	avail := utils.HumanReadableByteCount(stat.Bavail * uint64(stat.Bsize))
	used := utils.HumanReadableByteCount((stat.Blocks - stat.Bfree) * uint64(stat.Bsize))
	total := utils.HumanReadableByteCount(stat.Blocks * uint64(stat.Bsize))

	return &space{avail, used, total}
}

func (c *Config) Run(args *model.ModuleArgs) {
	outputBlock := model.NewBlock(moduleName, c.BaseConfig, args.Index)
	var outStr string
	t := template.Must(template.New("load").Parse(c.Format))

	for range time.NewTicker(c.Period).C {
		buf := bytes.NewBufferString(outStr)
		space := getSpace(c.Path)

		if err := t.Execute(buf, space); err == nil {
			outputBlock.FullText = buf.String()
		} else {
			outputBlock.FullText = fmt.Sprint(err)
		}

		args.OutCh <- outputBlock
	}
}
