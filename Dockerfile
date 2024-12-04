FROM golang:1.21

EXPOSE 8080

RUN mkdir /app
WORKDIR /app

ADD . .

ENTRYPOINT [ "go", "run", "/app/cmd/pr9/main.go" ]