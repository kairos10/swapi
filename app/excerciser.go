package main

import (
	"fmt"
        "github.com/kairos10/swapi/motor/wifi"
	//"time"
)


func main() {
	mounts := wifi.FindMounts()
	if len(mounts) > 0 {
		for i:=0; i<len(mounts); i++ {
			m := mounts[i]
			fmt.Printf("found: %s:%d, version[%s], time[%s]\n", m.UDPAddr.IP, m.UDPAddr.Port, m.MCversion, m.DiscoveryTime)

			err := m.RetrieveMountParameters()
			if err != nil {
				fmt.Println("RetrieveMountParameters error: ", err)
			}
		}

		cmds := []string { ":e1", ":", ":e2" }
		for _, cmd := range cmds {
			response, err := mounts[0].SendCmdSync(cmd)
			fmt.Printf("[%s]\t\tresponse[%s]\t\terr[%v]\n", cmd, response, err)
		}

		/*
		vi, err = mounts[0].SWgetPosition(wifi.AXIS_RA_AZ)
		if err != nil {
			fmt.Println("SWgetPosition error: ", err)
		} else {
			fmt.Println("SWgetPosition: ", vi)
		}
		err = mounts[0].SWsetPosition(wifi.AXIS_RA_AZ, 1234)
		if err != nil {
			fmt.Println("SWsetPosition error: ", err)
		} else {
			fmt.Println("SWgetPosition: done")
		}
		vi, err = mounts[0].SWgetPosition(wifi.AXIS_RA_AZ)
		if err != nil {
			fmt.Println("SWgetPosition error: ", err)
		} else {
			fmt.Println("SWgetPosition: ", vi)
		}
		//*/

		err := mounts[0].SWstopMotion(wifi.AXIS_BOTH)
		if err != nil {
			fmt.Println("SWstopMotion error: ", err)
		} else {
			fmt.Println("SWstopMotion done")
		}

		vx, err := mounts[0].SWgetMotorStatus(wifi.AXIS_RA_AZ)
		if err != nil {
			fmt.Println("SWgetMotorStatus error: ", err)
		} else {
			fmt.Println("SWgetMotorStatus: ", vx)
		}

		vy, err := mounts[0].SWgetExtendedInfo(wifi.AXIS_RA_AZ)
		if err != nil {
			fmt.Println("SWgetExtendedInfo error: ", err)
		} else {
			fmt.Println("SWgetExtendedInfo: ", vy)
		}

		/*
		err = mounts[0].SetSlewRate(wifi.AXIS_RA_AZ, -wifi.SLEW_SPEED_5, 1500 * time.Millisecond)
		if err != nil {
			fmt.Println("SetSlewRate error: ", err)
		} else {
			fmt.Println("SetSlewRate: done")
		}
		<- time.After(1 * time.Second)
		*/

		err = mounts[0].GoToPosition(wifi.AXIS_RA_AZ, 500)
		if err != nil {
			fmt.Println("GoToPosition error: ", err)
		} else {
			fmt.Println("GoToPosition: done")
		}

		vi, err := mounts[0].SWgetPosition(wifi.AXIS_RA_AZ)
		if err != nil {
			fmt.Println("SWgetPosition error: ", err)
		} else {
			fmt.Println("SWgetPosition: ", vi)
		}
		err = mounts[0].GoToRelativeIncrement(wifi.AXIS_RA_AZ, 2073600-10000)
		if err != nil {
			fmt.Println("GoToRelativeIncrement error: ", err)
		} else {
			fmt.Println("GoToRelativeIncrement: done")
		}
		vi, err = mounts[0].SWgetPosition(wifi.AXIS_RA_AZ)
		if err != nil {
			fmt.Println("SWgetPosition error: ", err)
		} else {
			fmt.Println("SWgetPosition: ", vi)
		}
	} else {
		fmt.Println("nothing found!")
	}
}
