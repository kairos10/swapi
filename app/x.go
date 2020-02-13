package main

import (
	"fmt"
        "github.com/kairos10/swapi/motor/wifi"
	"time"
	//"net"
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

		//m := mounts[0]
		//fmt.Println("CPR: ", m.MCParamCPR/2)
		m := new(wifi.Mount)
		m.Resolve("192.168.3.22", 0)
		fmt.Println("mount: ", m)

		// CPR:  1036800
		// eq 0 0 = ra=518399=CPR/4  dec= 1
		// ra=12h=265698 .. cpr/4  dec=90'=408513
		//
		// home: -90', 90'
		// eq: ra: -180@E .. -0'@W
		// ra= -1  dec= 518401
		// ra= 0  dec= 518399 // start=0,0; connect eq=0,518399
		//
		// home: 270', 90'
		// eq: ra: 180@E .. 360@W
		// ra= 2073599  dec= 518401
		//
		// az home: 0, 0
		
		ra, _ := m.SWgetPosition(wifi.AXIS_RA_AZ)
		dec, _ := m.SWgetPosition(wifi.AXIS_DEC_ALT)
		fmt.Println("ra=", ra, " dec=", dec)

		//_ = m.SWsetPosition(wifi.AXIS_RA_AZ, 0)
		//_ = m.SWsetPosition(wifi.AXIS_DEC_ALT, 0)

		///////////////////////////////////
		_ = m.SWsetExtendedAttr(wifi.AXIS_1, wifi.SW_EXTENDED_ATTR_DUAL_ENCODER_DISABLE)
		fmt.Println("disabled")
		cpr, _ := m.SWgetCountsPerRevolution(wifi.AXIS_RA); fmt.Println("CPR: ", cpr)
		ra, _ = m.SWgetPosition(wifi.AXIS_RA_AZ); dec, _ = m.SWgetPosition(wifi.AXIS_DEC_ALT);
		raExt, _ := m.SWgetPositionExt(wifi.AXIS_RA_AZ); decExt, _ := m.SWgetPositionExt(wifi.AXIS_DEC_ALT)
		fmt.Println("EXT/disabled ra=", ra, " dec=", dec, " raExt=", raExt, " decExt=", decExt)
		fmt.Println("wait")
		<- time.After(10 * time.Second)
		//
		fmt.Println("disabled")
		ra, _ = m.SWgetPosition(wifi.AXIS_RA_AZ); dec, _ = m.SWgetPosition(wifi.AXIS_DEC_ALT);
		raExt, _ = m.SWgetPositionExt(wifi.AXIS_RA_AZ); decExt, _ = m.SWgetPositionExt(wifi.AXIS_DEC_ALT)
		fmt.Println("EXT/disabled ra=", ra, " dec=", dec, " raExt=", raExt, " decExt=", decExt)
		fmt.Println("wait")
		<- time.After(10 * time.Second)
		//
		//
		_ = m.SWsetExtendedAttr(wifi.AXIS_1, wifi.SW_EXTENDED_ATTR_DUAL_ENCODER_ENABLE)
		fmt.Println("enabled")
		cpr, _ = m.SWgetCountsPerRevolution(wifi.AXIS_RA); fmt.Println("CPR: ", cpr)
		cpr, _ = m.SWgetCountsPerRevolution(wifi.AXIS_DEC); fmt.Println("CPR: ", cpr)
		cpr, _ = m.SWgetCountsPerRevolution(wifi.AXIS_BOTH); fmt.Println("CPR: ", cpr)
		ra, _ = m.SWgetPosition(wifi.AXIS_RA_AZ); dec, _ = m.SWgetPosition(wifi.AXIS_DEC_ALT);
		raExt, _ = m.SWgetPositionExt(wifi.AXIS_RA_AZ); decExt, _ = m.SWgetPositionExt(wifi.AXIS_DEC_ALT)
		fmt.Println("EXT/disabled ra=", ra, " dec=", dec, " raExt=", raExt, " decExt=", decExt)
		fmt.Println("wait")
		<- time.After(10 * time.Second)
		//
		fmt.Println("enabled")
		ra, _ = m.SWgetPosition(wifi.AXIS_RA_AZ); dec, _ = m.SWgetPosition(wifi.AXIS_DEC_ALT);
		raExt, _ = m.SWgetPositionExt(wifi.AXIS_RA_AZ); decExt, _ = m.SWgetPositionExt(wifi.AXIS_DEC_ALT)
		fmt.Println("EXT/disabled ra=", ra, " dec=", dec, " raExt=", raExt, " decExt=", decExt)
		///////////////////////////////////

	} else {
		fmt.Println("nothing found!")
	}
}
