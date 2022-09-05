package elasticsearch

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/jaysonhurd/employee-tasks/internal/pkg/postgres"
	"github.com/rs/zerolog"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type ElasticSearcher interface {
	LoadEmployeesFromPostgres() error
	LoadTasksFromPostgres() error
}

const (
	esCertPath = "config/http_ca.crt"
)

type elasticSearch struct {
	esconn *elasticsearch.Client
	pgconn postgres.Postgreser
	l      *zerolog.Logger
}

func New(
	esconn *elasticsearch.Client,
	pgconn postgres.Postgreser,
	l *zerolog.Logger,
) ElasticSearcher {
	return &elasticSearch{
		esconn: esconn,
		pgconn: pgconn,
		l:      l,
	}
}

func (e *elasticSearch) LoadEmployeesFromPostgres() error {
	var employees []postgres.Employee
	ctx := context.Background()
	employees, err := e.pgconn.AllEmployees()
	if err != nil {
		return err
	}

	i := 0
	for range employees {
		employeeJSON, err := json.Marshal(employees[i])
		if err != nil {
			return err
		}
		req := esapi.IndexRequest{
			Index:      "employee",
			DocumentID: strconv.Itoa(i + 1),
			Body:       bytes.NewReader(employeeJSON),
			Refresh:    "true",
		}
		_, err = req.Do(ctx, e.esconn)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *elasticSearch) LoadTasksFromPostgres() error {
	var tasks []postgres.Task
	ctx := context.Background()
	tasks, err := e.pgconn.AllTasks()
	if err != nil {
		return err
	}

	i := 0
	for range tasks {
		taskJSON, err := json.Marshal(tasks[i])
		if err != nil {
			return err
		}
		req := esapi.IndexRequest{
			Index:      "employee",
			DocumentID: strconv.Itoa(i + 1),
			Body:       bytes.NewReader(taskJSON),
			Refresh:    "true",
		}
		_, err = req.Do(ctx, e.esconn)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewESConnection(config elasticsearch.Config, l *zerolog.Logger) (*elasticsearch.Client, error) {
	var es *elasticsearch.Client
	var err error

	cert, err := getESCACert(esCertPath)
	if err != nil {
		return es, err
	}
	config.CACert = cert
	config.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	es, err = elasticsearch.NewClient(config)
	if err != nil {
		return es, err
	}

	es.Ping.WithHuman()
	l.Info().Msgf("elasticsearch connection established with " + config.Addresses[0])

	return es, nil
}

func getESCACert(path string) (cert []byte, err error) {

	_, err = os.Stat(path)
	if err != nil {
		return nil, err
	}
	cert, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return cert, nil
}
