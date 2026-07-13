FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/mintok-api ./cmd/api

FROM alpine:3.22

RUN addgroup -S mintok && adduser -S mintok -G mintok

USER mintok
WORKDIR /app

COPY --from=build /bin/mintok-api /app/mintok-api
COPY --from=build /app/migrations /app/migrations

EXPOSE 8080

ENTRYPOINT ["/app/mintok-api"]
