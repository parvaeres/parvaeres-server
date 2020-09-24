FROM golang:1.14 AS build
WORKDIR /go/src

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY pkg pkg
COPY main.go main.go
ENV CGO_ENABLED=0
RUN go build -a -installsuffix cgo -o parvaeres .

FROM scratch AS runtime
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/parvaeres ./
EXPOSE 8080/tcp
ENTRYPOINT ["./parvaeres"]
