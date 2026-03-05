FROM golang:1.25

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY src/ ./src/
COPY services/ ./services/

RUN CGO_ENABLED=0 GOOS=linux go build -o /notes_app

EXPOSE 8080

CMD ["/notes_app"]