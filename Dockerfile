FROM golang:1.23

WORKDIR /project

COPY go.mod ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

EXPOSE 8080

CMD ["/docker-gs-ping"]

ENTRYPOINT [ "go", "run", "." ]
