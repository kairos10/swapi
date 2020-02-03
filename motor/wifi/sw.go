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

func (mount *Mount) SWgetVersion(ax AXIS) (ret0 string, err0 error) {
	// :e[1,2,3]
	if ax!=AXIS_RA_AZ && ax!=AXIS_DEC_ALT && ax!=AXIS_BOTH {
		err0 = &cmdError{ERR01_AXIS, "AXIS parameter error"}
	} else {
		cmd := fmt.Sprintf(":e%d", ax)
		b, err := mount.SendCmdSync(cmd)
		if err != nil {
			err0 = err
		} else if len(b) == 7 {
			decBuf := make([]byte, 2)
			hex.Decode(decBuf[:], b[1:5])
			ret0 = fmt.Sprintf("%d.%d %s", decBuf[0], decBuf[1], b[5:7])
		} else {
			err0 = &cmdError{ERR02_RESP_LEN, "invalid response ["+string(b)+"]"}
		}
	}
	return
}

func (mount *Mount) SWgetCountsPerRevolution(ax AXIS) (ret0 int, err0 error) {
	// :a[1,2]
	if ax!=AXIS_RA_AZ && ax!=AXIS_DEC_ALT {
		err0 = &cmdError{ERR01_AXIS, "AXIS parameter error"}
	} else {
		cmd := fmt.Sprintf(":a%d", ax)
		b, err := mount.SendCmdSync(cmd)
		if err != nil {
			err0 = err
		} else if len(b) == 7 {
			decBuf := make([]byte, 3)
			hex.Decode(decBuf[:], b[1:7])
			ret0 = 256*(256*int(decBuf[2])+int(decBuf[1]))+int(decBuf[0])
		} else {
			err0 = &cmdError{ERR02_RESP_LEN, "invalid response ["+string(b)+"]"}
		}
	}
	return
}

func (mount *Mount) SWgetPosition(ax AXIS) (ret0 int, err0 error) {
	// :a[1,2]
	if ax!=AXIS_RA_AZ && ax!=AXIS_DEC_ALT {
		err0 = &cmdError{ERR01_AXIS, "AXIS parameter error"}
	} else {
		cmd := fmt.Sprintf(":a%d", ax)
		b, err := mount.SendCmdSync(cmd)
		if err != nil {
			err0 = err
		} else if len(b) == 7 {
			decBuf := make([]byte, 3)
			hex.Decode(decBuf[:], b[1:7])
			ret0 = 256*(256*int(decBuf[2])+int(decBuf[1]))+int(decBuf[0]) - 0x800000
		} else {
			err0 = &cmdError{ERR02_RESP_LEN, "invalid response ["+string(b)+"]"}
		}
	}
	return
}

func (mount *Mount) SWstopMotion(ax AXIS) (err0 error) {
	// :K[1,2,3]
	if ax!=AXIS_RA_AZ && ax!=AXIS_DEC_ALT && ax!=AXIS_BOTH {
		err0 = &cmdError{ERR01_AXIS, "AXIS parameter error"}
	} else {
		cmd := fmt.Sprintf(":K%d", ax)
		b, err := mount.SendCmdSync(cmd)
		if err != nil {
			err0 = err
		} else if len(b) != 1 {
			err0 = &cmdError{ERR02_RESP_LEN, "invalid response ["+string(b)+"]"}
		}
	}
	return
}

