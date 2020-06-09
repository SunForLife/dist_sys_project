FROM golang:1.13

ENV DOCKER_HOST=127.0.0.1:8181

WORKDIR /go/src/app

COPY lib ../lib

COPY upload/go.mod .
COPY upload/go.sum .
RUN go mod download

COPY upload/. .
RUN go install -v ./...

# Export necessary port
EXPOSE 8181

CMD ["upload"]