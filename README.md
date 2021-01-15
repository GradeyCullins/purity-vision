# Purity Vision Web Extension API Server
## About
Purity Vision is a project that utilizes the [Google Vision API](https://cloud.google.com/vision) to auto-filter images from web pages that are detected to have gore, nudity, or other explicit content. The filter is user-configurable and can be disabled with the click of a button. Browse the web with peace in mind.

## Project
This repository serves as the backend for the web extension. It interfaces with the Google Vision API and has a database that caches image requests in order to avoid redundant calls to Google's API. This is done to improve efficiency and reduce costs, as Google charges on a per request basis.

## Future
Eventually, this project is meant to be completely run within containers. As of now the Postgres database runs in Docker, and the API server runs locally using the user's local installation of Go.

## Requirements
### Required
 -  UNIX Operating System (Linux, MacOS, Window's Subsystem for Linux (WSL/WSL2))
 -  [golang](https://golang.org/dl/) (version must support [Go Modules](https://blog.golang.org/using-go-modules))
 - [docker](https://www.docker.com/)
 - [Google Cloud Account](https://cloud.google.com/)

### Optional
 - [direnv](https://direnv.net/) - useful for loading the environment variables needed to run the project

## Setup
### Google Account Credentials
Users will need to sign up for a Google Cloud account if they have not already and create a project with the Google Vision API enabled. In order to authenticate with the account, the user must create a Service Account under the GCP project, then export the credentials to a JSON a file that is saved locally on disk. The `GOOGLE_APPLICATION_CREDENTIALS` environment variable defined in the `.envrc` file must be an absolute path to this credential file for the image filtering to work.

### Environment
To connect to the database, the API server loads environment variables to setup the database credentials, point to the Google API credential file, among other things.

The required environment variables are listed in the `.envrc` file in the project repository. It is recommended to use [direnv](https://direnv.net/) to handle this, as it allows developers to utilize the .envrc file instead of manually exporting the environment variables.

**note**: make sure to fill in the `PURITY_DB_PASS` entry if using direnv, otherwise the database setup may fail.

### Database

With golang installed and Docker running, start the the database with the `start-db.sh` script.
```bash
./start-db.sh
```

### API Server
Run the API server with: 
```bash 
go run ./main.go
```
#### Example
Use curl to hit the *filter* endpoint:
```bash
curl -i localhost:8080/filter \
    -d '{"imgUriList": ["https://previews.124rf.com/images/valio84sl/valio84sl1311/valio84sl131100006/23554524-autumn-landscape-orange-trre.jpg"]}'
```
If everything is working, the response should look like:
```json
{
  "imgFilterResList": [
    {
      "imgURI": "https://www.allaboutbirds.org/news/wp-content/uploads/2020/03/THeron-Anderson-124505431.jpg",
      "error": "",
      "pass": true
    }
  ]
}
```
In this example, the image has passed the default filter rule which filters out images found to contain nudity.

## Support
Shoot me an email if you have any questions:
[gradeycullins@gmail.com](mailto:gradeycullins.com)