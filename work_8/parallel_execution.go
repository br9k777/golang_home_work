package parallelExecution

import (
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CustomTimeEncoder function of own formulating time for output to the log
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func init() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = CustomTimeEncoder
	config.Encoding = "console"
	var err error
	var logger *zap.Logger
	if logger, err = config.Build(); err != nil {
		zap.L().Fatal("Logger create error", zap.Error(err))
	}
	zap.ReplaceGlobals(logger)
}

//ParalelWorkMain основная функция в которую передаем slcie задач
//мыксимальное число выполняемых паралельно задач и максимальное число ошибок
func ParalelWorkMain(works []func() error, maxParalelWorks, maxErrors int) {
	fmt.Fprintf(os.Stdout, "Функция запущена со слайсом %d задач, паралельно будет выполнятся %d, до %d ошибок\n", len(works), maxParalelWorks, maxErrors)
	worksChanel := make(chan func() error)
	errorChanel := make(chan error, maxParalelWorks-1)
	chanelForEmergancyStop := make(chan interface{})

	var (
		work                  func() error
		workerWait, errorWait sync.WaitGroup
	)
	errorWait.Add(1)
	go func() {
		errorCounter(errorChanel, chanelForEmergancyStop, maxErrors)
		errorWait.Done()
	}()

	for i := 0; i < maxParalelWorks; i++ {
		workerWait.Add(1)
		go func() {
			worker(worksChanel, chanelForEmergancyStop, errorChanel)
			workerWait.Done()
		}()
	}
OUT:
	for _, work = range works {
		select {
		case worksChanel <- work:
		case <-chanelForEmergancyStop:
			break OUT
		}
	}
	close(worksChanel)
	workerWait.Wait()
	close(errorChanel)
	errorWait.Wait()
}

func worker(workChanel <-chan func() error, stopChanel <-chan interface{}, errorChanel chan<- error) {
	var ok bool
	var workFunc func() error
	for {
		select {
		case workFunc, ok = <-workChanel:
			if !ok {
				return
			}
			errorChanel <- workFunc()
		case <-stopChanel:
			return
		}
	}
}

func errorCounter(errorChanel <-chan error, stopChan chan<- interface{}, maximumNumberOfErrors int) {
	var workDone, errCount int
	var err error
	for err = range errorChanel {
		workDone++
		if err != nil {
			errCount++
		}
		if errCount == maximumNumberOfErrors-1 {
			fmt.Printf("Достигнуто максимальное число ошибок %d\n", maximumNumberOfErrors)
			close(stopChan)
			for range errorChanel {
			}
			return
		}
	}
}
