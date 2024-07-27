FROM golang:1.22-alpine AS builder

COPY . /app
WORKDIR /app

RUN go mod download
RUN go build -o ./bin/auth_service cmd/main.go

FROM alpine:3.20

WORKDIR /root

COPY --from=builder /app/bin/auth_service .
COPY --from=builder /app/${ENV_FILE} .
RUN apk add curl

CMD [ "./auth_service" ]
