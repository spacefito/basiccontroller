package of13

import (
	"github.com/contiv/ofnet/ofctrl"
	"github.com/shaleman/libOpenflow/openflow13"
)


type FlowModMessage struct {
	Sw *ofctrl.OFSwitch
	FlowMod  *openflow13.FlowMod
}

func (p *FlowModMessage) Dpid() Dpid {
	// Gets DPID from the sw which is included in the The flowmod is for
	var falsedpid Dpid
	falsedpid = 1
	return falsedpid

}

func (p *FlowModMessage) Send() {
	p.Sw.Send(p.FlowMod)
}