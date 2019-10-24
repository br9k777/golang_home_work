package logOtus

import (
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
)

var (
	log *zap.Logger
)

func init() {
	var err error
	if log, err = zap.NewProduction(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(log)
}

//HwAccepted структура с сообщением о том что ДЗ приянто
type HwAccepted struct {
	ID    int
	Grade int
}

//LogIt логируем что ДЗ принято
func (hw *HwAccepted) LogIt(w io.Writer) {
	fmt.Fprintf(w, "%s accepted %d %d\n", time.Now().Format("2006-01-02"), hw.ID, hw.Grade)
}

//HwSubmitted структура с сообщением о том что ДЗ отправлено студентом
type HwSubmitted struct {
	ID      int
	Code    string
	Comment string
}

//LogIt логируем что ДЗ отправлено
func (hw *HwSubmitted) LogIt(w io.Writer) {
	fmt.Fprintf(w, "%s submitted %d \"%s\"\n", time.Now().Format("2006-01-02"), hw.ID, hw.Comment)
}

//OtusEvent набор функций необходимых для логирования Otus событий
type OtusEvent interface {
	LogIt(w io.Writer)
}

//Log логируем событие Otus
func Log(e OtusEvent, w io.Writer) {
	e.LogIt(w)
	return
}
