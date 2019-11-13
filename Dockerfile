FROM golang:1.13 AS builder
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -o simplelb cmd/simplelb/main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /app/simplelb .
ENTRYPOINT [ "/root/simplelb" ]
