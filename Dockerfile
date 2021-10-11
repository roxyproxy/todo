FROM golang

WORKDIR /go/src/github.com/roxyproxy/todo

COPY . .

RUN go build -o todo


CMD ./todo