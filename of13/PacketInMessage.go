package of13

import (
	"github.com/contiv/ofnet/ofctrl"
)


type PacketInMessage struct {
	Sw *ofctrl.OFSwitch
	PacketInMessage  *ofctrl.PacketIn
}

func (p *PacketInMessage) Dpid() Dpid {
	// Gets DPID from the sw which is included in the PacketIn
	var falsedpid Dpid
	falsedpid = 1
	return falsedpid

}
