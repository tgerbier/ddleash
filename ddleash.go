package ddleash

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
)

type Config struct {
	Team     string
	Username string
	Password string
}

type DDLeash struct {
	cookieJar   http.CookieJar
	client      *http.Client
	hasLoggedIn bool
}

type Metric struct {
	Name string
}

var (
	ErrNotLoggedIn = errors.New("DDLeash not logged in")
)

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
		hasLoggedIn: false,
	}, nil
}

func (leash *DDLeash) Login() error {
	leash.hasLoggedIn = true
	return nil
}

func (leash *DDLeash) FetchAllMetrics() ([]Metric, error) {
	if !leash.hasLoggedIn {
		return nil, ErrNotLoggedIn
	}

	return []Metric{
		Metric{Name: "foo"},
		Metric{Name: "bar"},
		Metric{Name: "baz"},
	}, nil
}
