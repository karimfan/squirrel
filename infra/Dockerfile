# Build stage
FROM golang:1.20 AS build
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server

# Runtime stage
FROM gcr.io/distroless/base-debian11
COPY --from=build /app/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
