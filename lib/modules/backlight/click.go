package backlight

import (
	"fmt"

	"github.com/rumpelsepp/i3gostatus/lib/model"
)

var clickHandlers *model.ClickHandlers = &model.ClickHandlers{
	HandleWheelUp:   onWheelUp,
	HandleWheelDown: onWheelDown,
	// Handle clicks as aliases to mouse wheel events
	HandleLeftClick:  onWheelUp,
	HandleRightClick: onWheelDown,
}

func onWheelUp(args *model.ModuleArgs, block *model.I3BarBlock, data interface{}) {
	config := data.(*Config)

	incBrightness(5)
	block.FullText = fmt.Sprintf(config.Format, getBrightness())
	args.ClickEventCh <- block
}

func onWheelDown(args *model.ModuleArgs, block *model.I3BarBlock, data interface{}) {
	config := data.(*Config)

	// Set hard limit. Otherwise the display could be blacked out
	// completetly when using the mouse wheel. This is annoying.
	if b := getBrightness(); b <= 10 {
		setBrightness(10)
	} else {
		decBrightness(5)
		block.FullText = fmt.Sprintf(config.Format, getBrightness())
		args.ClickEventCh <- block
	}
}
