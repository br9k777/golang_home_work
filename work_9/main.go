package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
	// "github.com/spf13/viper"
)

const (
	defaultInFile  = "/tmp/test_in.txt"
	defaultOutFile = "/tmp/out.txt"
	defaultIBS     = 512
	defaultOBS
)

var (
	log           *zap.Logger
	input, output string
	ibs, obs      int
	offset        int64
)

func init() {
	var err error
	if log, err = zap.NewProduction(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(log)
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
	log.Sugar().Debug(`аргументы %v`, os.Args)
	var err error
	if input != `` && output != `` {
		fmt.Printf("Берем из файла %s копируем в файл %s. Смешение %d.\n", input, output, offset)
		fmt.Printf("Читаем за раз %d. Пишем за раз %d.\n", ibs, obs)
		// fmt.Printf("Запускаем копирование.\n")
		if err = CopyFile(input, output, offset, ibs, obs); err != nil {
			fmt.Fprintf(os.Stderr, `Во вреия копирования произошла ошибка %s`, err)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Вы не указали обязательные ключи: --in_file `путь до копируемого файла` и --out_file `путь до нового файла`\n")
		fmt.Fprintf(os.Stderr, "\tПример использования %s --in_file /tmp/test_in.txt --out_file /tmp/out.txt", os.Args[0])
		fmt.Fprintf(os.Stderr, " --ibs 512  --obs 512\n")
		fmt.Fprintf(os.Stderr, "Более подробно о ключах можно посмотреть выполнив %s -h\n", os.Args[0])
	}

}
