FROM golang:1.6

RUN mkdir -p /go/src/github.com/matthauglustaine/ddleash/
WORKDIR /go/src/github.com/matthauglustaine/ddleash/

COPY . /go/src/github.com/matthauglustaine/ddleash/

RUN go-wrapper download github.com/matthauglustaine/ddleash/cmd/ddleash
RUN go-wrapper install github.com/matthauglustaine/ddleash/cmd/ddleash

ENTRYPOINT [ "ddleash" ]
