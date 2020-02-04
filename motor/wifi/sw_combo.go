package wifi

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
