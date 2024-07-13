FROM golang:1.22-alpine as builder
WORKDIR /go/src/app
COPY . .
RUN go build -ldflags "-X 'internal/app.ApiVersion=1.0.0'"  -o /go/src/app/lite-api cmd/lite-api/main.go

FROM alpine
RUN mkdir /app
WORKDIR /app

COPY --from=builder /go/src/app/lite-api /app/lite-api
RUN chmod +x /app/lite-api

ENTRYPOINT ["/app/lite-api"]
CMD ["start"]