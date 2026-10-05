FROM golang:1.24-alpine AS build
WORKDIR /src
COPY contracts ./contracts
COPY schedule-service/go.mod schedule-service/go.sum ./schedule-service/
WORKDIR /src/schedule-service
RUN go mod download
COPY schedule-service .
RUN CGO_ENABLED=0 go build -o /service .
FROM scratch
COPY --from=build /service /service
ENTRYPOINT ["/service"]
