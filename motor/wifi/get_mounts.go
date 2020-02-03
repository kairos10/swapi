package wifi

import (
	"log"
	"fmt"
	"net"
	"time"
	"encoding/hex"
)

func FindMounts() []*Mount {
	lAddr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:11881")
	lConn, err := net.ListenUDP("udp", lAddr)
	if err != nil {
		log.Fatal("FindMounts err: ", err)
	}
	defer lConn.Close()

	mounts := make(map[string]*Mount)
	synCh1 := make(chan bool)
	synCh2 := make(chan bool)
	go func() {
		buf := make([]byte, 64)
		decodeBuf := make([]byte, 3)

		synCh1 <- true
		for cont:=true; cont; {
			lConn.SetReadDeadline(time.Now().Add(TIMEOUT_REPLY))
			n, fromAddr, err := lConn.ReadFromUDP(buf)
			if err==nil && !lAddr.IP.Equal(fromAddr.IP) {
				// ignore packets received from self
				mount := mounts[string(fromAddr.IP)]
				if mount == nil {
					mount = new(Mount)
					mount.UDPAddr = *fromAddr
					mount.DiscoveryTime = time.Now()
					mounts[string(fromAddr.IP)] = mount
				}
				if n > 7 {
					hex.Decode(decodeBuf, buf[1:5])
					mount.MCversion = fmt.Sprintf("%d.%d %s", decodeBuf[0], decodeBuf[1], buf[5:7])
				}
			}
			select {
			case <- synCh2:
				cont = false
			default:
			}
		}
		synCh1 <- true
	}()
	<-synCh1

	rAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:11880")
	if err != nil {
		log.Fatal("FindMounts err: ", err)
	}

	for i:=0; i<NUM_REPEAT_CMD; i++ {
		var cmd string
		if i < (NUM_REPEAT_CMD+1)/2 {
			cmd = ":e1" // get MC version
		} else {
			cmd = ":" 	// abort processing; this command will return "!0"
		}
		_, err = lConn.WriteToUDP([]byte(cmd+"\r"), rAddr)
		time.Sleep(TIMEOUT_REPLY)
	}
	synCh2 <- true
	<-synCh1

	aMounts := make([]*Mount, len(mounts))
	for _, m := range mounts {
		aMounts[len(aMounts)-1] = m
	}
	return aMounts
}
