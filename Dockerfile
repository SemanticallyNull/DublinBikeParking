# Stage 1: Build the frontend
FROM node:22-alpine AS build-frontend

WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build the Go binary
FROM golang:1.21-bookworm AS build-go

WORKDIR /go/src/app
ADD . /go/src/app
COPY --from=build-frontend /app/static/ /go/src/app/static/

RUN go build -o /go/bin/dublinbikeparking

# Stage 3: Minimal runtime image
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=build-go /go/bin/dublinbikeparking /app/dublinbikeparking
COPY --from=build-go /go/src/app/static/ /app/static/

CMD ["/app/dublinbikeparking"]
