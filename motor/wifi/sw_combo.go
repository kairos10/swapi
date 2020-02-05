package wifi

import (
	"time"
)

func (mount *Mount) RetrieveMountParameters() (err0 error) {
	switch {
	case true:
		version1, err0 := mount.SWgetVersion(AXIS_RA_AZ)
		if err0 != nil { break }
		version2, err0 := mount.SWgetVersion(AXIS_DEC_ALT)
		if err0 != nil { break }
		freq, err0 := mount.SWgetTimerFreq()
		if err0 != nil { break }
		cpr1, err0 := mount.SWgetCountsPerRevolution(AXIS_RA_AZ)
		if err0 != nil { break }
		cpr2, err0 := mount.SWgetCountsPerRevolution(AXIS_DEC_ALT)
		if err0 != nil { break }
		hsMult1, err0 := mount.SWgetHighSpeedRatio(AXIS_RA_AZ)
		if err0 != nil { break }
		hsMult2, err0 := mount.SWgetHighSpeedRatio(AXIS_DEC_ALT)
		if err0 != nil { break }

		if version1 != version2 || cpr1 != cpr2 || hsMult1 != hsMult2 {
			err0 = &cmdError{ERR05_NOT_SUPPORTED, "The mount has different parameters for each AXIS; not supported"}
		} else {
			mount.MCversion = version2
			mount.MCParamFrequency = freq
			mount.MCParamCPR = cpr2
			mount.MCParamHighSpeedMult = hsMult2
		}
	}
	return
}

type SLEW_SPEED float64
const (
	SLEW_SPEED_SIDERAL SLEW_SPEED	=	360.0/24/3600
	SLEW_SPEED_0			=	SLEW_SPEED_SIDERAL / 2
	SLEW_SPEED_1			=	SLEW_SPEED_SIDERAL * 1
	SLEW_SPEED_2			=	SLEW_SPEED_SIDERAL * 8
	SLEW_SPEED_3			=	SLEW_SPEED_SIDERAL * 16
	SLEW_SPEED_4			=	SLEW_SPEED_SIDERAL * 32
	SLEW_SPEED_5			=	SLEW_SPEED_SIDERAL * 64
	SLEW_SPEED_6			=	SLEW_SPEED_SIDERAL * 128
	SLEW_SPEED_7			=	SLEW_SPEED_SIDERAL * 400
	SLEW_SPEED_8			=	SLEW_SPEED_SIDERAL * 600
	SLEW_SPEED_9			=	SLEW_SPEED_SIDERAL * 800
)
// set slew rate (in deg/sec); usefull absolute values could be between 0.002 .. 3.3 deg/sec
// positive speed move the axis CW, while a negative speed results in a CCW rotation
// if duration == 0, the slew will continue after the function returns and needs to be stopped by calling SWstopMotion
func (mount *Mount) SetSlewRate(ax AXIS, speed SLEW_SPEED, duration time.Duration) (err0 error) {
	isCCW := speed<0.0
	if speed < 0 { speed = -speed }
	isHighSpeed := speed > 400*360/24/3600 // x400 sideral rate

	switch {
	case true:
		if mount.MCParamFrequency == 0 || mount.MCParamCPR == 0 || mount.MCParamHighSpeedMult == 0 {
			err0 = mount.RetrieveMountParameters()
		}
		if err0 != nil { break }
		if mount.MCParamFrequency == 0 || mount.MCParamCPR == 0 || mount.MCParamHighSpeedMult == 0 {
			err0 = &cmdError{ERR05_NOT_SUPPORTED, "The mount cannot be used to set the slew rate"}
			break
		}
		clockDivider := int(SLEW_SPEED(mount.MCParamFrequency * 360 / mount.MCParamCPR) / speed)
		if isHighSpeed {
			clockDivider *= mount.MCParamHighSpeedMult
		}
		err0 = mount.SWsetStepPeriod(ax, clockDivider)
		//fmt.Printf("speed=%f highSpeed=%v preset=%d", speed, isHighSpeed, clockDivider)
		if err0 != nil { break }

		var mm MotionMode
		mm.MmTrackingNotGoto = true
		mm.MmSpeedFast = isHighSpeed
		mm.MmSpeedMedium = false
		mm.MmSlowGoTo = false
		mm.IsCCW = isCCW
		mm.IsSouth = false
		mm.IsCoarseGoto = false
		err0 = mount.SWsetMotionMode(ax, mm)
		if err0 != nil { break }

		err0 = mount.SWstartMotion(ax)
		if err0 != nil {
			_ = mount.SWstopMotion(ax)
			break
		}

		if duration != 0 {
			<- time.After(duration)
			err0 = mount.SWstopMotion(ax)
			if err0 != nil {
				_ = mount.SWstopMotion(AXIS_BOTH)
				break
			}
		}
	}

	return
}

func (mount *Mount) StopMotor(ax AXIS) (err0 error) {
	if ax == AXIS_BOTH {
		_ = mount.SWstopMotion(ax)
		err0 = mount.StopMotor(AXIS_DEC_ALT)
		if err0 == nil {
			err0 = mount.StopMotor(AXIS_RA_AZ)
		}
		return
	}

	isRunning := true
	for x:=0; err0==nil && x<NUM_REPEAT_CMD; x++ {
		v, err0 := mount.SWgetMotorStatus(ax)
		if err0 == nil && !v.IsRunning {
			isRunning = false
			break
		}
		err0 = mount.SWstopMotion(ax)
	}
	if isRunning && err0==nil {
		err0 = &cmdError{ERR04_NA, "could not stop motor"}
	}
	return
}


