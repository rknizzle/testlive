package main

import (
	"github.com/rknizzle/testlive/pkg/api"
	"github.com/rknizzle/testlive/pkg/datastore"
	"github.com/rknizzle/testlive/pkg/datastore/inmemory"
	"github.com/rknizzle/testlive/pkg/scheduler"
)

var jobstore datastore.Datastore

func main() {
	// create a datastore
	// inmemory store is default until others are added
	jobstore = inmemory.New()

	s := scheduler.New(jobstore)
	// start periodically firing off job requests
	go s.Init(jobstore)

	// initialize REST API endpoints
	a := api.New(jobstore)
	a.Init()
}
