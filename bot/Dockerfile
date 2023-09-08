# Use an official Go image as a build stage
FROM golang:1.20-alpine AS build

WORKDIR /
COPY . /

# Install necessary packages and update CA certificates
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates
RUN go build -o /app/main ./bot/cmd/main.go

# Create a minimal runtime image
FROM scratch

WORKDIR /app

# Copy the built binary and CA certificates from the build stage
COPY --from=build /app/main /app/main
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 9013
ENTRYPOINT ["/app/main"]
