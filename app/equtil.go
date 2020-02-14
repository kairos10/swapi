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
	fmt.Printf("\t\tflip\t-\tPerform a meridian flip\n")
	//fmt.Printf("\t\ttrackSideral\t-\tStart tracking at 1x sideral speed\n")
	fmt.Println("")
}

func main() {
	cmd := ""
	if len(os.Args) < 2 {
		cmd = "help"
	} else if len(os.Args) == 2 {
		cmd = os.Args[1]
	}

	cmds := map[string]bool{"help":true, "setHome":true, "initEQ":true, "flip":true, "trackSideral":true}
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

		if cmd == "initEQ" {
			err = m.InitializeEQ()
			if err != nil {
				fmt.Println("eq init error: ", err)
			} else {
				fmt.Println("eq initialized: ")
			}
			os.Exit(0)
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
		} else if cmd == "setHome" {
			err = m.EqFlipMeridian(true)
			if err != nil {
				fmt.Println("flip error: ", err)
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
