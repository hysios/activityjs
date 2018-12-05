FROM node:8.14.0-jessie as web

COPY ./web /var/web
WORKDIR /var/web
RUN yarn install
RUN yarn build

FROM golang:1.11.2-stretch as serve

COPY ./serve /go/src/activityjs.io/serve
WORKDIR  /go/src/activityjs.io/serve
ENV GO111MODULE=off
RUN go get ./...
RUN go get -u github.com/gopherjs/gopherjs

EXPOSE 8080

COPY --from=web /var/web/build /go/src/activityjs.io/serve/static
RUN gopherjs build demo/main.go -o demo/main.js
ENV WORKDIR=/go
RUN mkdir /go/tmp
ENTRYPOINT [ "go", "run", "serve/main.go" ]