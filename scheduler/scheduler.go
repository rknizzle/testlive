// Package scheduler handles running jobs at the appropriate time and verification of responses

package scheduler

import (
	"fmt"
	"github.com/rknizzle/testlive/datastore"
	"github.com/rknizzle/testlive/job"
	"net/http"
	"sync"
	"time"
)

// periodically starts a batch of jobs
// Logic is very simple for now, just send out every request every 10 seconds.
// In a future update each job will have a set frequency and the scheduler will
// handle when each request should be sent
func Init(jobStore datastore.Datastore) {
	batchJobs(jobStore)
	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ticker.C:
			batchJobs(jobStore)
		}
	}
}

// used to match up the job with the result using the id
type result struct {
	id  string
	res http.Response
	err error
}

// sends out a request for each job
func batchJobs(jobStore datastore.Datastore) {
	ch := make(chan *result)

	var wg sync.WaitGroup
	jobs := jobStore.GetAll()
	// loop through all jobs
	for _, job := range jobs {
		wg.Add(1)
		go execute(job, ch, &wg)
	}

	// close the channel in the background
	go func() {
		wg.Wait()
		close(ch)
	}()

	// read from channel as they come in until its closed
	// the channel will close when all requests receive a response
	for res := range ch {
		job, err := jobStore.Get(res.id)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
		fmt.Println(job)
	}
}

// sends an HTTP request for a job
func execute(job *job.Job, ch chan<- *result, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{}
	r, err := http.NewRequest(job.HTTPMethod, job.URL, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(r)
	// if theres an error then create an empty http response to avoid dereferencing a nil value
	if err != nil {
		resp = &http.Response{}
	}

	result := &result{job.ID, *resp, err}

	// send http result to result channel
	ch <- result
}
