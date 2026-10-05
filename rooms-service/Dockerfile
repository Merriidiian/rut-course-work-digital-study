FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY internal ./internal
COPY rooms-service ./rooms-service
RUN CGO_ENABLED=0 go build -o /service ./rooms-service
FROM scratch
COPY --from=build /service /service
ENTRYPOINT ["/service"]
