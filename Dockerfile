# Start by building the application.
FROM golang:1.21-bookworm as build-go

WORKDIR /go/src/app
ADD . /go/src/app

RUN go build -o /go/bin/dublinbikeparking

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=build-go /go/bin/dublinbikeparking /app/dublinbikeparking
COPY --from=build-go /go/src/app/static/ /app/static/

CMD ["/app/dublinbikeparking"]
