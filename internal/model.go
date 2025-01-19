package internal

type Entry struct {
	StartedItems     int    `db:"started_items"`
	CompletedItems   int    `db:"completed_items"`
	CompletedCourses int    `db:"completed_courses"`
	StudyTime        int    `db:"study_time"`
	Date             string `db:"date"`
}
