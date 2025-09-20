# ---------- base ----------
FROM golang:1.24-alpine AS base
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# ---------- test stage ----------
FROM base AS test
RUN go test ./... -v

# ---------- build stage ----------
FROM base AS build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api ./cmd/api

# ---------- run stage ----------
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/api /app/api
COPY docs/api/v1/openapi.yaml /app/docs/api/v1/openapi.yaml
ENV HTTP_ADDR=:8080
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/api"]