package main

import (
	"encoding/json"
	"fmt"
	"github.com/jaysonhurd/employee-tasks/internal/pkg/elasticsearch"
	"github.com/jaysonhurd/employee-tasks/internal/pkg/postgres"
	"github.com/jaysonhurd/employee-tasks/pkg/tasks"
	"github.com/jaysonhurd/employee-tasks/pkg/tasks/models"
	"github.com/rs/zerolog"
	"io"
	"io/ioutil"
	"os"
	"time"
)

var (
	DBs postgres.Postgreser
)

const (
	// Making these constants for the purposes of the exercise.  These would either be configs passed in
	// via config.json or using flags if this were a production ready service.
	configLocation = "config/config.json"
	logFile        = "logs/employee-tasks.log"
)

func main() {
	// Set up logger
	logger, err := loggerSetup()

	// Load Config from file (config/config.json usually)
	config, err := loadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Set up PostgresConn connection
	dbconn, err := postgres.NewPostgresConnection(config.Postgres, &logger)
	if err != nil {
		logger.Fatal().Msgf("unable to connect to postgres")
		return
	}
	dbconn.SqlClient.Ping()

	DBs := postgres.New(dbconn, &logger)
	if err != nil {
		return
	}

	// Set up ElasticSearch connection
	es, err := elasticsearch.NewESConnection(config.ElasticConfig, &logger)
	if err != nil {
		logger.Fatal().Msgf(err.Error())
	}
	es.Ping.WithHuman()
	info := es.Info
	fmt.Println(info.WithHuman())
	ES := elasticsearch.New(es, DBs, &logger)

	err = ES.EmptyES()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = ES.LoadTasksFromPostgres()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = ES.LoadEmployeesFromPostgres()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	routes := tasks.New(DBs, ES, *config, &logger)

	err = routes.MakeRoutes()
	if err != nil {
		logger.Fatal().Msgf(err.Error())
		return
	}

}

// loggerSetup - could be more elaborate.  Basic logging for the purposes of the exercise.
func loggerSetup() (zerolog.Logger, error) {
	var err error
	var file *os.File
	var logger zerolog.Logger
	{
		file, err = os.OpenFile(
			logFile,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
		)
		if err != nil {
			panic(err)
		}
		logger = zerolog.New(io.MultiWriter(os.Stdout, file)).With().Timestamp().Logger()
		zerolog.TimeFieldFormat = time.RFC3339
		zerolog.TimestampFieldName = "timestamp"
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return logger, nil
}

// loadConfig - loads your config/config.json file
func loadConfig() (*models.Config, error) {
	var c *models.Config

	_, err := os.Stat(configLocation)
	if err != nil {
		return nil, err
	}
	f, err := ioutil.ReadFile(configLocation)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(f, &c)
	if err != nil {
		return nil, err
	}

	return c, nil

}
