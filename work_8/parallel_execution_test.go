package parallelExecution

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestParallelWOrk(t *testing.T) {
	totalWorkNumber := 50
	rand.Seed(int64(time.Now().Nanosecond()))
	works := make([]func() error, totalWorkNumber)
	for i := 0; i < totalWorkNumber; i++ {
		works[i] = func() error {
			wait := rand.Intn(1e3)
			time.Sleep(time.Duration(wait) * time.Millisecond)
			if wait > 800 {
				fmt.Fprintf(os.Stderr, "Работа выолнялась долго %d млисекунд и вернула ошибку \n", wait)
				return errors.New(`Долгое выолнение`)
			}
			fmt.Fprintf(os.Stdout, "Работы выолнялась %d млисекунд\n", wait)
			return nil
		}
	}
	ParalelWorkMain(works, 10, 3)
	ParalelWorkMain(works, 20, 20)
}
