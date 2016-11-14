package model

import (
	"github.com/pelletier/go-toml"
	"syscall"
	"time"
)

type I3BarHeader struct {
	Version     int  `json:"version"`
	StopSignal  int  `json:"stop_signal",omitempty`
	ContSignal  int  `json:"cont_signal",omitempty`
	ClickEvents bool `json:"click_events",omitempty`
}

type I3BarBlock struct {
	FullText            string `json:"full_text"`
	ShortText           string `json:"short_text,omitempty"`
	Color               string `json:"color,omitempty"`
	Background          string `json:"background,omitempty"`
	Border              string `json:"border,omitempty"`
	MinWidth            string `json:"min_width,omitempty"`
	Align               string `json:"align,omitempty"`
	Name                string `json:"name",omitempty`
	Instance            string `json:"instance,omitempty"`
	Urgent              bool   `json:"urgent,omitempty"`
	Separator           bool   `json:"separator,omitempty"`
	SeparatorBlockWidth string `json:"separator_block_width,omitempty"`
	Markup              string `json:"markup,omitempty"`
}

type I3BarBlockWrapper struct {
	I3BarBlock
	Index int
}

type I3ClickEvent struct {
	Name     string `json:"name"`
	Instance string `json:"instance"`
	Button   int    `json:"button"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
}

const DefaultPeriod = 1000

type BaseConfig struct {
	I3BarBlock
	Period time.Duration
}

type Module interface {
	Run(chan *I3BarBlockWrapper, int)
	ReadConfig(configTree *toml.TomlTree)
}

func NewHeader() *I3BarHeader {
	return &I3BarHeader{
		Version:     1,
		StopSignal:  int(syscall.SIGSTOP),
		ContSignal:  int(syscall.SIGCONT),
		ClickEvents: false,
	}
}

func NewBlock(name string, initValues BaseConfig, index int) *I3BarBlockWrapper {
	return &I3BarBlockWrapper{
		I3BarBlock: I3BarBlock{
			Name:                name,
			Color:               initValues.Color,
			Background:          initValues.Background,
			Border:              initValues.Border,
			MinWidth:            initValues.MinWidth,
			Align:               initValues.Align,
			Urgent:              initValues.Urgent,
			Separator:           initValues.Separator,
			SeparatorBlockWidth: initValues.SeparatorBlockWidth,
			Markup:              initValues.Markup,
		},
		Index: index,
	}
}

func (config *BaseConfig) ReadConfig(name string, configTree *toml.TomlTree) {
	config.Name = name
	config.Color = configTree.GetDefault(config.Name+".color", "").(string)
	config.Background = configTree.GetDefault(config.Name+".background", "").(string)
	config.Border = configTree.GetDefault(config.Name+".border", "").(string)
	config.MinWidth = configTree.GetDefault(config.Name+".min_width", "").(string)
	config.Align = configTree.GetDefault(config.Name+".align", "").(string)
	config.Urgent = configTree.GetDefault(config.Name+".urgent", false).(bool)
	config.Separator = configTree.GetDefault(config.Name+".separator", false).(bool)
	config.SeparatorBlockWidth = configTree.GetDefault(config.Name+".separator_block_width", "").(string)
	config.Markup = configTree.GetDefault(config.Name+".markup", "").(string)
	config.Period = time.Duration(configTree.GetDefault(config.Name+".period", int64(DefaultPeriod)).(int64)) * time.Millisecond
}
