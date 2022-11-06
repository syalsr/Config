FROM golang:latest


WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o config ./app/cmd/main.go
CMD ["./config"]

EXPOSE 8081 50051