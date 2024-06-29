FROM golang:1.22.0

RUN mkdir /usr/local/ama
WORKDIR /usr/local/ama

ADD ./ ./
RUN go get .
RUN go build .

RUN groupadd ama \
  && useradd -g ama ama \
  && mkdir -p /home/ama \
  && cp ./api /home/ama
USER ama
WORKDIR  /home/ama

ENTRYPOINT [ "./api" ]
