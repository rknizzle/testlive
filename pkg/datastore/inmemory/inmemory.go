// Package inmemory stores a collection of jobs in application memory
// This package should be used for testing only as the data will not persist
// across runs of the program

package inmemory

import (
	"errors"
	"github.com/google/uuid"
	"github.com/rknizzle/testlive/job"
)

// inmemory
type inmemory struct {
	jobCollection []*job.Job
}

// Initializes a new in memory datastore of jobs
func New() *inmemory {
	return &inmemory{}
}

// Add a new job to the datastore
func (i *inmemory) Create(j *job.Job) *job.Job {
	// generate a UUID for the job
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	j.ID = id.String()

	i.jobCollection = append(i.jobCollection, j)
	return j
}

// Updates job data
func (i *inmemory) Update(id string, j *job.Job) (*job.Job, error) {
	for index, k := range i.jobCollection {
		if k.ID == id {
			i.jobCollection[index] = j
			return k, nil
		}
	}

	return &job.Job{}, errors.New("Job could not be found for updating")
}

// Gets the job with the given ID
func (i *inmemory) Get(id string) (*job.Job, error) {
	for _, job := range i.jobCollection {
		if job.ID == id {
			return job, nil
		}
	}
	return &job.Job{}, errors.New("No job found")
}

// Gets all jobs
func (i *inmemory) GetAll() []*job.Job {
	return i.jobCollection
}
