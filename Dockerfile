FROM golang:1.17

WORKDIR /go/src/app

COPY go.mod ./
COPY go.sum ./
COPY *.go ./

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build
CMD [ "./qzcatch" ]
