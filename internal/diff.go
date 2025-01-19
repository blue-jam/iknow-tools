package internal

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"
)

func CmdDiff(ctx *cli.Context) error {
	beginDateStr := ctx.Args().Get(0)
	endDateStr := ctx.Args().Get(1)
	if beginDateStr == "" || endDateStr == "" {
		return cli.Exit("Please specify the begin and end dates", 1)
	}

	db, err := sqlx.Open("sqlite3", DefaultDBName)
	if err != nil {
		return fmt.Errorf("failed to open the database: %w", err)
	}

	beginEntry, endEntry, err := fetchDiffData(db, beginDateStr, endDateStr)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, "Metric\tBegin Date\tEnd Date\tDifference")
	fmt.Fprintf(w, "Date\t%s\t%s\t\n", beginEntry.Date, endEntry.Date)
	fmt.Fprintf(w, "Started Items\t%d\t%d\t%+d\n", beginEntry.StartedItems, endEntry.StartedItems, endEntry.StartedItems-beginEntry.StartedItems)
	fmt.Fprintf(w, "Completed Items\t%d\t%d\t%+d\n", beginEntry.CompletedItems, endEntry.CompletedItems, endEntry.CompletedItems-beginEntry.CompletedItems)
	fmt.Fprintf(w, "Completed Courses\t%d\t%d\t%+d\n", beginEntry.CompletedCourses, endEntry.CompletedCourses, endEntry.CompletedCourses-beginEntry.CompletedCourses)
	fmt.Fprintf(w, "Study Time\t%d\t%d\t%+d\n", beginEntry.StudyTime, endEntry.StudyTime, endEntry.StudyTime-beginEntry.StudyTime)
	w.Flush()

	return nil
}

func fetchDiffData(db *sqlx.DB, beginDateStr string, endDateStr string) (Entry, Entry, error) {
	var beginEntry Entry
	err := db.Get(&beginEntry, "SELECT * FROM entries WHERE date >= ? ORDER BY date ASC LIMIT 1", beginDateStr)
	if err != nil {
		return Entry{}, Entry{}, fmt.Errorf("failed to fetch the begin entry: %w", err)
	}
	var endEntry Entry
	err = db.Get(&endEntry, "SELECT * FROM entries WHERE date <= ? ORDER BY date DESC LIMIT 1", endDateStr)
	if err != nil {
		return Entry{}, Entry{}, fmt.Errorf("failed to fetch the end entry: %w", err)
	}
	return beginEntry, endEntry, nil
}
