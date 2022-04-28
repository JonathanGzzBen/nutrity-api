#!/usr/bin/env bash
cd ./api/v1
mockgen -source repository/users.go -destination repository/mocks/UsersRepository.go -package mocks