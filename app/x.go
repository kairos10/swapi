package main

import (
	"fmt"
        "github.com/kairos10/swapi/motor/wifi"
)


func main() {
	mounts := wifi.FindMounts()
	if len(mounts) > 0 {
		for i:=0; i<len(mounts); i++ {
			m := mounts[i]
			fmt.Printf("found: %s:%d, version[%s], time[%s]\n", m.UDPAddr.IP, m.UDPAddr.Port, m.MCversion, m.DiscoveryTime)
		}

		cmds := []string { ":e1", ":", ":e2", ":e3", ":e1", ":", ":b1", ":a1" }
		for _, cmd := range cmds {
			response, err := wifi.SendCmdSync(mounts[0], cmd)
			fmt.Printf("[%s]\t\tresponse[%s]\t\terr[%v]\n", cmd, response, err)
		}
	} else {
		fmt.Println("nothing found!")
	}
}
