package plugins

import (
	"net/http"
	"net/url"
	"path"
	"time"
)

type Headers map[string][]interface{}

type Pocs struct {
	Name         string
	Target       string
	Port         int
	Cookie       string
	Timeout      time.Duration
	Proxy        func(*http.Request) (*url.URL, error)
	Charset      string
	Chunked      bool
	PluginsInfos PluginsInfos
	PocResult    PocResult
}

type PocResult struct {
	Status  bool
	Message string
}

func (p *Pocs) UrlJoin(raw string, paths string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	u.Path = path.Join(u.Path, paths)
	return u.String()
}

func (p *Pocs) AppendUri(raw string, paths string) string {
	return path.Join(raw, paths)
}

func (p *Pocs) GetHostname(target string) string {
	raw, err := url.Parse(target)
	if err != nil {
		return ""
	}
	return raw.Hostname()
}

func (p *Pocs) InitHeaders() (headers Headers) {
	return
}

func (p *Pocs) GET(target string) {}

func (p *Pocs) POST(target string, postdata string, headers Headers) {}

func (p *Pocs) PUT(target string, putdata string, headers Headers) {}
