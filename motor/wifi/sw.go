package wifi

/*
Low level functions for controlling wifi-aware SW mounts
*/

import (
	"fmt"
	"encoding/hex"
)

type AXIS int

const (
	AXIS_RA_AZ   = 1
	AXIS_DEC_ALT = 2
	AXIS_BOTH    = 3
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
	if true &&
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
	v, err0 := mount.swSend('e', ax, nil)
	if err0 == nil {
		ret0 = fmt.Sprintf("[%d.%d %02X]", v&0xff, v&0xff00>>8, v&0xff0000>>16)
	}
	return
}

func (mount *Mount) SWgetCountsPerRevolution(ax AXIS) (int, error) {
	return mount.swSend('a', ax, nil)
}

func (mount *Mount) SWgetPosition(ax AXIS) (ret0 int, err0 error) {
	ret0, err0 = mount.swSend('j', ax, nil)
	if err0 == nil {
		ret0 -= 0x800000
	}
	return
}

func (mount *Mount) SWstopMotion(ax AXIS) (err0 error) {
	_, err0 = mount.swSend('K', ax, nil)
	return
}

func (mount *Mount) SWgetTimerFreq() (int, error) {
	return mount.swSend('b', 1, nil)
}

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
	v, err0 := mount.swSend('f', ax, nil)
	// 111 > 1110 > 1011
	// 611 > 6110 > 1061
	if err0 == nil {
		d1 := v&0xf0 >> 4
		d2 := v&0x0f
		d3 := v&0xf000 >> 12

		ret0.IsTracking = (d1&1 > 0)
		ret0.IsCCW = (d1&2 > 0)
		ret0.IsFast = (d1&4 > 0)

		ret0.IsRunning = (d2&1 > 0)
		ret0.IsBlocked = (d2&2 > 0)

		ret0.IsInitDone= (d3&1 > 0)
		ret0.IsLevelSwitchOn= (d3&2 > 0)
	}
	return
}

