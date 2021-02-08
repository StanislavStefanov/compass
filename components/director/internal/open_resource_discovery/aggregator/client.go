package aggregator

import (
	"context"
	"encoding/json"
	"github.com/kyma-incubator/compass/components/director/internal/open_resource_discovery"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type client struct {
	http.Client
}

func NewClient(timeout time.Duration) *client {
	return &client{
		Client: http.Client{
			Timeout: timeout,
		},
	}
}

func (c *client) FetchOpenDiscoveryDocuments(ctx context.Context, url string) ([]*open_resource_discovery.Document, error) {
	resp, err := c.Get(url + open_resource_discovery.WellKnownEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "error while fetching open resource discovery well-known configuration")
	}

	defer closeBody(ctx, resp.Body)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body")
	}

	config := open_resource_discovery.WellKnownConfig{}
	if err := json.Unmarshal(bodyBytes, &config); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling json body")
	}

	docs := make([]*open_resource_discovery.Document, 0, 0)
	for _, config := range config.OpenResourceDiscoveryV1.Documents {
		doc, err := c.FetchOpenDiscoveryDocument(ctx, url+config.URL)
		if err != nil {
			return nil, errors.Wrapf(err, "error fetching ORD document from: %s", url+config.URL)
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

func (c *client) FetchOpenDiscoveryDocument(ctx context.Context, documentURL string) (*open_resource_discovery.Document, error) {
	resp, err := c.Get(documentURL)
	if err != nil {
		return nil, err
	}
	defer closeBody(ctx, resp.Body)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading document body")
	}
	result := &open_resource_discovery.Document{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling document")
	}
	return result, nil
}

func closeBody(ctx context.Context, body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		log.C(ctx).WithError(err).Warnf("Got error on closing response body")
	}
}