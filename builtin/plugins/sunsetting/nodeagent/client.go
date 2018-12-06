package nodeagent

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Client struct {
	logger     logrus.FieldLogger
	httpClient *http.Client
}

type Instance struct {
	URI string
}

func NewClient(logger logrus.FieldLogger) (*Client, error) {
	var client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1,
		},
		Timeout: time.Second * 30,
	}

	return &Client{
		logger:     logger,
		httpClient: client,
	}, nil
}

func (c *Client) Get(nodeAgent Instance, path string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, nodeAgent.URI+path, nil)
	if err != nil {
		return "", err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(rb), nil
}
