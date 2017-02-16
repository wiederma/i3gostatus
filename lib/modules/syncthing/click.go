package syncthing

import (
	"os/exec"

	"github.com/rumpelsepp/i3gostatus/lib/model"
)

var clickHandlers = &model.ClickHandlers{
	HandleRightClick: onRightClick,
	HandleLeftClick:  onLeftClick,
}

func onRightClick(args *model.ModuleArgs, block *model.I3BarBlock, data interface{}) {
	config := data.(*Config)

	if isUp(config.STUrl) {
		stopSyncthing()
	} else {
		startSyncthing()
	}
}

func onLeftClick(args *model.ModuleArgs, block *model.I3BarBlock, data interface{}) {
	config := data.(*Config)
	exec.Command(xdgOpen, config.STUrl).Run()
}
