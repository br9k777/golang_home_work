package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "go-envdir"
	app.Version = "1.0.0"
	app.Compiled = time.Now()
	app.Usage = "runs another program with environment modified according to files in a specified directory."
	app.UsageText = fmt.Sprintf("%s путь_до_директории_с_файлами_переменных_окружения вызов_какой_либо_программы", app.Name)
	// cli.App
	writer := os.Stdout
	writerErr := os.Stderr
	app.Action = func(c *cli.Context) error {

		fmt.Fprintf(writer, "Всего аргументов %d :\n %s\n", c.NArg(), c.Args())
		if c.NArg() != 2 {
			fmt.Fprintf(writerErr, "Не правильный вызов программы, пример использования:\n%s\n", app.UsageText)
			return nil
		}
		var enviroments map[string]string
		var err error
		if enviroments, err = GetEnviromentsFromDir(c.Args().Get(0)); err != nil {
			return err
		}
		if err = RunProgragWirhEnviroments(enviroments, c.Args().Get(1)); err != nil {
			return err
		}
		return nil

	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(writerErr, "ERROR: %s", err)
	}

}
