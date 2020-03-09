FROM golang:1.13

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

# Export necessary port
EXPOSE 7171

CMD ["app"]