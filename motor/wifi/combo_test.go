package wifi_test

import "github.com/kairos10/swapi/motor/wifi"
import (
	"fmt"
	"log"
)

func ExampleMount_InitializeEQ() {
	m := new(wifi.Mount)
	m.Resolve("192.168.4.1", -1)

	err := m.InitializeEQ() // home position for EQ is [ra=0; dec=CPR/4];
	if err != nil {
		log.Println("The mount may have been initialized from another aplication; err: ", err)
	}
}

func ExampleMount_EqFlipMeridian() {
	m := new(wifi.Mount)
	m.Resolve("192.168.4.1", -1)

	_ = m.InitializeEQ()

	// move the mount away from the home position
	m.GoToRelativeIncrement(wifi.AXIS_RA, m.MCParamCPR/15)
	m.GoToRelativeIncrement(wifi.AXIS_DEC, -m.MCParamCPR/20)

	ra, _ := m.SWgetPosition(wifi.AXIS_RA_AZ)
	dec, _ := m.SWgetPosition(wifi.AXIS_DEC_ALT)
	fmt.Println("BEFORE FLIP: ra=", ra, " dec=", dec)
	m.EqFlipMeridian(true, true) // do the flip anyway, trusting that the mount was initialized in EQ mode; correct the RA position to account for the time spend to perform the flip
	ra, _ = m.SWgetPosition(wifi.AXIS_RA_AZ)
	dec, _ = m.SWgetPosition(wifi.AXIS_DEC_ALT)
	fmt.Println("AFTER FLIP: ra=", ra, " dec=", dec)
}
