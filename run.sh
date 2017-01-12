#!/bin/bash
#docker run -it try
docker run -it --rm -p 8000:3000 -v ~/work/golang/crawler-go:/go/src/crawler-go -w /go/src/crawler-go try