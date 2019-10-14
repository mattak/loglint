package main

import (
	"bufio"
	"github.com/mattak/loglint/internal"
	"github.com/urfave/cli"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func readall(r io.Reader) string {
	data, err := ioutil.ReadAll(r)

	if err != nil {
		log.Fatal(err)
	}

	return string(data)
}

func main() {
	app := cli.NewApp()
	app.Name = "loglint"
	app.Usage = "log linter by local rules"
	app.Description = "Analyze errors and warnings from log by simple rule file."
	app.Version = "0.1.2"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "execute, e",
			Value: "",
			Usage: "rule json content by argument",
		},
		cli.StringFlag{
			Name: "file, f",
			Value: ".loglint.json",
			Usage: "rule json file path",
		},
	}
	app.Action = func(c *cli.Context) error {
		var reader io.Reader

		if len(c.Args()) <= 0 {
			reader = os.Stdin
		} else {
			filename := c.Args()[0]
			file, err := os.Open(filename)

			if err != nil {
				log.Fatal(err)
			}

			reader = bufio.NewReader(file)
		}

		linter := internal.Linter{}

		if len(c.String("execute")) > 0 {
			linter.PrepareByContent(c.String("execute"))
		} else if len(c.String("file")) > 0 {
			linter.PrepareByFile(c.String("file"))
		} else {
			log.Fatal("Rule option error. please specify --execute or --file option.")
		}

		result := linter.Run(readall(reader))
		result.Print()

		return nil
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
