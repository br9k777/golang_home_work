package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"

	log "github.com/sirupsen/logrus"
	// "github.com/spf13/viper"
)

const (
	defaultInFile  = "/tmp/test_in.txt"
	defaultOutFile = "/tmp/out.txt"
	defaultIBS     = 512
	defaultOBS
)

var (
	input, output string
	ibs, obs      int
	offset        int64
)

func init() {
	// log.SetFormatter(&log.TextFormatter{})
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.00", FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	flag.StringVarP(&input, "in_file", "i", "", "ввести путь до файла источника `ПУТЬ_ДО_ФАЙЛА_ИСТОЧНИКА`")
	//
	flag.Lookup("in_file").NoOptDefVal = defaultInFile
	flag.StringVarP(&output, "out_file", "o", "", "ввести путь до файла источника `ПУТЬ_ДО_НОВОГО_ИСТОЧНИКА`")
	flag.Lookup("out_file").NoOptDefVal = defaultOutFile
	flag.Int64Var(&offset, "offset", 0, "смещение во входном файле относитльено начачла файла")
	flag.IntVar(&ibs, "ibs", defaultIBS, "Читает по bytes байт за раз.")
	flag.IntVar(&obs, "obs", defaultOBS, "Пишет по bytes байт за раз.")
}

func main() {
	flag.Parse()
	log.Debugf(`аргументы %v`, os.Args)

	if input != `` && output != `` {
		fmt.Printf("Берем из файла %s копируем в файл %s. Смешение %d.\n", input, output, offset)
		fmt.Printf("Читаем за раз %d. Пишем за раз %d.\n", ibs, obs)
		// fmt.Printf("Запускаем копирование.\n")
		CopyFile(input, output, offset, ibs, obs)
	} else {
		fmt.Fprintf(os.Stderr, "Вы не указали обязательные ключи: --in_file `путь до копируемого файла` и --out_file `путь до нового файла`\n")
		fmt.Fprintf(os.Stderr, "\tПример использования %s --in_file /tmp/test_in.txt --out_file /tmp/out.txt", os.Args[0])
		fmt.Fprintf(os.Stderr, " --ibs 512  --obs 512\n")
		fmt.Fprintf(os.Stderr, "Более подробно о ключах можно посмотреть выполнив %s -h\n", os.Args[0])
	}

}
