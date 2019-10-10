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
	app.Name = "unitylint"
	app.Usage = "unity linter by build log"
	app.Description = "Parse unity log to simplify errors and warnings."
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
		linter.Prepare(".loglint.json")
		result := linter.Run(readall(reader))
		result.Print()

		return nil
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
