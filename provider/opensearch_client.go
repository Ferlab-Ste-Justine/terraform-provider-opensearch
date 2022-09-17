package provider

import (
	"errors"
	"fmt"
	"io/ioutil"
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
	return ((*reqCon).CurrentEndpoint + 1) == len((*(*reqCon).Client).Endpoints)
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

	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	res, resErr := (*(*reqCon).Client).Client.Do(req)
	if resErr != nil {
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

	if res.StatusCode >= 400 {
		defer res.Body.Close()

		if (*reqCon).RetriesLeft == 0 {
			b, bErr := ioutil.ReadAll(res.Body)
			var errMsg string
			if bErr != nil {
				errMsg = "Could not extract error code response body"
			} else {
				errMsg = string(b)
			}

			errMsg = fmt.Sprintf("Request return code %d: %s", res.StatusCode, errMsg)
			return res, errors.New(errMsg)
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