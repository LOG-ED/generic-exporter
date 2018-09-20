FROM golang as builder
WORKDIR /go/src/github.com/LOG-ED/generic-exporter

# Install and run dep
RUN go get -u github.com/golang/dep/cmd/dep
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

# Copy the code and compile it
ADD . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /exporter ./cmd/exporter/

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /exporter ./
CMD ["./exporter"]