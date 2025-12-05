FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.sum go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o imphnen .

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /app

COPY --from=builder imphnen .

USER nonroot:nonroot

EXPOSE 8080

CMD["imphnen"]
