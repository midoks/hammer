#!/bin/sh

echo "test"

docker build -t hammer:v1 ./


# docker run -p 9090:80 -d --network=host --cap-add=SYS_PTRACE --name hammer hammer:v1