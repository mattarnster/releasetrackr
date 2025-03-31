FROM golang:1.24.1

WORKDIR /go/src/releasetrackr
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["releasetrackr"]
