# testlive
#### testlive is a continous testing service/tool for testing backend endpoints/features in development or production environments

#### testlive keeps a collection of jobs. A job represents an endpoint to be periodically requested and the response to be verified to monitor the working state of the endpoint

## Binary Installation
```
$ go get -u github.com/rknizzle/testlive/cmd/testlive
```

## Run binary locally in Docker container
```
docker run -it -d -p 8080:8080 rkneills/testlive:latest
```

## Usage
The binary starts an HTTP server running on localhost:8080  
Open localhost:8080 in the browser to access the status page and add jobs

## REST API
```
POST /jobs - Add a new job  
GET /jobs - Get all jobs and their statuses  
PUT /jobs - Update a job  
```

Job structure:
```
{
  "title": "jobTitle",
  "url": "http://test.com/feature",
  "httpMethod": "GET",
  "frequency": 30, (seconds)
  "response": {
    "statusCode": 200
  }
}
```
