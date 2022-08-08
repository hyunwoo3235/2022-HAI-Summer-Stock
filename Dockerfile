FROM golang:alpine as build

COPY . /app
WORKDIR /app

RUN GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w' -o main main.go

FROM alpine:latest
COPY --from=build /app /app
WORKDIR /app

CMD /app/main