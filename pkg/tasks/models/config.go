package models

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/lib/pq"
)

type Config struct {
	Postgres      PostgresConfig       `json:"postgres"`
	ElasticSearch string               `json:"elasticsearch"`
	ElasticConfig elasticsearch.Config `json:"elastic_config"`
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
	SSLMode  string `json:"sslmode"`
}

type Task struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Create_time string        `json:"create_time"`
	Owners      pq.Int32Array `json:"owners"`
	Private     bool          `json:"private"`
	Due_by      string        `json:"due_by"`
}

type Employee struct {
	Employee_id    int    `json:"employee_id"`
	Nickname       string `json:"nickname"`
	First_name     string `json:"first_name"`
	Last_name      string `json:"last_name"`
	Street_address string `json:"street_address"`
	City           string `json:"city"`
	State          string `json:"state"`
	Zip            string `json:"zip"`
}
