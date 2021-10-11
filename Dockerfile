FROM golang

WORKDIR /go/src/github.com/roxyproxy/todo

RUN go build -o todo

CMD ./todo