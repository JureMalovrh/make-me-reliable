FROM golang:latest AS build-env
WORKDIR /app
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.39.0
COPY . .
RUN go mod download

FROM build-env AS builder
RUN GOOS=linux CGO_ENABLED=0 go build -v -o reliable-api ./cmd

FROM alpine:latest AS production
WORKDIR /app
COPY --from=builder /app/reliable-api .
CMD [ "./reliable-api" ]
