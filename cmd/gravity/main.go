package main

import (
	"log"
	"os"

	"github.com/Gravity-Tech/gravity-core/cmd/gravity/commands"

	"github.com/urfave/cli/v2"
)

const (
	version           = "0.0.1"
	DefaultGravityDir = ".gravity"
)

func main() {
	app := &cli.App{
		Name:  "Gravity CLI",
		Usage: "the gravity command line interface",
		Commands: []*cli.Command{
			commands.LedgerCommand,
			commands.OracleCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
