package internal

import (
	"reflect"
	"testing"
)

func TestFetchDiffData(t *testing.T) {
	db := createTestDB()
	defer db.Close()

	err := initDB(db)
	if err != nil {
		t.Fatalf("failed to initialize the database: %v", err)
	}

	_, err = db.Exec(`
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

	begin, end, err := fetchDiffData(db, "2021-01-02", "2021-01-03")
	if err != nil {
		t.Fatalf("failed to fetch diff data: %v", err)
	}

	wantBegin := Entry{
		Date:             "2021-01-02",
		StartedItems:     200,
		CompletedItems:   100,
		CompletedCourses: 2,
		StudyTime:        20,
	}
	if !reflect.DeepEqual(begin, wantBegin) {
		t.Errorf("got: %v, want: %v", begin, wantBegin)
	}

	wantEnd := Entry{
		Date:             "2021-01-03",
		StartedItems:     300,
		CompletedItems:   150,
		CompletedCourses: 3,
		StudyTime:        30,
	}
	if !reflect.DeepEqual(end, wantEnd) {
		t.Errorf("got: %v, want: %v", end, wantEnd)
	}
}
