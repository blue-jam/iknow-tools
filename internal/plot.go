package internal

import (
	"gonum.org/v1/plot/vg"
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

	db := GetDBFromContext(ctx)

	var entries []Entry
	err := db.Select(&entries, "SELECT * FROM entries WHERE date BETWEEN ? AND ? ORDER BY date ASC", startDateStr, endDateStr)

	p := plot.New()

	p.X.Label.Text = "Date"
	p.Y.Label.Text = "Items"

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

	predictCompleted := ctx.Bool(PredictCompletedFlag.Name)
	if predictCompleted {
		startedTime, err := time.Parse("2006-01-02", entries[0].Date)
		if err != nil {
			return err
		}
		dataEndTime, err := time.Parse("2006-01-02", entries[len(entries)-1].Date)
		if err != nil {
			return err
		}
		rangeEndTime, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return err
		}

		completedItemsDiff := completedItems[len(completedItems)-1].Y - completedItems[0].Y
		totalDays := dataEndTime.Sub(startedTime).Hours() / 24
		rangeDays := rangeEndTime.Sub(startedTime).Hours() / 24

		predictedItems := make(plotter.XYs, 2)
		predictedItems[0].X = completedItems[0].X
		predictedItems[0].Y = completedItems[0].Y
		predictedItems[1].X = float64(rangeEndTime.Unix())
		predictedItems[1].Y = completedItemsDiff / totalDays * rangeDays
		pl, err := plotter.NewLine(predictedItems)
		if err != nil {
			return err
		}
		pl.LineStyle.Width = 2
		pl.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
		pl.LineStyle.Color = plotutil.Color(2)

		p.Add(pl)
	}

	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02"}

	err = p.Save(400, 300, "plot.png")

	return nil
}

func parseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}
