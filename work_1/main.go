package main

import (
	"fmt"
	"github.com/beevik/ntp"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"time"
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

func main() {

	app := cli.NewApp()
	app.Name = "hello go"
	app.Usage = "just do something"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "test-hello, i",
			Usage: "say hi",
		},
		cli.BoolFlag{
			Name:  "say-time, t",
			Usage: "hust say time now",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "set debug level for log and output stdout",
		},
		cli.StringFlag{
			Name:  "ntp-server, n",
			Value: `ntp3.stratum2.ru`,
			Usage: "use ntp server `HOST`",
		},
	}

	//

	app.Action = func(c *cli.Context) error {

		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
			log.SetOutput(os.Stdout)
			//} else {
			//	logFile, err := os.OpenFile(c.String("log-file"), os.O_CREATE|os.O_WRONLY, 0666)
			//	if err != nil {
			//		log.SetOutput(os.Stdout)
			//		log.Error("creat log failed", err.Error())
			//		log.Info("Failed to log to file, using default stderr")
			//	} else {
			//		log.SetOutput(logFile)
			//		log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.00", FullTimestamp: true})
			//	}
			//	defer logFile.Close()
		}

		log.Debugf(`Total arguments %s %d`, c.Args().Get(0), c.NArg())

		if c.Bool("test-hello") {
			fmt.Print(`Hello world`)
		}
		if c.Bool("say-time") {
			fmt.Printf("Time is %s", time.Now().String())
		}
		//в любом случае напишем время от NTP сервера
		if c.String(`ntp-server`) != `` {
			if timeFromServer, err := ntp.Time(c.String(`ntp-server`)); err == nil {
				fmt.Printf(`Получили время от NTP сервера %s раное %s`+"\n", c.String(`ntp-server`), timeFromServer.String())
			} else {
				log.Errorf("Не удалось получить время от NTP сервера %s\n%s\n", c.String(`ntp-server`), err)
			}
		}
		return nil
	}
	app.Run(os.Args)

}
