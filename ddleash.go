package ddleash

import (
	"net/http"
	"net/http/cookiejar"
)

type Config struct {
	Team     string
	Username string
	Password string
}

type DDLeash struct {
	cookieJar http.CookieJar
	client    *http.Client
}

type Metric struct {
	Name string
}

func New(config Config) (*DDLeash, error) {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &DDLeash{
		cookieJar: cookieJar,
		client: &http.Client{
			Jar: cookieJar,
		},
	}, nil
}

func (leash *DDLeash) Login() error {
	return nil
}

func (leash *DDLeash) FetchAllMetrics() ([]Metric, error) {
	return []Metric{
		Metric{Name: "foo"},
		Metric{Name: "bar"},
		Metric{Name: "baz"},
	}, nil
}
