package main

import (
	"encoding/json"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"github.com/rknizzle/testlive/job"
	"github.com/rknizzle/testlive/datastore"
	"github.com/rknizzle/testlive/datastore/inmemory"
	"github.com/rknizzle/testlive/scheduler"
)

func main() {
	// create a datastore
	// inmemory store is default until others are added
	jobStore := inmemory.New()

	// start periodically firing off job requests
	go scheduler.Init(jobStore)

	// initialize REST API endpoints
	initRestAPI(jobStore)
}

func initRestAPI(jobStore datastore.Datastore) {

	r := gin.Default()

	// ping
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	////////////////////////////////////////////
	// /jobs
	////////////////////////////////////////////

	// Get all jobs
	r.GET("/jobs", func(c *gin.Context) {
		c.JSON(200, jobStore.GetAll())
	})

	// Create a new job
	r.POST("/jobs", func(c *gin.Context) {
		// get request body
		var j job.Job
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}

		// create a new job object (j) based on the input body
		err = json.Unmarshal(body, &j)
		if err != nil {
			panic(err)
		}
		jobStore.Create(&j)
		c.JSON(200, j)
	})

	// Update a job
	r.PUT("/jobs/:id", func(c *gin.Context) {
		id := c.Param("id")

		var j job.Job
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			panic(err)
		}

		// create a new job object (j) based on the input body
		err = json.Unmarshal(body, &j)
		if err != nil {
			panic(err)
		}

		_, err = jobStore.Update(id, &j)
		if err != nil {
			panic(err)
		}
		c.JSON(200, j)
	})

	// Listen and serve on localhost
	r.Run()
}
