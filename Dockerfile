FROM golang:1.23-bookworm AS builder
ARG VERSION
ARG BUILD_HASH
ARG BUILD_DATE
ARG BUILD_PGO

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux
RUN go build -pgo=$BUILD_PGO -ldflags "-X main.Version=$VERSION -X main.BuildHash=$BUILD_HASH -X main.BuildDate=$BUILD_DATE" -o mox ./cmd/mox/*.go

FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=builder /app/mox .

ENTRYPOINT ["/mox"]
