package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

var outputString strings.Builder

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
func repeatAndWrite(OneRune string, repeat int32) {
	if OneRune == `` {
		return
	}
	if repeat == 0 {
		outputString.WriteString(OneRune)
	}
	var i int32
	log.Debugf(`write rune = %s ,repeat=%d`, OneRune, repeat)
	for i = 0; i < repeat; i++ {
		outputString.WriteString(OneRune)
	}
	OneRune = ``
	repeat = 0
}

func StringUnpack(inputString string) string {
	outputString.Reset()
	escapeCharacter := false
	var lastRune string
	var totalRepeat, digit int32
	for _, oneRune := range inputString {
		log.Debugf(`Read rune %s`, string(oneRune))
		switch {
		case string(oneRune) == `\`:
			repeatAndWrite(lastRune, totalRepeat)
			totalRepeat = 0
			if escapeCharacter {
				//ло этого уже был \
				lastRune = `\`
				escapeCharacter = false
			} else {
				lastRune = ``
				escapeCharacter = true
			}
		case escapeCharacter:
			repeatAndWrite(lastRune, totalRepeat)
			escapeCharacter = false
			totalRepeat = 0
			lastRune = string(oneRune)

		case isDigit(string(oneRune), &digit):
			totalRepeat = totalRepeat*10 + digit
		default:
			repeatAndWrite(lastRune, totalRepeat)
			escapeCharacter = false
			totalRepeat = 0
			lastRune = string(oneRune)
		}
	}
	repeatAndWrite(lastRune, totalRepeat)

	return outputString.String()
}

func init() {
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.00", FullTimestamp: true})
	log.SetOutput(os.Stderr)
}
