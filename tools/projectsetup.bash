#!/bin/sh

docker run -e POSTGRES_USER=myusername -e POSTGRES_PASSWORD=mypassword -p 5432:5432 -v /data:/home/jayson/go/src/github.com/jaysonhurd/employee-tasks/tools/postgres/data -d postgres

docker network create elastic

docker run -d --net elastic -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:8.4.0
docker run -d -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:8.4.0



