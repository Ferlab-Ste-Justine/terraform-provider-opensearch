package provider

import (
	"net/http"
	"net/url"
	"path"
	"strings"
)

type OpensearchClient struct {
	Client    *http.Client
	Endpoints []string
	Username  string
	Password  string
	Retries   int
}

func (cli *OpensearchClient) GetRequestContext() *RequestContext {
	return &RequestContext{
		Client: cli,
		CurrentEndpoint: 0,
		RetriesLeft: (*cli).Retries,
	}
}

type RequestContext struct {
	Client          *OpensearchClient
	CurrentEndpoint int
	RetriesLeft     int
}

func (reqCon *RequestContext) GetCurrentEndpoint() string {
	return (*(*reqCon).Client).Endpoints[(*reqCon).CurrentEndpoint]
}

func (reqCon *RequestContext) AtLastEndpoint() bool {
	return (*reqCon).CurrentEndpoint >= len((*(*reqCon).Client).Endpoints)
}

func (reqCon *RequestContext) Do(method string, urlPath string, body string) (*http.Response, error) {
	endpoint := reqCon.GetCurrentEndpoint()
	u, uErr := url.Parse(endpoint)
	if uErr != nil {
		return nil, uErr
	}
	u.Path = path.Join(u.Path, urlPath)

	r := strings.NewReader(body)
	req, reqErr := http.NewRequest(method, u.String(), r)
    if reqErr != nil {
		return nil, reqErr
	}

	if (*(*reqCon).Client).Username != "" {
		req.SetBasicAuth ((*(*reqCon).Client).Username, (*(*reqCon).Client).Password)
	}

	res, resErr := (*(*reqCon).Client).Client.Do(req)
	if resErr != nil {
		defer res.Body.Close()
		if (*reqCon).RetriesLeft == 0 {
			return res, resErr	
		}

		(*reqCon).RetriesLeft -= 1

		if reqCon.AtLastEndpoint() {
			(*reqCon).CurrentEndpoint = 0
		} else {
			(*reqCon).CurrentEndpoint += 1
		}
		
		return reqCon.Do(method, urlPath, body)
	}

	return res, nil
}