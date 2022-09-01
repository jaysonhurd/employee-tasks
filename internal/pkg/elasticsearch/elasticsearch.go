package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog"
)

func NewESConnection(config elasticsearch.Config, l *zerolog.Logger) (elasticsearch.Client, error) {
	es, err := elasticsearch.NewClient(config)
	if err != nil {
		return *es, err
	}
	es.Ping.WithHuman()
	l.Info().Msgf("elasticsearch connection established with " + config.Addresses[0])

	return *es, nil
}
