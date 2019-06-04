package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

func init() {
	// log.SetFormatter(&log.TextFormatter{})
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.00", FullTimestamp: true})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
	//для вывода на экран - предпологается что первый лог идет в файл, второй на экран - пока не удалось реализвать
	// logError := log.New()
	// logError.Formatter=&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.00", FullTimestamp: true}
	// logError.SetOutput(os.Stderr)
	// log.SetLevel(log.Warn)
}

const (
	jsonFilePath = `urlBase.json`
)

func main() {

	app := cli.NewApp()
	app.Name = "hello go"
	app.Usage = "just do something"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "shorten, s",
			Usage: "Дать короткую ссылку",
		},
		cli.StringFlag{
			Name:  "resolve, r",
			Usage: "Дать длинную ссылку",
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

		log.Debugf(`Total arguments %s %d`, c.Args().Get(0), c.NArg())

		if err := ReadUrlBaseFromJsonFile(jsonFilePath); err != nil {
			log.Error(err)
		}
		if c.String(`shorten`) != `` {
			var url Url
			log.Info(url.Shorten(c.String(`shorten`)))
			WriteUrlBaseToJsonFile(jsonFilePath)
		}
		if c.String(`resolve`) != `` {
			var url Url
			log.Info(url.Shorten(c.String(`resolve`)))
		}
		return nil
	}
	app.Run(os.Args)

}
