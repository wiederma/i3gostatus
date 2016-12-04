package pulseaudio

import (
	"fmt"
	"github.com/rumpelsepp/i3gostatus/lib/model"
)

var clickHandlers *model.ClickHandlers = &model.ClickHandlers{
	HandleRightClick: onRightClick,
	HandleWheelUp:    onWheelUp,
	HandleWheelDown:  onWheelDown,
}

func onRightClick(args *model.ModuleArgs, block *model.I3BarBlockWrapper, data interface{}) {
	toggleMute()

	// TODO: Remove duplicated code
	if isMuted() {
		block.Color = "#FF0000"
	} else {
		// Use standard color; unsetting it let the JSON
		// marshaller omit the relevant field.
		block.Color = ""
	}

	args.ClickEventCh <- block
}

func onWheelUp(args *model.ModuleArgs, block *model.I3BarBlockWrapper, data interface{}) {
	if err := increaseVolume(5); err != nil {
		return
	}

	config := data.(*Config)
	block.FullText = fmt.Sprintf(config.Format, getCurrentVolume())
	args.ClickEventCh <- block
}

func onWheelDown(args *model.ModuleArgs, block *model.I3BarBlockWrapper, data interface{}) {
	if err := decreaseVolume(5); err != nil {
		return
	}

	config := data.(*Config)
	block.FullText = fmt.Sprintf(config.Format, getCurrentVolume())
	args.ClickEventCh <- block
}
