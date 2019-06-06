package main

import (
	log "github.com/sirupsen/logrus"
	"strings"
)

var (
	//хранилище ссылок
	UrlBaseByShortName map[string]*Url
	UrlBaseByLongName  map[string]*Url
)

//проверяем являяется ли руна число и возвращем число
func isDigit(oneRune string, digit *int32) bool {
	// так как нелься исползовать массивы ии сторонниие библотеки типо strconv Itoa
	// будем использоватьь switch
	switch oneRune {
	case `0`:
		*digit = 0
		return true
	case "1":
		*digit = 1
		return true
	case "2":
		*digit = 2
		return true
	case "3":
		*digit = 3
		return true
	case "4":
		*digit = 4
		return true
	case "5":
		*digit = 5
		return true
	case "6":
		*digit = 6
		return true
	case "7":
		*digit = 7
		return true
	case "8":
		*digit = 8
		return true
	case "9":
		*digit = 9
		return true
	default:
		*digit = 0
		return false
	}
}
func (out *strings.Builder) repeatAndWrite(OneRune *rune, repeat *int32) {
	if OneRune == `` {
		return
	}
	for i := 0; i < repeat; i++ {
		out.WriteRune(OneRune)
	}
	OneRune = ``
	repeat = 0
}

func StringUnpack(inputString string) string {
	escapeCharacter := false
	var lastRune rune
	var totalRepeat, digit int32
	var outputString strings.Builder
	for _, oneRune := range inputString {
		switch {

		case oneRune == `\`:
			outputString.repeatAndWrite(lastRune, lastRuneIsDigit)
			if escapeCharacter {
				//ло этого уже был \
				lastRune = `\`
				escapeCharacter = false
			} else {
				lastRune = ``
				escapeCharacter = true
			}
		case escapeCharacter:
			outputString.repeatAndWrite(lastRune, lastRuneIsDigit)
			escapeCharacter = false
			totalRepeat = 0
			lastRune = oneRune

		case isDigit(oneRune, digit):
			totalRepeat = totalRepeat*10 + digit
		default:
			outputString.repeatAndWrite(lastRune, lastRuneIsDigit)
			escapeCharacter = false
			totalRepeat = 0
			lastRune = oneRune
		}
	}
	outputString.repeatAndWrite(lastRune, lastRuneIsDigit)

	return outputString.String()
}

func init() {
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.00", FullTimestamp: true})
	log.SetOutput(os.Stderr)
}
