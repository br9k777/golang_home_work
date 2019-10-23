package main

import (
	"fmt"
	"io/ioutil"
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

func TestGetWordsNumbers(t *testing.T) {
	// t.Skip()

	var content []byte
	var err error
	if content, err = ioutil.ReadFile(`words.txt`); err != nil {
		log.Error("Ошибка открытия файла", zap.Error(err))
		return
	}
	var words []string
	if words, err = GetSliceOfWords(content); err != nil {
		log.Error("Ошибка подсчета слов", zap.Error(err))
		return
	}
	for _, w := range words {
		fmt.Fprintf(os.Stdout, "%s\n", w)
	}
	return

}
