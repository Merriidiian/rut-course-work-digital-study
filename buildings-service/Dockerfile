FROM golang:1.24-alpine AS build
WORKDIR /src
COPY contracts ./contracts
COPY buildings-service/go.mod buildings-service/go.sum ./buildings-service/
WORKDIR /src/buildings-service
RUN go mod download
COPY buildings-service .
RUN CGO_ENABLED=0 go build -o /service .
FROM scratch
COPY --from=build /service /service
ENTRYPOINT ["/service"]
