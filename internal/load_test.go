package internal

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/urfave/cli/v2"
)

func TestCmdLoad(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../test/data/mock_iknow.html")
	}))
	defer ts.Close()

	// Create a test database
	db := createTestDB()
	defer db.Close()

	// Create a CLI context with the necessary arguments
	app := cli.NewApp()
	set := flag.NewFlagSet("test", 0)
	BaseURLFlag.Apply(set)
	set.Parse([]string{"-base-url", ts.URL, "userID"})

	ctx := cli.NewContext(app, set, nil)

	ctx.App.Metadata = make(map[string]interface{})
	SetDBToContext(ctx, db)

	// Run the CmdLoad function
	err := CmdLoad(ctx)
	if err != nil {
		t.Fatalf("CmdLoad failed: %v", err)
	}

	// Verify the data in the database
	var entry Entry
	err = db.Get(&entry, "SELECT * FROM entries LIMIT 1")
	if err != nil {
		t.Fatalf("failed to fetch the entry: %v", err)
	}

	dateStr := time.Now().Format("2006-01-02")

	want := Entry{
		StartedItems:     6491,
		CompletedItems:   4258,
		CompletedCourses: 73,
		StudyTime:        218,
		Date:             dateStr,
	}
	if !reflect.DeepEqual(entry, want) {
		t.Errorf("got: %v, want: %v", entry, want)
	}
}
