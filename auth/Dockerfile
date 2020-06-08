FROM golang:1.13

ENV DOCKER_HOST=127.0.0.1:9191

WORKDIR /go/src/auth

COPY lib ../lib

COPY auth/go.mod .
COPY auth/go.sum .
RUN go mod download

COPY auth/. .
RUN go install -v ./...
# RUN go test -v ./...

# Export necessary port
EXPOSE 9191

CMD ["auth"]