FROM golang:1.22 as base

FROM base as documentation
WORKDIR /app

COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g ./main/main.go -o ./main/docs

FROM base as builder

WORKDIR /app

COPY . .
# Instale o compilador protobuf e as ferramentas gRPC
RUN apt-get update && apt-get install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN protoc --go_out=. --go-grpc_out=. proto/*.proto

COPY --from=documentation /app/main/docs ./main/docs
RUN GOOS=linux GOARCH=amd64 go build -v -o server ./main

FROM base as production

WORKDIR /app
ENV GIN_MODE=release
RUN apt-get update && apt-get install -y curl
COPY --from=builder /app/server .


CMD ["./server"]