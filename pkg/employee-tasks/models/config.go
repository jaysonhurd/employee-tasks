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
