FROM golang:1.8-alpine

WORKDIR /go/src/releasetrackr
COPY . .

RUN apk add --no-cache git

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"]