FROM golang:1.24-alpine AS build
WORKDIR /src
COPY contracts ./contracts
COPY rooms-service/go.mod rooms-service/go.sum ./rooms-service/
WORKDIR /src/rooms-service
RUN go mod download
COPY rooms-service .
RUN CGO_ENABLED=0 go build -o /service .
FROM scratch
COPY --from=build /service /service
ENTRYPOINT ["/service"]
