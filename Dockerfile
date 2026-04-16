# Pin to a Go version that satisfies go.mod (bump if go.mod `go` directive changes).
FROM golang:alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /orbit ./cmd/orbit

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /orbit .
ENV PORT=8080
EXPOSE 8080
USER 65532:65532
CMD ["./orbit"]
