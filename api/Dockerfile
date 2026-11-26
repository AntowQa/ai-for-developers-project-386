FROM golang:1.26.1-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api

FROM alpine:3.22
RUN apk add --no-cache wget
COPY --from=build /out/api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
