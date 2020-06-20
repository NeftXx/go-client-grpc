#!/bin/bash
docker login
docker build -t go-client-grpc .
docker tag go-client-grpc $1/go-client-grpc
docker push $1/go-client-grpc
