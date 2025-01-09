package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

type Entry struct {
	StartedItems     int
	CompletedItems   int
	CompletedCourses int
	StudyTime        int
	date             string
}

func fetchEntry(siteURL string) (Entry, error) {
	// Fetch the HTML page with net/http
	resp, err := http.Get(siteURL)
	if err != nil {
		wErr := fmt.Errorf("failed to fetch the URL: %w", err)
		return Entry{}, wErr
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		wErr := fmt.Errorf("failed to parse the HTML: %w", err)
		return Entry{}, wErr
	}

	var sItemsStr string
	var cItemsStr string
	var cCoursesStr string
	var sTimeStr string
	doc.Find(".statistics").Each(func(i int, s *goquery.Selection) {
		sItemsStr = s.Find(".started_items dd").Text()
		cItemsStr = s.Find(".completed_items dd").Text()
		cCoursesStr = s.Find(".completed_courses dd").Text()
		sTimeStr = s.Find(".study_time dd").Text()
	})

	entry := Entry{}
	entry.StartedItems, err = strconv.Atoi(sItemsStr)
	if err != nil {
		wErr := fmt.Errorf("failed to parse started items: %w", err)
		return Entry{}, wErr
	}
	entry.CompletedItems, err = strconv.Atoi(cItemsStr)
	if err != nil {
		wErr := fmt.Errorf("failed to parse completed items: %w", err)
		return Entry{}, wErr
	}
	entry.CompletedCourses, err = strconv.Atoi(cCoursesStr)
	if err != nil {
		wErr := fmt.Errorf("failed to parse completed courses: %w", err)
		return Entry{}, wErr
	}

	// Removing the "時間" suffix using Regular Expressions
	sTimeStr = strings.Replace(sTimeStr, "時間", "", -1)
	entry.StudyTime, err = strconv.Atoi(sTimeStr)
	if err != nil {
		wErr := fmt.Errorf("failed to parse study time: %w", err)
		return Entry{}, wErr
	}

	return entry, nil
}

func cmdLoad(ctx *cli.Context) error {
	uID := ctx.Args().First()
	if uID == "" {
		return cli.Exit("Please specify the user ID", 1)
	}
	url := "https://iknow.jp/users/" + uID
	entry, err := fetchEntry(url)
	if err != nil {
		return err
	}

	// Extract date from system time
	entry.date = time.Now().Format("2006-01-02")

	db, err := sql.Open("sqlite3", "iknow.sqlite3")
	if err != nil {
		return fmt.Errorf("failed to open the database: %w", err)
	}
	defer db.Close()

	queries := []string{
		`CREATE TABLE IF NOT EXISTS entries (date CHAR(16) PRIMARY KEY, started_items INTEGER, completed_items INTEGER, completed_courses INTEGER, study_time INTEGER)`,
	}
	for _, query := range queries {
		_, err = db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to execute the DB init query: %w", err)
		}
	}
	_, err = db.Exec(
		`INSERT OR REPLACE INTO entries (date, started_items, completed_items, completed_courses, study_time) VALUES (?, ?, ?, ?, ?)`,
		entry.date, entry.StartedItems, entry.CompletedItems, entry.CompletedCourses, entry.StudyTime,
	)
	if err != nil {
		return fmt.Errorf("failed to insert the entry: %w", err)
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Commands = []*cli.Command{
		{
			Name:   "load",
			Usage:  "Load the statistics from the specified URL",
			Args:   true,
			Action: cmdLoad,
		},
	}
	app.Name = "iknow-tools"
	app.Usage = "A CLI tool for managing iKnow! statistics"
	app.Version = "0.1.0"
	_ = app.Run(os.Args)
}
