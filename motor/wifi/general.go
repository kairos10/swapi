package wifi

import (
	"net"
	"time"
	"fmt"
	"strings"
)

// default UDP communication port
const SW_UDP_PORT = 11880

// UDP communication parameters
const (
	TIMEOUT_REPLY  = 100 * time.Millisecond // reply timeout in [ms]
	NUM_REPEAT_CMD = 5                      // resend the [cmd] for how many times if there is no reply
)

/*
Information related to a networked motor controller.
The SW* methods mirror the low level commands supported by the MC
*/
type Mount struct {
	LoggerFunc func(string)
	UDPAddr       net.UDPAddr // MC Address
	DiscoveryTime time.Time // discovery time

	localConn *net.UDPConn // local UDP address used to receive responses from the MC

	isInit bool // true if the MC* parameters have been retrieved

	MCversion     string // MC version code, as reported by the mount

	MCParamFrequency int
	MCParamCPR int
	MCParamHighSpeedMult int
	MCParamT1Tracking1X int

	HasDualEncoder bool
	HasPPEC bool
	HasOriginalIndex bool
	HasEqAz bool
	HasPolarScopeLED bool
	MustSeparateStartAxis bool
	HasTorqueSelection bool

	isEqInit bool // InitializeEQ has been called on the mount, initializing DEC=0 and RA=CPR/4; for an already initialized mount, there is no way to tell wether the mount was initialized in AZ mode (0, 0) or EQ mode (CPR/4, 0)
}
func (m *Mount) String() (r string) {
	if !m.isInit { _ = m.RetrieveMountParameters() }
	r += fmt.Sprintf("Addr[%v] Ver[%s] DualEnc[%v] EqAz[%v] AxSepStart[%v]", m.UDPAddr, m.MCversion, m.HasDualEncoder, m.HasEqAz, m.MustSeparateStartAxis)
	return
}

func (m *Mount) log(s string) {
	if m.LoggerFunc != nil {
		m.LoggerFunc(s)
	}
}

// Initialize an existing Mount with a static IP address
func (m *Mount) Resolve(addr string, port int) (err0 error) {
	if port <= 0 { port = SW_UDP_PORT }
	if !strings.Contains(addr, ":") {
		addr = fmt.Sprintf("%s:%d", addr, port)
	}
	udpAddr, err0 := net.ResolveUDPAddr("udp", addr)
	if err0 == nil {
		m.UDPAddr = *udpAddr
		err0 = m.RetrieveMountParameters()
	}
	return
}

type cmdError struct {
	code byte
	desc string
}
func (e *cmdError) Error() string {
	return fmt.Sprintf("%d - %s", e.code, e.desc)
}

// various error codes
const (
	ERR01_AXIS	=	iota+100
	ERR02_RESP_LEN
	ERR03_PARAM
	ERR04_NA
	ERR05_NOT_SUPPORTED
	ERR06_VALUE_TOO_LARGE
	ERR07_ALREADY_INITIALIZED
	ERR08
	ERR09
)
