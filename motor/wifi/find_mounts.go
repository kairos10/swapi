package wifi

import (
	"log"
	"net"
	"time"
	"fmt"
)

func FindMounts() []*Mount {
	//lAddr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:11881")
	lAddr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:")
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

		synCh1 <- true
		for cont:=true; cont; {
			lConn.SetReadDeadline(time.Now().Add(TIMEOUT_REPLY))
			_, fromAddr, err := lConn.ReadFromUDP(buf)
			if err==nil && !lAddr.IP.Equal(fromAddr.IP) {
				// ignore packets received from self
				mount := mounts[string(fromAddr.IP)]
				if mount == nil {
					mount = new(Mount)
					mount.UDPAddr = *fromAddr
					mount.DiscoveryTime = time.Now()
					mounts[string(fromAddr.IP)] = mount
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

	rAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("255.255.255.255:%d", SW_UDP_PORT))
	if err != nil {
		log.Fatal("FindMounts err: ", err)
	}

	for i:=0; i<NUM_REPEAT_CMD; i++ {
		cmd := ":"
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
