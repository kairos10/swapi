package main

import (
	"fmt"
        "github.com/kairos10/swapi/motor/wifi"
	"time"
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

		cmds := []string { ":e1", ":", ":e2", ":e3", ":e1", ":", ":b1", ":a1" }
		for _, cmd := range cmds {
			response, err := mounts[0].SendCmdSync(cmd)
			fmt.Printf("[%s]\t\tresponse[%s]\t\terr[%v]\n", cmd, response, err)
		}

		vs, err := mounts[0].SWgetVersion(wifi.AXIS_RA_AZ)
		if err != nil {
			fmt.Println("SWgetVersion error: ", err)
		} else {
			fmt.Println("SWgetVersion: ", vs)
		}

		vi, err := mounts[0].SWgetCountsPerRevolution(wifi.AXIS_RA_AZ)
		if err != nil {
			fmt.Println("SWgetCountsPerRevolution error: ", err)
		} else {
			fmt.Println("SWgetCountsPerRevolution: ", vi)
		}

		vi, err = mounts[0].SWgetTimerFreq()
		if err != nil {
			fmt.Println("SWgetTimerFreq error: ", err)
		} else {
			fmt.Println("SWgetTimerFreq: ", vi)
		}

		vi, err = mounts[0].SWgetHighSpeedRatio(wifi.AXIS_RA_AZ)
		if err != nil {
			fmt.Println("SWgetHighSpeedRatio error: ", err)
		} else {
			fmt.Println("SWgetHighSpeedRatio: ", vi)
		}

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

		var mm wifi.MotionMode
		mm.MmTrackingNotGoto = true
		mm.MmSpeedFast = true
		mm.MmSpeedMedium = false
		mm.MmSlowGoTo = false
		mm.IsCCW = true
		mm.IsSouth = false
		mm.IsCoarseGoto = false
		err = mounts[0].SWsetMotionMode(wifi.AXIS_BOTH, mm)
		if err != nil {
			fmt.Println("SWsetMotionMode error: ", err)
		} else {
			fmt.Println("SWsetMotionMode done: ", mm)
		}

		err = mounts[0].SWstartMotion(wifi.AXIS_BOTH)
		if err != nil {
			fmt.Println("SWstartMotion error: ", err)
		} else {
			fmt.Println("SWstartMotion done")
		}
		time.Sleep(2*time.Second)

		err = mounts[0].SWstopMotion(wifi.AXIS_BOTH)
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

		err = mounts[0].SWstopMotion(wifi.AXIS_BOTH)
		if err != nil {
			fmt.Println("SWstopMotion error: ", err)
		} else {
			fmt.Println("SWstopMotion: done")
		}

		vy, err := mounts[0].SWgetExtendedInfo(wifi.AXIS_RA_AZ)
		if err != nil {
			fmt.Println("SWgetExtendedInfo error: ", err)
		} else {
			fmt.Println("SWgetExtendedInfo: ", vy)
		}

		err = mounts[0].SetSlewRate(wifi.AXIS_RA_AZ, -wifi.SLEW_SPEED_5, 1500 * time.Millisecond)
		if err != nil {
			fmt.Println("SetSlewRate error: ", err)
		} else {
			fmt.Println("SetSlewRate: done")
		}
	} else {
		fmt.Println("nothing found!")
	}
}
