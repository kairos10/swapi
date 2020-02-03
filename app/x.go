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
			response, err := mounts[0].SendCmdSync(cmd)
			fmt.Printf("[%s]\t\tresponse[%s]\t\terr[%v]\n", cmd, response, err)
		}

		vs, err := mounts[0].SWgetVersion(wifi.AXIS_RA_AZ)
		if err != nil {
			fmt.Println("getVer error: ", err)
		} else {
			fmt.Println("ver: ", vs)
		}

		vi, err := mounts[0].SWgetCountsPerRevolution(wifi.AXIS_RA_AZ)
		if err != nil {
			fmt.Println("SWgetCountsPerRevolution error: ", err)
		} else {
			fmt.Println("SWgetCountsPerRevolution: ", vi)
		}

		vi, err = mounts[0].SWgetPosition(wifi.AXIS_RA_AZ)
		if err != nil {
			fmt.Println("SWgetPosition error: ", err)
		} else {
			fmt.Println("SWgetPosition: ", vi)
		}

		err = mounts[0].SWstopMotion(wifi.AXIS_BOTH)
		if err != nil {
			fmt.Println(".SWstopMotion error: ", err)
		} else {
			fmt.Println("SWstopMotion done")
		}
	} else {
		fmt.Println("nothing found!")
	}
}
