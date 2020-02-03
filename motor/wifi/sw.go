package wifi

type AXIS int
const (
	AXIS_RA_AZ	=	1
	AXIS_DEC_ALT	=	2
	AXIS_BOTH	=	3
)

func SWgetVersion(AXIS) string {
	// :e[1,2,3]
	return ""
}

func SWgetCountsPerRevolution(AXIS) int {
	// :a[1,2]
	return 0
}
func SWgetTimerFreq(AXIS) int {
	// :b[1]
	return 0
}
