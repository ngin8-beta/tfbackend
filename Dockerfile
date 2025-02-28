FROM golang:1.24-alpine AS builder

WORKDIR /go/src/github.com/ngin8-beta/tfbackend

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o tfbackend ./cmd/tfbackend



FROM alpine:3

RUN adduser -D appuser

COPY --from=builder /go/src/github.com/ngin8-beta/tfbackend/tfbackend /tfbackend

RUN chown appuser:appuser /tfbackend
RUN mkdir /data && chown appuser:appuser /data

USER appuser

ENTRYPOINT ["/tfbackend"]