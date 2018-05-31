package bower

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Registry struct {
	BaseURL *url.URL
}

var DefaultRegistry = Registry{BaseURL: &url.URL{Scheme: "https", Host: "bower.herokuapp.com"}}

type LookupResponse struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (r Registry) Lookup(pkg string) (*LookupResponse, error) {
	url := r.BaseURL.ResolveReference(&url.URL{Path: "/packages/" + url.QueryEscape(pkg)})
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	var lr *LookupResponse
	err = json.NewDecoder(resp.Body).Decode(&lr)
	if err != nil {
		return nil, err
	}

	return lr, nil
}
