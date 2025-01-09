package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFetchEntry(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "data/mock_iknow.html")
	}))
	defer ts.Close()

	entry, err := fetchEntry(ts.URL)
	if err != nil {
		t.Errorf("failed to fetch the entry: %v", err)
	}

	want := Entry{
		StartedItems:     6491,
		CompletedItems:   4258,
		CompletedCourses: 73,
		StudyTime:        218,
	}
	if !reflect.DeepEqual(entry, want) {
		t.Errorf("got: %v, want: %v", entry, want)
	}
}
