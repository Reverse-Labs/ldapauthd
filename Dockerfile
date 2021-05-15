FROM golang:1.16.4-alpine3.13

RUN mkdir -p /go/src/github.com/Reverse-Labs/ldapauthd
WORKDIR /go/src/github.com/Reverse-Labs/ldapauthd
COPY . .
RUN go build -o /usr/bin/ldapauthd
CMD ["/usr/bin/ldapauthd", "serve"]