# employee-tasks
Demonstration of an employee task system

## Outline
This project uses Postgres as a persistent data store along with ElasticSearch for 
logging and runtime application alerts.  The Go service is an API which connects to this
to these data sources and provides interaction with the data (updating, inserting
deleting, etc.).

## Starting Prerequisite Databases
If you do not already have Postgres and ElasticSearch running, please 
follow the below instructions for each to run as a Docker container locally.
This assumes that you have `dockerd` running locally so that you can 
run Docker containers.

### Postgres
Install Postgres client (`psql`) if not done already
redhat & centos : `dnf install postgresql.x86_64`

Run a Docker container locally (you can change the `/data` volume location if you would like it somewhere else):
```azure
docker run -e POSTGRES_USER=myusername -e POSTGRES_PASSWORD=mypassword -p 5432:5432 -v /data:$HOME/data -d postgres
```
You will need to create the database first:
```azure
psql -d postgres -U myusername -h 0.0.0.0 -p 5432
drop database if exists workers;
create database workers;
```

You will then need to reconnect to the new DB and load the data:
```azure
\q
psql -d workers -U myusername -h 0.0.0.0 -p 5432
```
Then create tables and load the data by running the script `tools/postgres/populate_postgres.sql`


### ElasticSearch

Run the ElasticSearch Docker container:

```docker run -d --name es01 -p 9200:9200 -p 9300:9300 -e "http.publish_host=127.0.0.1" -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:8.4.0```

Put the `http_ca.crt` in the `config` directory so that the Go application can access:

```azure
cd config
docker cp :/usr/share/elasticsearch/config/certs/http_ca.crt .
chmod 644 http_ca.crt
```

Reset the password to ElasticSearch:

```azure
docker exec -it es01 /usr/share/elasticsearch/bin/elasticsearch-reset-password -u elastic
```
Take this password created and add it into the `config/config.json` field replacing `<PASSWORD_HERE>`:

```  "elastic_config": {
    "Addresses": ["https://127.0.0.1:9200"],
    "Username": "elastic",
    "Password": "<PASSWORD_HERE>"
  }
```

### Starting your service
`.idea` has been included here if you are using JetBrains Goland.  If not run:

//TODO: add `go build ...` HERE

