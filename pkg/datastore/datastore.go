// Package datastore is a data store abstraction in which all data store methods will be a subpackage of

package datastore

import (
	"github.com/rknizzle/testlive/pkg/job"
)

type Datastore interface {
	Create(*job.Job) *job.Job
	Update(string, *job.Job) (*job.Job, error)
	Get(string) (*job.Job, error)
	GetAll() []*job.Job
}
