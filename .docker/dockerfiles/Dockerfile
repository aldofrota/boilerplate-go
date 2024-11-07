FROM golang:1.23 AS builder

WORKDIR /usr/boilerplate-go/app

RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install gotest.tools/gotestsum@latest
RUN apt-get install -y git curl

COPY . .

EXPOSE 9001

CMD [ "tail", "-f", "/dev/null" ]