FROM golang:alpine AS builder
RUN apk add --no-cache --update \
        git \
        ca-certificates
ADD . /app
WORKDIR /app
COPY go.mod ./
RUN  go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /main .

FROM alpine
COPY --from=builder /main ./
RUN chmod +x ./main
ENTRYPOINT ["./main"]