package internal

import (
	"time"

	"github.com/urfave/cli/v2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
)

func CmdPlot(ctx *cli.Context) error {
	startDateStr := ctx.Args().Get(0)
	endDateStr := ctx.Args().Get(1)
	if startDateStr == "" || endDateStr == "" {
		return cli.Exit("Please specify the start and end dates", 1)
	}

	db, err := openDB(DefaultDBName)
	if err != nil {
		return err
	}

	var entries []Entry
	err = db.Select(&entries, "SELECT * FROM entries WHERE date BETWEEN ? AND ? ORDER BY date ASC", startDateStr, endDateStr)

	p := plot.New()

	p.X.Label.Text = "Date"

	startedItems := make(plotter.XYs, len(entries))
	completedItems := make(plotter.XYs, len(entries))

	for i, entry := range entries {
		d, err := parseDate(entry.Date)
		if err != nil {
			return err
		}
		x := float64(d.Unix())

		startedItems[i].X = x
		startedItems[i].Y = float64(entry.StartedItems - entries[0].StartedItems)
		completedItems[i].X = x
		completedItems[i].Y = float64(entry.CompletedItems - entries[0].CompletedItems)
	}

	s, err := plotter.NewLine(startedItems)
	if err != nil {
		return err
	}
	s.LineStyle.Width = 2
	s.LineStyle.Color = plotutil.Color(0)

	c, err := plotter.NewLine(completedItems)
	if err != nil {
		return err
	}
	c.LineStyle.Width = 2
	c.LineStyle.Color = plotutil.Color(1)

	p.Add(s, c)

	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02"}

	err = p.Save(400, 300, "plot.png")

	return nil
}

func parseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}
