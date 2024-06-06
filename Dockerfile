#FROM golang:1.21.0-alpine AS build-stage
#
#WORKDIR /app
#
#COPY go.mod go.sum ./
#RUN go mod download && go mod tidy && go clean --modcache
#
#COPY ./ ./
#
#RUN CGO_ENABLED=0 GOOS=linux go build -o ./ ./server.go
#
#
#FROM gcr.io/distroless/base-debian11 AS build-release-stage
#
#WORKDIR /
#
#COPY --from=build-stage app app
#COPY . .
#
#EXPOSE 8080
#
#USER nonroot:nonroot
#
#CMD ["/app"]


FROM golang:1.22-alpine AS builder

WORKDIR /app
RUN apk --no-cache add bash git make # gcc gettext musl-dev

# dependencies
COPY go.mod go.sum ./
RUN  go mod tidy && go mod download && go clean --modcache

# build
COPY . .
RUN go build -o ./bin/app ./server.go

# run
FROM alpine AS runner

COPY --from=builder /app/bin/app /
COPY .env .env
CMD ["/app"]