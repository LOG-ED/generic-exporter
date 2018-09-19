FROM golang as builder
RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR /go/src/github.com/LOG-ED/generic-exporter
ADD . /go/src/github.com/LOG-ED/generic-exporter
RUN dep ensure 

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o exporter ./cmd/exporter/

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/LOG-ED/generic-exporter/exporter .
CMD ["./exporter"]