FROM golang:1.23.0
WORKDIR /app
COPY . .

COPY go.mod go.sum ./
RUN go mod download

RUN go build -o filend .

ENTRYPOINT ["/app/filend"]