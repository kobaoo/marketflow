FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

COPY . .

RUN go build -o marketflow .

FROM alpine:3.20

RUN adduser -D appuser
USER appuser

COPY --from=builder /app/marketflow /marketflow

EXPOSE 8080

CMD ["./marketflow"]
