package internal

import "github.com/urfave/cli/v2"

var Flags = []cli.Flag{
	dbNameFlag,
}

var dbNameFlag = &cli.StringFlag{
	Name:  "db",
	Value: DefaultDBName,
	Usage: "The database file name",
}

var BaseURLFlag = &cli.StringFlag{
	Name:  "base-url",
	Value: "https://iknow.jp",
	Usage: "The base URL of the iKnow! website",
}

var PredictCompletedFlag = &cli.BoolFlag{
	Name:  "predict-completed",
	Usage: "Draw the predicted number of completed items",
}
