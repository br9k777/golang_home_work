package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
)

type wordStruct struct {
	count int
	word  string
}

//ReadWordsFromFile открываем файл достаем слова
func ReadWordsFromFile(filePath string) ([]string, error) {
	var file os.FileInfo
	var fileErr error
	if file, fileErr = os.Stat(filePath); os.IsNotExist(fileErr) || !file.Mode().IsRegular() {
		log.Errorf(`Not exist regular file %s`, filePath)
		return nil, fileErr
	}
	if fileErr != nil {
		log.Errorf(`Какая то новая злая ошибка прии чтении файла%s`, filePath)
		return nil, fileErr
	}
	var openedFile *os.File
	if openedFile, fileErr = os.Open(filePath); fileErr == nil {

		defer openedFile.Close()
		f := bufio.NewReader(openedFile)
		// сделаем задел сразу на 100 слов чтоб уменшить количество реолакаций
		var words = make([]string, 0, 100)
		var bufError error
		buf := make([]byte, 1024)
		for {
			//isPrefix возвращаемое вторым значением слшиком сложно для меня
			// пока будем сччиитать что слишком длиннных строк в файле нет
			buf, _, bufError = f.ReadLine()
			if bufError != nil {
				if bufError == io.EOF {
					break
				}
				return nil, bufError
			}
			//на уроке предлагали пакет string
			worsFromStrings := strings.Fields(string(buf))
			words = append(words, worsFromStrings...)

		}
		return words, nil
	}
	return nil, fileErr
}

//ShowHightFrequencyWords показать самые часто иссполььзуемые слова
func ShowHightFrequencyWords(filePath string) {
	var words []string
	var fileError error
	if words, fileError = ReadWordsFromFile(filePath); fileError != nil {
		log.Error(fileError)
		return
	}
	// а теперь пройдемся повторно учитывая знаки разделения
	var totalWordSum = make(map[string]*wordStruct, 100)
	for _, word := range words {
		newWords := regexp.MustCompile(`(\w|[а-яА-Я\_])+`).FindAllString(word, -1)
		for _, s := range newWords {
			// totalWordSum[s]++
			if w, ok := totalWordSum[s]; ok {
				w.count++
			} else {
				totalWordSum[s] = &wordStruct{
					count: 1,
					word:  s,
				}
			}
		}
	}
	//обратно в масив для сортировки
	var totalWordSumArray = make([]*wordStruct, 0, 300)
	for _, w := range totalWordSum {
		totalWordSumArray = append(totalWordSumArray, w)
	}
	sort.Slice(totalWordSumArray, func(i, j int) bool {
		return totalWordSumArray[i].count >= totalWordSumArray[j].count
	})
	for i, w := range totalWordSumArray {
		if i > 9 {
			break
		}
		fmt.Fprintf(os.Stderr, "%-3d Word: %s\tcount= %d\n", i+1, w.word, w.count)
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.00", FullTimestamp: true})
	log.SetOutput(os.Stderr)
}
