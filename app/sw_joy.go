package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kairos10/swapi/motor/wifi"
	"github.com/simulatedsimian/joystick"
)

type axisStatus struct {
	sync.RWMutex
	alt, az float64 // +/-[0-1] values
}

const (
	// JSPollInterval gives the joystick polling interval
	JSPollInterval = 100 * time.Millisecond
)

func initJs() (joystick.Joystick, error) {
	jsid := 0
	js, err0 := joystick.Open(jsid)
	switch {
	default:
		if err0 != nil {
			break
		}
		axisCount := js.AxisCount()
		if axisCount < 2 {
			err0 = errors.New("joystick " + js.Name() + " does not expose at least two axis")
			break
		}
	}
	return js, err0
}
func startPolling(f func(), tOut time.Duration) {
	ticker := time.NewTicker(JSPollInterval)
	var tOutChan <-chan time.Time
	if tOut != 0 {
		tOutChan = time.After(tOut)
	}
	for cont := true; cont; {
		if tOut != 0 {
			select {
			case <-ticker.C:
				f()
			case <-tOutChan:
				cont = false
				break
			}
		} else {
			<-ticker.C
			f()
		}
	}
}

func runJoy(axStat *axisStatus) (err0 error) {
	js, err := initJs()
	if err != nil {
		log.Printf("A joystick/gamepad has not been detected [%s]\n", err)
	}
	for js == nil {
		<-time.After(5 * time.Second)
		js, err = initJs()
		if err != nil {
			fmt.Printf(".")
		}
	}

	axisCount := js.AxisCount()
	axMinMax := make([]struct{ min, max, zero, zeroMin, zeroMax, selected int }, axisCount)

	fmt.Println("please leave all sticks centered for 3 seconds")
	startPolling(func() {
		jinfo, err := js.Read()
		if err != nil {
			log.Println("JOYSTICK error: ", err)
		}
		isFirstRun := true
		for axis := 0; axis < axisCount; axis++ {
			a := jinfo.AxisData[axis]
			aMinMax := &axMinMax[axis]
			if isFirstRun {
				aMinMax.zeroMin = a
				aMinMax.min = a
				aMinMax.zeroMax = a
				aMinMax.max = a
				isFirstRun = false
				continue
			}
			if a < aMinMax.zeroMin {
				aMinMax.zeroMin = a
			}
			if a > aMinMax.zeroMax {
				aMinMax.zeroMax = a
			}
		}
	}, 2*time.Second)
	fmt.Println("zero point evaluated")

	fmt.Println("please move one stick, then leave it free for 1 second")
	selectedAxis := make(map[int]bool)
	for hasMoved := false; hasMoved || len(selectedAxis) < 2; {
		hasMoved = false

		startPolling(func() {
			jinfo, err := js.Read()
			if err != nil {
				log.Println("JOYSTICK error: ", err)
			}
			for axis := 0; axis < axisCount; axis++ {
				a := jinfo.AxisData[axis]
				aMinMax := &axMinMax[axis]
				if a < aMinMax.min {
					aMinMax.min = a
					hasMoved = true
				}
				if a > aMinMax.max {
					aMinMax.max = a
					hasMoved = true
				}
				if aMinMax.min < aMinMax.zeroMin && aMinMax.max > aMinMax.zeroMax {
					selectedAxis[axis] = true
				}
			}
		}, 1*time.Second)

	}
	fmt.Println("max throw evaluated", selectedAxis)
	<-time.After(1 * time.Second)

	go func() {
		startPolling(func() {
			jinfo, err := js.Read()
			if err != nil {
				log.Println("JOYSTICK error: ", err)
			}
			axID := 0
			for axis := range selectedAxis {
				axID++
				a := jinfo.AxisData[axis]
				aMinMax := &axMinMax[axis]
				if a < aMinMax.min {
					aMinMax.min = a
				}
				if a > aMinMax.max {
					aMinMax.max = a
				}
				v := float64(0)
				if a < aMinMax.zeroMin {
					v = -float64(a) / float64(aMinMax.min-aMinMax.zeroMin)
				} else if a > aMinMax.zeroMax {
					v = float64(a) / float64(aMinMax.max-aMinMax.zeroMax)
				}
				axStat.Lock()
				if axID == 1 {
					axStat.alt = v
				} else if axID == 2 {
					axStat.az = v
				}
				axStat.Unlock()
			}
		}, 0)
	}()

	return
}

func main() {
	var crtAxStat axisStatus
	var swMount *wifi.Mount
	synCh := make(chan bool)

	// prefetch mounts
	go func() {
		for j := 0; j < 5; j++ {
			mounts := wifi.FindMounts()
			if len(mounts) > 0 {
				swMount = mounts[0]
				swMount.RetrieveMountParameters()
				fmt.Println("mount found")
				break
			}
		}
		synCh <- true
	}()

	err := runJoy(&crtAxStat)
	if err != nil {
		return
	}
	fmt.Println("joystick ok")

	<-synCh
	if swMount == nil {
		fmt.Println("no SW mount found on the network; retrying...")
	}
	for swMount == nil {
		mounts := wifi.FindMounts()
		if len(mounts) < 1 {
			fmt.Printf("*")
			<-time.After(5 * time.Second)
		} else {
			swMount = mounts[0]
		}
	}
	fmt.Println(swMount)

	swMount.StopMotor(wifi.AXIS_BOTH)
	defer swMount.StopMotor(wifi.AXIS_BOTH) // we don't want to accidentally leave the motors running on exit
	startPolling(func() {
		crtAxStat.Lock()
		axAlt := crtAxStat.alt
		axAz := crtAxStat.az
		crtAxStat.Unlock()

		if axAlt == 0 {
			swMount.SWstopMotion(wifi.AXIS_DEC_ALT)
		} else {
			speed := wifi.SLEW_SPEED(axAlt * float64(wifi.SLEW_SPEED_7-wifi.SLEW_SPEED_0))
			swMount.SetSlewRate(wifi.AXIS_DEC_ALT, speed, 0)
		}
		if axAz == 0 {
			swMount.SWstopMotion(wifi.AXIS_RA_AZ)
		} else {
			speed := wifi.SLEW_SPEED(axAz * float64(wifi.SLEW_SPEED_6-wifi.SLEW_SPEED_0))
			swMount.SetSlewRate(wifi.AXIS_RA_AZ, speed, 0)
		}
	}, 0)
	<-make(chan bool)
}
