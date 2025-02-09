package main

import (
	"github.com/blue-jam/iknow-tools/internal"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Commands = []*cli.Command{
		{
			Name:   "load",
			Usage:  "Load the statistics from the specified URL",
			Args:   true,
			Action: internal.CmdLoad,
			Flags: []cli.Flag{
				internal.BaseURLFlag,
			},
		},
		{
			Name:   "diff",
			Usage:  "Show the difference between the statistics of the specified dates",
			Args:   true,
			Action: internal.CmdDiff,
		},
		{
			Name:   "plot",
			Usage:  "Show the graph of the statistics",
			Args:   true,
			Action: internal.CmdPlot,
			Flags: []cli.Flag{
				internal.PredictCompletedFlag,
			},
		},
	}
	app.Flags = internal.Flags
	app.Name = "iknow-tools"
	app.Usage = "A CLI tool for managing iKnow! statistics"
	app.Version = "0.1.0"
	app.Before = internal.OpenDBFromContext
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
