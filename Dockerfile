FROM golang:1.18-alpine3.14

WORKDIR /build

COPY . .
RUN apk add --update make && apk add --update gcc && apk add --update musl-dev
RUN echo "foo"
ENV GOOS=linux GOARCH=amd64
RUN make build_all