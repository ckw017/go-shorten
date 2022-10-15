FROM golang

EXPOSE 8080

ARG package=github.com/ckw017/go-shorten
# ARG package=app
# ${PWD#$GOPATH/src/}

RUN mkdir -p /go/src/${package}
WORKDIR /go/src/${package}

COPY . /go/src/${package}
RUN go get -t -v ./...
RUN go install .

CMD ["go-shorten"]
