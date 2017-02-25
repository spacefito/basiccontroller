package controller

import (
	"github.com/shaleman/libOpenflow/openflow13"
	"github.com/contiv/ofnet/ofctrl"
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

//class definition starts here.
type OfApp struct {
	//PacketIns from switch to be picked up by this application go in here
	PacketInMessageChannel chan PacketInMessage

	// FlowMods to Switch get put in here
	FlowModMessageChannel chan FlowModMessage

	// List of switches this app is looking at.
	SwitchList map[DPID]ofctrl.OFSwitch
}



