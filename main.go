package main

import (

	"github.com/spacefito/basiccontroller/controller"
)



func main() {
	controllerApp := controller.NewController()

	go controllerApp.ProcessPacketInMessages()
	go controllerApp.ProcessFlowModMessages()

	controllerApp.Listen(":6633")
}
