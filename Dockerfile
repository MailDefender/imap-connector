# syntax=docker/dockerfile:1

FROM golang:1.25 AS build

ENV GOROOT=

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY --exclude=go.* . .

# Build
WORKDIR /app/cmd/imap-connector
RUN CGO_ENABLED=0 GOOS=linux go build -o /imap-connector



FROM alpine

COPY --from=build /imap-connector /imap-connector

EXPOSE 8080

# Run
CMD ["/imap-connector"]