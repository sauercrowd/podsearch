FROM golang:1.8
MAINTAINER sauercrowd <jonadev95@posteo.org>

ADD . /go/src/github.com/sauercrowd/podsearch
RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR /go/src/github.com/sauercrowd/podsearch
RUN dep ensure
RUN go install github.com/sauercrowd/podsearch
EXPOSE 8080
