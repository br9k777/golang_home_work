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
			Name:  "test",
			Usage: "set debug level for log and output stdout",
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
			fmt.Fprintf(os.Stderr, "Original string= %s\nUnpack string =%s\n", c.String(`unpack`), StringUnpack(c.String(`unpack`)))
		}
		if c.Bool("test") {
			fmt.Fprintf(os.Stderr, "Original string= %s\t|Unpack string =%s\n", `a4bc2d5e`, StringUnpack(`a4bc2d5e`))
			fmt.Fprintf(os.Stderr, "Original string= %s\t|Unpack string =%s\n", `abcd`, StringUnpack(`abcd`))
			fmt.Fprintf(os.Stderr, "Original string= %s\t|Unpack string =%s\n", `45`, StringUnpack(`45`))
			fmt.Fprintf(os.Stderr, "Original string= %s\t|Unpack string =%s\n", `qwe\4\5`, StringUnpack(`qwe\4\5`))
			fmt.Fprintf(os.Stderr, "Original string= %s\t|Unpack string =%s\n", `qwe\\5`, StringUnpack(`qwe\\5`))
		}
		return nil
	}
	app.Run(os.Args)

}
