# Ask Me Anything API

This code represents the API server that will handle requests from the front end(s) and serve data to and from the database.

This API is written in go because I have not touched Go in a while and wanted to get back into it.

It will call the database (GCP Firestore) and return data in JSON format. That is really all I have for now. The data model and all are kind of in my head but I will toss that in here too once I get a more solid grasp on the basic outline.

## Local Dev

The local development setup consists of 2 parts.

1. The Firestore Emulator
1. The ama api

### Firestore Emulator

To launch the firestore emulator run `make db`

It should stand up a local instance of firestore that can be used to test the application.

Exit with Ctrl+c

### The AMA API

To launch the API locally just run `make run` and the api will spin up.

Exit with Ctrl+c

## Deploying

The code will be built and deployed using Docker. The target compute platform will be a GCP Cloud Run service. The docker image will be pushed up to the GCP Artifact Registry in the target GCP project.

Take a look at the Dockerfile for details.

## Testing

I plan to set up unit testing and maybe somewhere some integration testing. I have not done golang unit tests before so I hope to learn some about how to do that here.
