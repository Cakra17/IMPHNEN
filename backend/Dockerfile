FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.sum go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o imphnen ./cmd/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /app

COPY --from=builder /app/imphnen .

USER nonroot:nonroot

EXPOSE 6969

CMD ["./imphnen"]
