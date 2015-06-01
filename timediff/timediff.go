package timediff

import (
	"log"
	"strconv"
	"time"
)

//Gets two times in future and past and returns string describing their estimate difference
func GetDifference(future time.Time, past time.Time) string {
	t1 := past.Unix()
	t2 := future.Unix()

	if d := t2 - t1; d == 0 {
		return "Just now"
	} else if d := t2 - t1; d >= 1 && d < 60 {
		return strconv.Itoa(int(d)) + " " + appendS("Second", d) + " ago"
	} else if d := (t2 - t1) / 60; d >= 1 && d < 60 {
		return strconv.Itoa(int(d)) + " " + appendS("Minute", d) + " ago"
	} else if d := (t2 - t1) / (60 * 60); d >= 1 && d < 24 {
		return strconv.Itoa(int(d)) + " " + appendS("Hour", d) + " ago"
	} else if d := (t2 - t1) / (60 * 60 * 24); d >= 1 && d < 7 {
		return strconv.Itoa(int(d)) + " " + appendS("Day", d) + " ago"
	} else if d := (t2 - t1) / (60 * 60 * 24 * 7); d >= 1 && d < 5 {
		return strconv.Itoa(int(d)) + " " + appendS("Week", d) + " ago"
	} else if d := (t2 - t1) / (60 * 60 * 24 * 30); d >= 1 && d < 12 {
		return strconv.Itoa(int(d)) + " " + appendS("Month", d) + " ago"
	} else if d := (t2 - t1) / (60 * 60 * 24 * 365); d >= 1 {
		return strconv.Itoa(int(d)) + " " + appendS("Year", d) + " ago"
	} else {
		log.Printf("Couldn't get proper difference between %v and %v", future, past)
		return strconv.Itoa(int(t2-t1)) + " Seconds ago"
	}
}

func appendS(s string, d int64) string {
	if d > 1 {
		return s + "s"
	} else {
		return s
	}
}
