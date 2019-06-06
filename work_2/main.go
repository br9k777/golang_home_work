package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

func init() {
	// log.SetFormatter(&log.TextFormatter{})
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.00", FullTimestamp: true})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
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
			Name:  "unpack, u",
			Usage: "Распакавать строку",
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

		if c.String(`unpack`) != `` {
			fmt.Fprintf(os.Stderr, "Original string= %s\nUnpack string =%s\n", c.String(`unpack`), fmt.StringUnpack(c.String(`unpack`)))
		}
		return nil
	}
	app.Run(os.Args)

}
