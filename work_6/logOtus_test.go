package logOtus

import (
	"fmt"
	"os"
	"testing"

	"go.uber.org/zap"
)

//не разобрался еще как лучше включать dev логирование для тестов
func TestSetLogger(t *testing.T) {
	var err error
	if log, err = zap.NewDevelopment(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func TestLog(t *testing.T) {
	// t.Skip()
	events := []interface{}{
		HwSubmitted{
			ID:      6715,
			Code:    `-------code-------`,
			Comment: `Pls see it`,
		},
		HwAccepted{
			ID:    6715,
			Grade: 3,
		},
		HwSubmitted{
			ID:      8431,
			Code:    `-------code2-------`,
			Comment: `Just did it`,
		},
		HwAccepted{
			ID:    8431,
			Grade: 4,
		},
		HwSubmitted{
			ID:      6124,
			Code:    `-------code3-------`,
			Comment: `Done`,
		},
		HwAccepted{
			ID:    6124,
			Grade: 5,
		},
	}
	var eOtus OtusEvent
	// var ok bool
	for _, e := range events {
		// e.(OtusEvent).LogIt(os.Stdout)
		switch e.(type) {
		case HwSubmitted:
			submit := e.(HwSubmitted)
			eOtus = &submit
		case HwAccepted:
			accept := e.(HwAccepted)
			eOtus = &accept
		}
		eOtus.LogIt(os.Stdout)
	}

	// var event OtusEvent = &HwSubmitted{
	// 	ID:      6715,
	// 	Code:    `-------code-------`,
	// 	Comment: `Pls see it`,
	// }
	// Log(event, os.Stdout)
	return
}
