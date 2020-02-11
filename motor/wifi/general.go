package wifi

import (
	"net"
	"time"
	"fmt"
)

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
