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

	poleSign := float64(1)
	if star.DEC < 0 { poleSign = -poleSign }
	distStarNCP := 3600*(90-poleSign*star.DEC) - poleSign*motionDecSec
	difRaRad := math.Atan2(motionRaSec, distStarNCP)
	difRaDeg := difRaRad * 180 / math.Pi
	crtRA = star.RA + difRaDeg

	crtDec = star.DEC + motionDecSec/3600
	// ? substract RA's contribution towards DEC
	dDecRa := math.Sqrt(star.DEC*star.DEC + difRaDeg*difRaDeg) - poleSign*star.DEC
	crtDec = crtDec - poleSign*dDecRa
	return
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
	now = time.Date(2020, time.July, 31, 1, 0, 0, 0, time.Local)
	theta := getLocalSideralTime(now, long)
	fmt.Println("NOW:", now.In(time.UTC))
	fmt.Println("Years since J2000", yearsSinceEPOCH)
	fmt.Println("local sideral time:", degreeToStr(theta, true))

	starCatalog := loadStarCatalog()

	fmt.Printf("%v|%v|%v|%v|%v|%v|%v|%v|%v|%v\n", "name", "RA J2000", "DEC J2000", "PM RA", "PM Dec", "Parallax", "VMag", "crt RA", "crt DEC", "atm refr.")
	//for _, name := range []string{"polaris", "yildun", "hd5848", "ngc188", "hd5914", "hd5621", "hip85822", "hip37391", "hip115746"} {
	for _, name := range []string{"polaris", "yildun", "hd5848", "ngc188", "hd5914", "hd5621", "hip37391", "hip115746", "miaplacidus"} {
		star := starCatalog[name]
		starRA, starDEC := star.getRaDec()
		starHA := theta - starRA
		alt, _ := getAltAz(starHA, star.DEC, lat)
		r := getAtmosphericRefraction(alt)
		fmt.Printf("%v|%v|%v|%v|%v|%v|%v|%v|%v|%v\n", name, star.RA, star.DEC, star.PMra, star.PMdec, star.Parallax, star.Vmag, starRA, starDEC, r)
		/*
		fmt.Println("\nSTAR: ", name, "ra0:", star.RA, "raNow:", starRA)
		fmt.Printf("PmRA[%v] PmDec[%v] J2000_RA[%v] J2000_DEC[%v] RA[%v] DEC[%v]\n", star.PMra, star.PMdec, degreeToStr(star.RA, true), degreeToStr(star.DEC, false), degreeToStr(starRA, true), degreeToStr(starDEC, false))
		fmt.Println("HA: ", degreeToStr(starHA, true))
		fmt.Println("AtmosphericRefraction(°):", r)
		fmt.Println("alt:", degreeToStr(alt, false), "adjAlt:", degreeToStr(alt+r, false), "az:", degreeToStr(az, false))
		//*/
	}
}
