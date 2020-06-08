FROM golang:1.13

ENV DOCKER_HOST=127.0.0.1:7171

WORKDIR /go/src/app

COPY lib ../lib

COPY app/go.mod .
COPY app/go.sum .
RUN go mod download

COPY app/. .
RUN go install -v ./...

# Export necessary port
EXPOSE 7171

CMD ["app"]