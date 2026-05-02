FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /bin/api ./cmd/api

FROM gcr.io/distroless/static AS final

COPY --from=build /bin/api /app/api

EXPOSE 3000

CMD ["/app/api"]
