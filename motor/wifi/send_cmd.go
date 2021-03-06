package wifi

import (
	"log"
	"fmt"
	"sync"
	"net"
	"time"
	"encoding/hex"
)

// use a different socket for each command, so the replies from different cmds do not get mixed up
func (mount *Mount) SendCmdSync(cmd string) (ret []byte, err error) {
	var numReplies struct {
		sync.RWMutex
		pending int
	}

	synCh := make(chan bool)
	localConn := mount.localConn
	if localConn == nil {
		localConn, err = net.ListenUDP("udp", nil)
		if err != nil {
			log.Fatal("SendCmdSync listen: ", err)
		}
	}

	go func() {
		buf := make([]byte, 64)

		synCh <- true // let WRITE continue

		// wait TIMEOUT_REPLY for each NUM_REPEAT_CMD
		localConn.SetReadDeadline(time.Now().Add(TIMEOUT_REPLY * (NUM_REPEAT_CMD+1)))
		n, _, err := localConn.ReadFromUDP(buf)
		if err != nil {
			//mount.log(fmt.Sprintln("read error: ", err))
		} else {
			numReplies.Lock()
			numReplies.pending--
			numReplies.Unlock()
			ret = buf[:n-1]
			//mount.log(fmt.Sprintln("YYY: ", string(cmd), " - ", string(ret)))
		}

		synCh <- true

	}()

	<- synCh // wait for READ to get ready

	for i, cont :=0, true; cont && i<NUM_REPEAT_CMD; i++ {
		_, err := localConn.WriteToUDP([]byte(cmd+"\r"), &mount.UDPAddr)
		if err != nil {
			log.Fatal("write to ", mount.UDPAddr.IP, ": ", err)
		}
		numReplies.Lock()
		numReplies.pending++
		numReplies.Unlock()
		select {
		case <- synCh: // READ is done
			cont = false
		case <- time.After(TIMEOUT_REPLY):
		}
	}

	if numReplies.pending == 0 {
		// reuse localAddr
		mount.localConn = localConn
	} else {
		mount.log(fmt.Sprintln("localConn not in sync; pending=", numReplies.pending))
		mount.localConn = nil
		defer localConn.Close()
	}

	if len(ret) == 0 {
		err = &cmdError{0, "cmd no reply"}
	} else if ret[0] == '!' {
		var code [1]byte
		hex.Decode(code[:], ret[1:])
		err = &cmdError{code[0], "remote error"}
	}
	return
}
