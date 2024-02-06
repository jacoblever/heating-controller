package common

import (
	"testing"
	"time"
)

type FakeClock struct {
	TimeNow time.Time
}

func (f *FakeClock) Now() time.Time {
	return f.TimeNow
}

func (f *FakeClock) SetTime(t time.Time) {
	f.TimeNow = t
}

func (f *FakeClock) SetTimeRFC3339(t *testing.T, str string) {
	newNow, err := time.Parse(time.RFC3339, str)
	if err != nil {
		t.Fatal(err)
	}
	f.SetTime(newNow)
}

func (f *FakeClock) Advance(t *testing.T, str string) {
	duration, err := time.ParseDuration(str)
	if err != nil {
		t.Fatal(err)
	}
	f.TimeNow = f.TimeNow.Add(duration)
}
