package main

import (
	"bufio"
	"regexp"
	"sort"
	"strings"
)

const numberOfMostCommonWordsToDisplay = 10

//GetSliceOfWords показать самые часто иссполььзуемые слова
func GetSliceOfWords(biteWithWords []byte) ([]string, error) {
	// var err error

	type word struct {
		count int
		word  string
	}
	// пройдемся учитывая знаки разделения
	var totalWordSum = make(map[string]*word, 100)
	scanner := bufio.NewScanner(strings.NewReader(string(biteWithWords)))

	for scanner.Scan() {
		newWords := regexp.MustCompile(`(\w|[а-яА-Я\_])+`).FindAllString(string(scanner.Text()), -1)
		for _, s := range newWords {
			// totalWordSum[s]++

			if w, ok := totalWordSum[s]; ok {
				w.count = w.count + 1
				continue
			}
			totalWordSum[s] = &word{
				count: 1,
				word:  s,
			}

		}
	}
	//обратно в масив для сортировки
	var totalWordSumArray = make([]*word, 0, len(totalWordSum))
	for _, w := range totalWordSum {
		totalWordSumArray = append(totalWordSumArray, w)
	}
	sort.Slice(totalWordSumArray, func(i, j int) bool {
		if totalWordSumArray[i].count == totalWordSumArray[j].count {
			return sort.StringsAreSorted([]string{totalWordSumArray[i].word, totalWordSumArray[j].word})
		}
		return totalWordSumArray[i].count > totalWordSumArray[j].count
	})
	// если слов больше нужного количество то возвращаем толко десять
	resultWords := make([]string, 0, numberOfMostCommonWordsToDisplay)
	for i := 0; i < len(totalWordSumArray) && i < numberOfMostCommonWordsToDisplay; i++ {
		resultWords = append(resultWords, totalWordSumArray[i].word)
	}
	return resultWords, nil
}
