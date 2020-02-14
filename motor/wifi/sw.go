package wifi

/*
Low level functions for controlling wifi-aware SW mounts
*/

import (
	"fmt"
	"encoding/hex"
)

// mount axis
type AXIS int

// mount axis
const (
	AXIS_RA_AZ   AXIS = 1
	AXIS_DEC_ALT AXIS = 2
	AXIS_BOTH    AXIS = 3
	AXIS_RA      = AXIS_RA_AZ
	AXIS_AZ      = AXIS_RA_AZ
	AXIS_1       = AXIS_RA_AZ
	AXIS_DEC     = AXIS_DEC_ALT
	AXIS_ALT     = AXIS_DEC_ALT
	AXIS_2       = AXIS_DEC_ALT
)

// byte position: (int) to motor controller representation
const (
	D1	=	0x0000f0
	D2	=	0x00000f
	D3	=	0x00f000
	D4	=	0x000f00
	D5	=	0xf00000
	D6	=	0x0f0000
)
// bit position: (int) to motor controller representation
const (
	D2_B0	=	1 << iota
	D2_B1
	D2_B2
	D2_B3

	D1_B0
	D1_B1
	D1_B2
	D1_B3

	D4_B0
	D4_B1
	D4_B2
	D4_B3

	D3_B0
	D3_B1
	D3_B2
	D3_B3
)

func (m *Mount) swSend(cmd byte, ax AXIS, cmdParam *int) (ret0 int, err0 error) {
	cmdAxisType := map[byte]byte{
		'b':1,'C':1,'D':1,'n':1,'N':1,
		'a':2,'c':2,'f':2,'g':2,'h':2,'H':2,'i':2,'j':2,'k':2,'m':2,
		'A':3,'B':3,'d':3,'e':3,'E':3,'F':3,'G':3,'I':3,'J':3,'K':3,'L':3,'M':3,'O':3,'P':3,'q':3,'Q':3,'r':3,'R':3,'s':3,'S':3,'T':3,'U':3,'V':3,'W':3,'z':3,
	}
	cmdParamLen := map[byte]byte{
		'a':0,'b':0,'c':0,'d':0,'D':0,'e':0,'f':0,'F':0,'g':0,'h':0,'i':0,'j':0,'J':0,'K':0,'L':0,'m':0,'n':0,'r':0,'s':0,'z':0,
		'B':1,'k':1,'O':1,'P':1,
		'A':2,'G':2,'N':2,'R':2,'V':2,
		'C':4,'Q':4,
		'E':6,'H':6,'I':6,'M':6,'q':6,'S':6,'T':6,'U':6,'W':6,
	}
	cmdResponseLen := map[byte]byte{
		'A':0,'B':0,'E':0,'F':0,'G':0,'H':0,'I':0,'J':0,'K':0,'L':0,'M':0,'O':0,'P':0,'R':0,'S':0,'T':0,'U':0,'V':0,'W':0,'z':0,
		'g':2,'r':2,
		'f':3,
		'a':6,'b':6,'c':6,'d':6,'D':6,'e':6,'h':6,'i':6,'j':6,'k':6,'m':6,'q':6,'s':6,

	}
	if 
		 (cmdAxisType[cmd] != 0) &&
		 (cmdAxisType[cmd] != 1 || ax!=AXIS_RA_AZ) &&
		 (cmdAxisType[cmd] != 2 || ax!=AXIS_RA_AZ && ax!=AXIS_DEC_ALT) &&
		 (cmdAxisType[cmd] != 3 || ax!=AXIS_RA_AZ && ax!=AXIS_DEC_ALT && ax!=AXIS_BOTH) &&
		true {
		err0 = &cmdError{ERR01_AXIS, "AXIS parameter error"}
	} else {
		var cmdStr string
		switch cmdParamLen[cmd] {
		case 0:
			if cmdParam != nil {
				err0 = &cmdError{ERR03_PARAM, "parameter error ["+string(*cmdParam)+"] for cmd ["+string(cmd)+"]"}
				return
			}
			cmdStr = fmt.Sprintf(":%c%d", cmd, ax)
		case 1:
			if *cmdParam != *cmdParam & 0xf {
				err0 = &cmdError{ERR03_PARAM, "parameter error ["+string(*cmdParam)+"] for cmd ["+string(cmd)+"]"}
				return
			}
			cmdStr = fmt.Sprintf(":%c%d%1X", cmd, ax, *cmdParam & 0xf)
		case 2:
			if *cmdParam != *cmdParam & 0xff {
				err0 = &cmdError{ERR03_PARAM, "parameter error ["+string(*cmdParam)+"] for cmd ["+string(cmd)+"]"}
				return
			}
			cmdStr = fmt.Sprintf(":%c%d%02X", cmd, ax, *cmdParam & 0xff)
		case 4:
			if *cmdParam != *cmdParam & 0xffff {
				err0 = &cmdError{ERR03_PARAM, "parameter error ["+string(*cmdParam)+"] for cmd ["+string(cmd)+"]"}
				return
			}
			cmdStr = fmt.Sprintf(":%c%d%02X%02X", cmd, ax, *cmdParam & 0xff, *cmdParam&0xff00 >> 8)
		case 6:
			cmdStr = fmt.Sprintf(":%c%d%02X%02X%02X", cmd, ax, *cmdParam & 0xff, *cmdParam&0xff00 >> 8, *cmdParam&0xff0000 >> 16)
		}

		b, err := m.SendCmdSync(cmdStr)
		if err != nil {
			err0 = err
		} else if len(b) != int(cmdResponseLen[cmd]) + 1 {
			err0 = &cmdError { ERR02_RESP_LEN, "["+"x"+"] invalid response ["+string(b)+"]" }
		} else {
			var decodeBuf [3]byte
			if len(b[1:]) /2 * 2 < len(b[1:]) {
				b = b[:len(b)+1]
				b[len(b)-1]='0'
			}
			hex.Decode(decodeBuf[:], b[1:])
			for x := ((len(b)+1)/2)-2; x>=0; x-- {
				ret0 = (ret0<<8) + int(decodeBuf[x])
			}
		}
		//fmt.Println("XXX: ", cmdStr, " -- ", string(b))
	}
	return
}

func (mount *Mount) SWgetVersion(ax AXIS) (ret0 string, err0 error) {
	v, err := mount.swSend('e', ax, nil)
	if err0=err; err0 == nil {
		ret0 = fmt.Sprintf("[%d.%d %02X]", v&0xff, v&0xff00>>8, v&0xff0000>>16)
	}
	return
}

func (mount *Mount) SWgetCountsPerRevolution(ax AXIS) (int, error) {
	return mount.swSend('a', ax, nil)
}

func (mount *Mount) SWgetT1Tracking1X() (int, error) {
	return mount.swSend('D', AXIS_1, nil)
}

func (mount *Mount) SWgetHighSpeedRatio(ax AXIS) (int, error) {
	return mount.swSend('g', ax, nil)
}

func (mount *Mount) SWgetPosition(ax AXIS) (ret0 int, err0 error) {
	ret0, err0 = mount.swSend('j', ax, nil)
	if err0 == nil {
		ret0 -= 0x800000
	}
	return
}

// get the position of the secondary  encoder. it's unclear wether the CPR for the secondary encoder is available.
func (mount *Mount) SWgetPositionExt(ax AXIS) (ret0 int, err0 error) {
	ret0, err0 = mount.swSend('d', ax, nil)
	if err0 == nil {
		ret0 -= 0x800000
	}
	return
}

func (mount *Mount) SWstopMotion(ax AXIS) (err0 error) {
	_, err0 = mount.swSend('K', ax, nil)
	return
}

func (mount *Mount) SWsetSwitch(ax AXIS, switchPos int) (err0 error) {
	_, err0 = mount.swSend('O', ax, &switchPos)
	return
}

func (mount *Mount) SWstartMotion(ax AXIS) (err0 error) {
	_, err0 = mount.swSend('J', ax, nil)
	return
}

func (mount *Mount) SWgetTimerFreq() (int, error) {
	return mount.swSend('b', 1, nil)
}

// current motor state, retrieved with the SWgetExtendedInfo method
type MotorStatus struct {
	IsTracking bool
	IsCCW bool
	IsFast bool
	IsRunning bool
	IsBlocked bool
	IsInitDone bool
	IsLevelSwitchOn bool
}
func (ms MotorStatus) String() string {
	ret := ""
	if ms.IsTracking { ret += "tracking" } else { ret += "goto" }
	if ms.IsCCW { ret += " CCW" } else { ret += " CW" }
	if ms.IsFast { ret += " fast" } else { ret += " slow" }
	if ms.IsRunning { ret += " running" } else { ret += " stopped" }
	if ms.IsBlocked { ret += " blocked" } else { ret += " normal" }
	if ms.IsInitDone { ret += " initDone" } else { ret += " initNot" }
	if ms.IsLevelSwitchOn { ret += " levelSwitchOn" } else { ret += " levelSwitchOff" }
	return ret
}
func (mount *Mount) SWgetMotorStatus(ax AXIS) (ret0 MotorStatus, err0 error) {
	v, err := mount.swSend('f', ax, nil)
	if err0=err; err0 == nil {
		ret0.IsTracking = (v&D1_B0 > 0)
		ret0.IsCCW = (v&D1_B1 > 0)
		ret0.IsFast = (v&D1_B2 > 0)

		ret0.IsRunning = (v&D2_B0 > 0)
		ret0.IsBlocked = (v&D2_B1 > 0)

		ret0.IsInitDone= (v&D3_B0 > 0)
		ret0.IsLevelSwitchOn= (v&D3_B1 > 0)
	}
	return
}


// decoded status information, retrieved with SWgetExtendedInfo
type ExtendedStatus struct {
	IsPecTrainingOn bool
	IsPecTrackingOn bool

	SupportDualEncoder bool
	SupportPPEC bool
	SupportOriginalIndex bool
	SupportEQAZ bool

	HasPolarScopeLED bool
	IsAxisSeparateStart bool
	HasTorqueSelection bool
}
func (es ExtendedStatus) String() string {
	ret := ""
	if es.IsPecTrainingOn { ret += " PecTrainingOn" } else { ret += " PecTrainingOff" }
	if es.IsPecTrackingOn { ret += " PecTrackingOn" } else { ret += " PecTrackingOff" }

	if es.SupportDualEncoder { ret += " SupportDualEncoder" } else { ret += " NoDualEncoder" }
	if es.SupportPPEC { ret += " SupportPPEC" } else { ret += " NoPPEC" }
	if es.SupportOriginalIndex { ret += " SupportOriginalIndex" } else { ret += " NoOriginalIndex" }
	if es.SupportEQAZ { ret += " SupportEqAz" } else { ret += " NoEqAz" }

	if es.HasPolarScopeLED { ret += " HasPolarScopeLED" } else { ret += " NoPolarScopeLED" }
	if es.IsAxisSeparateStart { ret += " IsAxisSeparateStart" } else { ret += " NoAxisSeparateStart" }
	if es.HasTorqueSelection { ret += " HasTorqueSelection" } else { ret += " NoTorqueSelection" }

	return ret
}
func (mount *Mount) SWgetExtendedInfo(ax AXIS) (ret0 ExtendedStatus, err0 error) {
	ext := 0x000001
	v, err := mount.swSend('q', ax, &ext)
	if err0=err; err0 == nil {
		ret0.IsPecTrainingOn = (v&D1_B0 > 0)
		ret0.IsPecTrackingOn = (v&D1_B1 > 0)

		ret0.SupportDualEncoder = (v&D2_B0 > 0)
		ret0.SupportPPEC = (v&D2_B1 > 0)
		ret0.SupportOriginalIndex = (v&D2_B2 > 0)
		ret0.SupportEQAZ = (v&D2_B3 > 0)

		ret0.HasPolarScopeLED = (v&D3_B0 > 0)
		ret0.IsAxisSeparateStart = (v&D3_B1 > 0)
		ret0.HasTorqueSelection = (v&D3_B2 > 0)
	}
	return
}

// commands for SWsetExtendedAttr()
type SW_EXTENDED_ATTR int
const (
	SW_EXTENDED_ATTR_PEC_TRAINING_START	SLEW_SPEED	= 0x000000
	SW_EXTENDED_ATTR_PEC_TRAINING_CANCEL	SLEW_SPEED	= 0x000001
	//
	SW_EXTENDED_ATTR_PEC_TRACKING_START	SLEW_SPEED	= 0x000002
	SW_EXTENDED_ATTR_PEC_TRACKING_CANCEL	SLEW_SPEED	= 0x000003
	//
	SW_EXTENDED_ATTR_DUAL_ENCODER_ENABLE	SLEW_SPEED	= 0x000004 // use both primary and secondary encoders; the CPR does not change
	SW_EXTENDED_ATTR_DUAL_ENCODER_DISABLE	SLEW_SPEED	= 0x000005
	//
	SW_EXTENDED_ATTR_FULL_TORQUE_ENABLE	SLEW_SPEED	= 0x000106
	SW_EXTENDED_ATTR_FULL_TORQUE_DISABLE	SLEW_SPEED	= 0x000006
	//
	SW_EXTENDED_ATTR_SLEW_STRIDE		SLEW_SPEED	= 0x000007
	SW_EXTENDED_ATTR_INDEX_POSITION_RESET	SLEW_SPEED	= 0x000008
	SW_EXTENDED_ATTR_FLUSH_TO_ROM		SLEW_SPEED	= 0x000009
)
func (mount *Mount) SWsetExtendedAttr(ax AXIS, speedId SW_EXTENDED_ATTR) (err0 error) {
	_, err0 = mount.swSend('W', ax, (*int)(&speedId))
	return
}

func (mount *Mount) SWsetGotoTargetRelative(ax AXIS, ticks int) (err0 error) {
	_, err0 = mount.swSend('H', ax, &ticks)
	return
}

func (mount *Mount) SWsetGotoTargetAbs(ax AXIS, tickPos int) (err0 error) {
	_, err0 = mount.swSend('S', ax, &tickPos)
	return
}

func (mount *Mount) SWsetBrakeIncrement(ax AXIS, ticks int) (err0 error) {
	_, err0 = mount.swSend('M', ax, &ticks)
	return
}

func (mount *Mount) SWsetInitializationDone(ax AXIS) (err0 error) {
	_, err0 = mount.swSend('F', ax, nil)
	return
}

func (mount *Mount) SWsetPosition(ax AXIS, numTicks int) (err0 error) {
	err0 = mount.StopMotor(ax)
	if err0 == nil {
		numTicks += 0x800000
		_, err0 = mount.swSend('E', ax, &numTicks)
	}
	return
}

func (mount *Mount) SWsetStepPeriod(ax AXIS, clockDivider int) (err0 error) {
	_, err0 = mount.swSend('I', ax, &clockDivider)
	return
}

type SW_AUTOGUIDE_SPEED_FRACTION_ID int
const (
	SW_AUTOGUIDE_SPEED_100PCT	SW_AUTOGUIDE_SPEED_FRACTION_ID = 0 // sw constant: autoguide speed id, 1x
	SW_AUTOGUIDE_SPEED_75PCT	SW_AUTOGUIDE_SPEED_FRACTION_ID = 1 // sw constant: autoguide speed id, 0.75x
	SW_AUTOGUIDE_SPEED_50PCT	SW_AUTOGUIDE_SPEED_FRACTION_ID = 2 // sw constant: autoguide speed id, 0.5x
	SW_AUTOGUIDE_SPEED_25PCT	SW_AUTOGUIDE_SPEED_FRACTION_ID = 3 // sw constant: autoguide speed id, 0.25x
	SW_AUTOGUIDE_SPEED_12PCT	SW_AUTOGUIDE_SPEED_FRACTION_ID = 4 // sw constant: autoguide speed id, 0.125x
)
func (mount *Mount) SWsetAutoguideSpeed(ax AXIS, speedId SW_AUTOGUIDE_SPEED_FRACTION_ID) (err0 error) {
	_, err0 = mount.swSend('P', ax, (*int)(&speedId))
	return
}

// motor motion parameters, passed to the mount with the SWgetExtendedInfo method
type MotionMode struct {
	MmTrackingNotGoto bool
	MmSpeedFast bool		// speedMedium takes precedence; speedLow is selected if (!speedMedium && !speedFast)
	MmSpeedMedium bool		// overrides speedFast
	MmSlowGoTo bool

	IsCCW bool			// CW otherwise
	IsSouth bool			// North otherwise
	IsCoarseGoto bool		// Normal otherwise
}
func (mm MotionMode) String() string {
	ret := ""
	if mm.MmTrackingNotGoto { ret += "tracking" } else { ret += "goto" }
	if mm.MmSpeedMedium { ret += " medSpeed" } else if mm.MmSpeedFast { ret += " fast" } else { ret += " slow" }
	if mm.MmSlowGoTo { ret += " slowGoto" }
	if mm.IsCCW { ret += " CCW" } else { ret += " CW" }
	if mm.IsSouth { ret += " South" } else { ret += " North" }
	if mm.IsCoarseGoto { ret += " isCoarseGoto" }
	return ret
}
//  Channel will always be set to Tracking Mode after stopped
func (mount *Mount) SWsetMotionMode(ax AXIS, mm MotionMode) (err0 error) {
	// 00 						: goto	fast CW
	// 10 D1_B0					: tracking slow CW
	// 11 D1_B0		D2_B0			: tracking slow CCW
	// 20 		D1_B1				: goto slow CW
	// 21 		D1_B1	D2_B0			: goto slow CCW
	// 30 D1_B0	D1_B1				: tracking fast CW
	// 31 D1_B0	D1_B1	D2_B0			: tracking fast CCW
	mode := 0
	if mm.MmTrackingNotGoto { mode |= D1_B0 }
	if mm.MmTrackingNotGoto == mm.MmSpeedFast { mode |= D1_B1 }
	if mm.MmSpeedMedium { mode |= D1_B2 }
	if mm.MmSlowGoTo{ mode |= D1_B3 }

	if mm.IsCCW { mode |= D2_B0 }
	if mm.IsSouth { mode |= D2_B1 }
	if mm.IsCoarseGoto { mode |= D2_B2 }

	//fmt.Printf("MOTION MODE: %02X\n", mode)
	_, err0 = mount.swSend('G', ax, &mode)
	return
}
