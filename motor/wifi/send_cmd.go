package wifi

import (
	"log"
	"sync"
	"net"
	"time"
)

// use a different socket for each command, so the replies from different cmds do not get mixed up
func SendCmdSync(mount *Mount, cmd string) (ret []byte) {
	var numReplies struct {
		sync.RWMutex
		pending int
	}

	synCh1 := make(chan bool)
	synCh2 := make(chan bool, 2)
	localConn, err := net.ListenUDP("udp", nil)
	if err != nil {
		log.Fatal("SendCmdSync listen: ", err)
	}
	defer localConn.Close()

	go func() {
		buf := make([]byte, 64)

		synCh1 <- true
		for cont:=true; cont; {
			localConn.SetReadDeadline(time.Now().Add(TIMEOUT_REPLY / 3))
			n, _, err := localConn.ReadFromUDP(buf)
			if err==nil {
				numReplies.Lock()
				numReplies.pending--
				numReplies.Unlock()
				synCh1 <- true
				cont = false
				ret = buf[:n-1]
			}

			// the read loop can exit before the writer is aware thet a reply was received; this is why we need 2 spaces on the synCh2 channel
			select {
                        case <- synCh2:
                                cont = false
                        default:
                        }

		}
		synCh1 <- true
	}()
	<- synCh1

	for i, cont :=0, true; cont && i<NUM_REPEAT_CMD; i++ {
		_, err := localConn.WriteToUDP([]byte(cmd+"\r"), &mount.UDPAddr)
		if err != nil {
			log.Fatal("write to ", mount.UDPAddr.IP, ": ", err)
		}
		numReplies.Lock()
		numReplies.pending++
		numReplies.Unlock()
		select {
		case <- synCh1:
			cont = false
		case <- time.After(TIMEOUT_REPLY * 3 / 2):
		}
	}
	synCh2 <- true
	<- synCh1

	if numReplies.pending == 0 {
		// reuse localAddr
		log.Printf("localConn %v still in sync\n", localConn.LocalAddr())
	} else {
		log.Println("localConn not in sync; pending=", numReplies.pending)
	}
	return
}
