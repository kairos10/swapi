package wifi

import (
	"net"
	"time"
	"fmt"
	"strings"
)

const SW_UDP_PORT = 11880

/*
SW* methods mirror the motor commands
*/
type Mount struct {
	UDPAddr       net.UDPAddr
	DiscoveryTime time.Time

	localConn *net.UDPConn

	isInit bool // true if the MC* parameters have been retrieved

	MCversion     string

	MCParamFrequency int
	MCParamCPR int
	MCParamHighSpeedMult int
	MCParamT1Tracking1X int

	HasDualEncoder bool
	HasPPEC bool
	HasOriginalIndex bool
	HasEqAz bool
	HasPolarScopeLED bool
	HasAxisSeparateStart bool
	HasTorqueSelection bool
}
func (m *Mount) String() (r string) {
	if !m.isInit { _ = m.RetrieveMountParameters() }
	r += fmt.Sprintf("Addr[%v] Ver[%s] DualEnc[%v] EqAz[%v] AxSepStart[%v]", m.UDPAddr, m.MCversion, m.HasDualEncoder, m.HasAxisSeparateStart, m.HasAxisSeparateStart)
	return
}

// m := new(wifi.Mount)
// m.Resolve("192.168.4.1", -1) // port defaults to 11880 [= SW_UDP_PORT]
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

const (
	TIMEOUT_REPLY  = 100 * time.Millisecond // reply timeout in [ms]
	NUM_REPEAT_CMD = 5                      // resend the [cmd] for how many times if there is no reply
)

type cmdError struct {
	code byte
	desc string
}
func (e *cmdError) Error() string {
	return fmt.Sprintf("%d - %s", e.code, e.desc)
}
const (
	ERR01_AXIS	=	iota+100
	ERR02_RESP_LEN
	ERR03_PARAM
	ERR04_NA
	ERR05_NOT_SUPPORTED
	ERR06_VALUE_TOO_LARGE
	ERR07
)
