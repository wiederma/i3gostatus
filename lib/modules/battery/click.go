package battery

import (
	"github.com/rumpelsepp/i3gostatus/lib/model"
	"strconv"
)

var clickHandlers = &model.ClickHandlers{
	HandleWheelUp:   onWheelUp,
	HandleWheelDown: onWheelDown,
}

func onWheelUp(args *model.ModuleArgs, block *model.I3BarBlock, data interface{}) {
	config := data.(*Config)
	curIndex := config.curIndex()

	if curIndex+1 > numberOfBatteries() {
		curIndex = 0
	} else {
		curIndex++
	}

	config.Index = strconv.Itoa(curIndex)
	block.FullText = config.renderTemplate()
	args.ClickEventCh <- block
}

func onWheelDown(args *model.ModuleArgs, block *model.I3BarBlock, data interface{}) {
	config := data.(*Config)
	curIndex := config.curIndex()

	if curIndex-1 < 0 {
		curIndex = numberOfBatteries()
	} else {
		curIndex--
	}

	config.Index = strconv.Itoa(curIndex)
	block.FullText = config.renderTemplate()
	args.ClickEventCh <- block
}
