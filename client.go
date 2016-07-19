package ddleash

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Account struct {
	Team     string
	User     string
	Password string
}

type Client struct {
	account    Account
	cookieJar  http.CookieJar
	httpClient *http.Client

	isLoggedIn bool
}

var (
	ErrNotLoggedIn           = errors.New("DDLeash not logged in")
	ErrDogweblCookieNotFound = errors.New("dogwebl cookie not found")
)

func New(account Account) *Client {
	cookieJar, _ := cookiejar.New(nil)

	return &Client{
		account:   account,
		cookieJar: cookieJar,
		httpClient: &http.Client{
			Jar: cookieJar,
		},
		isLoggedIn: false,
	}
}

func (c *Client) Login() error {
	dogwebl, err := c.fetchDogwebl()
	if err != nil {
		return err
	}

	form := url.Values{
		"username":              {c.account.User},
		"password":              {c.account.Password},
		"_authentication_token": {dogwebl},
	}
	loginUrl := urlForLogin(c.account.Team).String()
	resp, err := c.httpClient.PostForm(loginUrl, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf(
			"Invalid response for request %V with params %V: %V",
			loginUrl, form, resp,
		)
	}

	c.isLoggedIn = true
	return nil
}

func (c *Client) FetchAllMetricNames(window int) ([]string, error) {
	if !c.isLoggedIn {
		return nil, ErrNotLoggedIn
	}

	// Fetch all metric names
	metricListUrl := urlForMetricList(
		c.account.Team, window,
	).String()
	resp, err := c.httpClient.Get(metricListUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(
			"Invalid response for request %V: %V",
			metricListUrl, resp,
		)
	}

	// Decode the response
	var jsonResp struct{ Metrics []string }
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		return nil, err
	}

	return jsonResp.Metrics, nil
}

func (c *Client) FetchMetric(name string) (*Metric, error) {
	if !c.isLoggedIn {
		return nil, ErrNotLoggedIn
	}

	// Fetch the metric
	metricUrl := urlForMetric(
		c.account.Team, name,
	).String()
	resp, err := c.httpClient.Get(metricUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(
			"Invalid response for request %V: %V",
			metricUrl, resp,
		)
	}

	// Decode the response
	var jsonResp map[string]*Metric
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		return nil, err
	}

	jsonResp[name].Name = name
	return jsonResp[name], nil
}

func (c *Client) FetchMetricHostsTags(name string, window int) (*MetricHostsTags, error) {
	if !c.isLoggedIn {
		return nil, ErrNotLoggedIn
	}

	// Fetch the metric
	metricUrl := urlForMetricHostsTags(
		c.account.Team, name, window,
	).String()
	resp, err := c.httpClient.Get(metricUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(
			"Invalid response for request %V: %V",
			metricUrl, resp,
		)
	}

	// Decode the response
	var hostsTags MetricHostsTags
	if err := json.NewDecoder(resp.Body).Decode(&hostsTags); err != nil {
		return nil, err
	}

	return &hostsTags, nil
}

func (c *Client) fetchDogwebl() (string, error) {
	loginUrl := urlForLogin(c.account.Team).String()
	resp, err := c.httpClient.Get(loginUrl)
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

	ddCookies := c.cookieJar.Cookies(urlForRoot(c.account.Team))

	for _, ddCookie := range ddCookies {
		if ddCookie.Name == "dogwebl" {
			return ddCookie.Value, nil
		}
	}
	return "", ErrDogweblCookieNotFound
}
