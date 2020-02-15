package main

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/kairos10/swapi/motor/wifi"
)

func doHelp(av0 string) {
	fmt.Println("Usage: ", filepath.Base(av0), " <command>")
	fmt.Printf("\tAvailable commands\n")
	fmt.Printf("\t\thelp\t-\tDisplay this message\n")
	fmt.Printf("\t\tsetHome\t-\tSet Equatorial HOME to the current position\n")
	fmt.Printf("\t\tinitEQ\t-\tSet Equatorial HOME to the current position if the mount is not initialized, otherwise do nothing\n")
	fmt.Printf("\t\tflip\t-\tPerform a meridian flip, correcting the RA for the time spent to perform the flip\n")
	fmt.Printf("\t\tstop\t-\tStop motors\n")
	fmt.Printf("\t\ttrackSideral -\tStart tracking at 1x sideral speed (might need to stop the motors first)\n")
	fmt.Println("")
}

func main() {
	cmd := ""
	if len(os.Args) < 2 {
		cmd = "help"
	} else if len(os.Args) == 2 {
		cmd = os.Args[1]
	}

	cmds := map[string]bool{"help":true, "setHome":true, "initEQ":true, "flip":true, "trackSideral":true, "stop":true, "noop":true}
	if cmd == "help" {
		doHelp(os.Args[0])
		os.Exit(0)
	} else if !cmds[cmd] {
		fmt.Println("Command not recognized [", cmd, "]")
		doHelp(os.Args[0])
		os.Exit(1)
	}

	mounts := wifi.FindMounts()
	if len(mounts) > 0 {
		m := mounts[0]
		err := m.RetrieveMountParameters()
		m.LoggerFunc = func(s string) { fmt.Printf(s) }

		if cmd == "noop" {
			//
		} else if cmd == "initEQ" {
			err = m.InitializeEQ()
			if err != nil {
				fmt.Println("initEQ error: ", err)
				os.Exit(2)
			}
		} else if cmd == "setHome" {
			err = m.ReInitializeEQ()
			if err != nil {
				fmt.Println("setHome error: ", err)
				os.Exit(2)
			}
		} else if cmd == "flip" {
			err = m.EqFlipMeridian(true, true)
			if err != nil {
				fmt.Println("flip error: ", err)
				os.Exit(2)
			}
		} else if cmd == "stop" {
			err = m.StopMotor(wifi.AXIS_BOTH)
			if err != nil {
				fmt.Println("stop error: ", err)
				os.Exit(2)
			}
		} else if cmd == "trackSideral" {
			err = m.SetSlewRate(wifi.AXIS_RA_AZ, wifi.SLEW_SPEED_SIDERAL, 0)
			if err != nil {
				fmt.Println("slewSideral error: ", err)
				os.Exit(2)
			}
		}
	} else {
		fmt.Println("no MC found on the network")
		os.Exit(2)
	}
}
