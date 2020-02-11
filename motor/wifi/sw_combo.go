package wifi

import (
	"time"
	"fmt"
)

// ask the mount to provide various parameters
func (mount *Mount) RetrieveMountParameters() (err0 error) {
	if mount.isInit { return }
	switch {
	case true:
		version1, err0 := mount.SWgetVersion(AXIS_RA_AZ)
		if err0 != nil { break }
		version2, err0 := mount.SWgetVersion(AXIS_DEC_ALT)
		if err0 != nil { break }
		
		cpr1, err0 := mount.SWgetCountsPerRevolution(AXIS_RA_AZ)
		if err0 != nil { break }
		cpr2, err0 := mount.SWgetCountsPerRevolution(AXIS_DEC_ALT)
		if err0 != nil { break }
		
		hsMult1, err0 := mount.SWgetHighSpeedRatio(AXIS_RA_AZ)
		if err0 != nil { break }
		hsMult2, err0 := mount.SWgetHighSpeedRatio(AXIS_DEC_ALT)
		if err0 != nil { break }
		
		mount.MCParamFrequency, err0 = mount.SWgetTimerFreq()
		if err0 != nil { break }
		
		mount.MCParamT1Tracking1X, err0 = mount.SWgetT1Tracking1X()
		if err0 != nil { break }

		resExt, err0 := mount.SWgetExtendedInfo(AXIS_1)
		if err0 != nil { break }
		mount.HasDualEncoder = resExt.SupportDualEncoder
		mount.HasPPEC = resExt.SupportPPEC
		mount.HasOriginalIndex = resExt.SupportOriginalIndex
		mount.HasEqAz = resExt.SupportEQAZ
		mount.HasPolarScopeLED = resExt.HasPolarScopeLED
		mount.HasAxisSeparateStart = resExt.IsAxisSeparateStart
		mount.HasTorqueSelection = resExt.HasTorqueSelection

		if version1 != version2 || cpr1 != cpr2 || hsMult1 != hsMult2 {
			err0 = &cmdError{ERR05_NOT_SUPPORTED, "The mount has different parameters for each AXIS; not supported"}
		} else {
			mount.MCversion = version2
			mount.MCParamCPR = cpr2
			mount.MCParamHighSpeedMult = hsMult2
			mount.isInit = true
		}
	}
	return
}

func abs(x int) int { if x<0 { return -x } else { return x } }
func sign(x int) int { if x<0 { return -1 } else { return 1 } }

// normalize tick position value to [-CPR/2, CPR/2]
func normalizeTickPosition(pos int, CPR int) int {
        pos %= CPR // limit the number of ticks to only one rotation
        if abs(pos) > CPR/2 {
                pos -= sign(pos)*CPR
                return normalizeTickPosition(pos, CPR)
        }
        return pos
}

// optimize tick increment, by changing CW>CCW or CCW>CW if the absolute value of the increment exceeds CPR/2
func optimizeTickIncrement(incr int, CPR int) int {
        incr %= CPR // limit the number of ticks to only one rotation
        if incr > CPR/2 || incr < -CPR/2 {
                incr -= sign(incr) * CPR
        }
        return incr
}


// move axis to a specific position.
// The target position sent to the mount is normalized to [-CPR/2 ... +CPR/2]
func (mount *Mount) GoToPosition(ax AXIS, tgtPos int) (err0 error) {
	tgtPos = normalizeTickPosition(tgtPos, mount.MCParamCPR)
	switch {
	case true:
		// should not reach this error, but let's leave it here for now
		if abs(tgtPos) >  mount.MCParamCPR/2 {
			err0 = &cmdError{ERR06_VALUE_TOO_LARGE, "Value too large ["+string(tgtPos)+"] limit ["+string(mount.MCParamCPR/2)+"]"}
		}

		err0 = mount.StopMotor(ax)
		if err0 != nil { break }

		crtPos, err0 := mount.SWgetPosition(ax)
		if err0 != nil { break }

		err0 = mount.GoToRelativeIncrement(ax, tgtPos - crtPos)
		if err0 != nil { break }
	}

	return
}

// move the axis in the CW (positive increment) or CCW (negative increment) direction, for a given number of ticks.
// If the originalRelativeIncrement is greater than CPR/2 (in absolute value), the actual value sent to the mount is normalized to [-CPR/2 ... +CPR/2]
func (mount *Mount) GoToRelativeIncrement(ax AXIS, originalRelativeIncrement int) (err0 error) {
	relativeIncrement := optimizeTickIncrement(originalRelativeIncrement, mount.MCParamCPR)
	switch {
	case true:
		// the increment should be already optimized, but yeah
		if relativeIncrement >  mount.MCParamCPR/2 || relativeIncrement < - mount.MCParamCPR/2 {
			err0 = &cmdError{ERR06_VALUE_TOO_LARGE, "Value too large ["+string(relativeIncrement)+"] limit ["+string(mount.MCParamCPR/2)+"]"}
		}

		err0 = mount.StopMotor(ax)
		if err0 != nil { break }

		crtPos, err0 := mount.SWgetPosition(ax)
		targetPos := normalizeTickPosition(crtPos + originalRelativeIncrement, mount.MCParamCPR)
		if err0 != nil { break }
		//fmt.Printf("CURRENT POS[%v] increment[%v] target[%v]\n", crtPos, relativeIncrement, crtPos+relativeIncrement)

		isCCW := relativeIncrement<0
		if relativeIncrement<0 { relativeIncrement = -relativeIncrement }

		isHighSpeed := relativeIncrement >  mount.MCParamCPR*5/360 // use highSpeed if the increment exceeds 5degrees

		var mm MotionMode
		mm.MmTrackingNotGoto = false
		mm.MmSpeedFast = isHighSpeed
		mm.MmSpeedMedium = false
		mm.MmSlowGoTo = false
		mm.IsCCW = isCCW
		mm.IsSouth = false
		mm.IsCoarseGoto = false
		err0 = mount.SWsetMotionMode(ax, mm)
		if err0 != nil { break }

		err0 = mount.SWsetGotoTarget(ax, relativeIncrement)
		if err0 != nil { break }
		err0 = mount.SWsetBrakeIncrement(ax, 3500)

		err0 = mount.SWstartMotion(ax)
		if err0 != nil {
			_ = mount.SWstopMotion(ax)
			break
		}

		for {
			v, err0 := mount.SWgetMotorStatus(ax)
			if err0 != nil || !v.IsRunning { break }
			<- time.After(TIMEOUT_REPLY)
			//crtPos, _ := mount.SWgetPosition(ax)
			//fmt.Printf("GOTO crtPos=%d\n", crtPos)
		}
		crtPos, _ = mount.SWgetPosition(ax)
		if crtPos != targetPos {
			fmt.Printf("RELATIVE GOTO: initialTarget[%d] currentPos[%d] diff[%d]\n", targetPos, crtPos, targetPos-crtPos)
		}
		if err0 != nil { break }
	}
	return
}

// Slew speed in degrees/second
type SLEW_SPEED float64
const (
	SLEW_SPEED_SIDERAL SLEW_SPEED	=	360.0/24/3600			// ideal sideral slewing rate for EQ mounts; for AltAz each axis has to be calculated depending on the current position
	SLEW_SPEED_LUNAR SLEW_SPEED	=	(360.0 - 360/28)/24/3600	// in 28days the moon completes a full rotation, towards the East
	SLEW_SPEED_0			=	SLEW_SPEED_SIDERAL / 2		// sideral_speed/2
	SLEW_SPEED_1			=	SLEW_SPEED_SIDERAL * 1		// sideral speed
	SLEW_SPEED_2			=	SLEW_SPEED_SIDERAL * 8		// 8 * sideral_speed
	SLEW_SPEED_3			=	SLEW_SPEED_SIDERAL * 16		// 16 * sideral_speed
	SLEW_SPEED_4			=	SLEW_SPEED_SIDERAL * 32		// 32 * sideral_speed
	SLEW_SPEED_5			=	SLEW_SPEED_SIDERAL * 64		// 64 * sideral_speed
	SLEW_SPEED_6			=	SLEW_SPEED_SIDERAL * 128	// 128 * sideral_speed
	SLEW_SPEED_7			=	SLEW_SPEED_SIDERAL * 400	// 400 * sideral_speed
	SLEW_SPEED_8			=	SLEW_SPEED_SIDERAL * 600	// 600 * sideral_speed
	SLEW_SPEED_9			=	SLEW_SPEED_SIDERAL * 800	// 800 * sideral_speed
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
		if speed <= SLEW_SPEED_SIDERAL/5000 {
			// too slow, stop it
			err0 = mount.SWstopMotion(ax)
			break
		}

		if mount.MCParamFrequency == 0 || mount.MCParamCPR == 0 || mount.MCParamHighSpeedMult == 0 {
			err0 = mount.RetrieveMountParameters()
		}
		if err0 != nil { break }

		if mount.MCParamFrequency == 0 || mount.MCParamCPR == 0 || mount.MCParamHighSpeedMult == 0 {
			// to do: check against T1 from [Inquire 1X Tracking Period][:D1]; it should be ticks(SLEW_SPEED_SIDERAL)
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

// stop motion on the given axis
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

// cycle the camera trigger (On/Off), with a [duration] delay.
func (mount *Mount) SetPhotoSwitch(duration time.Duration) (err0 error) {
	switch {
	case true:
		err0 = mount.SWsetSwitch(AXIS_BOTH, 1)
		if err0 != nil { break }
		if duration == 0 {
			// a bulb mode photo needs to calls: 1+0 ... 1+0
			err0 = mount.SWsetSwitch(AXIS_BOTH, 0)
		} else {
			// in KEEP mode we don't want to wait for the duration of the exposure, so we only check for an error on the first OPEN command; meh
			go func() {
				<- time.After(duration)
				_ = mount.SWsetSwitch(AXIS_BOTH, 0)
				//fmt.Println("switch off")
			}()
		}
	}
	return
}
