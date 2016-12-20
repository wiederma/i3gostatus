package model

import (
	"syscall"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/satori/go.uuid"
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
	Index               int    `json:"-"`
}

type I3ClickEvent struct {
	Name     string `json:"name"`
	Instance string `json:"instance"`
	Button   int    `json:"button"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
}

const (
	_               = iota
	MouseButtonLeft = iota
	MouseButtonMiddle
	MouseButtonRight
	MouseWheelUp
	MouseWheelDown
)

const DefaultPeriod = 1000

type ModuleArgs struct {
	InCh         chan *I3ClickEvent
	OutCh        chan *I3BarBlock
	ClickEventCh chan *I3BarBlock
	Index        int
}

type BaseConfig struct {
	// TODO: Duplication in NewBlock()?
	I3BarBlock
	Period time.Duration
	Format string
}

func (c *BaseConfig) Parse(name string, configTree *toml.TomlTree) {
	c.Name = name
	c.Instance = uuid.NewV4().String()
	c.Color = configTree.GetDefault(c.Name+".color", "").(string)
	c.Background = configTree.GetDefault(c.Name+".background", "").(string)
	c.Border = configTree.GetDefault(c.Name+".border", "").(string)
	c.MinWidth = configTree.GetDefault(c.Name+".min_width", "").(string)
	c.Align = configTree.GetDefault(c.Name+".align", "").(string)
	c.Urgent = configTree.GetDefault(c.Name+".urgent", false).(bool)
	c.Separator = configTree.GetDefault(c.Name+".separator", false).(bool)
	c.SeparatorBlockWidth = configTree.GetDefault(c.Name+".separator_block_width", "").(string)
	c.Markup = configTree.GetDefault(c.Name+".markup", "").(string)
	c.Period = time.Duration(configTree.GetDefault(c.Name+".period", int64(DefaultPeriod)).(int64)) * time.Millisecond
	c.Format = configTree.GetDefault(c.Name+".format", "").(string)
}

type Module interface {
	Run(*ModuleArgs)
	ParseConfig(*toml.TomlTree)
}

func NewHeader(click_events bool) *I3BarHeader {
	return &I3BarHeader{
		Version:     1,
		StopSignal:  int(syscall.SIGSTOP),
		ContSignal:  int(syscall.SIGCONT),
		ClickEvents: click_events,
	}
}

func NewBlock(name string, initValues BaseConfig, index int) *I3BarBlock {
	return &I3BarBlock{
		Name:                name,
		Instance:            initValues.Instance,
		Color:               initValues.Color,
		Background:          initValues.Background,
		Border:              initValues.Border,
		MinWidth:            initValues.MinWidth,
		Align:               initValues.Align,
		Urgent:              initValues.Urgent,
		Separator:           initValues.Separator,
		SeparatorBlockWidth: initValues.SeparatorBlockWidth,
		Markup:              initValues.Markup,
		Index:               index,
	}
}
