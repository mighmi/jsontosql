FROM postgres:alpine

COPY ./scripts/tables.sql /docker-entrypoint-initdb.d

###

FROM golang:alpine as builder
#overcome tls: failed to verify certificate: x509
RUN apk update && apk upgrade && apk add --no-cache ca-certificates && \
update-ca-certificates

RUN mkdir ../home/app

WORKDIR ../home/app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o mybinary

####

FROM scratch

COPY --from=builder /home/app/mybinary .

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["./mybinary"]



# ------- how to get around
# // panic: Get "https://dummyjson.com/users": tls: failed to verify certificate: x509: certificate signed by unknown authority


# use CA bundle in Dockerfile along these lines
# # Copy the CA certificate into the image
# COPY ca.crt /app
# # Set the environment variable to point to the CA certificate file
# ENV SSL_CERT_FILE=/app/ca.crt

# # use Host's certificates
# # Mount the host's SSL certificate files into the container
# RUN mkdir -p /etc/ssl/certs
# COPY /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
