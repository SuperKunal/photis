FROM golang:1.14-alpine AS builder

RUN apk update && apk add git && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /go/src/photis

COPY . .

RUN go mod download
RUN go build -o main -ldflags '-w -s' .

EXPOSE 5000

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/photis/main .
CMD ["./main"]