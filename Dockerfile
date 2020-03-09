FROM golang:1.13

ENV DOCKER_HOST=127.0.0.1:7171

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

# Export necessary port
EXPOSE 7171

CMD ["app"]