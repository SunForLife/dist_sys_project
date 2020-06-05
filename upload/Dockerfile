FROM golang:1.13

ENV DOCKER_HOST=127.0.0.1:8181

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go install -v ./...

# Export necessary port
EXPOSE 8181

CMD ["upload"]