package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rknizzle/testlive/datastore"
	"github.com/rknizzle/testlive/datastore/inmemory"
	"github.com/rknizzle/testlive/job"
	"github.com/rknizzle/testlive/scheduler"
	"io/ioutil"
	"net/http"
)

func main() {
	// create a datastore
	// inmemory store is default until others are added
	jobStore := inmemory.New()

	s := &scheduler.Scheduler{}
	// start periodically firing off job requests
	go s.Init(jobStore)

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

	// load in html templates
	r.LoadHTMLGlob("templates/*")

	// load status page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "status.html", nil)
	})

	// load status page
	r.GET("/status", func(c *gin.Context) {
		c.HTML(http.StatusOK, "status.html", nil)
	})

	// load the job update form for the specified job
	r.GET("/edit/:id", func(c *gin.Context) {
		id := c.Param("id")

		j, err := jobStore.Get(id)
		if err != nil {
			panic(err)
		}

		c.HTML(http.StatusOK, "jobForm.html", j)
	})

	// load the form to create a new job
	r.GET("/new", func(c *gin.Context) {
		j := &job.Job{}

		c.HTML(http.StatusOK, "jobForm.html", j)
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
