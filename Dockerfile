FROM golang:1.6

#ADD ./src /go/src
#RUN go get -v hello
#RUN go build -o /hello hello

ADD ./src /go/src
#RUN go get -v golang.org/x/net
RUN go get -v github.com/codegangsta/negroni
RUN go get -v github.com/PuerkitoBio/goquery
RUN go get -v github.com/xyproto/mooseware
EXPOSE 3000

#CMD ["/hello"]

