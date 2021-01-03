# Create binary
FROM golang:1.15.6-alpine3.12 as build
WORKDIR /server
COPY ./build /server/build
COPY ./cmd /server/cmd
COPY ./pkg /server/pkg
COPY ./go.sum /server/go.sum
COPY ./go.mod /server/go.mod
COPY ./Makefile /server/Makefile
RUN apk add make
RUN make binary-locally

# Create image
FROM alpine:latest as prod
RUN apk add make
WORKDIR /server
COPY --from=build /server/bin/. /server/bin/.
COPY ./configs /server/configs
EXPOSE 8443
CMD [ "./bin/go-server", "--config", "configs/config.yaml" ]
