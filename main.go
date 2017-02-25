package main

import (

	"github.com/spacefito/basiccontroller/controller"
	log "github.com/Sirupsen/logrus"
)



func main() {
	controllerApp := controller.OfApp{nil, nil, nil}
	log.Infoln( controllerApp)
}