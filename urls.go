package ddleash

import (
	"fmt"
	"net/url"
	"strconv"
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

func urlForMetricList(team string, window int) *url.URL {
	baseUrl := urlForRoot(team)
	baseUrl.Path = "/metric/list"
	baseUrl.RawQuery = url.Values{
		"window": {strconv.Itoa(window)},
	}.Encode()

	return baseUrl
}

func urlForMetric(team string, name string) *url.URL {
	baseUrl := urlForRoot(team)
	baseUrl.Path = "/metric/metric_metadata"
	baseUrl.RawQuery = url.Values{
		"metrics[]": {name},
	}.Encode()

	return baseUrl
}

func urlForMetricHostsTags(team string, name string, window int) *url.URL {
	baseUrl := urlForRoot(team)
	baseUrl.Path = "/metric/hosts_and_tags"
	baseUrl.RawQuery = url.Values{
		"metric": {name},
		"window": {strconv.Itoa(window)},
	}.Encode()

	return baseUrl
}
