FROM golang:1.23.4-alpine3.21 as builder

RUN apk update
RUN apk add git

WORKDIR /app

COPY go.mod ./

RUN go mod download
RUN go install go.uber.org/mock/mockgen@latest
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.2
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.2
RUN go mod tidy

FROM golang:1.23.4-alpine3.21

RUN apk update
RUN apk add build-base bash
RUN apk add postgresql

COPY --from=builder /go /go

WORKDIR /app

COPY . .