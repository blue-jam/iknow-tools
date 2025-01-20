package internal

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestCmdDiff(t *testing.T) {
	const fromDate = "2021-01-02"
	const toDate = "2021-01-03"

	// Create a test database and insert test data
	db := createTestDB()
	defer db.Close()

	_, err := db.Exec(`
		INSERT INTO entries (date, started_items, completed_items, completed_courses, study_time)
		VALUES
			('2021-01-01', 100, 50, 1, 10),
			('2021-01-02', 200, 100, 2, 20),
			('2021-01-03', 300, 150, 3, 30),
			('2021-01-04', 400, 200, 4, 40)
	`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Create a CLI context with the necessary arguments
	app := cli.NewApp()
	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{fromDate, toDate})
	ctx := cli.NewContext(app, set, nil)

	ctx.App.Metadata = make(map[string]interface{})
	SetDBToContext(ctx, db)

	// Capture the output
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = old }()

	// Run the CmdDiff function
	err = CmdDiff(ctx)
	if err != nil {
		t.Fatalf("CmdDiff failed: %v", err)
	}

	// Read the output
	w.Close()
	var output bytes.Buffer
	_, _ = r.WriteTo(&output)

	// Verify the output
	expectedOutput := `
Metric               Begin Date    End Date      Difference
Date                 2021-01-02    2021-01-03    
Started Items        200           300           +100
Completed Items      100           150           +50
Completed Courses    2             3             +1
Study Time           20            30            +10
`[1:]
	if output.String() != expectedOutput {
		t.Errorf("got: %v, want: %v", output.String(), expectedOutput)
	}
}
