package end2end

import (
	"testing"
	"time"
)

type fakeClock struct {
	timeNow time.Time
}

func (f *fakeClock) Now() time.Time {
	return f.timeNow
}

func (f *fakeClock) SetTime(t time.Time) {
	f.timeNow = t
}

func (f *fakeClock) SetTimeRFC3339(t *testing.T, str string) {
	newNow, err := time.Parse(time.RFC3339, str)
	if err != nil {
		t.Fatal(err)
	}
	f.SetTime(newNow)
}

func (f *fakeClock) Advance(t *testing.T, str string) {
	duration, err := time.ParseDuration(str)
	if err != nil {
		t.Fatal(err)
	}
	f.timeNow = f.timeNow.Add(duration)
}
