// Package scheduler handles running jobs at the appropriate time and verification of responses

package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/rknizzle/testlive/pkg/datastore"
	"github.com/rknizzle/testlive/pkg/job"
	"net/http"
	"reflect"
	"sync"
	"time"
)

// Scheduler holds the collection of jobs and when to execute them
type Scheduler struct {
	jobs []*jobTimer
}

func New(jobstore datastore.Datastore) *Scheduler {
	// load in all the jobs with a start counter of 0 seconds
	var collection []*jobTimer
	jobs := jobstore.GetAll()
	for _, j := range jobs {
		jt := &jobTimer{0, j}
		collection = append(collection, jt)
	}
	return &Scheduler{collection}
}

// start scheduling counter
func (s *Scheduler) Init(jobstore datastore.Datastore) {
	// start the counter
	start := time.Now()
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			// for each tick of the ticker, make sure the jobs are synced up
			// with the datastore and check if any jobs should be executed
			count := int64(time.Since(start)) / 1e9
			c := int(count)
			go s.syncDatastore(jobstore, c)
			go s.startReadyJobs(jobstore, c)
		}
	}
}

// jobTimer contains a job as well as the count
// in which to trigger the execution of the job
type jobTimer struct {
	count int
	job   *job.Job
}

// used to match up the job with the result using the id
type result struct {
	id  string
	res http.Response
	err error
}

// Go through all jobTimers to find any jobs ready to be ran
func (s *Scheduler) startReadyJobs(jobstore datastore.Datastore, count int) {
	ch := make(chan *result)
	var wg sync.WaitGroup
	for _, i := range s.jobs {
		if i.count <= count {
			i.count += i.job.Frequency
			wg.Add(1)
			go execute(i.job, ch, &wg)
		}

	}
	// close the channel in the background
	go func() {
		wg.Wait()
		close(ch)
	}()

	// read from channel as they come in until its closed
	// the channel will close when all requests receive a response
	for res := range ch {
		job, err := jobstore.Get(res.id)
		if err != nil {
			panic(err)
		}
		// check if the response from the request matches the expected response
		verifyResponse(jobstore, job, res)
	}
}

// sync up the jobs in the datastore with the job data used by the scheduler
func (s *Scheduler) syncDatastore(jobstore datastore.Datastore, count int) {
	// add any new jobs
	// and update jobs that have been changed
	jobs := jobstore.GetAll()

	var match bool

	// loop through the jobs in the datastore
	for _, djob := range jobs {
		match = false
		// and loop through the jobs that the scheduler is currently working with
		for _, sjob := range s.jobs {
			// check if the job from the store already exists in the schedulers collection
			if djob.ID == sjob.job.ID {
				match = true
				// if the ids are the same but the objects are not equal anymore then reset the count
				// so the updated job will be executed immediately
				if !reflect.DeepEqual(djob, sjob.job) {
					sjob.count = count + 3
				}
				// set the job to the most recent version in the datastore
				sjob.job = djob
			}
		}
		// if the job does not already belong in the schedulers collection then the
		// job must be new and should be added to the collection
		if match != true {
			// add the datastore job to the schedulers collection
			// and make sure it is executed almost immediately after being added
			newTimer := &jobTimer{count + 3, djob}
			s.jobs = append(s.jobs, newTimer)
		}
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

// check that the status code of the http request matches the expected code
func verifyResponse(jobstore datastore.Datastore, job *job.Job, r *result) bool {

	fmt.Printf("Response status code: %d\n", r.res.StatusCode)
	fmt.Printf("Expected status code: %d\n", job.Response.StatusCode)

	// check if the status code matches the expected status code
	if r.res.StatusCode == job.Response.StatusCode {
		// if the status codes match, check if the response body matches the expected
		if bodiesAreEqual(r.res, job.Response.Body) {
			job.Status = "passing"
			jobstore.Update(job.ID, job)
			return true
		}
	}
	job.Status = "failing"
	jobstore.Update(job.ID, job)
	return false
}

func bodiesAreEqual(r http.Response, expectedBody interface{}) bool {
	var body interface{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Response body:")
	fmt.Println(body)

	fmt.Println("Expected body:")
	fmt.Println(expectedBody)

	eq := reflect.DeepEqual(body, expectedBody)
	if eq {
		fmt.Println("They're equal.")
		return true
	} else {
		fmt.Println("They're unequal.")
		return false
	}
}
