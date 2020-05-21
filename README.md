# Utsuru Concept 📷

This application is used for a photo repository as per [Shopify](https://www.shopify.com/)'s backend [coding challenge](https://docs.google.com/document/d/1ZKRywXQLZWOqVOHC4JkF3LqdpO3Llpfk_CkZPR8bjak/edit) for Fall 2020 internships.

## Why Utsuru Concept?

Utsuru comes from the Japanese word **映る** which means "to reflect/to project". In that sense, this service would be used for saving various reflections in the form of photography. 

Concept means that it is a proof-of-concept and that it may or may not be used for future usage.

## Prerequisites

- [Go](https://golang.org/)

- [MySQL](https://www.mysql.com/)

- [Google Cloud Platform](https://cloud.google.com/)

## How to get started

- Clone this repository
  
  `git clone https://github.com/zaclimon/UtsuruConcept`

- Install the packages
  
  `go mod download`

- Fill the environment variables inside `.env`.
  
  - For Google Cloud Storage, you will need to obtain your [service account's](https://cloud.google.com/docs/authentication/production) `.json` file.

- Build the project
  
  `go build UtsuruConcept`

- Execute the binary

  `./UtsuruConcept`

- Do a test call from using `curl http://localhost:5000/api/register`

- Enjoy! 🎉
  
  **Note**: Since a `Dockerfile` is available, it is also possible to use Docker Compose to build the whole development stack wihtout installing any external dependencies.

## Endpoints

All endpoints are to be used with a `POST` request method. As such, you could use a solution like [Postman](https://www.getpostman.com/) to try your solution.

- `/api/v1/register` To register a user it's body is a JSON composed of two fields which are `email` and `password`.
- `/api/v1/login` To authenticate a user it's body is a JSON composed of two fields which are `email` and `password`. It will return an `access_token` field containing the token which should be used for all further requests
- `/api/v1/upload` To upload an image. This time around, the body is a `form-data` and is composed of the `access_token` key which is the access token retrieved when authenticating the user and the `images` key which is a set of files. Please note that you will need to use [`multipart/form-data`](https://stackoverflow.com/a/4526286) as an encoding type to upload the images. It will return a list of all tge images that have been uploaded as well as images that have not been uploaded successfully.
- `/api/v1/search` has similar inputs to `/api/v1/upload` in which it only requests one image instead of one/more and it returns the metadata of all images that have been identified as "similar" to the one uploaded.

## Architecture

The document [`ARCHITECTURE.md`](./ARCHITECTURE.md) explains how the application should behave as well as the reasoning behind some of the choices made when designing it.

## Tests

Tests are also available for this project. Execute `go test UtsuruConcept/testing` in order to run them.