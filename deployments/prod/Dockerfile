FROM golang:1.23.4-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go mod tidy
RUN go build -o /app/company_crud/main ./cmd/company_crud/

FROM debian:bullseye-slim

RUN useradd -m appuser

COPY --from=builder /app/company_crud/main /app/company_crud/main
COPY app.env /app/app.env

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8000

WORKDIR /app/company_crud

CMD ["/app/company_crud/main"]