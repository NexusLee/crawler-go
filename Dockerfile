FROM golang:1.6

ADD ./src /go/src
RUN go get -v hello
RUN go build -o /hello hello
EXPOSE 3000

CMD ["/hello"]