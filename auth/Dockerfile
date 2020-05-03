FROM golang:1.13

ENV DOCKER_HOST=127.0.0.1:9191

WORKDIR /go/src/auth

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go install -v ./...

# Export necessary port
EXPOSE 9191

CMD ["auth"]