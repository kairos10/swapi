package wifi

import (
	"log"
	"fmt"
	"sync"
	"net"
	"time"
	"encoding/hex"
)

type cmdError struct {
	code byte
	desc string
}
func (e *cmdError) Error() string {
    return fmt.Sprintf("%d - %s", e.code, e.desc)
}

// use a different socket for each command, so the replies from different cmds do not get mixed up
func SendCmdSync(mount *Mount, cmd string) (ret []byte, err error) {
	var numReplies struct {
		sync.RWMutex
		pending int
	}

	synCh1 := make(chan bool)
	localConn := mount.localConn
	if localConn == nil {
		localConn, err = net.ListenUDP("udp", nil)
		if err != nil {
			log.Fatal("SendCmdSync listen: ", err)
		}
	}

	go func() {
		buf := make([]byte, 64)

		synCh1 <- true // let WRITE continue

		// wait TIMEOUT_REPLY for each NUM_REPEAT_CMD
		localConn.SetReadDeadline(time.Now().Add(TIMEOUT_REPLY * (NUM_REPEAT_CMD+1)))
		n, _, err := localConn.ReadFromUDP(buf)
		if err != nil {
			log.Println("read error: ", err)
		} else {
			numReplies.Lock()
			numReplies.pending--
			numReplies.Unlock()
			ret = buf[:n-1]
		}

		synCh1 <- true

	}()

	<- synCh1 // wait for READ to get ready

	for i, cont :=0, true; cont && i<NUM_REPEAT_CMD; i++ {
		_, err := localConn.WriteToUDP([]byte(cmd+"\r"), &mount.UDPAddr)
		if err != nil {
			log.Fatal("write to ", mount.UDPAddr.IP, ": ", err)
		}
		numReplies.Lock()
		numReplies.pending++
		numReplies.Unlock()
		select {
		case <- synCh1: // READ is done
			cont = false
		case <- time.After(TIMEOUT_REPLY * 3 / 2):
		}
	}

	if numReplies.pending == 0 {
		// reuse localAddr
		mount.localConn = localConn
	} else {
		log.Println("localConn not in sync; pending=", numReplies.pending)
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
