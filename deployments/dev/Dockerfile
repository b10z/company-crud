FROM golang:1.23.4-bullseye as builder

RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go mod tidy
RUN go build -o /app/company_crud/main ./cmd/company_crud/

FROM debian:bullseye-slim

COPY --from=builder /go/bin/dlv /usr/local/bin/dlv
COPY --from=builder /app/company_crud/main /app/company_crud/main

COPY app.env /app/app.env

EXPOSE 8000
EXPOSE 40000

WORKDIR /app/company_crud

CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "--continue=true", "exec", "main"]