# Start by building the application.
FROM golang:1.13-buster as build-go

WORKDIR /go/src/app
ADD . /go/src/app

RUN go build -o /go/bin/dublinbikeparking

FROM node:12.2.0-alpine as build-js

WORKDIR /go/src/app

# install and cache app dependencies
COPY static-vue/package.json /go/src/app/package.json
RUN npm install
RUN npm install @vue/cli -g

ADD static-vue/ /go/src/app/
RUN npm run build

# Now copy it into our base image.
FROM debian:buster

WORKDIR /app
COPY --from=build-go /go/bin/dublinbikeparking /app/dublinbikeparking
COPY --from=build-go /go/src/app/static/ /app/static/
COPY --from=build-js /go/src/app/dist/ /app/static-vue/dist/

CMD ["/app/dublinbikeparking"]
