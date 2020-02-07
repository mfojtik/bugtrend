package bugzilla

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/ddliu/go-httpclient"
)

type Client struct {
	apiKey string

	client *httpclient.HttpClient
}

func NewClient(apiKey string) Client {
	return Client{
		apiKey: apiKey,
		client: httpclient.Defaults(httpclient.Map{
			"Accept-Language": "en-us",
			"Accept":          "application/json",
			"Content-Type":    "application/json",
			"Host":            "bugzilla.redhat.com",
		}),
	}
}

func (c Client) withDefaultValues(options url.Values) url.Values {
	newOptions := options
	options["api_key"] = []string{c.apiKey}
	return newOptions
}

func (c Client) Search(values url.Values) (*SearchResult, error) {
	response, err := c.client.Get("https://bugzilla.redhat.com/rest/bug", c.withDefaultValues(values))
	if err != nil {
		return nil, err
	}
	responseBytes, err := response.ReadAll()
	if err != nil {
		return nil, err
	}
	ioutil.WriteFile("get.json", responseBytes, os.ModePerm)
	var result SearchResult
	if err := json.Unmarshal(responseBytes, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
