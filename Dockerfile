FROM golang:1.23

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY *.go ./
COPY src/ ./src/

RUN CGO_ENABLED=0 GOOS=linux go build -o /notes_app

EXPOSE 8080

CMD ["/bleh"]