package main

import (

	"github.com/spacefito/basiccontroller/controller"
	log "github.com/Sirupsen/logrus"
)



func main() {
	controllerApp := controller.NewOfApp()
	log.Infoln( controllerApp)
	controllerApp.SwitchConnected(nil)
	controllerApp.SwitchDisconnected(nil)
	controllerApp.PacketRcvd(nil, nil)
	controllerApp.SendFlowMod(nil, nil)
}