package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"
	"go.uber.org/zap"
)

var (
	log *zap.Logger
)

func init() {
	var err error
	if log, err = zap.NewProduction(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(log)
}

func main() {

	app := cli.NewApp()
	app.Name = "words"
	app.Usage = "show hight words frequency"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "words-file-path, f",
			Value: `words.txt`,
			Usage: "Изменить пппуть до файлоа ссо словами",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "set debug level for log and output stdout",
		},
	}

	defer log.Sync()
	var err error
	app.Action = func(c *cli.Context) error {

		if c.Bool("debug") {
			if log, err = zap.NewDevelopment(); err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
		}

		if c.String(`words-file-path`) != `` {
			var content []byte
			var err error
			if content, err = ioutil.ReadFile(c.String(`words-file-path`)); err != nil {
				return err
			}
			var words []string
			if words, err = GetSliceOfWords(content); err != nil {
				return err
			}
			for _, w := range words {
				fmt.Fprintf(os.Stdout, "%s\n", w)
			}
			return nil
		}
		return nil
	}

	if err = app.Run(os.Args); err != nil {
		log.Error(`Ошибка выпорлнения`, zap.Error(err))
	}

}
