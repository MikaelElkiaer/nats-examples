FROM golang:alpine AS builder
RUN mkdir /build
WORKDIR /build
COPY main.go go.mod go.sum ./
RUN go build -o client .

FROM alpine
COPY --from=builder /build/client /usr/bin/client
ENTRYPOINT [ "/usr/bin/client" ]
