FROM golang:1.21-alpine AS build
RUN apk add --no-cache git build-base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server main.go

# final stage
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=build /app/server /app/server
# если нужны файлы статические - копируем их
# COPY --from=build /app/uploads /app/uploads
USER nonroot:nonroot
EXPOSE 8082
CMD ["/app/server"]