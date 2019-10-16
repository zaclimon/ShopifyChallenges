# Utsuru Concept 📷

This application is used for a photo repository as per [Shopify](https://www.shopify.com/)'s backend [coding challenge](https://docs.google.com/document/d/1ZKRywXQLZWOqVOHC4JkF3LqdpO3Llpfk_CkZPR8bjak/edit) for Winter 2020 internships.

## Why Utsuru Concept?

Utsuru comes from the Japanese word **映る** which means "to reflect/to project". In that sense, this service would be used for saving various reflections in the form of photography. 

Concept means that it is a proof-of-concept and that it may or may not be used for future usage.

## Prerequisites

- [Python 3](https://www.python.org/)

- [PostgreSQL](https://www.postgresql.org/)

- [Google Cloud Platform](https://www.postgresql.org/)

## How to get started

- Clone this repository
  
  `git clone https://github.com/zaclimon/UtsuruConcept`

- Install the packages
  
  `pip install -r requirements.txt`

- Go to the `app/` directory, fill the environment variables inside `.env.sample` and rename it to `.env`.
  
  - For Google Cloud Storage, you will need to obtain your [service account's](https://cloud.google.com/docs/authentication/production) `.json` file.
  
  - For `SQLALCHEMY_DATABASE_URI`, please refer to [Flask's](https://flask-sqlalchemy.palletsprojects.com/en/2.x/config/) documentation for examples on how to connect your database.

- Start Flask
  
  `flask run`

- Do a test call from using `curl http://localhost:5000/api/register`

- Enjoy! 🎉
  
  **Note**: It is possible to also use a [virtual environment](https://docs.python.org/3/library/venv.html) for better isolation between your Python packages.

## Endpoints

All endpoints are to be used with a `POST` request method. As such, you could use a solution like [Postman](https://www.getpostman.com/) to try your solution.

- `/api/register` To register a user it's body is a JSON composed of two fields which are `email` and `password`.
- `/api/login` To authenticate a user it's body is a JSON composed of two fields which are `email` and `password`. It will return a `access_token` field containing the token which should be used for all further requests
- `/api/upload` To upload an image. This time around, the body is a `form-data` and is composed of the `token` key which is the access token retrieved when authenticating the user and the `images` key which is a set of files. Please note that you will need to use [`multipart/form-data`](https://stackoverflow.com/a/4526286) as an encoding type to upload the images. It will return the metadata of all images that have been uploaded as well as images that have not been uploaded successfully.
- `/api/search` is similar to `/api/upload` in which it only requests one image instead of one/more and it returns the metadata of an image that is similar to the one uploaded. 

## Architecture

The document [`ARCHITECTURE.MD`](./ARCHITECTURE.MD) explains how the application should behave as well as the reasoning behind some of the choices made when designing it.

## Tests

Tests are also available for this project. In order to do run them:

- Install the required testing packages
  
  `pip install -r requirements_test.txt`

- Run the tests
  
  `pytest`
