package wifi

import (
	"net"
	"time"
)

type Mount struct {
        UDPAddr        net.UDPAddr
        MCversion      string
        DiscoveryTime time.Time

	localConn *net.UDPConn
}

const (
        TIMEOUT_REPLY   =       100 * time.Millisecond	// reply timeout in [ms]
        NUM_REPEAT_CMD  =       5			// resend the [cmd] for how many times if there is no reply
)

