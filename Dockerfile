FROM golang:1.8

WORKDIR /go/src/releasetrackr
COPY . .

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"]