

FROM golang:1.16-alpine 
ENV PORT $PORT

WORKDIR /go/src/app
COPY . .
CMD go run server.go