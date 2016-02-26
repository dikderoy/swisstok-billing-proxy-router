FROM golang:1.6
MAINTAINER "Roman Bulgakov (roman.bulgakov@morepower.ru)"

RUN mkdir -p /go/src/app
RUN mkdir -p /go/build
WORKDIR /go/src/app

# this will ideally be built by the ONBUILD below ;)
CMD ["go-wrapper", "run"]

ONBUILD COPY . /go/src/app
ONBUILD RUN go-wrapper download
ONBUILD RUN go-wrapper install
