#!/bin/bash

go mod tidy

CGO_ENABLED=0 swag init -g ./main.go -o ./app/main/docs

air
