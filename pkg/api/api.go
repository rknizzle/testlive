package api

import (
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/rknizzle/testlive/pkg/datastore"
	"github.com/rknizzle/testlive/pkg/job"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type API struct {
	jobstore datastore.Datastore
}

func New(jobstore datastore.Datastore) *API {
	return &API{jobstore}
}

func sendErrorMessage(w http.ResponseWriter) {
	e := make(map[string]string)
	e["error"] = "true"
	e["message"] = "Operation failed"
	eData, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(eData)
}

// Renders the specified html page with the given template data
func renderPage(w http.ResponseWriter, pageName string, data interface{}) error {
	page, err := templatesBox.FindString(pageName)
	if err != nil {
		return err
	}

	t := template.New("page")
	t, err = t.Parse(string(page))
	if err != nil {
		return err
	}
	err = t.ExecuteTemplate(w, "page", data)
	if err != nil {
		return err
	}
	return nil
}

// ping
func (a *API) ping(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	j := make(map[string]string)
	j["message"] = "pong"
	jData, err := json.Marshal(j)
	if err != nil {
		sendErrorMessage(w)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

// load status page
func (a *API) status(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := renderPage(w, "status.html", nil)
	if err != nil {
		sendErrorMessage(w)
	}
}

// load the form to create a new job
func (a *API) newJob(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// load the form to create a new job
	j := &job.Job{}
	err := renderPage(w, "jobForm.html", j)
	if err != nil {
		sendErrorMessage(w)
	}
}

// load the job update form for the specified job
func (a *API) edit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	j, err := a.jobstore.Get(id)
	if err != nil {
		sendErrorMessage(w)
	}

	renderPage(w, "jobForm.html", j)
}

////////////////////////////////////////////
// /jobs
////////////////////////////////////////////

// Get all jobs
func (a *API) getJobs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jobs := a.jobstore.GetAll()
	j, err := json.Marshal(jobs)
	if err != nil {
		sendErrorMessage(w)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// Create a new job
func (a *API) createJob(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// get request body
	var j job.Job
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendErrorMessage(w)
	}

	// create a new job object (j) based on the input body
	err = json.Unmarshal(body, &j)
	if err != nil {
		sendErrorMessage(w)
	}
	a.jobstore.Create(&j)
	job, err := json.Marshal(j)
	if err != nil {
		sendErrorMessage(w)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(job)
}

// Update a job
func (a *API) updateJob(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	var j job.Job
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendErrorMessage(w)
	}

	// create a new job object (j) based on the input body
	err = json.Unmarshal(body, &j)
	if err != nil {
		sendErrorMessage(w)
	}

	_, err = a.jobstore.Update(id, &j)
	if err != nil {
		sendErrorMessage(w)
	}

	job, err := json.Marshal(j)
	if err != nil {
		sendErrorMessage(w)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(job)
}

// will this work or will I need to put it inside the API struct?
var templatesBox = packr.New("Templates", "./../../templates")

//func initRestAPI(jobStore datastore.Datastore) {
func (a *API) Init() {

	router := httprouter.New()
	router.GET("/", a.status)
	router.GET("/status", a.status)
	router.GET("/ping", a.ping)
	router.GET("/new", a.newJob)
	router.GET("/edit/:id", a.edit)
	router.GET("/jobs", a.getJobs)
	router.POST("/jobs", a.createJob)
	router.PUT("/jobs/:id", a.updateJob)

	fmt.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
