# Utsuru Concept 📷

This application is used for a photo repository as per [Shopify]([https://www.shopify.com/](https://www.shopify.com/)'s backend [coding challenge]([https://docs.google.com/document/d/1ZKRywXQLZWOqVOHC4JkF3LqdpO3Llpfk_CkZPR8bjak/edit](https://docs.google.com/document/d/1ZKRywXQLZWOqVOHC4JkF3LqdpO3Llpfk_CkZPR8bjak/edit)) for Winter 2020 internships.

## Why Utsuru Concept?

Utsuru comes from the Japanese word **映る** which means "to reflect/to project". In that sense, this service would be used for saving various reflections in the form of photography. 

Concept means that it is a proof-of-concept and that it may or may not be used for future usage.

## Prerequisites

- [Python 3]([https://www.python.org/](https://www.python.org/)

- [PostgreSQL]([https://www.postgresql.org/](https://www.postgresql.org/)

- [Google Cloud Platform]([https://cloud.google.com/](https://cloud.google.com/)

## How to get started

- Clone this repository
  
  `git clone https://github.com/zaclimon/UtsuruConcept`

- Install the packages
  
  `pip install -r requirements.txt`

- Fill the environment variables inside `.env.sample` and rename it to `.env`.
  
  - For Google Cloud Storage, you will need to obtain your [service account's]([https://cloud.google.com/docs/authentication/production](https://cloud.google.com/docs/authentication/production) `.json` file.
  
  - For `SQLALCHEMY_DATABASE_URI`, please refer to [Flask's]([https://flask-sqlalchemy.palletsprojects.com/en/2.x/config/](https://flask-sqlalchemy.palletsprojects.com/en/2.x/config/) documentation for examples on how to connect your database.

- Start Flask
  
  `flask run`

- Enjoy! 🎉
  
  **Note**: It is possible to also use a [virtual environment]([https://docs.python.org/3/library/venv.html](https://docs.python.org/3/library/venv.html) for better isolation between your Python packages.

## Architecture

The document [`ARCHITECTURE.MD`](./ARCHITECTURE.MD) explains how the application should behave as well as the reasoning behind some of the choices made when designing it.

## Tests

Tests are also available for this project. In order to do run them:

- Install the required testing packages
  
  `pip install -r requirements_test.txt`

- Run the tests
  
  `pytest`
