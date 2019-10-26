package parallelExecution

import (
	"fmt"
	"os"
	"sync"
)

//ParalelWorkMain основная функция в которую передаем slcie задач
//мыксимальное число выполняемых паралельно задач и максимальное число ошибок
func ParalelWorkMain(works []func() error, maxParalelWorks, maxErrors int) {
	fmt.Fprintf(os.Stdout, "Функция запущена со слайсом %d задач, паралельно будет выполнятся %d, до %d ошибок\n", len(works), maxParalelWorks, maxErrors)
	worksInProgess := make(chan func() error, maxParalelWorks-1)
	worksInProgessCounter := make(chan interface{}, maxParalelWorks-1)
	chanelForEmergancyStop := make(chan interface{})

	var errorsCount int
	var m sync.RWMutex

	go func() {
	OUT:
		for i, work := range works {
			select {
			case worksInProgess <- work: // отсюда будем брать задачи
				fmt.Fprintf(os.Stdout, "В очередь задач добавлена задача номер %d\n", i+1)
				worksInProgessCounter <- nil // а этот канал будет нас блокировать пока не выполнится одна из задач

			case <-chanelForEmergancyStop:
				fmt.Fprintf(os.Stderr, "Закончили передавать задачи в канал потому что слишком много ошибок\n")
				break OUT
			}
		}
		close(worksInProgess)
	}()

	for work := range worksInProgess {
		go func(workFunc func() error) {
			if err := workFunc(); err != nil {
				m.RLock()
				errorsCount++
				if errorsCount == maxErrors {
					fmt.Fprintf(os.Stderr, "Количество ошибок по работам равно %d выходим из программы\n", maxErrors)
					chanelForEmergancyStop <- nil
				}
				m.RUnlock()
			}
			<-worksInProgessCounter
		}(work)
	}
	for {
		if len(worksInProgessCounter) == 0 || errorsCount >= maxErrors {
			break
		}
	}
}
