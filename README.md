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

### Postgres & ElasticSearch Setup

Install Postgres client
`redhat, centos` : dnf install postgresql.x86_64

