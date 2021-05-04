package persistent

import (
	"bytes"
	"context"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ESWriter struct {
	es    *elasticsearch.Client
	index string
}

func NewESWriter(host, login, pwd, index string) (*ESWriter, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			host,
		},
		Username: login,
		Password: pwd,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	log.Println(elasticsearch.Version)
	log.Println(es.Info())

	return &ESWriter{
		es:    es,
		index: index,
	}, nil
}

func (w *ESWriter) Write(jsonData []byte, recordId string) {
	if err := w.saveToES(jsonData, recordId); err != nil {
		log.Printf("Failed to load data to ES with err: %s", err)
	}
}

func (w *ESWriter) saveToES(jsonData []byte, recordId string) error {
	ctx := context.Background()

	req := esapi.IndexRequest{
		Index:      w.index,
		DocumentID: recordId,
		Body:       bytes.NewReader(jsonData),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, w.es)
	if err != nil || res.StatusCode >= 300 || res.StatusCode < 200 {
		return err
	}

	log.Printf("Record with id \"%s\" send to ES.\n", recordId)
	return nil
}
