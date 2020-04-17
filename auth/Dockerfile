FROM golang:1.13

ENV DOCKER_HOST=127.0.0.1:9191

WORKDIR /go/src/auth
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

# Export necessary port
EXPOSE 9191

CMD ["auth"]