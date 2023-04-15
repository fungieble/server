FROM golang:alpine AS base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /cmd/alia/main cmd/alia/main.go

FROM alpine

WORKDIR /app

COPY --from=base /cmd/alia/main /app/

RUN chmod +x ./main

EXPOSE 8080

CMD ["./main"]