package controller

import (
	"github.com/shaleman/libOpenflow/openflow13"
	"github.com/contiv/ofnet/ofctrl"
	log "github.com/Sirupsen/logrus"
	"github.com/spacefito/basiccontroller/of13"
	"net"
	"strings"
	"github.com/shaleman/libOpenflow/util"
	"github.com/shaleman/libOpenflow/common"
	"time"
)


//"Class" definition starts here.
type Controller struct {
	//PacketIns from switch to be picked up by this application go in here
	PacketInMessageChannel chan of13.PacketInMessage

	// FlowMods to Switch get put in here
	FlowModMessageChannel chan of13.FlowModMessage

	// List of switches this app is looking at.
	SwitchList map[of13.Dpid]*ofctrl.OFSwitch

	// tcp connection (listener)
	listener *net.TCPListener
}

// "Constructor"
func NewController () *Controller {
	object := new(Controller)
	object.PacketInMessageChannel = make(chan of13.PacketInMessage)
	object.FlowModMessageChannel = make(chan of13.FlowModMessage)
	object.SwitchList = make(map[of13.Dpid]*ofctrl.OFSwitch)
	return object
}


func (c *Controller) SwitchConnected(sw *ofctrl.OFSwitch) {
	log.Infoln("Controller: Switch ", sw, " connected to app: ", c)
	c.InstallDefaultFlowsOnSW(sw)
}

func (c *Controller) SwitchDisconnected(sw *ofctrl.OFSwitch) {
	log.Infoln("Controller: Switch ", sw," disconnected from: ", c)
}

func (c *Controller) PacketRcvd(sw *ofctrl.OFSwitch, packet *ofctrl.PacketIn) {
	log.Infoln("Controller: Received packet ", packet, " from: ", sw)
	c.PacketInMessageChannel <- of13.PacketInMessage{sw, packet}
}

func (c *Controller) SendFlowMod(sw *ofctrl.OFSwitch, flowMod *openflow13.FlowMod) {
	log.Infoln("Controller: sending flomod: ",flowMod, " to switch: ",sw )
}

func (c *Controller) HandlePacketInMessage(packetInMessage of13.PacketInMessage) {
	c.SwitchList[packetInMessage.Dpid()] = packetInMessage.Sw
}


func (c *Controller) ProcessPacketInMessages() {
	for {
		packetInMessage := <-c.PacketInMessageChannel
                go c.HandlePacketInMessage(packetInMessage)
	}
}

func (o *Controller) ProcessFlowModMessages() {
// InstallFlowMod takes flowModRequest from FLowModRequestChannel and installs it on apropriate switch
	log.Infoln("processing flowmods")
	for {
		flowModMessage := <-o.FlowModMessageChannel
		log.Infoln("CONTROLLER: got a flowmod")
		flowModMessage.Send()
	}
}

//Listen on port
func (c *Controller) Listen (port string) {
	addr, _ := net.ResolveTCPAddr("tcp", port)

	var err error
	c.listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	defer c.listener.Close()

	log.Println("Listening for connections on", addr)
	for {
		conn, err := c.listener.AcceptTCP()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
			log.Fatal(err)
		}
		go c.handleConnection(conn)
	}
}

// Handle TCP connection from the swtich
func (c *Controller) handleConnection(conn net.Conn) {
	stream := util.NewMessageStream(conn, c)

	log.Println("New connection..")

	// Send ofp 1.3 Hello by default
	h, err := common.NewHello(4)
	if err != nil {
		return
	}
	stream.Outbound <- h

	for {
		select {
		// Send hello messge with latest protocol version.
		case msg := <-stream.Inbound:
			switch m := msg.(type) {
			//if we get a correct version reply we can continue,
			// if not we close connection
			case *common.Hello:
				if m.Version == openflow13.VERSION {
					log.Infoln("Received Openflow 1.3 Hello message")
					// Version negotiation is
					// considered complete. Create
					// new Switch and notifiy listening
					// applications.
					stream.Version = m.Version
					stream.Outbound <- openflow13.NewFeaturesRequest()
				} else {
					// Connection should be severed if controller
					// doesn't support switch version.
					log.Println("Received unsupported ofp version", m.Version)
					stream.Shutdown <- true
				}
			// After a vaild FeaturesReply has been received we
			// have all the information we need. Create a new
			// switch object and notify applications
			case *openflow13.SwitchFeatures:
				log.Printf("Received ofp1.3 Switch feature response: %+v", *m)

				// Create a new switch and handover the stream
				ofctrl.NewSwitch(stream, m.DPID, c)

				// Let switch instance handle all future messages..
				return

			// An error message may indicate a version mismatch. We
			// disconnect if an error occurs this early.
			case *openflow13.ErrorMsg:
				log.Warnf("Received ofp1.3 error msg: %+v", *m)
				stream.Shutdown <- true
			}
		case err := <-stream.Error:
			// The connection has been shutdown
		log.Println(err)
			return
		case <-time.After(time.Second *3):
			// This shouldn't happen. If it does, both the controller
			// and switch are no longer communicating. The TCPConn is
			// still established though.
			log.Warnln("Connection timed out")
			return
		}
	}
}

// Demux based on message version
func (c *Controller) Parse(b []byte) (message util.Message, err error) {
	switch b[0] {
	case openflow13.VERSION:
		message, err = openflow13.Parse(b)
	default:
		log.Errorf("Received unsupported openflow version: %d", b[0])
	}
	return
}


func (o *Controller) MultipartReply(sw *ofctrl.OFSwitch, rep *openflow13.MultipartReply){
	log.Infoln("NetMonitor: Received multipartReply from switch: ", sw.DPID())
}


func (c *Controller) InstallDefaultFlowsOnSW(sw *ofctrl.OFSwitch) {
	log.Infoln("NetMonitor: Installing default flows on ", sw.DPID())
	// For unmatched flows, the default is to send packet out to NORMAL,
	// send a copy of the header to the controller.
	// flow matches all ip packets (FlowMatch), sends enough of package to process ip header to controller (FlowAction),
	// and forwards the packet to NORMAL (FlowAction)

	//to NORMAL
	defaultFlowMod := openflow13.NewFlowMod()
	defaultFlowMod.Priority = 0

	outputActNormal := openflow13.NewActionOutput(openflow13.P_NORMAL)

	defaultOutputInstr := openflow13.NewInstrApplyActions()
	defaultOutputInstr.AddAction(outputActNormal, false)


	defaultFlowMod.AddInstruction(defaultOutputInstr)

	//Install it by  creating a FlowModRequest with sw and defaultFlowMod
	// and putting the FlowModRequest in the go channel ("queue")
	c.FlowModMessageChannel <- of13.FlowModMessage{sw, defaultFlowMod}



	//Next we need to create a ip (eth type 0x0800) flow to steal a copy of
	// packet headers to the controller
	ipOutputInstr := openflow13.NewInstrApplyActions()
	ipOutputInstr.AddAction(outputActNormal, false)
	outputActController := openflow13.NewActionOutput(openflow13.P_CONTROLLER)
	outputActController.MaxLen = 50
	ipOutputInstr.AddAction(outputActController, false)
	flowMod := openflow13.NewFlowMod()
	flowMod.Cookie = 1
		flowMod.Match.AddField(*openflow13.NewEthTypeField(0x0800))
		flowMod.AddInstruction(ipOutputInstr)

	//Install it by  creating a FlowModRequest with sw and defaultFlowMod
	// and putting the FlowModRequest in the go channel ("queue")
	c.FlowModMessageChannel <- of13.FlowModMessage{sw, flowMod}

}
