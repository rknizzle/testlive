// Package job represents a feature or endpoint to be tested.
// The result of the job is checked against an expected result
// to express a passing or failing feature

package job

type Job struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	URL        string   `json:"url"`
	HTTPMethod string   `json:"httpMethod"`
	Frequency  int      `json:"frequency"`
	Status     string   `json:"status"`
	Response   Response `json:"response"`
}

type Response struct {
	StatusCode int `json:"statusCode"`
}

func New(id string, title string, url string, httpMethod string, frequency int, status string, response Response) *Job {
	return &Job{id, title, url, httpMethod, frequency, status, response}
}
