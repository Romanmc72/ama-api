# Ask Me Anything API

This code represents the API server that will handle requests from the front end(s) and serve data to and from the database.

This API is written in go because I have not touched Go in a while and wanted to get back into it.

It will call the database (GCP Firestore) and return data in JSON format. That is really all I have for now. The data model and all are kind of in my head but I will toss that in here too once I get a more solid grasp on the basic outline.

## Local Dev

The local development setup consists of 2 parts.

1. The Firestore Emulator
1. The ama api

### Firestore Emulator

To launch the firestore emulator run `./build.sh db`

It should stand up a local instance of firestore that can be used to test the application.

Exit with Ctrl+c

### The AMA API

To launch the API locally just run `./build.sh run` and the api will spin up. Just make sure the database is also up if you intend to try and hit the API.

Exit with Ctrl+c

## Deploying

The code will be built and deployed using Docker. The target compute platform will be a GCP Cloud Run service. The docker image will be pushed up to the GCP Artifact Registry in the target GCP project.

Take a look at the Dockerfile for details and the associated CDK infrastructure repo.

## Testing

Run unit tests using the `./build.sh test` script and run integration tests after standing up the API and database using `./build.sh integ`.

There are some integration tests to simply set up the test environment for app development locally, those can be run with `./build.sh integsetup` and the tear down can be run with `./build.sh integteardown`.
