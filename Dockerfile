FROM golang:1.25.0-alpine3.22 AS build-stage

RUN mkdir -p /usr/local/ama
WORKDIR /usr/local/ama

COPY ./ ./
RUN go get .
RUN go build .

## STAGE: 2 ##
FROM alpine:3.22.1
LABEL maintainer="Roman Czerwinski romanczerwinski@r0m4n.com"

RUN addgroup ama \
    && adduser -h /home/ama -D -G ama ama

COPY --from=build-stage --chown=ama:ama /usr/local/ama/api /home/ama/api
USER ama
WORKDIR /home/ama
EXPOSE 8088

# The GCP Project ID this app will point at (default is dev)
ENV ENVIRONMENT_NAME='development'
ENV PROJECT_ID='ama-dev-414718'
ENV GO_LOG='info'

ENTRYPOINT [ "/home/ama/api" ]
