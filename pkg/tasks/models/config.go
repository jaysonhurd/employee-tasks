package models

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jaysonhurd/employee-tasks/internal/pkg/postgres"
)

type Config struct {
	Postgres      postgres.PostgresConfig `json:"postgres"`
	ElasticSearch string                  `json:"elasticsearch"`
	ElasticConfig elasticsearch.Config    `json:"elastic_config"`
}

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Create_time string `json:"create_time"`
	Owners      []int  `json:"owners"`
	Private     bool   `json:"private"`
	Due_by      string `json:"due_by"`
}
