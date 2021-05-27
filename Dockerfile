FROM golang:1.16.4-buster AS builder
RUN mkdir /workdir
WORKDIR /workdir
COPY ./srcs .
ENV GORARCH amd64
ENV GOOS linux
RUN go mod tidy
RUN go build -o app .

FROM debian:buster
ENV API_HOST 0.0.0.0:4000
COPY --from=builder /workdir/app .
CMD ["./app"]
