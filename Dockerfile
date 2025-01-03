FROM golang:1.23.0-alpine as golang

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM scratch

COPY --from=golang /app/main /app/main

EXPOSE 4000
EXPOSE 4001

CMD ["/app/main"]