FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -v -o unreliable-api ./main.go

FROM alpine:latest AS production
WORKDIR /app
COPY --from=builder /app/unreliable-api .
RUN chmod +x unreliable-api
CMD [ "./unreliable-api" ]
