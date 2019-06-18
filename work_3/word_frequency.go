package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"sort"

	log "github.com/sirupsen/logrus"
)

var (
	logErr log.Logger
)

type wordStruct struct {
	count int
	word  string
}

//ShowHightFrequencyWords основная функция читает файл находит слова, подсчитывает их
func ShowHightFrequencyWords(filePath string) error {

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		logErr.Error(err)
		return err
	}
	return nil
	words, err := GetSliceOfWords(content)
	if err != nil {
		logErr.Error(err)
		return err
	}
	for i, w := range words {
		if i > 9 {
			break
		}
		fmt.Fprintf(os.Stdout, "%-3d Word: %s\tcount= %d\n", i+1, w.Word, w.Count)
	}
	return nil
}

//ShowHightFrequencyWords показать самые часто иссполььзуемые слова
func GetSliceOfWords(biteWithWords []byte) ([]string, error) {
	// var words []string
	var bufError error
	var w struct {
		Count int
		Word  string
	}
	buf := make([]byte, 1024)
	// пройдемся учитывая знаки разделения
	var totalWordSum = make(map[string]interface{}, 100)
	reader := bufio.NewReader(biteWithWords)

	for {
		buf, _, bufError = reader.ReadLine()
		if bufError != nil {
			if bufError == io.EOF {
				break
			}
			return nil, bufError
		}
		newWords := regexp.MustCompile(`(\w|[а-яА-Я\_])+`).FindAllString(string(buf), -1)
		for _, s := range newWords {
			// totalWordSum[s]++

			if w, ok := totalWordSum[s]; ok {
				w.Count = w.Count + 1
			} else {
				totalWordSum[s] = &struct {
					Count int
					Word  string
				}{
					Count: 1,
					Word:  s,
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
	return totalWordSumArray, nil
}

func init() {
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.00", FullTimestamp: true})
	log.SetOutput(os.Stderr)
	logErr = log.New(os.Stderr, ``, 0)
}
