package utils

import (
	"testing"
	"time"
)

func TestDateCalculator(t *testing.T) {
	dateEndOW := time.Date(2021, 11, 1, 0, 0, 0, 0, time.UTC)
	wantMondayEOW := "2021-11-01"
	wantSundayEOW := "2021-11-07"

	gotMondayEOW, gotSundayEOW := dateCalculator(dateEndOW)

	dateEndOM := time.Date(2021, 11, 29, 0, 0, 0, 0, time.UTC)
	wantMondayEOM := "2021-11-29"
	wantSundayEOM := "2021-11-30"
	gotMondayEOM, gotSundayEOM := dateCalculator(dateEndOM)

	if gotMondayEOM != wantMondayEOM {
		t.Errorf("got %q want %q", gotMondayEOM, wantMondayEOM)
	}
	if gotSundayEOM != wantSundayEOM {
		t.Errorf("got %q want %q", gotSundayEOM, wantSundayEOM)
	}
	if gotMondayEOW != wantMondayEOW {
		t.Errorf("got %q want %q", gotMondayEOW, wantMondayEOW)
	}
	if gotSundayEOW != wantSundayEOW {
		t.Errorf("got %q want %q", gotSundayEOW, wantSundayEOW)
	}
}
