FROM golang:1.15.6-alpine3.12
WORKDIR /server
COPY ./cmd /server/cmd
COPY ./pkg /server/pkg
COPY ./test /server/test
COPY ./go.sum /server/go.sum
COPY ./go.mod /server/go.mod
COPY ./Makefile /server/Makefile
RUN apk add bash
RUN apk add make
RUN make test-locally
