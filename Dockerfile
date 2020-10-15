FROM golang:1.14 AS builder
ENV CGO_ENABLED 0
WORKDIR /go/src/app
ADD . .
RUN go build -o /auto-retain-pv

FROM alpine:3.12
COPY --from=builder /auto-retain-pv /auto-retain-pv
CMD ["/auto-retain-pv"]