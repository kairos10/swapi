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

	MCversion     string
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
	ERR03
	ERR04
)
