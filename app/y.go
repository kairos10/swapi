package main

import (
	"fmt"
        "github.com/kairos10/swapi/motor/wifi"
	"time"
)


func main() {
	m := new(wifi.Mount)
	m.Resolve("192.168.3.22", 0)

	ra, _ := m.SWgetPosition(wifi.AXIS_RA_AZ)
	dec, _ := m.SWgetPosition(wifi.AXIS_DEC_ALT)
	fmt.Println("ra=", ra, " dec=", dec)
	<- time.After(time.Second)

	_ = m.InitializeEQ()
	ra, _ = m.SWgetPosition(wifi.AXIS_RA_AZ)
	dec, _ = m.SWgetPosition(wifi.AXIS_DEC_ALT)
	fmt.Println("ra=", ra, " dec=", dec)

	m.GoToRelativeIncrement(wifi.AXIS_RA, m.MCParamCPR/15)
	m.GoToRelativeIncrement(wifi.AXIS_DEC, m.MCParamCPR/20)
	ra, _ = m.SWgetPosition(wifi.AXIS_RA_AZ)
	dec, _ = m.SWgetPosition(wifi.AXIS_DEC_ALT)
	fmt.Println("ra=", ra, " dec=", dec)

	m.EqFlipMeridian(true)
	ra, _ = m.SWgetPosition(wifi.AXIS_RA_AZ)
	dec, _ = m.SWgetPosition(wifi.AXIS_DEC_ALT)
	fmt.Println("ra=", ra, " dec=", dec)
}
