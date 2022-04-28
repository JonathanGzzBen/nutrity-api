#!/usr/bin/env bash
go get -u github.com/swaggo/swag/cmd/swag
cd ./api/v1
swag init