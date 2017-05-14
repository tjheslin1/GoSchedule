package model

// SubmitJob represents an incooming request to persist a job.
type SubmitJob struct {
	Name      string `json:"name"`
	StartTime int    `json:"start_time"`
	Interval  int    `json:"interval"`
	URL       string `json:"url"`
}
