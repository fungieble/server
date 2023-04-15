FROM golang:alpine AS base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /cmd/fungie/main cmd/fungie/main.go

FROM alpine

WORKDIR /app

COPY --from=base /cmd/fungie/main /app/

RUN chmod +x ./main

EXPOSE 8080

CMD ["./main"]