package controller

import (
	"github.com/shaleman/libOpenflow/openflow13"
	"github.com/contiv/ofnet/ofctrl"
	log "github.com/Sirupsen/logrus"
)


//Types:

type FlowModMessage struct {
	sw *ofctrl.OFSwitch
	flowMod *openflow13.FlowMod
}

type PacketInMessage struct {
	sw *ofctrl.OFSwitch
	PacketInMessage  *ofctrl.PacketIn
}

type DPID int64

//"Class" definition starts here.
type OfApp struct {
	//PacketIns from switch to be picked up by this application go in here
	PacketInMessageChannel chan PacketInMessage

	// FlowMods to Switch get put in here
	FlowModMessageChannel chan FlowModMessage

	// List of switches this app is looking at.
	SwitchList map[DPID]*ofctrl.OFSwitch
}

// "Constructor"
func NewOfApp () *OfApp {
	object := new(OfApp)
	object.PacketInMessageChannel = make(chan PacketInMessage)
	object.FlowModMessageChannel = make(chan FlowModMessage)
	object.SwitchList = make(map[DPID]*ofctrl.OFSwitch)
	return object
}


func (o *OfApp) SwitchConnected(sw *ofctrl.OFSwitch) {
	log.Infoln("OfApp: Switch ", sw, " connected to app: ", o)
}

func (o *OfApp) SwitchDisconnected(sw *ofctrl.OFSwitch) {
	log.Infoln("OfApp: Switch ", sw," disconnected from: ", o)
}

func (o *OfApp) PacketRcvd(sw *ofctrl.OFSwitch, packet *ofctrl.PacketIn) {
	log.Infoln("OfApp: Received packet ", packet, " from: ", sw)
}

func (o *OfApp) SendFlowMod(sw *ofctrl.OFSwitch, flowMod *openflow13.FlowMod) {
	log.Infoln("OfApp: sending flomod: ",flowMod, " to switch: ",sw )
}