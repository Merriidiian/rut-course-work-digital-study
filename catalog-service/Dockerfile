FROM golang:1.24-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /catalog .
FROM scratch
COPY --from=build /catalog /catalog
ENTRYPOINT ["/catalog"]
