FROM golang:latest

COPY . /usr/local/go/src/RBACapp/

WORKDIR /usr/local/go/src/RBACapp/

RUN go mod download && go build -o RBACapp .

EXPOSE 8000

ENTRYPOINT go run RBACapp