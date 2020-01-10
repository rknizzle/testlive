package main

import (
	"github.com/gin-gonic/gin"
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
	// Listen and serve on localhost
	r.Run()
}
