package model

import (
	"fmt"
)

type SimplifiedDate struct {
	Year   uint16
	Month  uint8
	Day    uint8
	Hour   uint8
	Minute uint8
}

func (a SimplifiedDate) Before(b SimplifiedDate) bool {
	return a.Year < b.Year || a.Year == b.Year && (a.Month < b.Month || a.Month == b.Month && (a.Day < b.Day || a.Day == b.Day && (a.Hour < b.Hour || a.Hour == b.Hour && a.Minute < b.Minute)))
}

func (a SimplifiedDate) Format() string {
	return fmt.Sprintf("%4d-%2d-%2dT%2d:%2d:00", a.Year, a.Month, a.Day, a.Hour, a.Minute)
}

func (a SimplifiedDate) AddMinute() SimplifiedDate {
	if a.Minute == 59 {
		a.Minute = 0
		a.Hour++
	} else {
		a.Minute++
	}

	return a
}

func (a SimplifiedDate) Unix() int32 {
	return (int32(a.Year)-1970)*365*24*60*60 + (int32(a.Month)-1)*31*24*60*60 + (int32(a.Day)-1)*24*60*60 + int32(a.Hour)*60*60 + int32(a.Minute)*60
}
