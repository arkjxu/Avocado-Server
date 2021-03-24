FROM golang:alpine

COPY dist dist

WORKDIR /go/dist

RUN apk update && apk --no-cache add tzdata

ENTRYPOINT [ "./avocado_server" ]