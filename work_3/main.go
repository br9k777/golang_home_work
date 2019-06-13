package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func init() {
	// log.SetFormatter(&log.TextFormatter{})
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.00", FullTimestamp: true})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
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

	//

	app.Action = func(c *cli.Context) error {

		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
			log.SetOutput(os.Stderr)
		}

		if c.String(`words-file-path`) != `` {
			ShowHightFrequencyWords(c.String(`words-file-path`))
		}
		return nil
	}
	app.Run(os.Args)

}
