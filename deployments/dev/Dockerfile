#Stage 1: Builder
FROM golang:1.23.4-alpine3.21 as builder

#Install necessary tools
RUN apk update && \
    apk add --no-cache git

#Install Delve debugger
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

RUN pwd
RUN ls

#Copy and download dependencies
COPY go.mod ./
RUN go mod download


#Build the application
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /app/company_crud/main ./cmd/company_crud/

#Stage 2: Final image
FROM alpine:3.21.0

#Install necessary runtime dependencies
RUN apk update && \
    apk add --no-cache ca-certificates

#Copy the binary from the builder stage
COPY --from=builder /go/bin/dlv /usr/local/bin/dlv
COPY --from=builder /app/company_crud/main /app/company_crud/main

COPY app.env /app/app.env

EXPOSE 8000
EXPOSE 40000

WORKDIR /app/company_crud

#Run the application with the Delve debugger
CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "--continue=true", "exec", "main"]