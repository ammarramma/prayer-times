package main

import (
	"fmt"
	"time"

	"github.com/mnadev/adhango/pkg/calc"
	"github.com/mnadev/adhango/pkg/data"
	"github.com/mnadev/adhango/pkg/util"
)

const (
	// Koordinat Lhokseumawe, Aceh
	LATITUDE  = 5.1801
	LONGITUDE = 97.1406
)

func main() {
	// Setup koordinat Lhokseumawe
	coords, err := util.NewCoordinates(LATITUDE, LONGITUDE)
	if err != nil {
		fmt.Printf("Error creating coordinates: %v\n", err)
		return
	}

	// Setup tanggal hari ini
	today := time.Now()
	// WIB timezone (UTC+7)
	loc, _ := time.LoadLocation("Asia/Jakarta")
	todayWIB := today.In(loc)
	dateComponents := data.NewDateComponents(todayWIB)

	// Setup calculation parameters untuk Kemenag Indonesia
	// Kemenag: Fajr 20°, Isha 18°, Madhab Shafii
	params := calc.NewCalculationParametersBuilder().
		SetFajrAngle(20.0).
		SetIshaAngle(18.0).
		SetMadhab(calc.SHAFI_HANBALI_MALIKI).
		Build()

	// Kalkulasi waktu shalat
	prayerTimes, err := calc.NewPrayerTimes(coords, dateComponents, params)
	if err != nil {
		fmt.Printf("Error calculating prayer times: %v\n", err)
		return
	}

	// Print hasil
	fmt.Println("===========================================")
	fmt.Println("      JADWAL WAKTU SHALAT")
	fmt.Println("      Lhokseumawe, Aceh")
	fmt.Printf("      %s\n", todayWIB.Format("Monday, 02 January 2006"))
	fmt.Println("===========================================")
	fmt.Printf("Subuh  : %s\n", prayerTimes.Fajr.In(loc).Format("15:04"))
	fmt.Printf("Terbit : %s\n", prayerTimes.Sunrise.In(loc).Format("15:04"))
	fmt.Printf("Dzuhur : %s\n", prayerTimes.Dhuhr.In(loc).Format("15:04"))
	fmt.Printf("Ashar  : %s\n", prayerTimes.Asr.In(loc).Format("15:04"))
	fmt.Printf("Maghrib: %s\n", prayerTimes.Maghrib.In(loc).Format("15:04"))
	fmt.Printf("Isya   : %s\n", prayerTimes.Isha.In(loc).Format("15:04"))
	fmt.Println("===========================================")
	fmt.Println("Metode: Kemenag RI (Fajr 20°, Isha 18°)")
	fmt.Println("Madhab Asr: Shafii")
	fmt.Println("===========================================")
}
