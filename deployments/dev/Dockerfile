#Stage 1: Builder
FROM golang:1.23.4-bullseye as builder

#Install Delve debugger
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

#Copy and download dependencies
COPY go.mod ./
RUN go mod download


#Build the application
COPY . .
RUN go mod tidy
RUN go build -o /app/company_crud/main ./cmd/company_crud/

#Stage 2: Final image
FROM debian:bullseye-slim

#Copy the binary from the builder stage
COPY --from=builder /go/bin/dlv /usr/local/bin/dlv
COPY --from=builder /app/company_crud/main /app/company_crud/main

COPY app.env /app/app.env

EXPOSE 8000
EXPOSE 40000

WORKDIR /app/company_crud

#Run the application with the Delve debugger
CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "--continue=true", "exec", "main"]