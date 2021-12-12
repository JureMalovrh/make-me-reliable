package http

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// GetCaller is a definition of HTTP caller that will use Get method for obtaining results
type GetCaller struct {
	Auth string
	URL  string
	c    *http.Client
}

// NewGetCaller returns new GetCaller
func NewGetCaller(auth, url string) *GetCaller {
	c := &http.Client{
		Timeout: 2 * time.Second,
	}

	if !strings.HasPrefix(url, "http") {
		url = fmt.Sprintf("http://%s", url)
	}

	return &GetCaller{
		Auth: auth,
		URL:  url,
		c:    c,
	}
}

// CallAPI will call provided path and params using GET request
func (g *GetCaller) CallAPI(ctx context.Context, path string, params map[string]string) (interface{}, error) {
	url := fmt.Sprintf("%s/%s", g.URL, path)

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	q := r.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}

	r.URL.RawQuery = q.Encode()

	r.Header.Add("Authorization", g.Auth)
	res, err := g.c.Do(r)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusOK {
		return string(body), nil
	}

	log.WithFields(log.Fields{
		"status": res.StatusCode,
		"body":   string(body),
	}).Error("Call to API failed")

	return nil, fmt.Errorf("API call failed")
}
