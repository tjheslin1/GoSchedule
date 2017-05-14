package model

// SubmitJob represents an incoming request to persist a job.
type SubmitJob struct {
	Name      string `json:"name"`
	StartTime int64  `json:"start_time"`
	Interval  int64  `json:"interval"`
	URL       string `json:"url"`
}
