package ddleash

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
)

type Config struct {
	Team     string
	Username string
	Password string
}

type DDLeash struct {
	config    Config
	cookieJar http.CookieJar
	client    *http.Client

	hasLoggedIn bool
}

type Metric struct {
	Name string
}

var (

	ErrNotLoggedIn           = errors.New("DDLeash not logged in")
	ErrDogweblCookieNotFound = errors.New("dogwebl cookie not found")
)

func urlForRoot(team string) *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s.datadoghq.com", team),
		Path:   "/",
	}
}

func urlForLogin(team string) *url.URL {
	baseUrl := urlForRoot(team)
	baseUrl.Path = "/account/login"
	baseUrl.RawQuery = url.Values{
		"redirect": {"f"},
	}.Encode()

	return baseUrl
}

func New(config Config) (*DDLeash, error) {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &DDLeash{
		config:    config,
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

func (leash *DDLeash) fetchDogwebl() (string, error) {
	loginUrl := urlForLogin(leash.config.Team).String()
	resp, err := leash.client.Get(loginUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf(
			"Invalid response for request %V: %V",
			loginUrl, resp,
		)
	}

	ddCookies := leash.cookieJar.Cookies(urlForRoot(leash.config.Team))

	for _, ddCookie := range ddCookies {
		if ddCookie.Name == "dogwebl" {
			return ddCookie.Value, nil
		}
	}
	return "", ErrDogweblCookieNotFound
}
