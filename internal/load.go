package internal

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli/v2"
)

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

func CmdLoad(ctx *cli.Context) error {
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
	entry.Date = time.Now().Format("2006-01-02")

	db, err := openDB(DefaultDBName)
	if err != nil {
		return fmt.Errorf("failed to open the database: %w", err)
	}
	defer db.Close()

	err = initDB(db)
	if err != nil {
		return err
	}
	_, err = db.Exec(
		`INSERT OR REPLACE INTO entries (date, started_items, completed_items, completed_courses, study_time) VALUES (?, ?, ?, ?, ?)`,
		entry.Date, entry.StartedItems, entry.CompletedItems, entry.CompletedCourses, entry.StudyTime,
	)
	if err != nil {
		return fmt.Errorf("failed to insert the entry: %w", err)
	}
	return nil
}
