package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

const (
	long = 25.61 // local longitude
	lat  = 45.65 // local latitude
)

// EPOCH is the start date for our catalog
var EPOCH = time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC)
var yearsSinceEPOCH = float64(time.Now().Sub(EPOCH)/time.Hour) / 24 / 365.25636

// getLocalSideralTime calculates LST for the specified time and longitude (East positive)
func getLocalSideralTime(now time.Time, eastLong float64) float64 {
	L0 := float64(99.967794687)
	L1 := float64(360.985647366286)
	L2 := float64(2.907879e-13)
	L3 := float64(-5.302e-22)

	fromEpoch := now.Sub(EPOCH)
	dJ := float64(fromEpoch) / float64(time.Hour*24)
	dJ = dJ + .5
	dJ2 := dJ * dJ
	dJ3 := dJ2 * dJ
	theta := L0 + L1*dJ + L2*dJ2 + L3*dJ3 + eastLong
	over := int(theta) / 360
	theta = theta - float64(over*360)
	if theta < 0 {
		theta = theta + 360
	}
	return theta
}

// degreeToStr formats an angle° as H/M/S (if isTimeFmt is true), or as °/′/″
func degreeToStr(deg float64, isTimeFmt bool) (ret string) {
	if isTimeFmt {
		deg = deg / 15
	}
	h := int(deg)
	deg = deg - float64(h)
	deg = deg * 60
	m := int(deg)
	deg = deg - float64(m)
	deg = deg * 60
	//s := int(deg)
	//deg = deg - float64(s)
	s := deg

	if isTimeFmt {
		ret = fmt.Sprintf("%02dh%02dm%.1fs", h, m, s)
	} else {
		ret = fmt.Sprintf("%d°%d′%.1f″", h, m, s)
	}
	return
}

// getAltAz calculates the AltAz coordinates for a given hour angle and declination, at the specified latitude
func getAltAz(haDeg, decDeg, latDeg float64) (altDeg, azDeg float64) {
	ha, dec, lat := haDeg*math.Pi/180, decDeg*math.Pi/180, latDeg*math.Pi/180
	haSin, haCos := math.Sin(ha), math.Cos(ha)
	decSin, decCos := math.Sin(dec), math.Cos(dec)
	latSin, latCos := math.Sin(lat), math.Cos(lat)

	altSin := decSin*latSin + decCos*latCos*haCos
	alt := math.Asin(altSin)
	altCos := math.Cos(alt)
	altDeg = alt * 180 / math.Pi

	aCos := (decSin - altSin*latSin) / (altCos * latCos)
	a := math.Acos(aCos) * 180 / math.Pi
	if haSin < 0 {
		azDeg = a
	} else {
		azDeg = 360 - a
	}

	return
}

// getAtmosphericRefraction calculates the refraction R in degrees for a given true altitude, using the formula developed by Sæmundsson
func getAtmosphericRefraction(altDeg float64) (rDeg float64) {
	p := 10.3/(altDeg+5.11) + altDeg
	rAMin := 1.02 / math.Tan(p*math.Pi/180)
	rDeg = rAMin / 60
	return
}

type starInfo struct {
	RA, DEC     float64 // RA/Dec @J2000 in degrees
	Vmag        float32 // visual magnitude
	PMra, PMdec float64 // proper motion RA/Dec, expressed in milliarcseconds per year
	Parallax    float64 // parallax in mas
}

func (star starInfo) getRaDec() (crtRA, crtDec float64) {
	motionRaSec := yearsSinceEPOCH * star.PMra / 1000
	motionDecSec := yearsSinceEPOCH * star.PMdec / 1000

	distStarNCP := 3600*(90-star.DEC) - motionDecSec
	difRaDeg := math.Atan2(motionRaSec, distStarNCP) * 180 / math.Pi
	//fmt.Println("motionRA_sec", motionRaSec, "motionDec_sec", motionDecSec, "to NCP:", distStarNCP, "difRA:", difRaDeg)
	//fmt.Println(star.RA, degreeToStr(star.RA, true), star.DEC, degreeToStr(star.DEC, false))
	//fmt.Println(star.RA+difRaDeg, degreeToStr(star.RA+difRaDeg, true), star.DEC+motionDecSec/3600, degreeToStr(star.DEC+motionDecSec/3600, false))

	crtRA = star.RA + difRaDeg
	crtDec = star.DEC + motionDecSec/3600
	return
}
func (star starInfo) getDEC() float64 {
	return star.DEC + yearsSinceEPOCH*star.PMdec/1000/3600
}

type starCatalog map[string]*starInfo

func loadStarCatalog() starCatalog {
	starCatalog := make(starCatalog)
	catalogFileName := "starCatalog.json"
	file, err := os.Open(catalogFileName)
	if err != nil {
		log.Fatalf("error opening star catalog '%v' [%v]\n", catalogFileName, err)
	}
	dec := json.NewDecoder(file)
	err = dec.Decode(&starCatalog)
	file.Close()
	if err != nil {
		log.Fatalf("error decoding star catalog '%v' [%v]\n", catalogFileName, err)
	}

	return starCatalog
}

func main() {
	now := time.Now()
	theta := getLocalSideralTime(now, long)
	fmt.Println("local sideral time:", degreeToStr(theta, true))

	starCatalog := loadStarCatalog()

	star := starCatalog["polaris"]

	starRA, _ := star.getRaDec()
	starHA := theta - starRA
	fmt.Println("starRA:", degreeToStr(starRA, true), "star HourAngle: ", degreeToStr(starHA, true))

	alt, az := getAltAz(starHA, star.DEC, lat)
	fmt.Println("alt:", degreeToStr(alt, false), "az:", degreeToStr(az, false))
	r := getAtmosphericRefraction(alt)
	fmt.Println("AtmosphericRefraction(°):", r)
	fmt.Println("adjusted alt:", degreeToStr(alt-r, false))
}
